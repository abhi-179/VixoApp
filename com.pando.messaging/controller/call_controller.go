package controller

import (
	"context"
	"log"
	"net/http"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	call "pandoMessagingWalletService/com.pando.messaging/service"
	"strconv"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

type CallController struct {
	Usecase call.CallUsecase
}

/***********************************************Save Call logs**************************************/
// SaveCallLogs
//@Tags call service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Save call logs
// @Description save call logs
// @Accept  json
// @Produce  json
// @Param call body models.CallDetailDto true "Add call"
// @Success 200 {object} interface{}
// @Header 200 {string} Token "qwerty"
// @Failure 404  "Not Found"
// @Router /api/v1/call/save_call_details [post]
func (r *CallController) SaveCallLogs(c echo.Context) error {
	var s models.CallDetail
	c.Bind(&s)
	callLogs := models.CallDetail{
		CallerId:     s.CallerId,
		CallDuration: s.CallDuration,
		Filehash:     s.Filehash,
		AwsUrl:       s.AwsUrl,
		StartTime:    s.StartTime,
		EndTime:      s.EndTime,
		User_ids:     s.User_ids,
		IsAudioCall:  s.IsAudioCall,
		IsMissedCall: s.IsMissedCall,
		IsGroupCall:  s.IsGroupCall,
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
	err := v.Struct(callLogs)
	var result []string
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			res := e.Translate(trans)
			result = append(result, res)
		}
		logger.Logger.WithError(err).WithField("error", result).Error("Validation error")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: 400, ValidationErrors: result})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request from controller part of save call details is ", callLogs)
	authResponse, _ := r.Usecase.SaveCallLogs(ctx, callLogs)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/**********************************************Fetch Call logs**************************************/
// FetchAllCallLogs
//@Tags call service
// @Security ApiKeyAuth
// @Security ApiKeyAuth Authorization
// @Summary Fetch all call logs
// @Description fetch all call logs
// @Accept  json
// @Produce  json
// @Param user_id query int64 true "user_id"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Success 200 {object} interface{}
// @Failure 404  "Not Found"
// @Router /api/v1/call/fetch_all_call_details [get]
func (r *CallController) FetchAllCallLogs(c echo.Context) error {
	user_id1, _ := strconv.Atoi(c.QueryParam("user_id"))
	user_id := int64(user_id1)
	source := c.QueryParam("source")
	query := c.Request().URL.Query()
	if user_id == 0 {
		logger.Logger.Error("User id is missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: http.StatusBadRequest, Msg: "Please enter user id."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request from controller part of fetch all call logs is ", logrus.Fields{"user_id": user_id, "source": source})
	authResponse, _ := r.Usecase.FetchAllCallLogs(ctx, user_id, source, query)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/**********************************************Fetch MissedCall*************************************/
// FetchMissedCallLogs
//@Tags call service
// @Security ApiKeyAuth
// @Security ApiKeyAuth Authorization
// @Summary Fetch missed call logs
// @Description fetch missed call logs
// @Accept  json
// @Produce  json
// @Param user_id query int64 true "user_id"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Success 200 {object} interface{}
// @Failure 404  "Not Found"
// @Router /api/v1/call/fetch_missed_call_details [get]
func (r *CallController) FetchMissedCallLogs(c echo.Context) error {
	user_id1, _ := strconv.Atoi(c.QueryParam("user_id"))
	user_id := int64(user_id1)
	source := c.QueryParam("source")
	query := c.Request().URL.Query()
	if user_id == 0 {
		logger.Logger.Error("User id is missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: http.StatusBadRequest, Msg: "Please enter user id."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request from controller part of fetch missed call logs is ", logrus.Fields{"user_id": user_id, "source": source})
	authResponse, _ := r.Usecase.FetchMissedCallLogs(ctx, user_id, source, query)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/***********************************************Delete Call Logs***********************************/
// DeleteCallLogs
//@Tags call service
// @Security ApiKeyAuth
// @Security ApiKeyAuth Authorization
// @Summary Delete call logs
// @Description Delete call logs
// @Accept  json
// @Produce  json
// @Param user_id query int64 true "user_id"
// @Success 200 {object} interface{}
// @Failure 404  "Not Found"
// @Router /api/v1/call/delete_call_details [get]
func (r *CallController) DeleteCallLogs(c echo.Context) error {
	user_id1, _ := strconv.Atoi(c.QueryParam("user_id"))
	user_id := int64(user_id1)
	id1, _ := strconv.Atoi(c.QueryParam("id"))
	id := int64(id1)
	if user_id == 0 || id == 0 {
		logger.Logger.Error("User id and id is missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: 400, Msg: "Please provide user_id and id"})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request from controller part of delete call logs is ", logrus.Fields{"user_id": user_id, "id": id})
	authResponse, _ := r.Usecase.DeleteCallLogs(ctx, user_id, id)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}
