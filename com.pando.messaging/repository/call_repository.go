package repository

import (
	"context"
	"net/http"
	"net/url"
	"pandoMessagingWalletService/com.pando.messaging/config"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	"strconv"
	"time"

	"github.com/imdario/mergo"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type callRepository struct {
	DBConn *gorm.DB
}

func NewcallRepository(conn *gorm.DB) CallRepository {
	return &callRepository{
		DBConn: conn,
	}
}

/*********************************************Save Call logs****************************************/
func (r *callRepository) SaveCallLogs(ctx context.Context, flow models.CallDetail, callduration time.Time, starttime time.Time, endtime time.Time) (*models.Response, error) {
	logger.Logger.Info("Enter in to save call logs repository part.")
	callUser := models.CallDetails{}
	var userids []int64
	userids = append(userids, flow.User_ids...)
	user := []models.User{}
	callerId := r.DBConn.Where("id = ?", flow.CallerId).First(&models.User{})
	if callerId.Error != nil {
		logger.Logger.WithError(callerId.Error).WithField("error", callerId.Error).Error("CallerId not found.")
		return &models.Response{Status: false, ResponseCode: http.StatusNotFound, Msg: "CallerId not found please give valid user ids."}, nil
	}
	check := r.DBConn.Raw("select user_ids,call_id from call_details where user_ids in (?)", flow.User_ids).Find(&callUser)
	if check.RowsAffected == 0 {

		call_id := time.Now().UnixNano() / 100000
		callLogs := models.CallDetails{
			CallId:       call_id,
			CallerId:     flow.CallerId,
			CallDuration: callduration,
			Filehash:     flow.Filehash,
			AwsUrl:       flow.AwsUrl,
			StartTime:    starttime,
			EndTime:      endtime,
			IsAudioCall:  flow.IsAudioCall,
			IsMissedCall: flow.IsMissedCall,
			IsGroupCall:  flow.IsGroupCall,
			User_ids:     pq.Int64Array(flow.User_ids),
		}

		check := r.DBConn.Where("id in (?)", userids).Find(&user)
		if check.RowsAffected < int64(len(flow.User_ids)) {
			logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("Users are not found.")
			return &models.Response{Status: false, ResponseCode: http.StatusNotFound, Msg: "These users are not found please give valid user ids."}, nil
		}
		db := r.DBConn.Create(&callLogs)
		if db.Error != nil {
			logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Call logs is not saved.")
			return &models.Response{Status: false, Msg: "Call logs is not saved.", ResponseCode: http.StatusBadRequest}, nil
		}
		result := models.CallDetail{
			ID:           callLogs.ID,
			CallId:       callLogs.CallId,
			CreatedAt:    float64(callLogs.CreatedAt.UnixNano() / 1000000),
			CallerId:     callLogs.CallerId,
			CallDuration: float64(callLogs.CallDuration.UnixNano() / 1000000),
			Filehash:     callLogs.Filehash,
			AwsUrl:       callLogs.AwsUrl,
			StartTime:    float64(callLogs.StartTime.UnixNano() / 1000000),
			EndTime:      float64(callLogs.EndTime.UnixNano() / 1000000),
			User_ids:     callLogs.User_ids,
			IsAudioCall:  callLogs.IsAudioCall,
			IsMissedCall: callLogs.IsMissedCall,
			IsGroupCall:  callLogs.IsGroupCall,
		}
		logger.Logger.Info("Call logs saved successfully. ", callLogs)
		return &models.Response{Status: true, Msg: "Call logs saved successfully.", ResponseCode: http.StatusOK, Details: &result}, nil
	}

	call := models.CallDetails{
		CallId:       callUser.CallId,
		CallerId:     flow.CallerId,
		CallDuration: callduration,
		Filehash:     flow.Filehash,
		AwsUrl:       flow.AwsUrl,
		StartTime:    starttime,
		EndTime:      endtime,
		IsAudioCall:  flow.IsAudioCall,
		IsMissedCall: flow.IsMissedCall,
		IsGroupCall:  flow.IsGroupCall,
		User_ids:     pq.Int64Array(flow.User_ids),
	}

	check = r.DBConn.Where("id in (?)", userids).Find(&user)
	if check.RowsAffected < int64(len(flow.User_ids)) {
		logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("Users are not found.")
		return &models.Response{Status: false, ResponseCode: http.StatusNotFound, Msg: "These users are not found please give valid user ids."}, nil
	}
	db := r.DBConn.Create(&call)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Call logs is not saved.")
		return &models.Response{Status: false, Msg: "Call logs is not saved.", ResponseCode: http.StatusBadRequest}, nil
	}
	result := models.CallDetail{
		ID:           call.ID,
		CreatedAt:    float64(call.CreatedAt.UnixNano() / 1000000),
		CallId:       call.CallId,
		CallerId:     call.CallerId,
		CallDuration: float64(call.CallDuration.UnixNano() / 1000000),
		Filehash:     call.Filehash,
		AwsUrl:       call.AwsUrl,
		StartTime:    float64(call.StartTime.UnixNano() / 1000000),
		EndTime:      float64(call.EndTime.UnixNano() / 1000000),
		User_ids:     call.User_ids,
		IsAudioCall:  call.IsAudioCall,
		IsMissedCall: call.IsMissedCall,
		IsGroupCall:  call.IsGroupCall,
	}
	logger.Logger.Info("Call logs saved successfully. ", call)
	return &models.Response{Status: true, Msg: "Call logs saved successfully.", ResponseCode: http.StatusOK, Details: &result}, nil
}

/*********************************************FetchAllCallLogs*****************************************/
func (r *callRepository) FetchAllCallLogs(ctx context.Context, user_id int64, source string, query url.Values) (*models.Response, error) {
	logger.Logger.Info("Enter in to fetch all call logs repository part.")
	userdetails := []models.User{}
	var userid []int64
	user_id1 := strconv.Itoa(int(user_id))
	calldetail := make([]models.CallDetail, 0)
	pagination := config.GeneratePaginationFromRequest(query)
	offset := (pagination.Page) * pagination.Limit
	queryBuider := r.DBConn.Limit(pagination.Limit).Offset(offset)

	test, _ := queryBuider.Raw("select id,EXTRACT(EPOCH FROM created_at::timestamptz(3))*1000,call_id,caller_id,EXTRACT(EPOCH FROM call_duration::timestamptz(3))*1000,filehash,aws_url,"+
		"EXTRACT(EPOCH FROM start_time::timestamptz(3))*1000,EXTRACT(EPOCH FROM end_time::timestamptz(3))*1000,user_ids,deleted_by_user_ids,is_audio_call,is_missed_call,is_group_call from call_details where "+user_id1+" = any(user_ids)"+
		"except select id, EXTRACT(EPOCH FROM created_at::timestamptz(3))*1000,call_id,caller_id,EXTRACT(EPOCH FROM call_duration::timestamptz(3))*1000,filehash,aws_url,EXTRACT(EPOCH FROM start_time::timestamptz(3))*1000,"+
		"EXTRACT(EPOCH FROM end_time::timestamptz(3))*1000,user_ids,deleted_by_user_ids,is_audio_call,is_missed_call,is_group_call from call_details where "+user_id1+" = any(deleted_by_user_ids)"+
		"except select distinct on(c.id)c.id, EXTRACT(EPOCH FROM c.created_at::timestamptz(3))*1000,call_id,caller_id,EXTRACT(EPOCH FROM call_duration::timestamptz(3))*1000,filehash,aws_url,"+
		"EXTRACT(EPOCH FROM start_time::timestamptz(3))*1000,EXTRACT(EPOCH FROM end_time::timestamptz(3))*1000,user_ids,deleted_by_user_ids,is_audio_call,is_missed_call,is_group_call "+
		"from call_details as c inner join blocked_contacts as b on "+user_id1+" = b.blocker_id and b.blockee_id = any(c.user_ids) and c.created_at > b.created_at where "+user_id1+" = any(c.user_ids) and array_length(c.user_ids,1)=2 order by 2 desc limit  ? offset ?", pagination.Limit, offset).Rows()
	defer test.Close()
	for test.Next() {
		f := models.CallDetail{}
		if err := test.Scan(&f.ID, &f.CreatedAt, &f.CallId, &f.CallerId, &f.CallDuration, &f.Filehash, &f.AwsUrl, &f.StartTime, &f.EndTime, &f.User_ids, &f.Deleted_by_user_ids, &f.IsAudioCall, &f.IsMissedCall, &f.IsGroupCall); err != nil {
			return nil, err
		}
		calldetail = append(calldetail, f)
	}
	for i := 0; i < len(calldetail); i++ {
		userid = append(userid, calldetail[i].User_ids...)
	}
	db := r.DBConn.Table("users").Select("id,EXTRACT(EPOCH FROM created_at::timestamptz(3))*1000,username,profile_pic_url").Where("id in (?)", userid).Scan(&userdetails)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Users are not found.")
		return &models.Response{Status: false, Msg: "Users are not found.", ResponseCode: http.StatusNotFound}, nil
	}
	var x map[int64]interface{}
	for i := 0; i < len(userdetails); i++ {
		m1 := map[int64]interface{}{userdetails[i].ID: userdetails[i]}
		mergo.Merge(&x, m1)
	}
	if len(calldetail) == 0 {
		logger.Logger.WithError(test.Err()).WithField("error", test.Err()).Error("Call logs not found.")
		return &models.Response{Status: false, Msg: "Call logs not found.", ResponseCode: http.StatusNotFound}, nil
	}
	if source == "Web" {
		logger.Logger.Info("Call logs found. ", calldetail)
		return &models.Response{Status: true, Msg: "Call logs found.", ResponseCode: http.StatusOK, CallDetails: &calldetail, User_details: &userdetails}, nil
	} else {
		logger.Logger.Info("Call logs found. ", calldetail)
		return &models.Response{Status: true, Msg: "Call logs found.", ResponseCode: http.StatusOK, CallDetails: &calldetail, UserDetails: &x}, nil
	}
}

/*********************************************FetchAllCallLogs******************************************/
func (r *callRepository) FetchMissedCallLogs(ctx context.Context, user_id int64, source string, query url.Values) (*models.Response, error) {
	logger.Logger.Info("Enter in to fetch missed call logs repository part")
	userdetails := []models.User{}
	var userid []int64
	user_id1 := strconv.Itoa(int(user_id))
	calldetail := []models.CallDetail{}
	pagination := config.GeneratePaginationFromRequest(query)
	offset := (pagination.Page) * pagination.Limit
	queryBuider := r.DBConn.Limit(pagination.Limit).Offset(offset)
	test, _ := queryBuider.Raw("select id,EXTRACT(EPOCH FROM created_at::timestamptz(3))*1000,call_id,caller_id,EXTRACT(EPOCH FROM call_duration::timestamptz(3))*1000,filehash,aws_url,"+
		"EXTRACT(EPOCH FROM start_time::timestamptz(3))*1000,EXTRACT(EPOCH FROM end_time::timestamptz(3))*1000,user_ids,deleted_by_user_ids,is_audio_call,is_missed_call,is_group_call from call_details where "+user_id1+" = any(user_ids) and is_missed_call is true "+
		"except select id, EXTRACT(EPOCH FROM created_at::timestamptz(3))*1000,call_id,caller_id,EXTRACT(EPOCH FROM call_duration::timestamptz(3))*1000,filehash,aws_url,EXTRACT(EPOCH FROM start_time::timestamptz(3))*1000,"+
		"EXTRACT(EPOCH FROM end_time::timestamptz(3))*1000,user_ids,deleted_by_user_ids,is_audio_call,is_missed_call,is_group_call from call_details where "+user_id1+" = any(deleted_by_user_ids)"+
		"except select distinct on(c.id)c.id, EXTRACT(EPOCH FROM c.created_at::timestamptz(3))*1000,call_id,caller_id,EXTRACT(EPOCH FROM call_duration::timestamptz(3))*1000,filehash,aws_url,"+
		"EXTRACT(EPOCH FROM start_time::timestamptz(3))*1000,EXTRACT(EPOCH FROM end_time::timestamptz(3))*1000,user_ids,deleted_by_user_ids,is_audio_call,is_missed_call,is_group_call "+
		"from call_details as c inner join blocked_contacts as b on "+user_id1+" = b.blocker_id and b.blockee_id = any(c.user_ids) and c.created_at > b.created_at where "+user_id1+" = any(c.user_ids) and array_length(c.user_ids,1)=2 order by 2 desc limit  ? offset ?", pagination.Limit, offset).Rows()
	defer test.Close()
	for test.Next() {
		f := models.CallDetail{}
		if err := test.Scan(&f.ID, &f.CreatedAt, &f.CallId, &f.CallerId, &f.CallDuration, &f.Filehash, &f.AwsUrl, &f.StartTime, &f.EndTime, &f.User_ids, &f.Deleted_by_user_ids, &f.IsAudioCall, &f.IsMissedCall, &f.IsGroupCall); err != nil {
			return nil, err
		}
		calldetail = append(calldetail, f)
	}
	for i := 0; i < len(calldetail); i++ {
		userid = append(userid, calldetail[i].User_ids...)
	}
	db := r.DBConn.Table("users").Select("id,EXTRACT(EPOCH FROM created_at::timestamptz(3))*1000,username,profile_pic_url").Where("id in (?)", userid).Scan(&userdetails)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Users are not found.")
		return &models.Response{Status: false, Msg: "Users are not found.", ResponseCode: http.StatusNotFound}, nil
	}
	var x map[int64]interface{}
	for i := 0; i < len(userdetails); i++ {
		m1 := map[int64]interface{}{userdetails[i].ID: userdetails[i]}
		mergo.Merge(&x, m1)
	}
	if len(calldetail) == 0 {
		logger.Logger.WithError(test.Err()).WithField("error", test.Err()).Error("Call logs not found.")
		return &models.Response{Status: false, Msg: "Call logs not found.", ResponseCode: http.StatusNotFound}, nil
	}
	if source == "Web" {
		logger.Logger.Info("Call logs found. ", calldetail)
		return &models.Response{Status: true, Msg: "Call logs found.", ResponseCode: http.StatusOK, CallDetails: &calldetail, User_details: &userdetails}, nil
	} else {

		logger.Logger.Info("Call logs found. ", calldetail)
		return &models.Response{Status: true, Msg: "Call logs found.", ResponseCode: http.StatusOK, CallDetails: &calldetail, UserDetails: &x}, nil
	}
}

/*********************************************Delete Call Logs*******************************************/
func (r *callRepository) DeleteCallLogs(ctx context.Context, user_id int64, id int64) (*models.Response, error) {
	logger.Logger.Info("Enter in to delete call logs repository part.")
	userId := strconv.Itoa(int(user_id))
	calldetails := models.CallDetails{}
	db := r.DBConn.Where(userId+" = any(user_ids) and id = ?", id).Find(&calldetails)
	if db.RowsAffected == 0 {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("This user id not present in user list.")
		return &models.Response{Status: false, ResponseCode: 404, Msg: "This user id not present in users list."}, nil
	}
	check := r.DBConn.Where(userId+" = any(deleted_by_user_ids) and id = ?", id).Find(&calldetails)
	if check.RowsAffected == 0 {
		db := r.DBConn.Exec("update call_details set deleted_by_user_ids = array_append(deleted_by_user_ids, ?) where id = ?", user_id, id)
		if db.Error != nil {
			logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Call logs not deleted.")
			return &models.Response{Status: false, Msg: "Call logs not deleted.", ResponseCode: 400}, nil
		}
		logger.Logger.Info("Call logs deleted.")
		return &models.Response{Status: true, Msg: "Call logs deleted.", ResponseCode: http.StatusOK}, nil
	}
	logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("Call logs already deleted.")
	return &models.Response{Status: false, Msg: "Call logs already deleted.", ResponseCode: http.StatusNotFound}, nil
}
