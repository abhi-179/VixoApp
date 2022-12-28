package repository

import (
	"context"
	"net/http"
	"pandoMessagingWalletService/com.pando.messaging/config"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	"strconv"
	"time"

	"github.com/carlescere/scheduler"
	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

type chatRepository struct {
	DBConn *gorm.DB
}

var BatchSize int
var ReportCount1, ReportCount2, ReportCount3, ReportCount4 int64

func NewchatRepository(conn *gorm.DB, conf *config.Config) ChatRepository {
	BatchSize = conf.Batch_size
	ReportCount1 = conf.ReportCount1
	ReportCount2 = conf.ReportCount2
	ReportCount3 = conf.ReportCount3
	ReportCount4 = conf.ReportCount4
	wh := &chatRepository{
		DBConn: conn,
	}
	// deleteStatusFunc := func() {
	// 	wh.deleteStatusScheduler()
	// }
	updategroupstatus := func() {
		wh.updategroupstatusScheduler()
	}
	if _, err := scheduler.Every(1).Day().Run(updategroupstatus); err != nil {
		logrus.Error("Error while starting scheduler")
	}
	// if _, err := scheduler.Every(1).Minutes().Run(deleteStatusFunc); err != nil {
	// 	logrus.Error("Error while starting scheduler")
	// }
	return wh
}

// func (r *chatRepository) deleteStatusScheduler() {
// 	s := models.Statuses{}
// 	db := r.DBConn.Raw("delete from status where TO_TIMESTAMP(created_at) <= now() - interval '24 hours'").Find(&s)
// 	if db.RowsAffected == 0 {
// 		logrus.Error("Status not deleted.")
// 		logger.Logger.Error("Status not deleted.")
// 	}
// 	logrus.Info("Status deleted.")
// 	logger.Logger.Info("Status deleted after 24 hours.")
// }
func (r *chatRepository) updategroupstatusScheduler() {
	s := models.Groups{}
	u := models.User{}
	db := r.DBConn.Raw("update groups set status = 'ACTIVE' where status = 'TEMPORARY_BLOCKED_1' and  blocked_time  <= now() - interval '24 hours' or status = 'TEMPORARY_BLOCKED_7' and  blocked_time  <= now() - interval '7 days' or status = 'TEMPORARY_BLOCKED_30' and  blocked_time  <= now() - interval '30 days'").Find(&s)
	if db.RowsAffected == 0 {
		logrus.Error("Group status not updated.")
		logger.Logger.Error("Group Status not updated.")
	}
	user := r.DBConn.Raw("update users set status = 'ACTIVE' where status = 'TEMPORARY_BLOCKED_1' and  status_updated_at  <= now() - interval '24 hours' or status = 'TEMPORARY_BLOCKED_7' and  status_updated_at  <= now() - interval '7 days' or status = 'TEMPORARY_BLOCKED_30' and  status_updated_at  <= now() - interval '30 days'").Find(&u)
	if user.RowsAffected == 0 {
		logrus.Error("User status not updated.")
		logger.Logger.Error("User Status not updated.")
	}
	logrus.Info("Group status updated.")
	logger.Logger.Info("Group status updated.")

}

/******************************************Post Status******************************************/
func (r *chatRepository) PostStatus(ctx context.Context, status models.Status) (*models.Response, error) {
	logger.Logger.Info("Enter in to post status repository part.")
	user := models.User{}
	check := r.DBConn.Where("id = ?", status.User_Id).First(&user)
	if check.Error != nil {
		logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("User id is not found.")
		return &models.Response{Status: false, ResponseCode: 404, Msg: "User id is not found."}, nil
	}
	if len(status.File_hash) == 0 {
		db := r.DBConn.Table("status").Create(&models.Statuses{
			ID:         status.ID,
			CreatedAt:  status.CreatedAt,
			User_Id:    status.User_Id,
			Username:   user.Username,
			StatusType: status.StatusType,
			Message:    status.Message,
		})
		if db.Error != nil {
			logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Status is not uploaded.")
			return &models.Response{Status: false, ResponseCode: 400, Msg: "Status is not uploaded."}, nil
		}
		return &models.Response{Status: true, ResponseCode: 200, Msg: "Status is uploaded."}, nil
	} else {
		for i := 0; i < len(status.File_hash); i++ {
			db := r.DBConn.Table("status").Create(&models.Statuses{
				ID:         status.ID,
				CreatedAt:  status.CreatedAt,
				User_Id:    status.User_Id,
				Username:   user.Username,
				StatusType: status.StatusType,
				File_hash:  status.File_hash[i].File_hash,
				AwsUrl:     status.AwsUrl[i].AwsUrl,
			})
			if db.Error != nil {
				logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Status is not uploaded.")
				return &models.Response{Status: false, ResponseCode: 400, Msg: "Status is not uploaded."}, nil
			}
		}
		logger.Logger.Info("Status is uploaded.")
		return &models.Response{Status: true, ResponseCode: 200, Msg: "Status is uploaded."}, nil
	}
}

/************************************************Fetch Status***********************************/
func (r *chatRepository) FetchStatus(ctx context.Context, user_id int64) (*models.Response, error) {
	logger.Logger.Info("Enter in to fetch status repository part.")
	status := []models.Statuses{}
	var UsersId []int64
	check := r.DBConn.Raw("select friend_id1 from friends where friend_id2 = ? and friendship_status = 'ACCEPTED' union select friend_id2 from friends where friend_id1 = ? and friendship_status = 'ACCEPTED'", user_id, user_id).Find(&UsersId)
	if check.Error != nil {
		logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("Friends not found for this user id.")
		return &models.Response{Status: false, ResponseCode: 404, Msg: "Friends not found for this user id."}, nil
	}
	db := r.DBConn.Table("status").Where("user_id = ?", user_id).Find(&status)
	if db.RowsAffected == 0 {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Status not found for this user id.")
		return &models.Response{Status: false, ResponseCode: 404, Msg: "Status not found for this user id."}, nil
	}
	status1 := []models.Statuses{}
	db1 := r.DBConn.Table("status").Where("user_id in (?)", UsersId).Find(&status1)
	if db1.Error != nil {
		logger.Logger.WithError(db1.Error).WithField("error", db1.Error).Error("Status not found for this user id.")
		return &models.Response{Status: false, ResponseCode: 404, Msg: "Status not found for this user id."}, nil
	}
	x := map[int64][]models.Statuses{}
	for i := 0; i < len(status1); i++ {
		value, ok := x[status1[i].User_Id]
		if ok {
			value = append(value, status1[i])
			x[status1[i].User_Id] = value
		} else {
			x[status1[i].User_Id] = []models.Statuses{status1[i]}
		}
	}
	logger.Logger.Info("Status found")
	return &models.Response{Status: true, ResponseCode: 200, Msg: "Status found.", Statues: &status, Statues1: &x}, nil
}

/******************************************Delete Status***************************************/
func (r *chatRepository) DeleteStatus(ctx context.Context, user_id int64, status_id []int64) (*models.Response, error) {
	logger.Logger.Info("Enter in to delete status repository part.")
	status := []models.Statuses{}
	check := r.DBConn.Table("status").Where("user_id = ?", user_id).First(&status)
	if check.Error != nil {
		logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("Status not found for this user id.")
		return &models.Response{Status: false, ResponseCode: 404, Msg: "Status not found for this user id."}, nil
	}
	db := r.DBConn.Table("status").Where("id in (?)", status_id).Find(&status).Delete(&status)
	if db.RowsAffected == 0 {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Status not found.")
		return &models.Response{Status: false, ResponseCode: 404, Msg: "Status not found."}, nil
	}
	logger.Logger.Info("Status deleted successfully.")
	return &models.Response{Status: true, ResponseCode: 200, Msg: "Status deleted successfully."}, nil

}

/******************************************SearchStatusByUsername*******************************/
func (r *chatRepository) SearchStatusByUsername(ctx context.Context, username string) (*models.Response, error) {
	status1 := []models.Statuses{}
	db1 := r.DBConn.Table("status").Where("username ILIKE ?", username+"%").Find(&status1)
	if db1.RowsAffected == 0 {
		logger.Logger.WithError(db1.Error).WithField("error", db1.Error).Error("Status not found for this user username.")
		return &models.Response{Status: false, ResponseCode: 404, Msg: "Status not found for this username."}, nil
	}
	x := map[string][]models.Statuses{}
	for i := 0; i < len(status1); i++ {
		value, ok := x[status1[i].Username]
		if ok {
			value = append(value, status1[i])
			x[status1[i].Username] = value
		} else {
			x[status1[i].Username] = []models.Statuses{status1[i]}
		}
	}
	logger.Logger.Info("Status found")
	return &models.Response{Status: true, ResponseCode: 200, Msg: "Status found.", Status2: &x}, nil
}

/********************************************Report Chat **********************************/
func (r *chatRepository) ReportChat(ctx context.Context, flow models.Reports) (*models.Response, error) {
	logger.Logger.Info("Enter in to report chat repository part.")
	if flow.GroupId == 0 {
		if flow.Reportee_Id == 0 || flow.Reporter_Id == 0 {
			logger.Logger.Error("Reporter id and Reportee id is missing.")
			return &models.Response{Status: false, ResponseCode: http.StatusBadRequest, Msg: "Please enter reporter id and reportee id."}, nil
		}
		user := r.DBConn.Where("id in (?)", []int64{flow.Reporter_Id, flow.Reportee_Id}).Find(&[]models.User{})
		if user.RowsAffected < 2 {
			logger.Logger.WithError(user.Error).WithField("error", user.Error).Error("Please check your reporter or reportee ids.")
			return &models.Response{Status: false, Msg: "Please check your reporter or reportee ids.", ResponseCode: 400}, nil
		}
		check := r.DBConn.Where(models.Reports{Reporter_Id: flow.Reporter_Id, Reportee_Id: flow.Reportee_Id}).FirstOrCreate(&flow)
		if check.Error != nil {
			logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("You have already reported this chat.")
			return &models.Response{Status: false, Msg: "You have already reported this chat.", ResponseCode: 400}, nil
		}
		var count int64
		db := r.DBConn.Model(&[]models.Reports{}).Where("reportee_id = ?", flow.Reportee_Id).Count(&count)
		if db.Error != nil {
			logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("count.")
			return &models.Response{Status: false, Msg: "Count.", ResponseCode: 400}, nil
		}
		switch count {
		case ReportCount1:
			{
				update := r.DBConn.Table("users").Where("id = ?", flow.Reportee_Id).First(&models.User{}).Updates(map[string]interface{}{"status": "TEMPORARY_BLOCKED_1", "status_updated_at": time.Now()})
				if update.RowsAffected == 0 {
					logger.Logger.WithError(update.Error).WithField("error", update.Error).Error("User status is not updated.")
					return &models.Response{Status: false, Msg: "User status is not updated.", ResponseCode: 400}, nil
				}
			}
		case ReportCount2:
			{
				update := r.DBConn.Table("users").Where("id = ?", flow.Reportee_Id).First(&models.User{}).Updates(map[string]interface{}{"status": "TEMPORARY_BLOCKED_7", "status_updated_at": time.Now()})
				if update.RowsAffected == 0 {
					logger.Logger.WithError(update.Error).WithField("error", update.Error).Error("User status is not updated.")
					return &models.Response{Status: false, Msg: "User status is not updated.", ResponseCode: 400}, nil
				}

			}
		case ReportCount3:
			{
				update := r.DBConn.Table("users").Where("id = ?", flow.Reportee_Id).First(&models.User{}).Updates(map[string]interface{}{"status": "TEMPORARY_BLOCKED_30", "status_updated_at": time.Now()})
				if update.RowsAffected == 0 {
					logger.Logger.WithError(update.Error).WithField("error", update.Error).Error("User status is not updated.")
					return &models.Response{Status: false, Msg: "User status is not updated.", ResponseCode: 400}, nil
				}
			}
		case ReportCount4:
			{
				update := r.DBConn.Table("users").Where("id = ?", flow.Reportee_Id).First(&models.User{}).Updates(map[string]interface{}{"status": "INACTIVE", "status_updated_at": time.Now()})
				if update.RowsAffected == 0 {
					logger.Logger.WithError(update.Error).WithField("error", update.Error).Error("User status is not updated.")
					return &models.Response{Status: false, Msg: "User status is not updated.", ResponseCode: 400}, nil
				}
			}
		}
		logger.Logger.Info("Chat Reported.")
		return &models.Response{Status: true, Msg: "Report sent.", ResponseCode: 200}, nil
	} else {
		if flow.GroupId == 0 || flow.Reporter_Id == 0 {
			logger.Logger.Error("Reporter id and Group id is missing.")
			return &models.Response{Status: false, ResponseCode: http.StatusBadRequest, Msg: "Please enter reporter id and group id."}, nil
		}
		user := r.DBConn.Where("id = ?", flow.Reporter_Id).First(&models.User{})
		if user.RowsAffected == 0 {
			logger.Logger.WithError(user.Error).WithField("error", user.Error).Error("Reporter id is not found.")
			return &models.Response{Status: false, Msg: "Reporter id is not found.", ResponseCode: 400}, nil
		}
		group := r.DBConn.Where("id= ?", flow.GroupId).First(&models.Groups{})
		if group.RowsAffected == 0 {
			logger.Logger.WithError(group.Error).WithField("error", group.Error).Error("Group id is not found.")
			return &models.Response{Status: false, Msg: "Group id is not found.", ResponseCode: 400}, nil
		}
		check := r.DBConn.Where(models.Reports{Reporter_Id: flow.Reporter_Id, GroupId: flow.GroupId}).FirstOrCreate(&flow)
		if check.Error != nil {
			logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("You have already reported this chat.")
			return &models.Response{Status: false, Msg: "You have already reported this chat.", ResponseCode: 400}, nil
		}
		var count int64
		db := r.DBConn.Model(&[]models.Reports{}).Where("reportee_id = ?", flow.Reportee_Id).Count(&count)
		if db.Error != nil {
			logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("count.")
			return &models.Response{Status: false, Msg: "Count.", ResponseCode: 400}, nil
		}
		switch count {
		case ReportCount1:
			{
				update := r.DBConn.Table("groups").Where("id = ?", flow.GroupId).First(&models.Groups{}).Updates(map[string]interface{}{"status": "TEMPORARY_BLOCKED_1", "blocked_time": time.Now()})
				if update.RowsAffected == 0 {
					logger.Logger.WithError(update.Error).WithField("error", update.Error).Error("User status is not updated.")
					return &models.Response{Status: false, Msg: "User status is not updated.", ResponseCode: 400}, nil
				}
			}
		case ReportCount2:
			{
				update := r.DBConn.Table("groups").Where("id = ?", flow.GroupId).First(&models.Groups{}).Updates(map[string]interface{}{"status": "TEMPORARY_BLOCKED_7", "blocked_time": time.Now()})
				if update.RowsAffected == 0 {
					logger.Logger.WithError(update.Error).WithField("error", update.Error).Error("User status is not updated.")
					return &models.Response{Status: false, Msg: "User status is not updated.", ResponseCode: 400}, nil
				}
			}
		case ReportCount3:
			{
				update := r.DBConn.Table("groups").Where("id = ?", flow.GroupId).First(&models.Groups{}).Updates(map[string]interface{}{"status": "TEMPORARY_BLOCKED_30", "blocked_time": time.Now()})
				if update.RowsAffected == 0 {
					logger.Logger.WithError(update.Error).WithField("error", update.Error).Error("User status is not updated.")
					return &models.Response{Status: false, Msg: "User status is not updated.", ResponseCode: 400}, nil
				}
			}
		case ReportCount4:
			{
				g := models.Groups{}
				update := r.DBConn.Table("groups").Where("id = ?", flow.GroupId).First(&g)
				if update.RowsAffected == 0 {
					logger.Logger.WithError(update.Error).WithField("error", update.Error).Error("User status is not updated.")
					return &models.Response{Status: false, Msg: "User status is not updated.", ResponseCode: 400}, nil
				}
				group := r.DBConn.Where("id = ?", flow.GroupId).First(&g).Delete(&g)
				if group.Error != nil {
					logger.Logger.WithError(group.Error).WithField("error", group.Error).Error("Group is not deleted.")
					return &models.Response{Status: false, Msg: "Group is not deleted.", ResponseCode: 400}, nil
				}

			}
		}
		logger.Logger.Info("Chat Reported.")
		return &models.Response{Status: true, Msg: "Report sent.", ResponseCode: 200}, nil
	}
}

/***************************************Save Blocked contacts Details***********************/
func (r *chatRepository) SaveBlockUserDetails(ctx context.Context, flow models.BlockedContacts) (*models.Response, error) {

	check := r.DBConn.Table("users").Where("id in (?)", []int64{flow.BlockeeId, flow.BlockerId}).Find(&[]models.User{})
	if check.RowsAffected < 2 {
		logger.Logger.Info("Please enter valid user id.")
		return &models.Response{Status: false, ResponseCode: 400, Msg: "Please enter valid user id."}, nil
	}
	err := r.DBConn.Where("blockee_id = ? AND blocker_id = ?", flow.BlockeeId, flow.BlockerId).Find(&models.BlockedContacts{})
	if err.RowsAffected != 0 {
		logger.Logger.Info("You have already blocked this person.")
		return &models.Response{Status: true, ResponseCode: 200, Msg: "You have aleady blocked this person."}, nil
	}
	db := r.DBConn.Create(&flow)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("User is not blocked.")
		return &models.Response{Status: false, Msg: "User is not blocked.", ResponseCode: 400}, nil
	}
	logger.Logger.Info("Users is blocked successfully.")
	return &models.Response{Status: true, ResponseCode: 200, Msg: "User is blocked successfully."}, nil
}

/****************************************Fetch Blocked Contacts Details************************/
func (r *chatRepository) FetchBlockedUserDetails(ctx context.Context, user_id int64) (*models.Response, error) {
	block := []models.BlockedContacts{}
	re := []models.BlockedContact{}
	db := r.DBConn.Where("blockee_id = ? or blocker_id = ?", user_id, user_id).Order("created_at DESC").Find(&block)
	if db.RowsAffected == 0 {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("User is not found.")
		return &models.Response{Status: false, Msg: "User is not found.", ResponseCode: 404}, nil
	}
	var x map[string]interface{}
	for i := 0; i < len(block); i++ {
		result := models.BlockedContact{
			ID:        block[i].ID,
			CreatedAt: float64(block[i].CreatedAt.UnixNano() / 1000000),
			BlockeeId: block[i].BlockeeId,
			BlockerId: block[i].BlockerId,
		}
		re = append(re, result)
		bloker_id := strconv.Itoa(int(re[i].BlockerId))
		blockee_id := strconv.Itoa(int(re[i].BlockeeId))
		m1 := map[string]interface{}{bloker_id + "-" + blockee_id: re[i]}
		mergo.Merge(&x, m1)
	}
	logger.Logger.Info("blocked contact found.", &x)
	return &models.Response{Status: true, ResponseCode: 200, Msg: "Blocked contacts found.", BlockedContactDetails: &x}, nil
}

/****************************************Fetch Blocked Contacts Details**************************/
func (r *chatRepository) FetchBlockedContactDetails(ctx context.Context, user_id int64) (*models.Response, error) {
	data := []models.BlockedContacts{}
	re := []models.BlockedContact{}
	db := r.DBConn.Where("blocker_id = ?", user_id).Find(&data)
	if db.RowsAffected == 0 {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Blocked contacts not found.")
		return &models.Response{Status: false, Msg: "Blocked contacts not found.", ResponseCode: 404}, nil
	}
	for i := 0; i < len(data); i++ {
		result := models.BlockedContact{
			ID:        data[i].ID,
			CreatedAt: float64(data[i].CreatedAt.UnixNano() / 1000000),
			BlockeeId: data[i].BlockeeId,
		}
		re = append(re, result)
	}
	return &models.Response{Status: true, ResponseCode: 200, Msg: "Blocked contacts found successfully.", BlockedContacts: &re}, nil
}

/***************************************Unblock Blocked contacts***********************/
func (r *chatRepository) Unblock_user(ctx context.Context, flow models.BlockedContacts) (*models.Response, error) {
	err := r.DBConn.Where("blockee_id = ? AND blocker_id = ?", flow.BlockeeId, flow.BlockerId).First(&models.BlockedContacts{}).Delete(&models.BlockedContacts{})
	if err.Error != nil {
		logger.Logger.WithError(err.Error).WithField("error", err.Error).Error("User is not found in blocked list.")
		return &models.Response{Status: false, ResponseCode: 400, Msg: "User is not found in blocked list."}, nil
	}
	logger.Logger.Info("User is unblocked successfully.")
	return &models.Response{Status: true, ResponseCode: 200, Msg: "User is unblocked successfully."}, nil
}

/**********************************************Fetch Wallpaper Details*********************************/
func (r *chatRepository) FetchWallpapersDetails(ctx context.Context) (*models.Response, error) {
	m := []models.Wallpapers{}
	n := []models.Wallpaper{}
	o := []models.Wallpaper{}
	p := []models.Wallpaper{}
	q := []models.Wallpaper{}
	db := r.DBConn.Find(&m)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Wallpaper not found.")
		return &models.Response{Status: false, Msg: "Wallpaper not found.", ResponseCode: 400}, nil
	}
	for i := 0; i < len(m); i++ {
		switch m[i].WallpaperType {
		case "dark":
			{
				media := models.Wallpaper{
					ID:            m[i].ID,
					CreatedAt:     float64(m[i].CreatedAt.UnixNano() / 1000000),
					WallpaperURL:  m[i].WallpaperURL,
					WallpaperType: m[i].WallpaperType,
				}
				p = append(p, media)
			}
		case "bright":
			{
				media := models.Wallpaper{
					ID:            m[i].ID,
					CreatedAt:     float64(m[i].CreatedAt.UnixNano() / 1000000),
					WallpaperURL:  m[i].WallpaperURL,
					WallpaperType: m[i].WallpaperType,
				}
				n = append(n, media)
			}
		case "light":
			{
				media := models.Wallpaper{
					ID:            m[i].ID,
					CreatedAt:     float64(m[i].CreatedAt.UnixNano() / 1000000),
					WallpaperURL:  m[i].WallpaperURL,
					WallpaperType: m[i].WallpaperType,
				}
				o = append(o, media)
			}
		case "pattern":
			{
				media := models.Wallpaper{
					ID:            m[i].ID,
					CreatedAt:     float64(m[i].CreatedAt.UnixNano() / 1000000),
					WallpaperURL:  m[i].WallpaperURL,
					WallpaperType: m[i].WallpaperType,
				}
				q = append(q, media)
			}
		}
	}
	if len(n) == 0 && len(o) == 0 && len(p) == 0 && len(q) == 0 {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Wallpaper not Found.")
		return &models.Response{Status: false, Msg: "Wallpaper not found.", ResponseCode: 400}, nil
	}
	s := models.Data{
		Dark:    p,
		Bright:  n,
		Light:   o,
		Pattern: q,
	}
	logger.Logger.Info("Wallpaper found.")
	return &models.Response{Status: true, ResponseCode: 200, Msg: "Wallpapers found.", Data1: &s}, nil

}

/***********************************************Save Group Chat Setting*************************/
func (r *chatRepository) SaveGroupChatSetting(ctx context.Context, flow models.ChatSettings) (*models.Response, error) {
	user := r.DBConn.Where("id =?", flow.User_Id).First(&models.User{})
	if user.RowsAffected == 0 {
		logger.Logger.WithError(user.Error).WithField("error", user.Error).Error("User is not found.")
		return &models.Response{Status: false, Msg: "User is not found.", ResponseCode: 400}, nil
	}
	check := r.DBConn.Where("user_id = ?", flow.User_Id).First(&models.ChatSettings{}).Update("group_chat_type", flow.Group_Chat_Type)
	if check.RowsAffected == 0 {
		db := r.DBConn.Create(&flow)
		if db.Error != nil {
			logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Group chat setting is not saved.")
			return &models.Response{Status: false, Msg: "Group chat setting is not saved.", ResponseCode: 400}, nil
		}
		logger.Logger.Info("Group chat setting is saved successfully.")
		return &models.Response{Status: true, Msg: "Group chat setting is saved successfully.", ResponseCode: 200}, nil
	}
	logger.Logger.Info("Group chat setting is saved successfully.")
	return &models.Response{Status: true, Msg: "Group chat setting is saved successfully.", ResponseCode: 200}, nil

}

/************************************************Fetch Group Chat Setting*************************/
func (r *chatRepository) FetchGroupChatSetting(ctx context.Context, user_id int64) (*models.Response, error) {
	chat := models.ChatSettings{}
	db := r.DBConn.Where("user_id = ?", user_id).First(&chat)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Group chat setting is not found.")
		return &models.Response{Status: false, Msg: "Group chat setting is not found.", ResponseCode: 400}, nil
	}
	logger.Logger.Info("Group chat setting is found.")
	return &models.Response{Status: true, Msg: "Group chat setting is found.", ResponseCode: 200, GroupChatSetting: &chat.Group_Chat_Type}, nil
}
