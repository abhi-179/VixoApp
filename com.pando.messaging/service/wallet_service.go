package service

//	"encoding/json"
import (
	"context"
	"fmt"
	"net/url"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	repo "pandoMessagingWalletService/com.pando.messaging/repository"
)

type walletUsecase struct {
	repository repo.WalletRepository
}

func NewWalletUsecase(repo repo.WalletRepository) WalletUsecase {
	return &walletUsecase{
		repository: repo,
	}
}

/***************************************Create Wallet**************************************/
func (r *walletUsecase) CreateWallet(ctx context.Context, password models.WalletReq) (*models.Response, error) {
	logger.Logger.Info("Request received from service part of create wallet.")
	return r.repository.CreateWallet(ctx, password)
}

/****************************************GetBalance*****************************************/
func (r *walletUsecase) GetBalance(ctx context.Context, walletId string) (*models.Response, error) {
	logger.Logger.Info("Request received from service part of get balance.")
	return r.repository.GetBalance(ctx, walletId)
}

/****************************************AddToken*****************************************/
func (r *walletUsecase) AddToken(ctx context.Context, flow map[string]interface{}) (*models.Response, error) {
	walletId := fmt.Sprintf("%v", flow["wallet_id"])
	amount := fmt.Sprintf("%v", flow["amount"])
	logger.Logger.Info("Request received from service part of add token.")
	return r.repository.AddToken(ctx, walletId, amount)
}

/****************************************AddToken*****************************************/
func (r *walletUsecase) RequestToken(ctx context.Context, req models.RequestTokens) (*models.Response, error) {
	logger.Logger.Info("Request received from service part of request token.")
	return r.repository.RequestToken(ctx, req)
}

/****************************************Reject Request*****************************************/
func (r *walletUsecase) RejectRequest(ctx context.Context, requestId int64, requestType string) (*models.Response, error) {
	logger.Logger.Info("Request received from service part of reject request.")
	return r.repository.RejectRequest(ctx, requestId, requestType)
}

/******************************************Send Token********************************************/
func (r *walletUsecase) SendToken(ctx context.Context, send map[string]interface{}) (*models.Response, error) {
	from := fmt.Sprintf("%v", send["from"])
	to := fmt.Sprintf("%v", send["to"])
	amount := fmt.Sprintf("%v", send["amount"])
	password := fmt.Sprintf("%v", send["password"])
	message := fmt.Sprintf("%v", send["message"])
	logger.Logger.Info("Request received from service part of send token.")
	return r.repository.SendToken(ctx, from, to, amount, password, message)
}

/*******************************************Get Transactions***********************************/
func (r *walletUsecase) GetTransactions(ctx context.Context, walletId string, query url.Values) (*models.Response, error) {
	logger.Logger.Info("Request received from service part of get transactions api.")
	return r.repository.GetTransactions(ctx, walletId, query)
}

/*****************************************ViewSpendAnalytics************************************/
func (r *walletUsecase) ViewSpendAnalytics(ctx context.Context, walletId string) (*models.Response, error) {
	logger.Logger.Info("Request received from service part of view spend analytics api.")
	return r.repository.ViewSpendAnalytics(ctx, walletId)
}

/******************************************Wallet Statement************************************/
func (r *walletUsecase) WalletStatement(ctx context.Context, flow map[string]interface{}) (string, error) {
	walletId := fmt.Sprintf("%v", flow["wallet_id"])
	startDate := fmt.Sprintf("%v", flow["start_date"])
	endDate := fmt.Sprintf("%v", flow["end_date"])
	totalMonths := fmt.Sprintf("%v", flow["total_months"])
	email := fmt.Sprintf("%v", flow["email"])
	queryType := fmt.Sprintf("%v", flow["query_type"])
	logger.Logger.Info("Request received from service part of wallet statement api.")
	return r.repository.WalletStatement(ctx, walletId, startDate, endDate, totalMonths, email, queryType)
}

/*****************************************Get Wallet Id ************************************/
func (r *walletUsecase) GetWalletId(ctx context.Context, userId int64) (*models.Response, error) {
	logger.Logger.Info("Request received from service part of get wallet id api.")
	return r.repository.GetWalletId(ctx, userId)
}

/*****************************************RecentTransactions******************************/
func (r *walletUsecase) RecentTransactions(ctx context.Context, walletId string, query url.Values) (*models.Response, error) {
	logger.Logger.Info("Request received from service part of recent transactions api.")
	return r.repository.RecentTransactions(ctx, walletId, query)
}

/****************************************Show all Requests*********************************/
func (r *walletUsecase) ShowPendingRequests(ctx context.Context, walletId string, userId int64, query url.Values) (*models.Response, error) {
	logger.Logger.Info("Request received from service part of show all request api.")
	return r.repository.ShowPendingRequests(ctx, walletId, userId, query)
}

/***************************************SendTokenToAdmin***********************************/
func (r *walletUsecase) SendTokenToAdmin(ctx context.Context, id, concertId int64, amount, password string) (*models.Response, error) {
	logger.Logger.Info("Request received from service part of send token to admin api.")
	return r.repository.SendTokenToAdmin(ctx, id, concertId, amount, password)
}

/****************************************Show all Requests*********************************/
func (r *walletUsecase) ShowOwnTokenRequests(ctx context.Context, walletId string, userId int64, query url.Values) (*models.Response, error) {
	logger.Logger.Info("Request received from service part of show own token request api.")
	return r.repository.ShowOwnTokenRequests(ctx, walletId, userId, query)
}
