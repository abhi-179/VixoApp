package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	constants "pandoMessagingWalletService/com.pando.messaging/constants"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	wallet "pandoMessagingWalletService/com.pando.messaging/service"
	"strconv"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

type WalletController struct {
	Usecase wallet.WalletUsecase
}

/*******************************************Create Wallet********************************/
// @Tags wallet service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary create wallet
// @Description create wallet
// @Accept  json
// @Produce  json
// @Param save review body models.WalletReq true "create wallet"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/wallet/create_wallet [post]
func (r *WalletController) CreateWallet(c echo.Context) error {
	var pass models.WalletReq
	c.Bind(&pass)
	id1, _ := strconv.Atoi(c.Request().Header.Get("id"))
	logger.Logger.Info("id ", id1)
	id := int64(id1)
	if id != pass.UserId {
		logger.Logger.Error("Wrong auth token.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter correct auth token."})
	}
	password := models.WalletReq{
		UserId:   pass.UserId,
		Password: pass.Password,
	}
	if password.Password == "" || password.UserId == 0 {
		logger.Logger.Error("Password or user id is missing.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter password or user id both."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of create wallet is ", logrus.Fields{"password": password.Password, "user_id": password.UserId})
	authResponse, _ := r.Usecase.CreateWallet(ctx, password)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*****************************************Get Balance*************************************/
// @Tags wallet service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Get balance
// @Description get balance
// @Accept  json
// @Produce  json
// @Param wallet_id query string true "wallet_id"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/wallet/get_balance [get]
func (r *WalletController) GetBalance(c echo.Context) error {
	walletId := c.QueryParam("wallet_id")
	if walletId == "" {
		logger.Logger.Error(constants.WalletIDError)
		return c.JSON(400, &models.Response{Status: false, Msg: constants.WalletIdReq})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of get balance is ", logrus.Fields{"wallet_id": walletId})
	authResponse, _ := r.Usecase.GetBalance(ctx, walletId)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/******************************************Add Token******************************************/
func (r *WalletController) AddToken(c echo.Context) error {
	var at map[string]interface{}
	c.Bind(&at)
	if at == nil {
		logger.Logger.Error("Wallet id or amount is missing.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter wallet id and amount."})
	}
	if at["wallet_id"] == "" || at["amount"] == "" {
		logger.Logger.Error("Wallet id or amount is missing.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter wallet id or amount both."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of add token is ", logrus.Fields{"wallet_id": at["wallet_id"], "amount": at["amount"]})
	authResponse, _ := r.Usecase.AddToken(ctx, at)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*******************************************Request Token************************************/
// @Tags wallet service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Request token
// @Description request token
// @Accept  json
// @Produce  json
// @Param request token body models.RequestTokensDto true "request token"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/wallet/request_token [post]
func (r *WalletController) RequestToken(c echo.Context) error {
	var at models.RequestTokens
	c.Bind(&at)
	id1, _ := strconv.Atoi(c.Request().Header.Get("id"))
	logger.Logger.Info("id ", id1)
	id := int64(id1)
	if id != at.RequestedByUserId {
		logger.Logger.Error("Wrong auth token.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter correct auth token."})
	}
	translator := en.New()
	uni := ut.New(translator, translator)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}
	v := validator.New()
	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}
	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
	err := v.Struct(at)
	var result []string
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			res := e.Translate(trans)
			result = append(result, res)
		}
		logger.Logger.WithError(err).WithField("error", result).Error("Validation error")
		return c.JSON(400, &models.Response{Status: false, ValidationErrors: result})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of request token is ", at)
	authResponse, _ := r.Usecase.RequestToken(ctx, at)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/********************************************Reject Request***********************************/
// @Tags wallet service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Reject/Accept request
// @Description reject/accept request
// @Accept  json
// @Produce  json
// @Param request_id query string true "request_id"
// @Param request_type query string true "request_type"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/wallet/reject_or_accept_request [get]
func (r *WalletController) RejectRequest(c echo.Context) error {
	requestId1, _ := strconv.Atoi(c.QueryParam("request_id"))
	requestId := int64(requestId1)
	requestType := c.QueryParam("request_type")
	if requestId == 0 || requestType == "" {
		logger.Logger.Error("Rquest id and request type are missing.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter request id and request type."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of reject request is ", logrus.Fields{"request_id": requestId, "request_type": requestType})
	authResponse, _ := r.Usecase.RejectRequest(ctx, requestId, requestType)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/**********************************************Send Token************************************/
// @Tags wallet service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Send Token
// @Description send token
// @Accept  json
// @Produce  json
// @Param send token body models.SendTokenDto true "send token"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/wallet/send_token [post]
func (r *WalletController) SendToken(c echo.Context) error {
	var send map[string]interface{}
	c.Bind(&send)
	if send == nil {
		logger.Logger.Error(constants.MissingReq)
		return c.JSON(400, &models.Response{Status: false, Msg: constants.MissingReq})
	}
	if send["from"] == "" || send["to"] == "" || send["amount"] == "" || send["password"] == "" {
		logger.Logger.Error("Either from,to,amount and password is missing.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter from, to, amount and password."})
	}
	if send["amount"] == 0 {
		logger.Logger.Error("Amount is zero.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Amount should be greater than zero."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of send token is ", logrus.Fields{"from": send["from"], "to": send["to"], "amount": send["amount"], "password": send["password"]})
	authResponse, _ := r.Usecase.SendToken(ctx, send)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*********************************************GetTransactions*******************************/
// @Tags wallet service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Get transactions
// @Description get transactions
// @Accept  json
// @Produce  json
// @Param wallet_id query string true "wallet_id"
// @Param page query int true "page"
// @Param limit query int true "limit"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/wallet/get_transactions [get]
func (r *WalletController) GetTransactions(c echo.Context) error {
	walletId := c.QueryParam("wallet_id")
	query := c.Request().URL.Query()
	if walletId == "" {
		logger.Logger.Error(constants.WalletIDError)
		return c.JSON(400, &models.Response{Status: false, Msg: constants.WalletIdReq})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of get transactions is ", logrus.Fields{"wallet_id": walletId, "query": query})
	authResponse, _ := r.Usecase.GetTransactions(ctx, walletId, query)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/******************************************view spend analytics**********************************/
// @Tags wallet service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary View spend analytics
// @Description view spend analytics
// @Accept  json
// @Produce  json
// @Param wallet_id query string true "wallet_id"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/wallet/view_spend_analytics [get]
func (r *WalletController) ViewSpendAnalytics(c echo.Context) error {
	walletId := c.QueryParam("wallet_id")
	if walletId == "" {
		logger.Logger.Error(constants.WalletIDError)
		return c.JSON(400, &models.Response{Status: false, Msg: constants.WalletIdReq})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of view spend analytics is ", logrus.Fields{"wallet_id": walletId})
	authResponse, _ := r.Usecase.ViewSpendAnalytics(ctx, walletId)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*******************************************Wallet Statement***********************************/

// @Tags wallet service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Wallet statement
// @Description wallet statement
// @Accept  json
// @Produce  json
// @Param wallet statement body models.WalletStatementDto true "wallet statement"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/wallet/wallet_statement [post]
func (r *WalletController) WalletStatement(c echo.Context) error {
	var req map[string]interface{}
	c.Bind(&req)
	if req == nil {
		logger.Logger.Error(constants.MissingReq)
		return c.JSON(400, &models.Response{Status: false, Msg: constants.MissingReq})
	}
	if req["query_type"] != "EMAIL" {
		if req["query_type"] != "DOWNLOAD" {
			logger.Logger.Error("Query type is wrong.")
			return c.JSON(400, &models.Response{Status: false, Msg: "Please use only EMAIL or DOWNLOAD in query_type."})
		}
	}
	if req["wallet_id"] == "" {
		logger.Logger.Error(constants.WalletIDError)
		return c.JSON(400, &models.Response{Status: false, Msg: constants.WalletIdReq})
	}
	if req["total_months"] == "" {
		if req["start_date"] == "" && req["end_date"] == "" {
			logger.Logger.Error("start_date and end_date are missing.")
			return c.JSON(400, &models.Response{Status: false, Msg: "Please enter start_date and end_date."})
		}
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of wallet statement is ", req)
	authResponse, err := r.Usecase.WalletStatement(ctx, req)
	if authResponse == "" {
		return c.JSON(http.StatusOK, &models.Response{Status: true, Msg: "You don't have any transactions for the given time period."})
	}
	if err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{Status: false, Msg: authResponse})
	}
	defer os.Remove(authResponse)
	if req["query_type"] == "DOWNLOAD" {
		return c.Attachment(authResponse, "")
	}
	return c.JSON(http.StatusOK, &models.Response{Status: true, Msg: "Mail sent to your registered email."})
}

/****************************************Get WalletId by userid***************************/
// @Tags wallet service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Get wallet id
// @Description get wallet id
// @Accept  json
// @Produce  json
// @Param user_id query int64 true "user_id"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/wallet/get_wallet_id [get]
func (r *WalletController) GetWalletId(c echo.Context) error {
	userid, _ := strconv.Atoi(c.QueryParam("user_id"))
	userId := int64(userid)
	id1, _ := strconv.Atoi(c.Request().Header.Get("id"))
	logger.Logger.Info("id ", id1)
	id := int64(id1)
	if userId == 0 {
		logger.Logger.Error("User id is missing.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter user id."})
	}
	if id != userId {
		logger.Logger.Error("Wrong auth token.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter correct auth token."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of get wallet id is ", logrus.Fields{"user_id": userId})
	authResponse, _ := r.Usecase.GetWalletId(ctx, userId)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/****************************************Recent Transactions***********************************/
// @Tags wallet service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Recent Transactions
// @Description recent transactions
// @Accept  json
// @Produce  json
// @Param wallet_id query string true "wallet_id"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/wallet/get_recent_transactions [get]
func (r *WalletController) RecentTransactions(c echo.Context) error {
	walletId := c.QueryParam("wallet_id")
	query := c.Request().URL.Query()
	if walletId == "" {
		logger.Logger.Error(constants.WalletIDError)
		return c.JSON(400, &models.Response{Status: false, Msg: constants.WalletIdReq})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of get recent transaction is ", logrus.Fields{"wallet_id": walletId})
	authResponse, _ := r.Usecase.RecentTransactions(ctx, walletId, query)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/**************************************ShowPendingRequests************************************/
// @Tags wallet service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Show pending requests
// @Description show pending requests
// @Accept  json
// @Produce  json
// @Param wallet_id query string true "wallet_id"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/wallet/show_pending_requests [get]
func (r *WalletController) ShowPendingRequests(c echo.Context) error {
	walletId := c.QueryParam("wallet_id")
	query := c.Request().URL.Query()
	id1, _ := strconv.Atoi(c.Request().Header.Get("id"))
	logger.Logger.Info("id ", id1)
	id := int64(id1)
	if walletId == "" {
		logger.Logger.Error(constants.WalletIDError)
		return c.JSON(400, &models.Response{Status: false, Msg: constants.WalletIdReq})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of show all requests is ", logrus.Fields{"wallet_id": walletId, "user_id": id})
	authResponse, _ := r.Usecase.ShowPendingRequests(ctx, walletId, id, query)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/************************************SendTokenToAdminByArtist****************************/
func (r *WalletController) SendTokenToAdmin(c echo.Context) error {
	var req map[string]interface{}
	c.Bind(&req)
	id1, _ := strconv.Atoi(c.Request().Header.Get("id"))
	logger.Logger.Info("id ", id1)
	id := int64(id1)
	if id == 0 {
		logger.Logger.Error("id not found from auth token.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter auth token."})
	}
	concert_id1, _ := strconv.Atoi(fmt.Sprintf("%v", req["concert_id"]))
	concertId := int64(concert_id1)
	if concertId == 0 {
		logger.Logger.Error("Concert id is missing.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter concert id."})
	}
	amount := fmt.Sprintf("%v", req["amount"])
	password := fmt.Sprintf("%v", req["password"])
	if amount == "" || password == "" {
		logger.Logger.Error("Amount and password are missing.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter amount and password."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of send token to admin is ", logrus.Fields{"user_id": id, "concert_id": concertId, "amount": amount, "password": password})
	authResponse, _ := r.Usecase.SendTokenToAdmin(ctx, id, concertId, amount, password)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/**************************************ShowOwnTokenRequests************************************/
// @Tags wallet service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Show own token requests
// @Description show own token requests
// @Accept  json
// @Produce  json
// @Param wallet_id query string true "wallet_id"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/wallet/show_own_token_requests [get]
func (r *WalletController) ShowOwnTokenRequests(c echo.Context) error {
	walletId := c.QueryParam("wallet_id")
	query := c.Request().URL.Query()
	id1, _ := strconv.Atoi(c.Request().Header.Get("id"))
	logger.Logger.Info("id ", id1)
	id := int64(id1)
	if walletId == "" {
		logger.Logger.Error(constants.WalletIDError)
		return c.JSON(400, &models.Response{Status: false, Msg: constants.WalletIdReq})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of show own token requests is ", logrus.Fields{"wallet_id": walletId, "user_id": id})
	authResponse, _ := r.Usecase.ShowOwnTokenRequests(ctx, walletId, id, query)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}
