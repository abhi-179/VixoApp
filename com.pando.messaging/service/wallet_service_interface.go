package service

import (
	"context"
	"net/url"
	models "pandoMessagingWalletService/com.pando.messaging/model"
)

type WalletUsecase interface {
	CreateWallet(ctx context.Context, password models.WalletReq) (*models.Response, error)
	GetBalance(ctx context.Context, walletId string) (*models.Response, error)
	AddToken(ctx context.Context, flow map[string]interface{}) (*models.Response, error)
	RequestToken(ctx context.Context, req models.RequestTokens) (*models.Response, error)
	RejectRequest(ctx context.Context, requestId int64, requestType string) (*models.Response, error)
	SendToken(ctx context.Context, send map[string]interface{}) (*models.Response, error)
	GetTransactions(ctx context.Context, walletId string, query url.Values) (*models.Response, error)
	ViewSpendAnalytics(ctx context.Context, walletId string) (*models.Response, error)
	WalletStatement(ctx context.Context, flow map[string]interface{}) (string, error)
	GetWalletId(ctx context.Context, userId int64) (*models.Response, error)
	RecentTransactions(ctx context.Context, walletId string, query url.Values) (*models.Response, error)
	ShowPendingRequests(ctx context.Context, walletId string, userId int64, query url.Values) (*models.Response, error)
	SendTokenToAdmin(ctx context.Context, id, concertId int64, amount, password string) (*models.Response, error)
	ShowOwnTokenRequests(ctx context.Context, walletId string, userId int64, query url.Values) (*models.Response, error)
}
