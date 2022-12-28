package service

import (
	"context"
	"net/url"
	models "pandoMessagingWalletService/com.pando.messaging/model"
)

type CallUsecase interface {
	SaveCallLogs(ctx context.Context, flow models.CallDetail) (*models.Response, error)
	FetchAllCallLogs(ctx context.Context, user_id int64, source string, query url.Values) (*models.Response, error)
	FetchMissedCallLogs(ctx context.Context, user_id int64, source string, query url.Values) (*models.Response, error)
	DeleteCallLogs(ctx context.Context, user_id int64, id int64) (*models.Response, error)
}
