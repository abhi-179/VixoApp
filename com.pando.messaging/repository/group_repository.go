package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"pandoMessagingWalletService/com.pando.messaging/config"
	"pandoMessagingWalletService/com.pando.messaging/config/kafka"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	"strconv"
	"strings"
	"time"

	"github.com/imdario/mergo"
	"github.com/lib/pq"

	"gorm.io/gorm"
)

type groupRepository struct {
	DBConn *gorm.DB
}

var AwsURL string

func NewgroupRepository(conn *gorm.DB, conf *config.Config) GroupRepository {
	AwsURL = conf.AwsURL
	return &groupRepository{
		DBConn: conn,
	}
}

/***********************************************Create Group*****************************************/
func (r *groupRepository) Create_Group(ctx context.Context, flow models.Group) (*models.Response, error) {
	logger.Logger.Info("Enter in to create group repository part.")
	var admin_id, pending_user_id, normal_user_id []int64
	user := models.User{}
	// var wg sync.WaitGroup
	// wg.Add(1)
	admin_id = append(admin_id, flow.Admin_ids...)
	chat_id := time.Now().UnixNano() / 1000000
	admin := r.DBConn.Where("id in (?) ", admin_id).First(&user)
	if admin.Error != nil {
		logger.Logger.WithError(admin.Error).WithField("error", admin.Error).Error("Please use valid user id in place of admin id.")
		return &models.Response{Status: false, ResponseCode: 400, Msg: "Please use valid user id in place of admin id."}, nil
	}
	for i := 0; i < len(flow.User_ids); i++ {
		user1 := models.User{}
		checker := r.DBConn.Where("id in (?)", flow.User_ids[i]).First(&user1)
		userid := strings.Trim(strings.Join(strings.Split(fmt.Sprint(flow.User_ids[i]), " "), ","), "[]")
		if checker.Error != nil {
			logger.Logger.WithError(checker.Error).WithField("error", checker.Error).Error("User id " + userid + " is not found.")
			return &models.Response{Status: false, ResponseCode: 404, Msg: "User id " + userid + " is not found."}, nil
		}
		chatsetting := models.ChatSettings{}
		private_user := r.DBConn.Where("user_id in (?)", flow.User_ids[i]).First(&chatsetting)
		if private_user.RowsAffected != 0 {
			message := models.Notifications{
				SenderUserId:      flow.Admin_ids[0],
				ReceiverUserId:    user1.ID,
				Message:           "You are invited by " + user.Username + " please accept to join the group " + flow.Group_name + ".",
				SenderUsername:    user.Username,
				ReceiverUsername:  user1.Username,
				NotificationTitle: "Group Invitation",
				NotificationType:  "GROUP_CHAT_INVITATION",
				Profile_Pic_Url:   flow.Profile_Pic_Url,
			}
			data, _ := json.Marshal(&message)
			go kafka.Push(context.Background(), nil, data)
			pending_user_id = append(pending_user_id, flow.User_ids[i])
		} else {
			normal_user_id = append(normal_user_id, flow.User_ids[i])
		}
	}

	grp := models.Groups{
		Group_name:        flow.Group_name,
		Admin_ids:         flow.Admin_ids,
		TotalUsers:        int64(len(normal_user_id)),
		Status:            "ACTIVE",
		ChatId:            chat_id,
		Subject_Timestamp: time.Unix(0, int64(flow.Subject_Timestamp)*int64(time.Millisecond)),
		Subject_Owner_Id:  flow.Subject_Owner_Id,
		Profile_Pic_Url:   flow.Profile_Pic_Url,
		User_ids:          normal_user_id,
		Pending_User_ids:  pending_user_id,
	}
	te := r.DBConn.Table("groups").Create(&grp)
	if te.Error != nil {
		logger.Logger.WithError(te.Error).WithField("err", te.Error).Errorf("Group is not created.")
		return &models.Response{Status: false, Msg: "Group is not created", ResponseCode: 400}, nil
	}
	result := models.Group{
		ID:                grp.ID,
		CreatedAt:         float64(grp.CreatedAt.UnixNano() / 1000000),
		Group_name:        grp.Group_name,
		Admin_ids:         grp.Admin_ids,
		ChatId:            grp.ChatId,
		TotalUsers:        grp.TotalUsers,
		Status:            grp.Status,
		Subject_Timestamp: float64(grp.Subject_Timestamp.UnixNano() / 1000000),
		Subject_Owner_Id:  grp.Subject_Owner_Id,
		Profile_Pic_Url:   grp.Profile_Pic_Url,
		User_ids:          grp.User_ids,
		Pending_User_ids:  grp.Pending_User_ids,
	}
	//wg.Wait()
	logger.Logger.Info("Group is created successfully.")
	return &models.Response{Status: true, Msg: "Group is created successfully", ResponseCode: http.StatusOK, Groups: &result}, nil
}

/********************************************Add users to Group************************************/
func (r *groupRepository) AddUserToGroup(ctx context.Context, flow models.Groups) (*models.Response, error) {
	logger.Logger.Info("Enter in to add users to group repository part.")
	g := models.Groups{}
	user1 := models.User{}
	var pending_user_id, normal_user_id []int64
	err := r.DBConn.Where("id = ?", flow.ID).First(&g)
	if err.Error != nil {
		logger.Logger.WithError(err.Error).WithField("error", err.Error).Error("Group is not found.")
		return &models.Response{Status: false, ResponseCode: 404, Msg: "Group is not found."}, nil
	}
	if len(flow.User_ids)+len(g.User_ids) > 200 {
		return &models.Response{Status: false, ResponseCode: 400, Msg: "You can only add 200 peoples."}, nil
	}
	admin := r.DBConn.Where("admin_ids && (?)", flow.Admin_ids).First(&g)
	if admin.Error != nil {
		logger.Logger.WithError(admin.Error).WithField("error", admin.Error).Error("You are not an admin.")
		return &models.Response{Status: false, ResponseCode: 400, Msg: "You are not an admin."}, nil
	}
	for i := 0; i < len(flow.User_ids); i++ {
		user := models.User{}
		checker := r.DBConn.Where("id = ?", flow.User_ids[i]).First(&user)
		userid := strings.Trim(strings.Join(strings.Split(fmt.Sprint(flow.User_ids[i]), " "), ","), "[]")
		if checker.Error != nil {
			logger.Logger.WithError(checker.Error).WithField("error", checker.Error).Error("User id " + userid + " is not found.")
			return &models.Response{Status: false, ResponseCode: 404, Msg: "User id " + userid + " is not found."}, nil
		}
		chatsetting := models.ChatSettings{}
		private_user := r.DBConn.Where("user_id in (?)", flow.User_ids[i]).First(&chatsetting)
		if private_user.RowsAffected != 0 {
			message := models.Notifications{
				SenderUserId:      user1.ID,
				ReceiverUserId:    user.ID,
				Message:           "You are invited by " + user1.Username + " please accept to join the group " + flow.Group_name + ".",
				SenderUsername:    user1.Username,
				ReceiverUsername:  user.Username,
				NotificationTitle: "Group Invitation",
				NotificationType:  "GROUP_CHAT_INVITATION",
				Profile_Pic_Url:   g.Profile_Pic_Url,
			}
			data, _ := json.Marshal(&message)
			fmt.Println(data)
			go kafka.Push(context.Background(), nil, data)
			pending_user_id = append(pending_user_id, flow.User_ids[i])
		} else {
			normal_user_id = append(normal_user_id, flow.User_ids[i])
		}
	}
	test := r.DBConn.Where("id = ? and user_ids && (?)", flow.ID, flow.User_ids).First(&models.Group{})
	if test.RowsAffected != 0 {
		logger.Logger.WithError(test.Error).WithField("error", test.Error).Error("These users are already in the group.")
		return &models.Response{Status: false, ResponseCode: 400, Msg: "These users are already in the group."}, nil
	}
	db := r.DBConn.Exec("update groups set user_ids = array_cat(user_ids, ?),pending_user_ids = array_cat(pending_user_ids, ?),total_users = ? where id = ?", []pq.Int64Array{normal_user_id}, []pq.Int64Array{pending_user_id}, g.TotalUsers+int64(len(normal_user_id)), flow.ID).Find(&g)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Users not added to group.")
		return &models.Response{Status: false, Msg: "Users not added to group.", ResponseCode: 400}, nil
	}
	data := models.Group{
		ID:                g.ID,
		CreatedAt:         float64(g.CreatedAt.UnixNano() / 1000000),
		Group_name:        g.Group_name,
		Admin_ids:         g.Admin_ids,
		ChatId:            g.ChatId,
		TotalUsers:        g.TotalUsers,
		Status:            g.Status,
		Subject_Timestamp: float64(g.Subject_Timestamp.UnixNano() / 1000000),
		Subject_Owner_Id:  g.Subject_Owner_Id,
		Profile_Pic_Url:   g.Profile_Pic_Url,
		User_ids:          g.User_ids,
		Pending_User_ids:  g.Pending_User_ids,
	}
	logger.Logger.WithField("success", g).Info("Users added to group successfully.")
	return &models.Response{Status: true, Msg: "Users added to group successfully", ResponseCode: http.StatusOK, Groups: &data}, nil
}

/*********************************************RemoveUsersFromGroup************************************/
func (r *groupRepository) RemoveUsersFromGroup(ctx context.Context, group_id int64, user_id int64, admin_id int64) (*models.Response, error) {
	logger.Logger.Info("Enter in to remove users from group repository part.")
	g := models.Groups{}
	if db := r.DBConn.Where("id = ?", group_id).First(&g).Error; db != nil {
		logger.Logger.WithError(db).WithField("err", db).Errorf("Group is not found.")
		return &models.Response{Status: false, Msg: "Group is not found", ResponseCode: 404}, nil
	}
	var result bool = false
	for _, x := range g.Admin_ids {
		if x == admin_id {
			result = true
			break
		}
	}
	userId := strconv.Itoa(int(user_id))
	if result {
		err := r.DBConn.Where(userId+" = any(user_ids) and id = ?", group_id).Find(&g)
		if err.RowsAffected == 0 {
			logger.Logger.WithError(err.Error).WithField("error", err.Error).Error("This user id is not in group.")
			return &models.Response{Status: false, Msg: "This user id is not in the group.", ResponseCode: 404}, nil
		}
		db := r.DBConn.Exec("update groups set user_ids = array_remove(user_ids, ?),admin_ids = array_remove(admin_ids, ?), total_users = ? where id = ?", user_id, user_id, g.TotalUsers-1, group_id)
		if db.Error != nil {
			logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("User is not removed from group.")
			return &models.Response{Status: false, Msg: "User is not removed from group.", ResponseCode: 400}, nil
		}
		logger.Logger.WithField("success", g).Info("User is removed from group successfully.")
		return &models.Response{Status: true, Msg: "User is removed from group successfully.", ResponseCode: 200}, nil
	}
	logger.Logger.Error("You are not an admin.")
	return &models.Response{Status: false, Msg: "You are not an admin.", ResponseCode: 400}, nil
}

/*********************************************LeaveGroup**********************************************/
func (r *groupRepository) LeaveGroup(ctx context.Context, group_id int64, user_id int64) (*models.Response, error) {
	logger.Logger.Info("Enter in to leave group repository part.")
	g := models.Groups{}
	if db := r.DBConn.Where("id = ?", group_id).First(&g).Error; db != nil {
		logger.Logger.WithError(db).WithField("err", db).Errorf("Group is not found.")
		return &models.Response{Status: false, Msg: "Group is not found", ResponseCode: 404}, nil
	}
	userId := strconv.Itoa(int(user_id))
	err := r.DBConn.Where(userId+" = any(user_ids) and id = ?", group_id).Find(&g)
	if err.RowsAffected == 0 {
		logger.Logger.WithError(err.Error).WithField("error", err.Error).Error("You are not in the group.")
		return &models.Response{Status: false, Msg: "You are not in the group.", ResponseCode: 404}, nil
	}
	if len(g.Admin_ids) > 1 {
		db := r.DBConn.Exec("update groups set user_ids = array_remove(user_ids, ?), admin_ids = array_remove(admin_ids, ?), total_users = ? where id = ?", user_id, user_id, g.TotalUsers-1, group_id)
		if db.Error != nil {
			logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("You haven't leaved the group successfully.")
			return &models.Response{Status: false, Msg: "You haven't leaved the group successfully.", ResponseCode: 400}, nil
		}
		logger.Logger.Info("You have left the group successfully.")
		return &models.Response{Status: true, Msg: "You have left the group successfully.", ResponseCode: 200}, nil
	} else {
		if g.TotalUsers > 1 {
			for i := 0; i < 1; {
				if user_id == g.Admin_ids[i] {
					var nextAdminId int64
					if g.User_ids[i] == user_id {
						nextAdminId = g.User_ids[i+1]
					} else {
						nextAdminId = g.User_ids[i]
					}
					db := r.DBConn.Exec("update groups set user_ids = array_remove(user_ids, ?), admin_ids = array_remove(admin_ids, ?), total_users = ? where id = ?", user_id, user_id, g.TotalUsers-1, group_id)
					if db.Error != nil {
						logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("You haven't leaved the group successfully.")
						return &models.Response{Status: false, Msg: "You haven't leaved the group successfully.", ResponseCode: 400}, nil
					}
					update := r.DBConn.Exec("update groups set admin_ids = array_append(admin_ids,?)", nextAdminId)
					if update.Error != nil {
						logger.Logger.WithError(update.Error).WithField("error", update.Error).Error("New admin is not created.")
						return &models.Response{Status: false, Msg: "You haven't leaved the group successfully.", ResponseCode: 400}, nil
					}
					logger.Logger.Info("You have left the group successfully.")
					return &models.Response{Status: true, Msg: "You have left the group successfully.", ResponseCode: 200, NewAdminId: nextAdminId}, nil

				} else {
					db := r.DBConn.Exec("update groups set user_ids = array_remove(user_ids, ?),total_users = ? where id = ?", user_id, g.TotalUsers-1, group_id)
					if db.Error != nil {
						logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("You haven't leaved the group successfully.")
						return &models.Response{Status: false, Msg: "You haven't leaved the group successfully.", ResponseCode: 400}, nil
					}
					logger.Logger.Info("You have left the group successfully.")
					return &models.Response{Status: true, Msg: "You have left the group successfully.", ResponseCode: 200}, nil
				}
			}
		} else {
			db := r.DBConn.Exec("update groups set user_ids = array_remove(user_ids, ?), admin_ids = array_remove(admin_ids, ?),total_users = ? where id = ?", user_id, user_id, g.TotalUsers-1, group_id)
			if db.Error != nil {
				logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("You haven't leaved the group successfully.")
				return &models.Response{Status: false, Msg: "You haven't leaved the group successfully.", ResponseCode: 400}, nil
			}
			logger.Logger.Info("You have left the group successfully.")
			return &models.Response{Status: true, Msg: "You have left the group successfully.", ResponseCode: 200}, nil

		}
	}
	return &models.Response{Status: false, Msg: "Something went wrong", ResponseCode: 400}, nil
}

/*****************************************EditGroupInfo********************************************/
func (r *groupRepository) EditGroupInfo(ctx context.Context, group_id int64, new_group_name string, profile_pic_url string) (*models.Response, error) {
	logger.Logger.Info("Enter in to edit group name repository part.")
	g := models.Groups{}
	err := r.DBConn.Where("id = ?", group_id).First(&g)
	if err.Error != nil {
		logger.Logger.WithError(err.Error).WithField("error", err.Error).Error("Group is not found.")
		return &models.Response{Status: false, ResponseCode: 404, Msg: "Group is not found."}, nil
	}
	if db := r.DBConn.Where("id = ?", group_id).Find(&g).Updates(map[string]interface{}{"group_name": new_group_name, "profile_pic_url": profile_pic_url}).Error; db != nil {
		logger.Logger.WithError(db).WithField("err", db).Errorf("Group is not found.")
		return &models.Response{Status: false, Msg: "Group is not found", ResponseCode: 404}, nil
	}
	logger.Logger.Info("Group info is updated successfully.")
	return &models.Response{Status: true, Msg: "Group info is updated successfully.", ResponseCode: 200}, nil

}

/*****************************************SearchUsersInGroup****************************************/
func (r *groupRepository) SearchUsersInGroup(ctx context.Context, group_id string, username string) (*models.Response, error) {
	logger.Logger.Info("Request received from SearchUsersInGroup repository")
	gu := models.Groups{}
	if db := r.DBConn.Where("id = ? AND user_ids ILIKE ?", group_id, username+"%").Find(&gu).Error; db != nil {
		logger.Logger.WithError(db).WithField("err", db).Errorf("Group is not found.")
		return &models.Response{Status: false, Msg: "Group is not found", ResponseCode: 404}, nil
	}
	return &models.Response{Status: true, Msg: "Users found in group.", ResponseCode: 200}, nil
}

/*********************************************Delete Group*******************************************/
func (r *groupRepository) DeleteGroup(ctx context.Context, group_id int64) (*models.Response, error) {
	logger.Logger.Info("Enter in to delete group repository part.")
	g := models.Groups{}
	if err := r.DBConn.Where("id = ?", group_id).First(&g).Delete(&g).Error; err != nil {
		logger.Logger.WithError(err).WithField("err", err).Errorf("Group is not found.")
		return &models.Response{Status: false, Msg: "Group is not found", ResponseCode: 404}, nil
	}
	logger.Logger.Info("Group is deleted successfully.")
	return &models.Response{Status: true, Msg: "Group is deleted successfully.", ResponseCode: 200}, nil
}

/*********************************************UploadProfilePhoto*****************************************/
func (r *groupRepository) UploadGroupProfilePhoto(ctx context.Context, file multipart.File, handler *multipart.FileHeader, filename string) (*models.Response, error) {
	url := AwsURL + filename
	logger.Logger.Info("Tenant/repository/tenant.go:Upload Profile photo")
	return &models.Response{Status: true, Msg: "Profile Photo uploaded successfully.", ResponseCode: http.StatusOK, URL: url}, nil
}

/*******************************************Make Or remove Admin********************************/
func (r *groupRepository) MakeOrRemoveAdmin(ctx context.Context, user_id int64, group_id int64, new_admin_id int64, method_type string) (*models.Response, error) {
	g := models.Groups{}
	check := r.DBConn.Where("id = ?", group_id).First(&g)
	if check.Error != nil {
		logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("Group is not found")
		return &models.Response{Status: false, Msg: "Group is not found.", ResponseCode: 400}, nil
	}

	if new_admin_id == user_id {
		logger.Logger.Info("User id and new admin id are same.")
		return &models.Response{Status: false, Msg: "You can not make or remove yourself as admin.", ResponseCode: 400}, nil
	}
	admninid := strconv.Itoa(int(user_id))
	userid := strconv.Itoa(int(new_admin_id))
	check2 := r.DBConn.Where(admninid + " = any(admin_ids)").First(&g)
	if check2.RowsAffected == 0 {
		logger.Logger.WithError(check2.Error).WithField("error", check2.Error).Error("You are not an admin.")
		return &models.Response{Status: false, Msg: "You are not an admin.", ResponseCode: 400}, nil
	}
	check3 := r.DBConn.Where(userid + " = any(user_ids)").First(&g)
	if check3.RowsAffected == 0 {
		logger.Logger.WithError(check3.Error).WithField("error", check3.Error).Error("User is not present in group.")
		return &models.Response{Status: false, Msg: "User is not present in group.", ResponseCode: 400}, nil
	}
	if method_type == "Create-admin" {
		var result bool = false
		for _, x := range g.Admin_ids {
			if x == new_admin_id {
				result = true
				break
			}
		}
		if result {
			logger.Logger.Info("This user is already an admin.")
			return &models.Response{Status: false, Msg: "This user is already an admin.", ResponseCode: 400}, nil
		}
		db := r.DBConn.Exec("update groups set admin_ids = array_cat(admin_ids, ?) where id = ?", pq.Int64Array{new_admin_id}, group_id)
		if db.Error != nil {
			logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Admin not added to list.")
			return &models.Response{Status: false, Msg: "Admin not created.", ResponseCode: 400}, nil
		}
		logger.Logger.Info("Admin created successfully.")
		return &models.Response{Status: true, Msg: "Admin created successfully.", ResponseCode: http.StatusOK}, nil
	} else if method_type == "Remove-admin" {
		db := r.DBConn.Exec("update groups set admin_ids = array_remove(admin_ids, ?) where id = ?", new_admin_id, group_id)
		if db.Error != nil {
			logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Admin not removed.")
			return &models.Response{Status: false, Msg: "Admin not removed.", ResponseCode: 400}, nil
		}
		logger.Logger.Info("Admin removed successfully.")
		return &models.Response{Status: true, Msg: "Admin removed successfully.", ResponseCode: http.StatusOK}, nil
	}

	return &models.Response{Status: false, ResponseCode: 400, Msg: "Please provide method_type"}, nil
}

/********************************************Get Group Details************************************/
func (r *groupRepository) GetGroupDetails(ctx context.Context, group_id int64) (*models.Response, error) {
	data := models.Groups{}
	user := []models.User{}
	var id []int64
	db := r.DBConn.Where("id = ?", group_id).Find(&data)
	if db.RowsAffected == 0 {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Group details not found.")
		return &models.Response{Status: false, Msg: "Group details not found.", ResponseCode: 400}, nil
	}
	for _, v := range data.User_ids {
		id = append(id, v)
	}
	err := r.DBConn.Select("id,EXTRACT(EPOCH FROM created_at::timestamptz(3))*1000,username,profile_pic_url").Where("id in (?)", id).Find(&user)
	if err.RowsAffected == 0 {
		logger.Logger.WithError(err.Error).WithField("error", err.Error).Error("User details not found.")
		return &models.Response{Status: false, Msg: "User details not found.", ResponseCode: 400}, nil
	}
	var x map[int64]interface{}
	for i := 0; i < len(user); i++ {
		m1 := map[int64]interface{}{user[i].ID: user[i]}
		mergo.Merge(&x, m1)
	}
	result := models.Group{
		ID:                data.ID,
		CreatedAt:         float64(data.CreatedAt.UnixNano() / 1000000),
		Group_name:        data.Group_name,
		Admin_ids:         data.Admin_ids,
		TotalUsers:        data.TotalUsers,
		Status:            data.Status,
		Subject_Timestamp: float64(data.Subject_Timestamp.UnixNano() / 1000000),
		Subject_Owner_Id:  data.Subject_Owner_Id,
		Profile_Pic_Url:   data.Profile_Pic_Url,
		ChatId:            data.ChatId,
		User_ids:          data.User_ids,
	}
	return &models.Response{Status: true, ResponseCode: 200, Msg: "Group details found.", Groups: &result, UserDetails: &x}, nil
}

/******************************************AcceptAndDeclineGroupInvitation****************************/
func (r *groupRepository) AcceptAndDeclineGroupInvitation(ctx context.Context, group models.AcceptGroupInvitationDto) (*models.Response, error) {
	grp := models.Groups{}
	check := r.DBConn.Where("id = ?", group.GroupId).Find(&grp)
	if check.RowsAffected == 0 {
		logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("Group is not found.")
		return &models.Response{Status: false, Msg: "User is not found.", ResponseCode: 400}, nil
	}
	if group.Type == "ACCEPT" {
		db := r.DBConn.Exec("update groups set pending_user_ids = array_remove(pending_user_ids,?),user_ids = array_cat(user_ids,?),total_users where id = ?", pq.Int64Array{group.UserId}, pq.Int64Array{group.UserId}, grp.TotalUsers+1, group.GroupId)
		if db.Error != nil {
			logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("User is not updated.")
			return &models.Response{Status: false, Msg: "User is not updated.", ResponseCode: 400}, nil
		}
		return &models.Response{Status: true, ResponseCode: 200, Msg: "You have added to the group successfully."}, nil

	} else if group.Type == "DECLINE" {
		update := r.DBConn.Exec("update groups set pending_user_ids = array_remove(pending_user_ids,?) where id = ?", pq.Int64Array{group.UserId}, group.GroupId)
		if update.Error != nil {
			logger.Logger.WithError(update.Error).WithField("error", update.Error).Error("User is not updated.")
			return &models.Response{Status: false, Msg: "User is not updated.", ResponseCode: 400}, nil
		}
		return &models.Response{Status: true, ResponseCode: 200, Msg: "You have declined the invitation."}, nil
	} else {
		return &models.Response{Status: false, ResponseCode: 400, Msg: "Please provide method_type either ACCEPT or DECLINE."}, nil
	}
}

/**************************************Get all group details of user***************************/
func (r *groupRepository) GetAllGroupDetailsOfUser(ctx context.Context, user_id int64) (*models.Response, error) {
	groupdetails := make([]models.Group, 0)
	userdetails := make([]models.UserDtos, 0)
	var userid []int64
	userId := strconv.Itoa(int(user_id))
	rows, err := r.DBConn.Table("groups").Select("id, EXTRACT(EPOCH FROM created_at::timestamptz(3))*1000,group_name,chat_id,admin_ids,total_users,status,EXTRACT(EPOCH FROM subject_timestamp::timestamptz(3))*1000,subject_owner_id,profile_pic_url,user_ids").Where(userId + " = any(user_ids)").Order("created_at DESC").Rows()
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error("Group details not found.")
		return &models.Response{Status: false, Msg: "Group details is not found.", ResponseCode: http.StatusBadRequest}, nil
	}
	defer rows.Close()
	for rows.Next() {
		f := models.Group{}
		if err := rows.Scan(&f.ID, &f.CreatedAt, &f.Group_name, &f.ChatId, &f.Admin_ids, &f.TotalUsers, &f.Status, &f.Subject_Timestamp, &f.Subject_Owner_Id, &f.Profile_Pic_Url, &f.User_ids); err != nil {
			return nil, err
		}
		groupdetails = append(groupdetails, f)
	}
	for i := 0; i < len(groupdetails); i++ {
		userid = append(userid, groupdetails[i].User_ids...)
	}
	test := r.DBConn.Raw("select u.id,EXTRACT(EPOCH FROM u.created_at::timestamptz(3))*1000,u.username,u.profile_pic_url,u.email, f.friendship_status from users u left join friends f on (u.id=f.friend_id1 AND f.friend_id2=?) OR (f.friend_id1=? AND u.id=f.friend_id2) where u.id in(?)", user_id, user_id, userid).Order("created_at DESC").Scan(&userdetails)
	if test.Error != nil {
		logger.Logger.WithError(err).WithField("error", err).Error("Chat details not found.")
		return &models.Response{Status: false, Msg: "User details is not found.", ResponseCode: http.StatusBadRequest}, nil
	}
	if len(groupdetails) == 0 && len(userdetails) == 0 {
		logger.Logger.Error("User details not found.")
		return &models.Response{Status: false, Msg: "User details is not found.", ResponseCode: http.StatusBadRequest}, nil
	}
	var y map[int64]interface{}

	for i := 0; i < len(userdetails); i++ {
		var friend bool
		if userdetails[i].ID != user_id {

			if userdetails[i].Friendship_status == "ACCEPTED" {
				friend = true
			} else {
				friend = false
			}
			result := models.UserDtos{
				ID:              userdetails[i].ID,
				CreatedAt:       userdetails[i].CreatedAt,
				Username:        userdetails[i].Username,
				Email:           userdetails[i].Email,
				Profile_Pic_Url: userdetails[i].Profile_Pic_Url,
				Is_friend:       friend,
			}
			m2 := map[int64]interface{}{userdetails[i].ID: result}
			mergo.Merge(&y, m2)
		}
	}
	logger.Logger.Info("Group details found.")
	return &models.Response{Status: true, Msg: "Group details found.", ResponseCode: http.StatusOK, GroupDetail: &groupdetails, UserDetails: &y}, nil
}
