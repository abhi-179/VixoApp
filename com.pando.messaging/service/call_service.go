package service

import (
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	repo "pandoMessagingWalletService/com.pando.messaging/repository"

	"context"
	"net/url"
	"time"
)

//	"encoding/json"

type callUsecase struct {
	repository repo.CallRepository
}

func NewcallUsecase(repo repo.CallRepository) CallUsecase {
	return &callUsecase{
		repository: repo,
	}
}

/*******************************************SaveCallLogs*******************************************/
func (r *callUsecase) SaveCallLogs(ctx context.Context, flow models.CallDetail) (*models.Response, error) {
	callduration := time.Unix(0, int64(flow.CallDuration)*int64(time.Millisecond))
	starttime := time.Unix(0, int64(flow.StartTime)*int64(time.Millisecond))
	endtime := time.Unix(0, int64(flow.EndTime)*int64(time.Millisecond))
	logger.Logger.Info("Request received from save call logs service part.")
	return r.repository.SaveCallLogs(ctx, flow, callduration, starttime, endtime)
}

/**********************************************Fetch All Call Logs**********************************/
func (r *callUsecase) FetchAllCallLogs(ctx context.Context, user_id int64, source string, query url.Values) (*models.Response, error) {
	logger.Logger.Info("Request received from fetch all call logs service part.")
	return r.repository.FetchAllCallLogs(ctx, user_id, source, query)
}

/**********************************************FetchMissedCallLogs************************************/
func (r *callUsecase) FetchMissedCallLogs(ctx context.Context, user_id int64, source string, query url.Values) (*models.Response, error) {
	logger.Logger.Info("Request received from fetch missed call logs service part.")
	return r.repository.FetchMissedCallLogs(ctx, user_id, source, query)
}

/**************************************************Delete Call Logs*********************************/
func (r *callUsecase) DeleteCallLogs(ctx context.Context, user_id int64, id int64) (*models.Response, error) {
	logger.Logger.Info("Request received from delete call logs service part.")
	return r.repository.DeleteCallLogs(ctx, user_id, id)
}
