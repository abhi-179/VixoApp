package controller

import (
	"context"
	"log"
	"net/http"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	chat "pandoMessagingWalletService/com.pando.messaging/service"
	"strconv"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

type ChatController struct {
	Usecase chat.ChatUsecase
}

/******************************************Post Status****************************************/
func (r *ChatController) PostStatus(c echo.Context) error {
	var ps models.Status
	c.Bind(&ps)
	status := models.Status{}
	if ps.StatusType == "text" {
		status = models.Status{
			ID:         ps.ID,
			CreatedAt:  ps.CreatedAt,
			User_Id:    ps.User_Id,
			StatusType: ps.StatusType,
			Message:    ps.Message,
		}
		if status.Message == "" {
			logger.Logger.Error("Message is missing.")
			return c.JSON(400, models.Response{Status: false, ResponseCode: 400, Msg: "Please enter message."})
		}

	} else if ps.StatusType == "media" {
		var filehash []models.FileHash
		var awsurl []models.AwsUrl
		for i := 0; i < len(ps.File_hash); i++ {
			filehash = append(filehash, ps.File_hash[i])
			awsurl = append(awsurl, ps.AwsUrl[i])
		}
		status = models.Status{
			ID:         ps.ID,
			CreatedAt:  ps.CreatedAt,
			User_Id:    ps.User_Id,
			StatusType: ps.StatusType,
			File_hash:  filehash,
			AwsUrl:     awsurl,
		}
		if len(status.File_hash) == 0 {
			logger.Logger.Error("Filehash is missing.")
			return c.JSON(400, models.Response{Status: false, ResponseCode: 400, Msg: "Please enter filehash."})
		}
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
	err := v.Struct(status)
	var result []string
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			res := e.Translate(trans)
			result = append(result, res)
		}
		logger.Logger.WithError(err).WithField("error", err).Error("Validation error")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: 400, ValidationErrors: result})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request from controller part of post status is ", status)
	authResponse, _ := r.Usecase.PostStatus(ctx, status)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/******************************************Fetch Status****************************************/
func (r *ChatController) FetchStatus(c echo.Context) error {
	userid1, _ := strconv.Atoi(c.QueryParam("user_id"))
	user_id := int64(userid1)
	if user_id == 0 {
		logger.Logger.Error("User id is missing")
		return c.JSON(400, models.Response{Status: false, ResponseCode: 400, Msg: "Please enter user id."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request from controller part of fetch status is ", logrus.Fields{"user_id": user_id})
	authResponse, _ := r.Usecase.FetchStatus(ctx, user_id)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/***************************************Delete Status*******************************************/
func (r *ChatController) DeleteStatus(c echo.Context) error {
	status_id := c.Request().URL.Query()
	c.Bind(&status_id)
	user_id1, _ := strconv.Atoi(c.QueryParam("user_id"))
	user_id := int64(user_id1)

	newstring := strings.Split(status_id.Get("status_id"), ",")
	var t2 = []int64{}
	for _, i := range newstring {
		j, _ := strconv.Atoi(i)
		t2 = append(t2, int64(j))
	}
	if len(t2) == 0 {
		logger.Logger.Error("Message id is missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: 400, Msg: "Please enter message id."})
	}
	logger.Logger.Debug(t2)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request from controller part of delete status is ", logrus.Fields{"status_id": status_id, "user_id": user_id})
	authResponse, _ := r.Usecase.DeleteStatus(ctx, user_id, t2)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/********************************************Search Status by username***************************/
func (r *ChatController) SearchStatusByUsername(c echo.Context) error {
	username := c.QueryParam("username")
	if username == "" {
		logger.Logger.Error("Username is missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: 400, Msg: "Please enter username."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from SearchStatusByUsername controller is", username)
	authresponse, _ := r.Usecase.SearchStatusByUsername(ctx, username)
	if authresponse == nil {
		return c.JSON(http.StatusUnauthorized, authresponse)
	}
	return c.JSON(http.StatusOK, authresponse)
}

/***************************************Report chat********************************************/
// @Tags chat service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary report chat
// @Description report chat
// @Accept  json
// @Produce  json
// @Param report chat body models.ReportDto true "report chat"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/chat/report_chat [post]
func (r *ChatController) ReportChat(c echo.Context) error {
	var m models.Reports
	c.Bind(&m)
	re := models.Reports{
		Reporter_Id: m.Reporter_Id,
		Reportee_Id: m.Reportee_Id,
		GroupId:     m.GroupId,
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request from controller part of report chat is ", logrus.Fields{"reporter_id": re.Reporter_Id, "reportee_id": re.Reportee_Id})
	authResponse, _ := r.Usecase.ReportChat(ctx, re)

	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/********************************************Save Block contact Details**************************/
// @Tags chat service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary SaveBlockUserDetails
// @Description SaveBlockUserDetails
// @Accept  json
// @Produce  json
// @Param SaveBlockUserDetails body models.BlockedContactDto true "SaveBlockUserDetails"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/chat/save_blocked_contacts [post]
func (r *ChatController) SaveBlockUserDetails(c echo.Context) error {
	var m models.BlockedContacts
	c.Bind(&m)
	re := models.BlockedContacts{
		BlockeeId: m.BlockeeId,
		BlockerId: m.BlockerId,
	}
	if m.BlockeeId == 0 || m.BlockerId == 0 {
		logger.Logger.Error("Blockee id and Blocker id are missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: http.StatusBadRequest, Msg: "Please enter blockee and blocker ids."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request from controller part of save block user details is ", re)
	authResponse, _ := r.Usecase.SaveBlockUserDetails(ctx, re)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*************************************Fetch Blocked contacts Details***************************************/
// @Tags chat service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary FetchBlockedUserDetails
// @Description FetchBlockedUserDetails
// @Accept  json
// @Produce  json
// @Param user_id query int64 true "user_id"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/chat/fetch_blocked_users [get]
func (r *ChatController) FetchBlockedUserDetails(c echo.Context) error {
	user_id1, _ := strconv.Atoi(c.QueryParam("user_id"))
	user_id := int64(user_id1)
	if user_id == 0 {
		logger.Logger.Error("User id is missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: http.StatusBadRequest, Msg: "Please enter user id."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request from controller part of fetch blocked contacts details is ", logrus.Fields{"user_id": user_id})
	authResponse, _ := r.Usecase.FetchBlockedUserDetails(ctx, user_id)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*************************************Fetch Blocked contacts Details for specific users***************************************/
// @Tags chat service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary FetchBlockedContactDetails
// @Description FetchBlockedContactDetails
// @Accept  json
// @Produce  json
// @Param user_id query int64 true "user_id"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/chat/fetch_blocked_contacts [get]
func (r *ChatController) FetchBlockedContactDetails(c echo.Context) error {
	user_id1, _ := strconv.Atoi(c.QueryParam("user_id"))
	user_id := int64(user_id1)
	if user_id == 0 {
		logger.Logger.Error("User id is missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: http.StatusBadRequest, Msg: "Please enter user id."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request from controller part of fetch blocked contacts is ", logrus.Fields{"user_id": user_id})
	authResponse, _ := r.Usecase.FetchBlockedContactDetails(ctx, user_id)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*************************************Unblock user***************************************/
// @Tags chat service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Unblock_user
// @Description Unblock_user
// @Accept  json
// @Produce  json
// @Param Unblock_user body models.BlockedContactDto true "Unblock_user"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/chat/unblock_user [post]
func (r *ChatController) Unblock_user(c echo.Context) error {
	var m models.BlockedContacts
	c.Bind(&m)
	re := models.BlockedContacts{
		BlockeeId: m.BlockeeId,
		BlockerId: m.BlockerId,
	}
	if m.BlockeeId == 0 || m.BlockerId == 0 {
		logger.Logger.Error("Blockee id and Blocker id are missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: http.StatusBadRequest, Msg: "Please enter blockee and blocker ids."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request from controller part of unblock blocked contacts is ", re)
	authResponse, _ := r.Usecase.Unblock_user(ctx, re)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*************************************Fetch Wallpapers Details***************************************/
// @Tags chat service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary FetchWallpapersDetails
// @Description FetchWallpapersDetails
// @Accept  json
// @Produce  json
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/chat/fetch_wallpaper_details [get]
func (r *ChatController) FetchWallpapersDetails(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	authResponse, _ := r.Usecase.FetchWallpapersDetails(ctx)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*****************************************SaveGroupChatSetting**********************************/
// @Tags chat service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary SaveGroupChatSetting
// @Description SaveGroupChatSetting
// @Accept  json
// @Produce  json
// @Param SaveGroupChatSetting body models.ChatSettingDto true "SaveGroupChatSetting"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/chat/save_group_chat_setting [post]
func (r *ChatController) SaveGroupChatSetting(c echo.Context) error {
	var g models.ChatSettings
	c.Bind(&g)
	chat := models.ChatSettings{
		User_Id:         g.User_Id,
		Group_Chat_Type: g.Group_Chat_Type,
	}

	if g.User_Id == 0 || g.Group_Chat_Type == "" {
		logger.Logger.Error("User id and group chat type is missing.")
		return c.JSON(400, models.Response{Status: false, ResponseCode: 400, Msg: "Please enter user id and group_chat_type."})
	}
	if g.Group_Chat_Type == "PUBLIC" {
		ctx := c.Request().Context()
		if ctx == nil {
			ctx = context.Background()
		}
		logger.Logger.Info("Request received from save chat setting is ", chat)
		authResponse, _ := r.Usecase.SaveGroupChatSetting(ctx, chat)
		if authResponse == nil {
			return c.JSON(http.StatusUnauthorized, authResponse)
		}
		return c.JSON(http.StatusOK, authResponse)
	} else {
		logger.Logger.Error("Group chat type is not correct.")
		return c.JSON(400, models.Response{Status: false, ResponseCode: 400, Msg: "Please enter correct group chat type i.e PUBLIC."})
	}
}

/********************************************Fetch Group chat setting*************************/
// @Tags chat service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary FetchGroupChatSetting
// @Description FetchGroupChatSetting
// @Accept  json
// @Produce  json
// @Param user_id query int64 true "user_id"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/chat/fetch_group_chat_setting [get]
func (r *ChatController) FetchGroupChatSetting(c echo.Context) error {
	user_id1, _ := strconv.Atoi(c.QueryParam("user_id"))
	user_id := int64(user_id1)
	if user_id == 0 {
		logger.Logger.Error("User id is missing.")
		return c.JSON(400, models.Response{Status: false, ResponseCode: 400, Msg: "Please enter user id."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from fetch group chat setting is ", logrus.Fields{"user_id": user_id})
	authResponse, _ := r.Usecase.FetchGroupChatSetting(ctx, user_id)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}
