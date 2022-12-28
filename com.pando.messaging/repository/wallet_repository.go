package repository

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/smtp"
	"net/url"
	"os"
	"pandoMessagingWalletService/com.pando.messaging/config"
	"pandoMessagingWalletService/com.pando.messaging/config/kafka"
	constants "pandoMessagingWalletService/com.pando.messaging/constants"
	"pandoMessagingWalletService/com.pando.messaging/enums"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"

	"gorm.io/gorm"
)

var WalletStatementFileUrl, LogoUrl string

type walletRepository struct {
	DBConn *gorm.DB
}

func NewWalletRepository(conn *gorm.DB, conf *config.Config) WalletRepository {
	WalletStatementFileUrl = conf.WalletStatementFileUrl
	LogoUrl = conf.LogoUrl
	return &walletRepository{
		DBConn: conn,
	}
	// SendTokenToArtist := func() {
	// 	wh.SendTokenToArtist()
	// }
	// if _, err := scheduler.Every(1).Minutes().Run(SendTokenToArtist); err != nil {
	// 	logrus.Error("Error while starting scheduler")
	// }
	//return wh
}

/*****************************************Create Wallet**********************************/
func (r *walletRepository) CreateWallet(ctx context.Context, password models.WalletReq) (*models.Response, error) {
	u := models.User{}
	user := r.DBConn.Where("id = ?", password.UserId).First(&u)
	if user.RowsAffected == 0 {
		logger.Logger.WithError(user.Error).WithField("error", user.Error).Error("Please enter correct user id.")
		return &models.Response{Status: false, Msg: "Please enter correct user id."}, nil
	}
	check := r.DBConn.Where("user_id = ?", password.UserId).First(&models.WalletDetails{})
	if check.RowsAffected != 0 {
		logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("You can not create more than one wallet.")
		return &models.Response{Status: false, Msg: "You can not create more than one wallet."}, nil
	}
	value, _ := json.Marshal(password)
	req, _ := http.NewRequest("POST", WalletApiUrl+"/api/new-key", bytes.NewBuffer(value))
	req.Header.Set(constants.ContentType, constants.ApplicationJson)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error(constants.WalletNotCreated)
		return nil, nil
	}
	var result models.Createwallet
	data, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(data, &result)
	if !result.Success {
		logger.Logger.Error(constants.WalletNotCreated)
		return &models.Response{Status: false, Msg: constants.WalletNotCreated}, nil
	}
	wallet := models.WalletDetails{
		WalletId: result.Data.Result.Address,
		UserId:   password.UserId,
		Username: u.Username,
		Balance:  "0",
	}
	db := r.DBConn.Create(&wallet)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error(constants.WalletNotCreated)
		return &models.Response{Status: false, Msg: constants.WalletNotCreated}, nil
	}
	from := SmtpSenderEmail
	toList := []string{u.Email}
	msg := []byte("To: " + u.Email + "\r\n" +
		"Subject: Wallet creation\r\n\r\n" + "Dear " + u.Username + "\nYour wallet is created successfully. The details are below : \nWalletId: " + wallet.WalletId + "\nPassword: " + password.Password + "\n \n*NOTE:- Please save this password for further operation and do not share it with others and please save it at your side because you can not change or update your password.")
	auth := smtp.PlainAuth("", SmtpUser, SmtpPass, SmtpHost)
	fmt.Println(SmtpHost, SmtpPass, SmtpSenderEmail, SmtpUser, u.Email)
	go smtp.SendMail(SmtpHost+":587", auth, from, toList, msg)
	return &models.Response{Status: true, Msg: "Wallet is created successfully.", WalletAddress: result.Data.Result.Address}, nil
}

/******************************************Get Balance****************************************/
func (r *walletRepository) GetBalance(ctx context.Context, walletId string) (*models.Response, error) {
	balreq := models.GetBalanceReq{
		ID:      1,
		Jsonrpc: "2.0",
		Method:  "pando.GetAccount",
		Param: []models.Param{
			{Address: walletId},
		},
	}
	value, _ := json.Marshal(balreq)
	req, _ := http.NewRequest("POST", WalletApiUrl+"/api/node", bytes.NewBuffer(value))
	req.Header.Set(constants.ContentType, constants.ApplicationJson)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error("Balance is not found.")
		return &models.Response{Status: false, Msg: "Error while fetching balance from blockchain."}, nil
	}
	var result models.Createwallet
	data, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(data, &result)
	oldBal, _ := strconv.ParseFloat(result.Data.Result.Coins.PTXWei, 64)
	newBal := float64(oldBal) / float64(1000000000000000000)
	b := math.Floor(newBal*100) / 100
	bal := fmt.Sprintf("%g", b)
	update := r.DBConn.Where("wallet_id = ?", walletId).Find(&models.WalletDetails{}).Update("balance", bal)
	if update.Error != nil {
		logger.Logger.WithError(err).WithField("error", err).Error("Balance is not updated.")
		return &models.Response{Status: false, Msg: "Wallet balance is not updated."}, nil
	}
	if !result.Success {
		return &models.Response{Status: true, Msg: "Wallet balance", Balance: "0"}, nil
	}
	return &models.Response{Status: true, Msg: "Wallet balance", Balance: bal}, nil
}

/***********************************************Add Token***************************************/
func (r *walletRepository) AddToken(ctx context.Context, walletId string, amount string) (*models.Response, error) {
	value, _ := json.Marshal(map[string]string{"to": walletId, "amount": amount})
	req, _ := http.NewRequest("POST", WalletApiUrl+"/api/transfer", bytes.NewBuffer(value))
	req.Header.Set(constants.ContentType, constants.ApplicationJson)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error("Token is not added to wallet.")
		return &models.Response{Status: false, Msg: "Token not added to wallet."}, nil
	}
	var result models.Createwallet
	data, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(data, &result)
	if !result.Success {
		logger.Logger.Error("Transaction is failed and token not added to wallet.")
		return &models.Response{Status: false, Msg: "Token not added to wallet."}, nil
	}
	db := r.DBConn.Create(&models.Transactions{
		CreatedAt:           time.Now().UTC(),
		SenderWalletId:      result.Data.Block.Proposer,
		ReceiverWalletId:    walletId,
		Amount:              amount,
		Message:             "Token added",
		Time:                strconv.Itoa(int(time.Now().Unix())),
		Status:              enums.Confirmed,
		TransactionHash:     result.Data.Result.Hash,
		TransactionCategory: enums.Received,
	})
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Transaction is not created.")
		return &models.Response{Status: false, Msg: "Transaction is not created."}, nil
	}
	return &models.Response{Status: true, Msg: "Token added successfully.", TransactionHash: result.Data.Result.Hash}, nil
}

/**********************************************Request Token*************************************/
func (r *walletRepository) RequestToken(ctx context.Context, req models.RequestTokens) (*models.Response, error) {
	reqToken := models.RequestTokens{}
	users := []models.User{}
	realAmount, _ := strconv.ParseFloat(req.Amount, 64)
	if realAmount == 0 {
		logger.Logger.Error("Amount is less than zero.")
		return &models.Response{Status: false, Msg: "Amount should be greater than zero."}, nil
	}
	check := r.DBConn.Where("requested_by_wallet_id = ? and requested_from_wallet_id = ? and request_status = 'PENDING'", req.RequestedByWalletId, req.RequestedFromWalletId).Find(&reqToken)
	if check.RowsAffected != 0 {
		logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("Request is sent already.")
		return &models.Response{Status: false, Msg: "You have already requested token from this user."}, nil
	}
	user := r.DBConn.Where("id in (?)", []int64{req.RequestedByUserId, req.RequestedFromUserId}).Find(&users)
	if user.RowsAffected < 2 {
		logger.Logger.WithError(user.Error).WithField("error", user.Error).Error("Users not found.")
		return &models.Response{Status: false, Msg: "User you have selected not found."}, nil
	}
	var sender, receiver, profile_pic_url string
	if users[0].ID == req.RequestedByUserId {
		receiver = users[1].Username
		sender = users[0].Username
		profile_pic_url = users[0].Profile_Pic_Url

	} else {
		receiver = users[0].Username
		sender = users[1].Username
		profile_pic_url = users[1].Profile_Pic_Url

	}
	db := r.DBConn.Create(&models.RequestTokens{
		RequestedFromWalletId: req.RequestedFromWalletId,
		RequestedFromUserId:   req.RequestedFromUserId,
		RequestedByUserId:     req.RequestedByUserId,
		RequestedByWalletId:   req.RequestedByWalletId,
		RequestStatus:         "PENDING",
		Amount:                req.Amount,
		Message:               req.Message,
	})
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Request has not sent.")
		return &models.Response{Status: false, Msg: "Request has not sent to the user."}, nil
	}
	amount, _ := strconv.ParseFloat(req.Amount, 64)
	var Message string
	if amount > 1 {
		Message = sender + " has requested " + req.Amount + " tokens from you."
	} else {
		Message = sender + " has requested " + req.Amount + " token from you."
	}
	message := models.Notifications{
		ReceiverUserId:    req.RequestedFromUserId,
		SenderUserId:      req.RequestedByUserId,
		SenderUsername:    sender,
		ReceiverUsername:  receiver,
		Message:           Message,
		NotificationTitle: enums.TokenRequest,
		NotificationType:  constants.RequestToken,
		Profile_Pic_Url:   profile_pic_url,
	}
	data, _ := json.Marshal(&message)
	go kafka.Push(context.Background(), nil, data)
	return &models.Response{Status: true, Msg: "Request has sent successfully."}, nil
}

/*****************************************Reject Request*************************************/
func (r *walletRepository) RejectRequest(ctx context.Context, requestId int64, requestType string) (*models.Response, error) {
	req := &models.RequestTokens{}
	if requestType == "ACCEPT" {
		db := r.DBConn.Where("id = ?", requestId).Find(&req).Update("request_status", enums.ACCEPTED)
		if db.Error != nil {
			logger.Logger.WithError(db.Error).WithField("error", db.Error).Error(constants.RequestStatusNotUpdated)
			return &models.Response{Status: false, Msg: constants.RequestStatusNotUpdated}, nil
		}
		return &models.Response{Status: true, Msg: "Request has accepted."}, nil
	} else if requestType == "REJECT" {
		db := r.DBConn.Where("id = ?", requestId).Find(&req).Update("request_status", enums.REJECTED)
		if db.Error != nil {
			logger.Logger.WithError(db.Error).WithField("error", db.Error).Error(constants.RequestStatusNotUpdated)
			return &models.Response{Status: false, Msg: constants.RequestStatusNotUpdated}, nil
		}
		return &models.Response{Status: true, Msg: "Request has rejected."}, nil
	} else {
		return &models.Response{Status: false, Msg: "Please use request type either ACCEPT or REJECT."}, nil
	}
}

/*****************************************SendToken********************************************/
func (r *walletRepository) SendToken(ctx context.Context, from string, to string, amount string, password string, message string) (*models.Response, error) {
	wallet := []models.WalletDetails{}
	users := models.User{}
	//first we unlock the wallet
	value, _ := json.Marshal(map[string]string{"address": from, "password": password})
	req, _ := http.NewRequest("POST", WalletApiUrl+"/api/unlock-wallet", bytes.NewBuffer(value))
	req.Header.Set(constants.ContentType, constants.ApplicationJson)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error(constants.ErrorWalletUnlocking)
		return &models.Response{Status: false, Msg: constants.ErrorWalletUnlocking}, nil
	}
	var result models.Createwallet
	data, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(data, &result)
	if result.Data.Error.Code != 0 {
		logger.Logger.Error(constants.WrongWalletPassword)
		return &models.Response{Status: false, Msg: constants.WrongWalletPassword}, nil
	}
	//Now we send the token to user
	amount1, _ := strconv.ParseFloat(amount, 64)
	amountInt := amount1 * 1000000000000000000
	newAmount := fmt.Sprintf("%f", amountInt)
	realAmount1 := strings.Split(newAmount, ".")
	sendTokenReq := models.GetBalanceReq{
		ID:      1,
		Jsonrpc: "2.0",
		Method:  "pandocli.Send",
		Param: []models.Param{
			{ChainId: "pandonet", From: from, To: to, PandoWei: "0", PTXWei: realAmount1[0], Fee: "10000000000000000", Async: true},
		},
	}
	value, _ = json.Marshal(sendTokenReq)
	req, _ = http.NewRequest("POST", WalletApiUrl+"/api/send", bytes.NewBuffer(value))
	req.Header.Set(constants.ContentType, constants.ApplicationJson)
	client = &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error(constants.SendTokenErrorUser)
		return &models.Response{Status: false, Msg: constants.SendTokenErrorUser}, nil
	}
	//var results models.CreateWallet
	data, _ = ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(data, &result)
	insufficient := strings.Contains(result.Message, "Insufficient")
	if insufficient {
		logger.Logger.Error("Insufficient balance.")
		return &models.Response{Status: false, Msg: "You have insufficient balance."}, nil
	}
	if !result.Success || result.Data.Result.Hash == "" {
		logger.Logger.Info(result.Message, result.Data.Result.Hash)

		logger.Logger.Error(constants.SendTokenErrorUser)
		return &models.Response{Status: false, Msg: constants.SendTokenErrorUser}, nil
	}

	//Now again lock the wallet
	value, _ = json.Marshal(map[string]string{"address": from})
	req, _ = http.NewRequest("POST", WalletApiUrl+"/api/lock-wallet", bytes.NewBuffer(value))
	req.Header.Set(constants.ContentType, constants.ApplicationJson)
	client = &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error(constants.ErrorWalletLocking)
		return &models.Response{Status: false, Msg: constants.ErrorWalletLocking}, nil
	}
	//var result models.Createwallet
	data, _ = ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(data, &result)
	if !result.Success {
		logger.Logger.Error(constants.WalletNotLocked)
		return &models.Response{Status: false, Msg: constants.WalletNotLocked}, nil
	}
	check := r.DBConn.Where("wallet_id in (?)", []string{from, to}).Find(&wallet)
	if check.RowsAffected == 0 {
		logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("Wallet details is not found.")
		return &models.Response{Status: false, Msg: "Wallet details is not found."}, nil
	}

	var sender, receiver string
	var sender_userid, receiver_userid int64
	if wallet[0].WalletId == from {
		sender = wallet[0].Username
		receiver = wallet[1].Username
		sender_userid = wallet[0].UserId
		receiver_userid = wallet[1].UserId
	} else {
		sender = wallet[1].Username
		receiver = wallet[0].Username
		sender_userid = wallet[1].UserId
		receiver_userid = wallet[0].UserId
	}
	user := r.DBConn.Where("id = ?", sender_userid).Find(&users)
	if user.RowsAffected == 0 {
		logger.Logger.WithError(user.Error).WithField("error", user.Error).Error("Your details have not found.")
		return &models.Response{Status: false, Msg: "Your details have not found."}, nil
	}
	db := r.DBConn.Create(&models.Transactions{
		CreatedAt:           time.Now().UTC(),
		SenderWalletId:      from,
		SenderUsername:      sender,
		ReceiverWalletId:    to,
		ReceiverUsername:    receiver,
		Amount:              amount,
		Message:             message,
		Time:                strconv.Itoa(int(time.Now().Unix())),
		Status:              enums.Confirmed,
		TransactionHash:     result.Data.Result.Hash,
		TransactionCategory: enums.Sent,
	})
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error(constants.TransactionNotSaved)
		return &models.Response{Status: false, Msg: constants.TransactionNotSaved}, nil
	}
	Amount, _ := strconv.ParseFloat(amount, 64)
	var Message string
	if Amount > 1 {
		Message = sender + " has sent you " + amount + " tokens."
	} else {
		Message = sender + " has sent you " + amount + " token."
	}
	messages := models.Notifications{
		ReceiverUserId:    receiver_userid,
		SenderUserId:      sender_userid,
		SenderUsername:    sender,
		ReceiverUsername:  receiver,
		Message:           Message,
		NotificationTitle: enums.TokenReceived,
		NotificationType:  constants.WalletNotificationType,
		Profile_Pic_Url:   users.Profile_Pic_Url,
	}
	data1, _ := json.Marshal(&messages)
	go kafka.Push(context.Background(), nil, data1)
	return &models.Response{Status: true, Msg: "Token has sent successfully.", TransactionHash: result.Data.Result.Hash}, nil
}

/***************************************Get Transactions******************************************/
func (r *walletRepository) GetTransactions(ctx context.Context, walletId string, query url.Values) (*models.Response, error) {
	t := []models.Transactions{}
	re := []models.TransactionsDto{}
	var location string
	userLocation := r.DBConn.Table("users as u").Select("u.time_zones").Joins("join wallet_details as w on w.user_id=u.id").Where("w.wallet_id=?", walletId).Find(&location)
	if userLocation.Error != nil {
		logger.Logger.WithError(userLocation.Error).WithField("error", userLocation.Error).Error("Location is not found.")
		return &models.Response{Status: true, Msg: "Location is not found."}, nil
	}
	pagination := config.GeneratePaginationFromRequest(query)
	offset := (pagination.Page) * pagination.Limit
	queryBuider := r.DBConn.Limit(pagination.Limit).Offset(offset)
	db := queryBuider.Where("LOWER(sender_wallet_id) = LOWER(?) or LOWER(receiver_wallet_id) = LOWER(?) order by time desc", walletId, walletId).Find(&t)
	go func() {
		req, _ := http.NewRequest("GET", WalletApiUrl+"/api/get-history?wallet="+walletId, nil)
		req.Header.Set(constants.ContentType, constants.ApplicationJson)
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			logger.Logger.WithError(err).WithField("error", err).Error("Token is not added to account.")
			return
		}
		var result models.Createwallet
		data, _ := ioutil.ReadAll(res.Body)
		_ = json.Unmarshal(data, &result)
		for i := 0; i < len(result.Data.Data); i++ {
			tx := models.Transactions{}
			db := r.DBConn.Where("transaction_hash = ?", result.Data.Data[i].Hash).Find(&tx)
			if db.RowsAffected != 0 {
				if tx.Status != result.Data.Data[i].Status {
					update := r.DBConn.Where("transaction_hash = ?", result.Data.Data[i].Hash).Find(&tx).Update("status", result.Data.Data[i].Status)
					if update.Error != nil {
						logger.Logger.WithError(update.Error).WithField("error", update.Error).Error("Status is not updated for transaction.")
						return
					}
				}
			} else {
				amount := result.Data.Data[i].Data.Output[0].Coins.PTXWei
				amount1, _ := strconv.ParseFloat(amount, 64)
				amountInt := (amount1) / 1000000000000000000
				newAmount := fmt.Sprintf("%v", amountInt)
				db := r.DBConn.Create(&models.Transactions{
					CreatedAt:           time.Now().UTC(),
					SenderWalletId:      result.Data.Data[i].Data.Input[0].Address,
					SenderUsername:      constants.ExternalAccount,
					ReceiverWalletId:    result.Data.Data[i].Data.Output[0].Address,
					ReceiverUsername:    "Unknown",
					Amount:              newAmount,
					Message:             "",
					Status:              result.Data.Data[i].Status,
					TransactionHash:     result.Data.Data[i].Hash,
					TransactionCategory: enums.TokenTransfer,
					Time:                result.Data.Data[i].Timestamp,
				})
				if db.Error != nil {
					logger.Logger.WithError(db.Error).WithField("error", db.Error).Error(constants.TransactionNotSaved)
					return
				}
			}
		}
	}()
	if db.RowsAffected == 0 {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Transactions are not found.")
		return &models.Response{Status: true, Msg: "Transactions are not found."}, nil
	}
	for i := 0; i < len(t); i++ {
		amount1, _ := strconv.ParseFloat(t[i].Amount, 64)
		amountInt := math.Floor(amount1*100) / 100
		newAmount := fmt.Sprintf("%v", amountInt)
		if walletId == t[i].SenderWalletId {
			result := models.TransactionsDto{
				Date:            config.TimeConvert("", t[i].CreatedAt).Format("Jan 2, 2006 15:04:05"),
				Amount:          newAmount,
				Username:        t[i].ReceiverUsername,
				TransactionHash: t[i].TransactionHash,
				WalletId:        t[i].ReceiverWalletId,
				Type:            "DR",
			}
			re = append(re, result)
		} else {
			result := models.TransactionsDto{
				Date:            config.TimeConvert("", t[i].CreatedAt).Format("Jan 2, 2006 15:04:05"),
				Amount:          newAmount,
				Username:        t[i].SenderUsername,
				TransactionHash: t[i].TransactionHash,
				WalletId:        t[i].SenderWalletId,
				Type:            "CR",
			}
			re = append(re, result)
		}

	}
	return &models.Response{Status: true, Msg: "Transactions are found.", Transactions: &re}, nil
}

/***********************************************View Spend Analytics******************************/
func (r *walletRepository) ViewSpendAnalytics(ctx context.Context, walletId string) (*models.Response, error) {
	t := []models.SpendAnalyticsReq{}
	results := []models.SpendAnalytics{}
	SpendLists := make([]models.SpendList, 0)
	month := time.Now().Month()
	currMon := int(month)
	total := 0.0
	i := 0
LOOP:
	for i < 6 {
		db := r.DBConn.Raw("select transaction_category as title ,round(CAST(sum(amount::float) as numeric),2) as amount from transactions where sender_wallet_id = ? and extract(month from created_at) = ? and transaction_category in('SENT','BOOKING','CONCERT CREATION') group by transaction_category", walletId, currMon).Find(&t)
		if db.RowsAffected == 0 {
			i = i + 1
			currMon = currMon - 1
			goto LOOP
		}
		for j := 0; j < len(t); j++ {
			SpendList := models.SpendList{
				Title:  t[j].Title,
				Amount: t[j].Amount,
			}
			SpendLists = append(SpendLists, SpendList)
			total = math.Floor((total+t[j].Amount)*100) / 100
		}
		result := models.SpendAnalytics{
			Month:     time.Month(currMon).String(),
			Total:     total,
			SpendList: SpendLists,
		}
		results = append(results, result)
		currMon = currMon - 1
		SpendLists = make([]models.SpendList, 0)
		total = 0
		i++
		if currMon == 0 {
			currMon = 12
		}
	}
	if len(results) == 0 {
		logger.Logger.Error("Have no data to show.")
		return &models.Response{Status: true, Msg: "Have no data to show."}, nil
	}
	return &models.Response{Status: true, Msg: "Data Found.", Data: &results}, nil
}

/************************************Wallet statement*****************************************/
func (r *walletRepository) WalletStatement(ctx context.Context, walletId string, startDate string, endDate string, totalMonths string, email string, queryType string) (string, error) {
	trans := []models.Transactions{}
	var from_date, to_date, location string
	userLocation := r.DBConn.Table("users as u").Select("u.time_zones").Joins("join wallet_details as w on w.user_id=u.id").Where("w.wallet_id=?", walletId).Find(&location)
	if userLocation.Error != nil {
		logger.Logger.WithError(userLocation.Error).WithField("error", userLocation.Error).Error("Location is not found.")
		return "", userLocation.Error
	}
	switch totalMonths {
	case "1":
		{
			start := time.Now().AddDate(0, 0, 1)
			from_date = start.AddDate(0, -1, -1).String()
			to_date = start.String()
		}
	case "3":
		{
			start := time.Now().AddDate(0, 0, 1)
			from_date = start.AddDate(0, -3, -1).String()
			to_date = start.String()
		}
	case "6":
		{
			start := time.Now().AddDate(0, 0, 1)
			from_date = start.AddDate(0, -6, -1).String()
			to_date = start.String()
		}
	case "12":
		{
			start := time.Now().AddDate(0, 0, 1)
			from_date = start.AddDate(-1, 0, -1).String()
			to_date = start.String()
		}
	default:
		{
			from_date = startDate
			layOut := "2006-01-02"
			timeStamp, _ := time.Parse(layOut, endDate)
			to_date = timeStamp.AddDate(0, 0, 1).String()

		}
	}
	fromdate := strings.SplitAfter(from_date, " ")
	todate := strings.SplitAfter(to_date, " ")
	db := r.DBConn.Where("(created_at  BETWEEN SYMMETRIC ? AND ?) AND (sender_wallet_id = ? or receiver_wallet_id = ?) order by created_at desc", fromdate[0], todate[0], walletId, walletId).Find(&trans)
	if db.RowsAffected == 0 {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Transactions not found.")
		return "", nil
	}
	var username string
	data := [][]string{{"Date", "Username", "WalletAddress", "Txhash", "Amount"}}
	for i := 0; i < len(trans); i++ {
		var data1 []string
		if walletId == trans[i].SenderWalletId {
			username = trans[i].SenderUsername
			data1 = []string{config.TimeConvert(location, trans[i].CreatedAt).Format(constants.DateFormat), trans[i].ReceiverUsername, trans[i].ReceiverWalletId, trans[i].TransactionHash, trans[i].Amount + " DR"}
		} else {
			username = trans[i].ReceiverUsername
			data1 = []string{config.TimeConvert(location, trans[i].CreatedAt).Format(constants.DateFormat), trans[i].SenderUsername, trans[i].SenderWalletId, trans[i].TransactionHash, trans[i].Amount + " CR"}
		}
		data = append(data, data1)
	}
	metadata := strconv.Itoa(time.Now().Nanosecond())
	filename := "wallet_statement" + metadata + ".csv"
	file, err := os.Create(filename)
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error("CSV file has not created.")
		return "CSV file has not created", err
	}
	defer file.Close()
	defer os.Remove(WalletStatementFileUrl + filename)
	writer := csv.NewWriter(file)
	defer writer.Flush()

	result := writer.WriteAll(data)
	if result != nil {
		logger.Logger.WithError(result).WithField("error", result).Error("Error while writing data to csv file.")
		return "Error while writing data to csv file", err
	}
	datap := loadCSV(WalletStatementFileUrl + filename)
	pass := config.RandomNumber(6)
	// Then we create a new PDF document and write the title and the current date.
	pdf := r.newReport(walletId, fromdate[0], todate[0], queryType, pass, username)

	pdf = image(pdf)
	if pdf.Err() {
		logger.Logger.WithError(pdf.Error()).WithField("error", pdf.Error()).Error("Error while adding image.")
		return "Error while adding image.", pdf.Error()
	}
	// After that, we create the table header and fill the table.
	pdf = header(pdf, datap[0])
	pdf = table(pdf, datap[1:])
	pdf.SetFont("Times", "", 12)
	pdf.CellFormat(60, 10, "* This is a system generated statement. Hence, it does not require any signature.", "", 0, "L", false, 0, "")
	page := strconv.Itoa(pdf.PageCount())
	pdf.SetLeftMargin(245)
	pdf.CellFormat(40, 10, "current page:", "", 0, "R", false, 0, "")
	pdf.Cell(40, 10, page)
	// And we should take the opportunity and beef up our report with a nice logo.
	filename1 := "wallet_statement" + metadata + ".pdf"
	// And finally, we write out our finished record to a file.
	err = savePDF(pdf, filename1)
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error("Error while saving pdf.")
		return "Error while saving pdf.", err
	}
	if queryType == "EMAIL" {
		message := "Dear Customer!\n\nThis is your wallet statement from " + fromdate[0] + " to " + fromdate[0] + " and this is an encrypted pdf you have to enter password to open it.\nYour password for this pdf is " + pass + " . Please do not disclose it with others."
		subject := "Wallet statement from " + fromdate[0] + " to " + todate[0]
		config.SendMail(WalletStatementFileUrl+filename1, email, message, subject)
		defer os.Remove(WalletStatementFileUrl + filename1)
		return "Mail sent on your registered email.", nil
	} else if queryType == "DOWNLOAD" {
		return WalletStatementFileUrl + filename1, nil
	}
	return "Please enter correct query_type", nil
}

/****************************************Functions used for statement****************************/
func loadCSV(path string) [][]string {
	f, err := os.Open(path)
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error("Error while opening csv file.")
	}
	defer f.Close()
	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error("Error while reading csv file.")
	}
	return rows
}
func (r *walletRepository) newReport(walletId, fromDate, toDate, queryType, pass, username string) *gofpdf.Fpdf {
	bal, _ := r.GetBalance(context.Background(), walletId)
	pdf := gofpdf.New("L", "mm", "A4", "")
	if queryType == "EMAIL" {
		pdf.SetProtection(gofpdf.CnProtectPrint, "123", pass)
	}
	pdf.AddPage()

	pdf.SetFont("Times", "B", 22)
	pdf.SetTopMargin(30)
	pdf.Ln(30)
	pdf.CellFormat(40, 10, "Wallet Statement", "", 0, "L", false, 0, "")
	pdf.Ln(10)
	pdf.SetFont("Times", "B", 15)
	pdf.Cell(40, 10, "Wallet id:")
	pdf.Cell(40, 10, walletId)
	pdf.SetLeftMargin(215)
	pdf.CellFormat(40, 10, "Date:", "", 0, "R", false, 0, "")
	pdf.Cell(40, 10, time.Now().Format("Jan 2, 2006"))
	pdf.SetLeftMargin(10)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Username:")
	pdf.Cell(40, 10, username)
	pdf.Ln(10)
	pdf.Cell(40, 10, "Balance:")
	pdf.Cell(40, 10, bal.Balance+" Pando")
	pdf.Ln(10)
	pdf.Cell(40, 10, "Transaction date:")
	pdf.Cell(40, 10, "  From  "+fromDate+"  to  "+toDate)
	pdf.SetLeftMargin(5)
	pdf.Ln(20)
	pdf.SetLeftMargin(5)

	return pdf
}

func header(pdf *gofpdf.Fpdf, hdr []string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "B", 14)
	pdf.SetTopMargin(20)
	pdf.SetFillColor(240, 240, 240)
	width := []float64{20, 37, 90, 120, 20}
	for i, str := range hdr {
		pdf.CellFormat(width[i], 7, str, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)
	return pdf
}

func table(pdf *gofpdf.Fpdf, tbl [][]string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "", 10)
	pdf.SetFillColor(255, 255, 255)
	align := []string{"C", "C", "C", "C", "C"}
	width := []float64{20, 37, 90, 120, 20}
	for _, line := range tbl {
		for i, str := range line {
			pdf.CellFormat(width[i], 7, str, "1", 0, align[i], false, 0, "")
		}
		pdf.Ln(-1)
	}
	return pdf
}
func image(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	pdf.ImageOptions(LogoUrl, 100, 10, 60, 20, false, gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}, 0, "")
	return pdf
}

func savePDF(pdf *gofpdf.Fpdf, filename string) error {
	return pdf.OutputFileAndClose(filename)
}

/*******************************************Get Wallet Id************************************/
func (r *walletRepository) GetWalletId(ctx context.Context, userId int64) (*models.Response, error) {
	wallet := models.WalletDto{}
	db := r.DBConn.Raw("select u.username,u.profile_pic_url,w.wallet_id from users as u join wallet_details as w on u.id=w.user_id where u.id = ?", userId).First(&wallet)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Wallet id is not found.")
		return &models.Response{Status: true, Msg: "Wallet id is not found."}, nil
	}
	result := models.UserDto{
		ID:              userId,
		Username:        wallet.Username,
		Profile_Pic_Url: wallet.Profile_Pic_Url,
		WalletId:        wallet.WalletId,
	}
	return &models.Response{Status: true, Msg: "Wallet id found.", UserDetail: &result}, nil
}

/*****************************************Recent Transactions********************************/
func (r *walletRepository) RecentTransactions(ctx context.Context, walletId string, query url.Values) (*models.Response, error) {
	trans := []models.Transactions{}
	var walletIds []string
	pagination := config.GeneratePaginationFromRequest(query)
	offset := (pagination.Page) * pagination.Limit
	queryBuider := r.DBConn.Limit(pagination.Limit).Offset(offset)
	db := queryBuider.Where("sender_wallet_id = ? or receiver_wallet_id = ? and transaction_category in ('SENT','RECEIVED') order by created_at desc", walletId, walletId).Find(&trans)
	if db.RowsAffected == 0 {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Recent transactions are not found.")
		return &models.Response{Status: true, Msg: "Recent transactions are not found."}, nil
	}
	for i := 0; i < len(trans); i++ {
		if trans[i].SenderWalletId != walletId {
			walletIds = append(walletIds, trans[i].SenderWalletId)
		} else if trans[i].ReceiverWalletId != walletId {
			walletIds = append(walletIds, trans[i].ReceiverWalletId)
		}
	}
	User := []models.UserDto{}
	user := r.DBConn.Raw("select u.id,u.username,u.profile_pic_url,w.wallet_id from users as u left join wallet_details as w on w.user_id = u.id where u.id in(select user_id from wallet_details where wallet_id in (?))", walletIds).Find(&User)
	if user.RowsAffected == 0 {
		logger.Logger.WithError(user.Error).WithField("error", user.Error).Error("User info not found.")
		return &models.Response{Status: false, Msg: "User information not found."}, nil
	}
	sort.Slice(User, func(i, j int) bool {
		return trans[i].CreatedAt.String() < trans[j].CreatedAt.String()
	})
	return &models.Response{Status: true, Msg: "Recent transactions have found.", Users: &User}, nil
}

/*******************************************ShowPendingRequests**********************************/
func (r *walletRepository) ShowPendingRequests(ctx context.Context, walletId string, userId int64, query url.Values) (*models.Response, error) {
	req := []models.RequestTokenDto{}
	pagination := config.GeneratePaginationFromRequest(query)
	offset := (pagination.Page) * pagination.Limit
	queryBuider := r.DBConn.Limit(pagination.Limit).Offset(offset)
	db := queryBuider.Table("request_tokens as r").Select("r.id,r.created_at,request_status,requested_from_user_id,requested_by_user_id,requested_from_wallet_id,requested_by_wallet_id,amount, username,profile_pic_url").Joins("left join users on users.id=r.requested_by_user_id").Where("(requested_from_wallet_id = ? and requested_from_user_id = ?) and request_status = 'PENDING'", walletId, userId).Find(&req)
	if db.RowsAffected == 0 {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Request is not found.")
		return &models.Response{Status: true, Msg: "You don't have any token requests."}, nil
	}
	return &models.Response{Status: true, Msg: "Token requests have found.", RequestToken: &req}, nil
}

/****************************************SendTokenToAdmin*********************************/
func (r *walletRepository) SendTokenToAdmin(ctx context.Context, id, concertId int64, amount, password string) (*models.Response, error) {
	wallet := models.UserWalletInfo{}
	db := r.DBConn.Table("wallet_details").Select("username,wallet_id").Where("user_id = ?", id).Find(&wallet)
	if db.RowsAffected == 0 {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Wallet id not found.")
		return &models.Response{Status: false, Msg: "Wallet id not found."}, nil
	}
	value, _ := json.Marshal(map[string]string{"address": wallet.WalletId, "password": password})
	req, _ := http.NewRequest("POST", WalletApiUrl+"/api/unlock-wallet", bytes.NewBuffer(value))
	req.Header.Set(constants.ContentType, constants.ApplicationJson)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error(constants.ErrorWalletUnlocking)
		return &models.Response{Status: false, Msg: constants.ErrorWalletUnlocking}, nil
	}
	var result models.Createwallet
	data, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(data, &result)
	if result.Data.Error.Code != 0 {
		logger.Logger.Error(constants.WrongWalletPassword)
		return &models.Response{Status: false, Msg: constants.WrongWalletPassword}, nil
	}
	//Now we send the token to user

	amount1, _ := strconv.ParseFloat(amount, 64)
	amountInt := amount1 * 1000000000000000000
	newAmount := fmt.Sprintf("%f", amountInt)
	realAmount := strings.Split(newAmount, ".")
	sendTokenReq := models.GetBalanceReq{
		ID:      1,
		Jsonrpc: "2.0",
		Method:  "pandocli.Send",
		Param: []models.Param{
			{ChainId: "pandonet", From: wallet.WalletId, To: AdminWalletId, PandoWei: "0", PTXWei: realAmount[0], Fee: "10000000000000000", Async: true},
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
	//var result models.Createwallet
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
	value, _ = json.Marshal(map[string]string{"address": wallet.WalletId})
	req, _ = http.NewRequest("POST", WalletApiUrl+"/api/lock-wallet", bytes.NewBuffer(value))
	req.Header.Set(constants.ContentType, constants.ApplicationJson)
	client = &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error(constants.ErrorWalletLocking)
		return &models.Response{Status: false, Msg: constants.ErrorWalletLocking}, nil
	}
	//var result models.Createwallet
	data, _ = ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(data, &result)
	if !result.Success {
		logger.Logger.Error(constants.WalletNotLocked)
		return &models.Response{Status: false, Msg: constants.WalletNotLocked}, nil
	}
	db1 := r.DBConn.Create(&models.Transactions{
		CreatedAt:           time.Now().UTC(),
		SenderWalletId:      wallet.WalletId,
		SenderUsername:      wallet.Username,
		ReceiverWalletId:    AdminWalletId,
		ReceiverUsername:    constants.VixoAdminAccount,
		Amount:              amount,
		Message:             "",
		Time:                strconv.Itoa(int(time.Now().Unix())),
		Status:              enums.Confirmed,
		TransactionHash:     result.Data.Result.Hash,
		TransactionCategory: enums.ConcertCreation,
	})
	if db1.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error(constants.TransactionNotSaved)
		return &models.Response{Status: false, Msg: constants.TransactionNotSaved}, nil
	}
	updateConcert := r.DBConn.Where("id=?", concertId).First(&models.Concerts{}).Update("status", "EDIT_MODE")
	if updateConcert.Error != nil {
		logger.Logger.WithError(updateConcert.Error).WithField("error", updateConcert.Error).Error("concert status has not updated.")
		return &models.Response{Status: false, Msg: "Concert status has not updated."}, nil
	}
	return &models.Response{Status: true, Msg: "Token has sent successfully."}, nil
}

/*******************************************ShowOwnTokenRequests**********************************/
func (r *walletRepository) ShowOwnTokenRequests(ctx context.Context, walletId string, userId int64, query url.Values) (*models.Response, error) {
	req := []models.RequestTokenDto{}
	pagination := config.GeneratePaginationFromRequest(query)
	offset := (pagination.Page) * pagination.Limit
	queryBuider := r.DBConn.Limit(pagination.Limit).Offset(offset)
	db := queryBuider.Table("request_tokens as r").Select("r.id,r.created_at,request_status,requested_from_user_id,requested_by_user_id,requested_from_wallet_id,requested_by_wallet_id,amount, username,profile_pic_url").Joins("left join users on users.id=r.requested_from_user_id").Where("requested_by_wallet_id = ? and requested_by_user_id = ?", walletId, userId).Find(&req)
	if db.RowsAffected == 0 {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Request is not found.")
		return &models.Response{Status: true, Msg: "You have not requested any token."}, nil
	}
	return &models.Response{Status: true, Msg: "Token requests have found.", RequestToken: &req}, nil
}
