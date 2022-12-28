package controller

import (
	"context"
	"log"
	"mime/multipart"
	"net/http"
	"pandoMessagingWalletService/com.pando.messaging/constants"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	ipfs "pandoMessagingWalletService/com.pando.messaging/service"
	"strconv"
	"strings"

	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type IPFSController struct {
	Usecase ipfs.IPFSUsecase
}

/*************************************Upload file to IPFS and AWS*************************/
// @Tags ipfs service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary UploadFile
// @Description UploadFile
// @Accept  mpfd
// @Produce  json
// @Param file formData file true "file"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/ipfs [post]
func (r *IPFSController) UploadFile(c echo.Context) error {

	files, handler, err := c.Request().FormFile("file")
	if err != nil {
		logger.Logger.Error("Please provide data in file field.")
		return c.JSON(http.StatusBadRequest, &models.Response{Status: false, ResponseCode: 400, Msg: "Please provide data in file field."})
	}
	if strings.ContainsAny(handler.Filename, "!@#$%^&*+=") {
		return c.JSON(400, &models.Response{Status: false, ResponseCode: 400, Msg: "file name is not correct."})
	}
	//defer files.Close()
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from upload file to ipfs is ", logrus.Fields{"files": files})
	authResponse, _ := r.Usecase.UploadFile(ctx, files, handler)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	defer files.Close()
	return c.JSON(http.StatusOK, authResponse)
}

/************************************Save Backup filehash*********************************/
// @Tags ipfs service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary SaveBackupFilehash
// @Description SaveBackupFilehash
// @Accept  json
// @Produce  json
// @Param SaveBackupFilehash body models.BackupDto true "SaveBackupFilehash"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/ipfs/save_backup_filehash [post]
func (r *IPFSController) SaveBackupFilehash(c echo.Context) error {
	var req models.Backups
	c.Bind(&req)
	id1, _ := strconv.Atoi(c.Request().Header.Get("id"))
	logger.Logger.Info("id ", id1)
	id := int64(id1)
	if id != req.UserId {
		logger.Logger.Error("Wrong auth token.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter correct auth token."})
	}
	if req.BackupNature != constants.Daily && req.BackupNature != constants.Weekly && req.BackupNature != constants.Monthly {
		logger.Logger.Error("Please use correct backup nature from these ['DAILY','WEEKLY','MONTHLY']")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: 400, Msg: "Please use correct backup nature from these ['DAILY','WEEKLY','MONTHLY']"})
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
	err := v.Struct(req)
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
	logger.Logger.Info("Request received from save backup filehash is ", req)
	authResponse, _ := r.Usecase.SaveBackupFilehash(ctx, req)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/******************************************Get Backup Filehash****************************/
// @Tags ipfs service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary GetBackupFilehash
// @Description GetBackupFilehash
// @Accept  json
// @Produce  json
// @Param user_id query int64 true "user_id"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/ipfs/get_backup_filehash [get]
func (r *IPFSController) GetBackupFilehash(c echo.Context) error {
	userId1, _ := strconv.Atoi(c.QueryParam("user_id"))
	userId := int64(userId1)
	id1, _ := strconv.Atoi(c.Request().Header.Get("id"))
	logger.Logger.Info("id ", id1)
	id := int64(id1)
	if id != userId {
		logger.Logger.Error("Wrong auth token.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter correct auth token."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from get backup filehash is ", logrus.Fields{"user_id": userId})
	authResponse, _ := r.Usecase.GetBackupFilehash(ctx, userId)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/**************************************UploadFileToIPFS**********************************/
// @Tags ipfs service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary UploadFileToIPFS
// @Description UploadFileToIPFS
// @Accept  mpfd
// @Produce  json
// @Param file formData file true "file"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/ipfs [post]
func (r *IPFSController) UploadFileToIPFS(c echo.Context) error {
	data := c.FormValue("file")
	var files multipart.File
	var handler *multipart.FileHeader
	var err error
	if data == "" {
		files, handler, err = c.Request().FormFile("file")
		if err != nil {
			logger.Logger.Error("Please provide data in file field.")
			return c.JSON(http.StatusBadRequest, &models.Response{Status: false, ResponseCode: 400, Msg: "Please provide data in file field."})
		}
		if strings.ContainsAny(handler.Filename, "!@#$%^&*+=") {
			return c.JSON(400, &models.Response{Status: false, ResponseCode: 400, Msg: "file name is not correct."})
		}
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from upload file to ipfs is ", logrus.Fields{"files": files})
	authResponse, _ := r.Usecase.UploadFileToIPFS(ctx, files, handler, data)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}
