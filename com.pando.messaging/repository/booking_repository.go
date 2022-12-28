package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/smtp"
	"pandoMessagingWalletService/com.pando.messaging/config"
	"pandoMessagingWalletService/com.pando.messaging/config/kafka"
	constants "pandoMessagingWalletService/com.pando.messaging/constants"
	"pandoMessagingWalletService/com.pando.messaging/enums"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/carlescere/scheduler"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type bookingRepository struct {
	DBConn *gorm.DB
}

var AdminWalletId, SmtpHost, SmtpPass, SmtpSenderEmail, SmtpUser, TicketTemplatePath, WalletApiUrl string

func NewBookingRepository(conn *gorm.DB, conf *config.Config) BookingRepository {
	AdminWalletId = conf.AdminWalletId
	SmtpHost = conf.SmtpHost
	SmtpPass = conf.SmtpPass
	SmtpSenderEmail = conf.SmtpSenderEmail
	SmtpUser = conf.SmtpUser
	TicketTemplatePath = conf.TicketTemplatePath
	WalletApiUrl = conf.WalletApiUrl
	wh := &bookingRepository{
		DBConn: conn,
	}
	sendTokenToArtist := func() {
		wh.SendTokenToArtist()
		wh.SendTokenToUserAfterTicketCancel()
	}
	if _, err := scheduler.Every().Day().At("23:00:00").Run(sendTokenToArtist); err != nil {
		logrus.Error("Error while starting scheduler")
	}
	return wh
}

/*****************************************BookTicket**************************************/
func (r *bookingRepository) BookTicket(ctx context.Context, bookTicket models.BookingDetailDto, message string) (*models.Response, error) {
	concert := models.ConcertInfo{}
	user := models.UserInfo{}
	//find the user info
	check := r.DBConn.Table("concerts").Select("title,concert_date,ticket_price").Where("id = ?", bookTicket.ConcertID).Find(&concert)
	if check.RowsAffected == 0 {
		logger.Logger.WithError(check.Error).WithField("error", check.Error).Error(constants.ConcertDetailError)
		return &models.Response{Status: false, Msg: constants.ConcertDetailError}, nil
	}
	userInfo := r.DBConn.Table("users as u").Select("u.username,u.email,u.profile_pic_url,w.wallet_id").Joins("left join wallet_details as w on u.id=w.user_id").Where("u.id = ?", bookTicket.UserId).Find(&user)
	if userInfo.RowsAffected == 0 {
		logger.Logger.WithError(userInfo.Error).WithField("error", userInfo.Error).Error("User details is not found.")
		return &models.Response{Status: false, Msg: "User details is not found."}, nil
	}
	//first we unlock the wallet
	value, _ := json.Marshal(map[string]string{"address": user.WalletId, "password": bookTicket.WalletPassword})
	req, _ := http.NewRequest("POST", WalletApiUrl+"/api/unlock-wallet", bytes.NewBuffer(value))
	req.Header.Set(constants.ContentType, constants.ApplicationJson)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error("Error while unlocking the wallet.")
		return &models.Response{Status: false, Msg: "Error while unlocking the wallet."}, nil
	}
	var result models.Createwallet
	data, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(data, &result)
	if result.Data.Error.Code != 0 {
		logger.Logger.Error("Wallet has not unlocked please check your password.")
		return &models.Response{Status: false, Msg: "Wallet has not unlocked please check your password."}, nil
	}
	//Now we send the token to user
	amount := concert.TicketPrice * float64(bookTicket.TotalTickets)
	amountInt := amount * 1000000000000000000
	newAmount := fmt.Sprintf("%f", amountInt)
	realAmount := strings.Split(newAmount, ".")
	sendTokenReq := models.GetBalanceReq{
		ID:      1,
		Jsonrpc: "2.0",
		Method:  "pandocli.Send",
		Param: []models.Param{
			{ChainId: "pandonet", From: user.WalletId, To: AdminWalletId, PandoWei: "0", PTXWei: realAmount[0], Fee: "10000000000000000", Async: true},
		},
	}
	value, _ = json.Marshal(sendTokenReq)
	req, _ = http.NewRequest("POST", WalletApiUrl+"/api/send", bytes.NewBuffer(value))
	req.Header.Set(constants.ContentType, constants.ApplicationJson)
	client = &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error(constants.SendTokenError)
		return &models.Response{Status: false, Msg: constants.SendTokenError}, nil
	}
	data, _ = ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(data, &result)
	insufficient := strings.Contains(result.Message, "Insufficient")
	if insufficient {
		logger.Logger.Error("Insufficient balance.")
		return &models.Response{Status: false, Msg: "You have insufficient balance."}, nil
	}
	if !result.Success || result.Data.Result.Hash == "" {
		logger.Logger.Error(constants.SendTokenError)
		return &models.Response{Status: false, Msg: constants.SendTokenError}, nil
	}
	//Now again lock the wallet
	value, _ = json.Marshal(map[string]string{"address": user.WalletId})
	req, _ = http.NewRequest("POST", WalletApiUrl+"/api/lock-wallet", bytes.NewBuffer(value))
	req.Header.Set(constants.ContentType, constants.ApplicationJson)
	client = &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error("Error while locking the wallet.")
		return &models.Response{Status: false, Msg: "Error while locking the wallet."}, nil
	}
	//var result models.Createwallet
	data, _ = ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(data, &result)
	if !result.Success {
		logger.Logger.Error("Wallet is not locked.")
		return &models.Response{Status: false, Msg: "Wallet has not locked."}, nil
	}
	var tickets []string
	for i := 0; i < int(bookTicket.TotalTickets); i++ {
		ticketCode := config.RandomNumber(10)
		db := r.DBConn.Create(&models.BookingDetails{
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UserId:      bookTicket.UserId,
			TicketCode:  ticketCode,
			ConcertID:   bookTicket.ConcertID,
			TicketPrice: int64(concert.TicketPrice),
			Status:      enums.Booked,
		})
		if db.Error != nil {
			logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Booking details is not saved.")
			return &models.Response{Status: false, Msg: "Booking details is not saved."}, nil
		}
		tickets = append(tickets, ticketCode)
	}
	db1 := r.DBConn.Create(&models.Transactions{
		CreatedAt:           time.Now().UTC(),
		SenderWalletId:      user.WalletId,
		SenderUsername:      user.Username,
		ReceiverWalletId:    AdminWalletId,
		ReceiverUsername:    constants.VixoAdminAccount,
		Amount:              fmt.Sprintf("%v", amount),
		Message:             "",
		Time:                strconv.Itoa(int(time.Now().Unix())),
		Status:              enums.Confirmed,
		TransactionHash:     result.Data.Result.Hash,
		TransactionCategory: enums.Booking,
	})
	if db1.Error != nil {
		logger.Logger.WithError(db1.Error).WithField("error", db1.Error).Error("Transaction is not saved.")
		return &models.Response{Status: false, Msg: "Transaction is not saved."}, nil
	}

	Notification := models.Notifications{
		ReceiverUserId:    bookTicket.UserId,
		SenderUserId:      0,
		SenderUsername:    constants.AdminUsername,
		ReceiverUsername:  user.Username,
		Message:           "You have booked " + strconv.Itoa(int(bookTicket.TotalTickets)) + " " + message + " of " + concert.Title,
		NotificationTitle: enums.BookingDetails,
		NotificationType:  constants.Ticket,
	}
	data1, _ := json.Marshal(&Notification)
	go kafka.Push(context.Background(), nil, data1)
	go r.SendTicketToMaiL(ctx, tickets, user.Email, bookTicket.ConcertID, bookTicket.TotalTickets)
	return &models.Response{Status: true, Msg: "Ticket is booked successfully."}, nil
}

/*****************************************My Bookings****************************************/
func (r *bookingRepository) SendTicketToMaiL(ctx context.Context, ticketCode []string, email string, concertId, totalTickets int64) {
	concert := models.Concerts{}
	concertDetail := r.DBConn.Table("concerts").Where("id = ?", concertId).First(&concert)
	if concertDetail.Error != nil {
		logger.Logger.WithError(concertDetail.Error).WithField("error", concertDetail.Error).Error(constants.ConcertDetailError)
		//return &models.Response{Status: false, , Msg: "Concert detail is not found."}, nil
	}
	var timezones string
	userLocation := r.DBConn.Table("users").Select("time_zones").Where("email = ?", email).Find(&timezones)
	if userLocation.RowsAffected == 0 {
		logger.Logger.WithError(userLocation.Error).WithField("error", userLocation.Error).Error("User timezone is not found.")
		//return &models.Response{Status: false, , Msg: "Concert detail is not found."}, nil
	}
	var viewType string
	if concert.Viewing_experience == "VIEW_2D" {
		viewType = "2D"
	} else {
		viewType = "VR"
	}
	date := strings.SplitAfter(config.TimeConvert(timezones, concert.ConcertDate).String(), " ")
	From := SmtpSenderEmail
	to := email
	Subject := concert.Title + " ticket details"
	auth := smtp.PlainAuth("", SmtpUser, SmtpPass, SmtpHost)
	t, _ := template.ParseFiles(TicketTemplatePath)
	var body bytes.Buffer
	headers := "MIME-version:1.0;\nContent-Type: text/html;"
	body.Write([]byte(fmt.Sprintf("To: "+email+"\r\n"+
		"Subject:"+Subject+"\n%s\n\n", headers)))

	t.Execute(&body, struct {
		Title              string
		TicketCode         []string
		ArtistName         string
		Date               string
		Time               string
		Language           string
		TotalTickets       int64
		TicketPrice        float64
		Viewing_experience string
	}{
		Title:              concert.Title,
		TicketCode:         ticketCode,
		ArtistName:         concert.ArtistName,
		Date:               date[0],
		Time:               date[1],
		Language:           concert.Language,
		TotalTickets:       totalTickets,
		TicketPrice:        math.Floor(concert.TicketPrice*10000) / 10000,
		Viewing_experience: viewType,
	})
	if err := smtp.SendMail(SmtpHost+":587", auth, From, []string{to}, body.Bytes()); err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error("Mail is not sent")
		//return &models.Response{, Msg: "Mail is not sent."}, nil
	}
	logger.Logger.Info("Mail is sent.")
	//return &models.Response{Status: true, Msg: "Mail sent", }, nil
}

/******************************************SendTokenToArtist*************************************/
func (r *bookingRepository) SendTokenToArtist() {
	var concertId []int64
	check := r.DBConn.Table("concerts").Select("id").Where("concert_date <= now() - interval '48 hours' and status = 'COMPLETED' and payment_status = 'PENDING'").Find(&concertId)
	if check.RowsAffected == 0 {
		logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("Concert is not found.")
	}
	type detail struct {
		username     string
		user_id      int64
		wallet_id    string
		show_type    string
		title        string
		total_amount float64
	}
	walletId := []detail{}
	rows, _ := r.DBConn.Raw("select w.username,w.user_id,w.wallet_id,c.show_type,c.title, sum(b.ticket_price)::float as total_amount from concerts as c left join wallet_details as w on w.user_id = c.user_id left join booking_details as b on b.concert_id=c.id where c.id in(?) group by w.wallet_id,w.username,w.user_id,c.show_type,c.title", concertId).Rows()
	defer rows.Close()
	for rows.Next() {
		f := detail{}
		if err := rows.Scan(&f.username, &f.user_id, &f.wallet_id, &f.show_type, &f.total_amount); err != nil {
			logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("Concert is not found.")
			//return nil, err
		}
		walletId = append(walletId, f)
	}
	for i := 0; i < len(walletId); i++ {
		totalAmount := (walletId[i].total_amount + float64(constants.SecurityMoney)) * 1000000000000000000
		amount := fmt.Sprintf("%f", totalAmount)
		value, _ := json.Marshal(map[string]string{"to": walletId[i].wallet_id, "amount": amount})
		req, _ := http.NewRequest("POST", constants.TokenTransferUrl, bytes.NewBuffer(value))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			logger.Logger.WithError(err).WithField("error", err).Error(constants.TokenTransferError)
		}
		var result models.Createwallet
		data, _ := ioutil.ReadAll(res.Body)
		_ = json.Unmarshal(data, &result)
		if !result.Success {
			logger.Logger.Error(constants.Transactionfailed)
		}
		db := r.DBConn.Create(&models.Transactions{
			CreatedAt:           time.Now().UTC(),
			SenderWalletId:      result.Data.Block.Proposer,
			ReceiverWalletId:    walletId[i].wallet_id,
			Amount:              amount,
			Message:             "Token added",
			Time:                strconv.Itoa(int(time.Now().Unix())),
			Status:              enums.Confirmed,
			TransactionHash:     result.Data.Result.Hash,
			TransactionCategory: enums.Received,
		})
		if db.Error != nil {
			logger.Logger.WithError(db.Error).WithField("error", db.Error).Error(constants.TransactionNotSaved)
		}
		update := r.DBConn.Where("id in (?)", concertId).Find(&[]models.Concerts{}).Update("payment_status", "TRANSFERRED_TO_ARTIST")
		if update.Error != nil {
			logger.Logger.WithError(update.Error).WithField("error", update.Error).Error("Concert payment status not updated.")
		}
		message := models.Notifications{
			ReceiverUserId:    walletId[i].user_id,
			SenderUserId:      0,
			SenderUsername:    constants.AdminUsername,
			ReceiverUsername:  walletId[i].username,
			Message:           "Your " + walletId[i].title + " " + walletId[i].show_type + " payment have been credited",
			NotificationTitle: enums.ConcertPaymentReceived,
			NotificationType:  constants.WalletNotificationType,
		}
		data1, _ := json.Marshal(&message)
		go kafka.Push(context.Background(), nil, data1)
	}

}

/***************************************Refund Token to user******************************************/
func (r *bookingRepository) RefundToken(ctx context.Context, concertId int64) (*models.Response, error) {
	walletInfo := []models.TicketInfo{}
	query := r.DBConn.Table("wallet_details as w").Select("w.wallet_id,w.username,w.user_id,sum(b.ticket_price)::text as ticket_price").Joins("join booking_details as b on b.user_id=w.user_id").Where("b.concert_id = ? and b.status in('BOOKED','SENT') group by w.wallet_id,w.username,w.user_id", concertId).Find(&walletInfo)
	if query.Error != nil {
		logger.Logger.WithError(query.Error).WithField("error", query.Error).Error("Artist's wallet id not found.")
		return &models.Response{Status: false, Msg: "Artist's wallet id not found."}, nil
	}
	var title string
	db := r.DBConn.Table("concerts").Select("title").Where("id=?", concertId).First(&title)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error(constants.ConcertDetailError)
		return &models.Response{Status: false, Msg: constants.ConcertDetailError}, nil
	}
	go func() {
		for i := 0; i < len(walletInfo); i++ {
			value, _ := json.Marshal(map[string]string{"to": walletInfo[i].WalletId, "amount": walletInfo[i].TicketPrice})
			req, _ := http.NewRequest("POST", constants.TokenTransferUrl, bytes.NewBuffer(value))
			req.Header.Set(constants.ContentType, constants.ApplicationJson)
			client := &http.Client{}
			res, _ := client.Do(req)
			var result models.Createwallet
			data, _ := ioutil.ReadAll(res.Body)
			_ = json.Unmarshal(data, &result)
			if !result.Success {
				logger.Logger.Error(constants.Transactionfailed)
			}
			db := r.DBConn.Create(&models.Transactions{
				CreatedAt:           time.Now().UTC(),
				SenderWalletId:      result.Data.Block.Proposer,
				ReceiverWalletId:    walletInfo[i].WalletId,
				Amount:              walletInfo[i].TicketPrice,
				Message:             "Token Refunded",
				Time:                strconv.Itoa(int(time.Now().Unix())),
				Status:              enums.Confirmed,
				TransactionHash:     result.Data.Result.Hash,
				TransactionCategory: enums.Received,
			})
			if db.Error != nil {
				logger.Logger.WithError(db.Error).WithField("error", db.Error).Error(constants.TransactionNotSaved)
			}
			update := r.DBConn.Where("concert_id = ?", concertId).Find(&models.BookingDetails{}).Update("status", "REFUNDED")
			if update.Error != nil {
				logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Concert status is not updated.")
			}
			message := models.Notifications{
				ReceiverUserId:    walletInfo[i].UserId,
				SenderUserId:      0,
				SenderUsername:    constants.AdminUsername,
				ReceiverUsername:  walletInfo[i].Username,
				Message:           "Admin has refunded you token of event " + title,
				NotificationTitle: enums.Refund,
				NotificationType:  constants.WalletNotificationType,
			}
			data1, _ := json.Marshal(&message)
			go kafka.Push(context.Background(), nil, data1)
		}
	}()
	return &models.Response{Status: true, Msg: "Refund has initiated."}, nil
}

/*********************************************ViewTickets***************************************/
func (r *bookingRepository) ViewTickets(ctx context.Context, concertId, userId int64) (*models.Response, error) {
	bookingDetail := []models.TicketDetail{}
	rows, err := r.DBConn.Raw("select distinct b.ticket_code,b.status,s.time_zones, u.username, b.receiver_id,s.username, b.user_id from booking_details as b left join users as u on u.id=b.receiver_id left join users as s on s.id=b.user_id where (b.concert_id = ?) and (b.user_id = ? or b.receiver_id = ?)", concertId, userId, userId).Rows()
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error("Booking details are not found for this concert.")
		return &models.Response{Status: true, Msg: "Booking details are not found for this concert.."}, nil
	}
	defer rows.Close()
	for rows.Next() {
		result := models.TicketDetail{}
		rows.Scan(&result.TicketCode, &result.Status, &result.TimeZone, &result.ReceiverUsername, &result.ReceiverUserId, &result.SenderUsername, &result.SenderUserId)
		if result.Status != "SENT" {
			result = models.TicketDetail{
				TicketCode: result.TicketCode,
				Status:     result.Status,
				TimeZone:   result.TimeZone,
			}
		}
		result = models.TicketDetail{
			TicketCode:       result.TicketCode,
			Status:           result.Status,
			ReceiverUsername: result.ReceiverUsername,
			ReceiverUserId:   result.ReceiverUserId,
			SenderUsername:   result.SenderUsername,
			SenderUserId:     result.SenderUserId,
			TimeZone:         result.TimeZone,
		}
		bookingDetail = append(bookingDetail, result)
	}
	concertDetail := models.Concerts{}
	db := r.DBConn.Table("concerts").Where("id = ?", concertId).First(&concertDetail)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Concert details not found.")
		return &models.Response{Status: false, Msg: "Concert details not found."}, nil
	}
	result := &models.ConcertDto{
		Id:                  concertDetail.Id,
		CreatedAt:           concertDetail.CreatedAt,
		UpdatedAt:           concertDetail.UpdatedAt,
		ArtistName:          concertDetail.ArtistName,
		ConcertDate:         config.TimeConvert(bookingDetail[0].TimeZone, concertDetail.ConcertDate),
		Description:         concertDetail.Description,
		Thumbnail_file_hash: concertDetail.Thumbnail_file_hash,
		TicketPrice:         math.Floor(concertDetail.TicketPrice*10000) / 10000,
		Title:               concertDetail.Title,
		UserId:              concertDetail.UserId,
		Viewing_experience:  concertDetail.Viewing_experience,
		Awss3url:            concertDetail.Awss3url,
		ShowType:            concertDetail.ShowType,
		Language:            concertDetail.Language,
		Status:              concertDetail.Status,
		PaymentStatus:       concertDetail.PaymentStatus,
		TicketDetail:        bookingDetail,
	}
	return &models.Response{Status: true, Msg: "Booking details have found.", ConcertDetails: result}, nil
}

/*****************************************VerifyTicketCode**********************************/
func (r *bookingRepository) VerifyTicketCode(ctx context.Context, userId, concertId int64, ticketCode string) (*models.Response, error) {
	bookingDetail := models.BookingDetails{}
	liveStreamInfo := models.VerifyTicketCodeDto{}
	db := r.DBConn.Where("ticket_code = ? AND concert_id= ?", ticketCode, concertId).First(&bookingDetail)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error(constants.TicketNotFound)
		return &models.Response{Status: false, Msg: constants.TicketNotFound}, nil
	}
	concert := r.DBConn.Table("concerts").Select("artist_name, concert_date,stream_id,live_streaming_url,viewing_experience,liked").Joins("left join concert_like_dislike on concerts.id=concert_like_dislike.concert_id").Where("concerts.id= ?", concertId).Find(&liveStreamInfo)
	if concert.Error != nil {
		logger.Logger.WithError(concert.Error).WithField("error", concert.Error).Error(constants.ConcertDetailError)
		return &models.Response{Status: false, Msg: constants.ConcertDetailError}, nil
	}
	// if time.Now().After(liveStreamInfo.ConcertDate) {
	// 	logger.Logger.Error("Concert date is over.")
	// 	return &models.Response{Status: false, Msg: "Concert is not available now."}, nil
	// }
	if bookingDetail.ViewerId == userId {
		logger.Logger.Info("Ticke is verified with the same user")
		return &models.Response{Status: true, Msg: "Ticket has verified successfully.", LiveStreamInfo: &liveStreamInfo}, nil
	} else if bookingDetail.ViewerId == 0 {
		update := r.DBConn.Where("ticket_code = ? AND concert_id = ?", ticketCode, concertId).Find(&bookingDetail).Updates(map[string]interface{}{"status": "UTILIZED", "viewer_id": userId})
		if update.Error != nil {
			logger.Logger.WithError(update.Error).WithField("error", update.Error).Error("Ticket status is not updated.")
			return &models.Response{Status: false, Msg: constants.TicketStatusNotUpdated}, nil
		}
		return &models.Response{Status: true, Msg: "Ticket has verified successfully.", LiveStreamInfo: &liveStreamInfo}, nil
	} else {
		return &models.Response{Status: true, Msg: "Ticket is already used."}, nil
	}
}

/*****************************************SendTicket**********************************/
func (r *bookingRepository) SendTicket(ctx context.Context, senderUserId, receiverUserId int64, ticketCode string) (*models.Response, error) {
	var title string
	users := []models.User{}
	db := r.DBConn.Table("booking_details as b").Select("title").Joins("left join concerts as c on b.concert_id=c.id").Where("b.ticket_code = ? AND b.user_id= ?", ticketCode, senderUserId).Find(&title)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error(constants.TicketNotFound)
		return &models.Response{Status: false, Msg: constants.TicketNotFound}, nil
	}
	userInfo := r.DBConn.Where("id in (?)", []int64{senderUserId, receiverUserId}).Find(&users)
	if userInfo.RowsAffected < 2 {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Users are not found.")
		return &models.Response{Status: false, Msg: "Please enter valid sender and receiver id."}, nil
	}
	var sender_username, receiver_username string
	if users[0].ID == senderUserId {
		sender_username = users[0].Username
		receiver_username = users[1].Username
	} else {
		receiver_username = users[0].Username
		sender_username = users[1].Username
	}
	update := r.DBConn.Where("user_id = ? and ticket_code = ?", senderUserId, ticketCode).First(&models.BookingDetails{}).Updates(map[string]interface{}{"receiver_id": receiverUserId, "status": enums.Sent, "viewer_id": receiverUserId})
	if update.Error != nil {
		logger.Logger.WithError(update.Error).WithField("error", update.Error).Error("Booking status has not updated.")
		return &models.Response{Status: false, Msg: "Booking status has not updated."}, nil
	}
	message := models.Notifications{
		ReceiverUserId:    receiverUserId,
		SenderUserId:      senderUserId,
		SenderUsername:    sender_username,
		ReceiverUsername:  receiver_username,
		Message:           sender_username + " has sended you " + title + " concert ticket",
		NotificationTitle: enums.ConcertTicketReceived,
		NotificationType:  constants.Ticket,
	}
	data1, _ := json.Marshal(&message)
	go kafka.Push(context.Background(), nil, data1)
	return &models.Response{Status: true, Msg: "You have send ticket to " + receiver_username}, nil

}

/******************************************CancelTicket**************************************/
func (r *bookingRepository) CancelTicket(ctx context.Context, userId, concertId int64, ticketCode string) (*models.Response, error) {
	concertDetail := models.ConcertDetail{}
	bookingDetail := models.BookingDetails{}
	concert := r.DBConn.Table("concerts").Select("title,concert_date,ticket_price").Where("id = ?", concertId).Find(&concertDetail)
	if concert.Error != nil {
		logger.Logger.WithError(concert.Error).WithField("error", concert.Error).Error(constants.ConcertDetailError)
		return &models.Response{Status: false, Msg: constants.ConcertDetailError}, nil
	}
	diff := time.Until(concertDetail.ConcertDate).Hours()
	if diff < 24 {
		logger.Logger.Error("Concert date is less than 24 hours.")
		return &models.Response{Status: false, Msg: "You can't cancel the ticket within 24 hours."}, nil
	}
	userInfo := models.UserInnfo{}
	wallet := r.DBConn.Table("users as u").Select("w.wallet_id,u.username,u.email").Joins("join wallet_details as w on u.id=w.user_id").Where("user_id = ?", userId).Find(&userInfo)
	if wallet.Error != nil {
		logger.Logger.WithError(wallet.Error).WithField("error", wallet.Error).Error("User info not found.")
		return &models.Response{Status: false, Msg: "Your information not found."}, nil
	}
	db := r.DBConn.Where("user_id = ? AND ticket_code = ?", userId, ticketCode).First(&bookingDetail)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Ticket has not found.")
		return &models.Response{Status: false, Msg: "Ticket has not found."}, nil
	}
	if bookingDetail.Status == "CANCELLED" {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Ticket is already cancelled.")
		return &models.Response{Status: true, Msg: "Ticket is already cancelled."}, nil
	}
	update := db.Update("status", "CANCELLED")
	if update.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error(constants.TicketStatusNotUpdated)
		return &models.Response{Status: false, Msg: constants.TicketStatusNotUpdated}, nil
	}

	message := models.Notifications{
		ReceiverUserId:    userId,
		SenderUserId:      0,
		SenderUsername:    constants.VixoAdminAccount,
		ReceiverUsername:  userInfo.Username,
		Message:           "You have cancelled the ticket of " + concertDetail.Title,
		NotificationTitle: constants.TicketCancel,
		NotificationType:  constants.Ticket,
	}
	data, _ := json.Marshal(&message)
	go kafka.Push(context.Background(), nil, data)
	from := SmtpSenderEmail
	password := SmtpPass
	toList := []string{userInfo.Email}
	port := "587"
	msg := "Dear " + userInfo.Username + "\nYour ticket " + ticketCode + " is cancelled of " + concertDetail.Title + " and your token will be refunded within 48 hours."
	body := []byte("To: " + userInfo.Email + "\r\n" +
		"Subject: Ticket cancelled\r\n\r\n" + msg)
	auth := smtp.PlainAuth("", SmtpUser, password, SmtpHost)
	go smtp.SendMail(SmtpHost+":"+port, auth, from, toList, body)
	return &models.Response{Status: true, Msg: "Ticket has cancelled."}, nil
}

/*****************************************SendTokenToUser*************************************/
func (r *bookingRepository) SendTokenToUserAfterTicketCancel() {
	var bookingInfo []models.BookingInfo
	bookingStatus := models.BookingDetails{}
	check := r.DBConn.Table("booking_details as b").Select("b.wallet_id,b.ticket_price::text,w.username,b.user_id,c.title,b.ticket_code").Joins("join wallet_details as w on w.user_id=b.user_id").Joins("join concerts as c on c.id=b.concert_id").Where("b.updated_at <= now() - interval '48 hours' and b.status = 'CANCELLED'").Find(&bookingInfo)
	if check.RowsAffected == 0 {
		logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("Concert is not found.")
	}
	for i := 0; i < len(bookingInfo); i++ {
		value, _ := json.Marshal(map[string]string{"to": bookingInfo[i].WalletId, "amount": bookingInfo[i].TicketPrice})
		req, _ := http.NewRequest("POST", constants.TokenTransferUrl, bytes.NewBuffer(value))
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			logger.Logger.WithError(err).WithField("error", err).Error(constants.TokenTransferError)
			//return &models.Response{Status: false, , Msg: "Token not added to wallet."}, nil
		}
		var result models.Createwallet
		data, _ := ioutil.ReadAll(res.Body)
		_ = json.Unmarshal(data, &result)
		if !result.Success {
			logger.Logger.Error(constants.Transactionfailed)
			//return &models.Response{Status: false, , Msg: "Token not added to wallet."}, nil
		}
		db := r.DBConn.Create(&models.Transactions{
			CreatedAt:           time.Now().UTC(),
			SenderWalletId:      AdminWalletId,
			SenderUsername:      constants.VixoAdminAccount,
			ReceiverUsername:    bookingInfo[i].Username,
			ReceiverWalletId:    bookingInfo[i].WalletId,
			Amount:              bookingInfo[i].TicketPrice,
			Message:             "Concert refund is received.",
			Time:                strconv.Itoa(int(time.Now().Unix())),
			Status:              enums.Confirmed,
			TransactionHash:     result.Data.Result.Hash,
			TransactionCategory: enums.Received,
		})
		if db.Error != nil {
			logger.Logger.WithError(db.Error).WithField("error", db.Error).Error(constants.TransactionNotSaved)
			//return &models.Response{Status: false, , Msg: constants.TransactionNotSaved}, nil
		}
		update := r.DBConn.Where("wallet_id=?", bookingInfo[i].WalletId).First(&bookingStatus).Update("status", enums.Refunded)
		if update.Error != nil {
			logger.Logger.WithError(update.Error).WithField("error", update.Error).Error("Ticket status is not updated ")
		}
		message := models.Notifications{
			ReceiverUserId:    bookingInfo[i].UserId,
			SenderUserId:      0,
			SenderUsername:    constants.VixoAdminAccount,
			ReceiverUsername:  bookingInfo[i].Username,
			Message:           "Your token have refunded for cancelling the ticket " + bookingInfo[i].TicketCode + " of" + bookingInfo[i].Title,
			NotificationTitle: enums.TokenReceivedOnCancellation,
			NotificationType:  constants.WalletNotificationType,
		}
		data1, _ := json.Marshal(&message)
		go kafka.Push(context.Background(), nil, data1)
	}
}


