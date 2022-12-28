package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"pandoMessagingWalletService/com.pando.messaging/config"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	group "pandoMessagingWalletService/com.pando.messaging/service"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

type GroupController struct {
	Usecase group.GroupUsecase
	Config  *config.Config
}

//var conf *config.Config

/*********************************************Create Group******************************************/
// @Tags group service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Create_Group
// @Description Create_Group
// @Accept  json
// @Produce  json
// @Param Create_Group body models.GroupDto true "Create_Group"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/group/create_group [post]
func (r *GroupController) Create_Group(c echo.Context) error {
	var f models.Group
	c.Bind(&f)
	newFlow := models.Group{
		Group_name:        f.Group_name,
		Admin_ids:         f.Admin_ids,
		Subject_Timestamp: f.Subject_Timestamp,
		Subject_Owner_Id:  f.Subject_Owner_Id,
		Profile_Pic_Url:   f.Profile_Pic_Url,
		User_ids:          f.User_ids,
	}
	logger.Logger.Info("Requests in controller part for create group is ", newFlow)
	translator := en.New()
	uni := ut.New(translator, translator)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}
	v := validator.New()
	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		logger.Logger.Error(err)
		log.Fatal(err)
	}
	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
	err := v.Struct(newFlow)
	var result []string
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			res := e.Translate(trans)
			result = append(result, res)
		}
		logger.Logger.WithError(err).WithField("err", err).Errorf("validation_errors", result)
		return c.JSON(400, &models.Response{Status: false, ResponseCode: 400, ValidationErrors: result})

	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Requests in controller part for create group is ", newFlow)
	authResponse, _ := r.Usecase.Create_Group(ctx, newFlow)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/***********************************************Add Users to Group*********************************/
// @Tags group service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary AddUserToGroup
// @Description AddUserToGroup
// @Accept  json
// @Produce  json
// @Param AddUserToGroup body models.AddUserToGroupDto true "AddUserToGroup"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/group/add_users_to_group [post]
func (r *GroupController) AddUserToGroup(c echo.Context) error {
	var f models.Groups
	c.Bind(&f)
	newFlow := models.Groups{
		ID:         f.ID,
		Admin_ids:  f.Admin_ids,
		TotalUsers: int64(len(f.User_ids)),
		User_ids:   f.User_ids,
	}
	if f.ID == 0 {
		logger.Logger.WithField("error", f.ID).Error("User id is missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: http.StatusBadRequest, Msg: "Please enter group id"})
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
	err := v.Struct(newFlow)
	var result []string
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			res := e.Translate(trans)
			result = append(result, res)
		}
		logger.Logger.WithError(err).WithField("err", err).Errorf("validation_errors", result)
		return c.JSON(400, &models.Response{Status: false, ResponseCode: 400, ValidationErrors: result})

	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Requests in controller part for add users in group is ", newFlow)
	authResponse, _ := r.Usecase.AddUserToGroup(ctx, newFlow)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/**********************************************Remove Users from group*****************************/
// @Tags group service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary RemoveUsersFromGroup
// @Description RemoveUsersFromGroup
// @Accept  json
// @Produce  json
// @Param group_id query int64 true "group_id"
// @Param user_id query int64 true "user_id"
// @Param admin_id query int64 true "admin_id"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/group/remove_user_from_group [get]
func (r *GroupController) RemoveUsersFromGroup(c echo.Context) error {
	group_id1, _ := strconv.Atoi(c.QueryParam("group_id"))
	user_id1, _ := strconv.Atoi(c.QueryParam("user_id"))
	admin_id1, _ := strconv.Atoi(c.QueryParam("admin_id"))
	admin_id := int64(admin_id1)
	user_id := int64(user_id1)
	group_id := int64(group_id1)
	if group_id1 == 0 || user_id1 == 0 || admin_id1 == 0 {
		logger.Logger.Error("Group id, admin id, and user id are missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: http.StatusBadRequest, Msg: "Please mention all feilds group id or user_id and admin id."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request in controller part for remove users from group is ", logrus.Fields{"group_id": group_id, "admin_id": admin_id, "user_id": user_id})
	authResponse, _ := r.Usecase.RemoveUsersFromGroup(ctx, group_id, user_id, admin_id)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/**********************************************LeaveGroup********************************************/
// @Tags group service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary LeaveGroup
// @Description LeaveGroup
// @Accept  json
// @Produce  json
// @Param group_id query int64 true "group_id"
// @Param user_id query int64 true "user_id"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/group/leave_group [get]
func (r *GroupController) LeaveGroup(c echo.Context) error {
	group_id1, _ := strconv.Atoi(c.QueryParam("group_id"))
	user_id1, _ := strconv.Atoi(c.QueryParam("user_id"))
	user_id := int64(user_id1)
	group_id := int64(group_id1)
	if group_id1 == 0 || user_id1 == 0 {
		logger.Logger.Error("Group id and user id is missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: http.StatusBadRequest, Msg: "Please enter group id and user id."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request in controller part for leave group is ", logrus.Fields{"group_id": group_id, "user_id": user_id})
	authResponse, _ := r.Usecase.LeaveGroup(ctx, group_id, user_id)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/***********************************************Edit Group Info*************************************/
// @Tags group service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary EditGroupInfo
// @Description EditGroupInfo
// @Accept  json
// @Produce  json
// @Param group_id query int64 true "group_id"
// @Param new_group_name query string true "new_group_name"
// @Param profile_pic_url query string true "profile_pic_url"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/group/edit_group_info [post]
func (r *GroupController) EditGroupInfo(c echo.Context) error {
	group_id1, _ := strconv.Atoi(c.QueryParam("group_id"))
	new_group_name := c.QueryParam("new_group_name")
	profile_pic_url := c.QueryParam("profile_pic_url")
	group_id := int64(group_id1)
	if group_id1 == 0 || new_group_name == "" || profile_pic_url == "" {
		logger.Logger.Error("Group_id, new_group_name and profile_pic_url is missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: 400, Msg: "Please enter group_id, new_group_name and profile_pic_url"})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request in controller part of edit group name is ", logrus.Fields{"group_id": group_id, "group_name": new_group_name, "profile_pic_url": profile_pic_url})
	authResponse, _ := r.Usecase.EditGroupInfo(ctx, group_id, new_group_name, profile_pic_url)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/**********************************************SearchUsersInGroup***********************************/

func (r *GroupController) SearchUsersInGroup(c echo.Context) error {
	group_id := c.QueryParam("group_id")
	username := c.QueryParam("username")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	authResponse, _ := r.Usecase.SearchUsersInGroup(ctx, group_id, username)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/**********************************************Delete Group*****************************************/
// @Tags group service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary DeleteGroup
// @Description DeleteGroup
// @Accept  json
// @Produce  json
// @Param group_id query int64 true "group_id"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/group/delete_group [get]
func (r *GroupController) DeleteGroup(c echo.Context) error {
	group_id1, _ := strconv.Atoi(c.QueryParam("group_id"))
	if group_id1 == 0 {
		logger.Logger.Error("Group id is missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: http.StatusBadRequest, Msg: "Please enter group id."})
	}
	group_id := int64(group_id1)
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request from controller part in delete group is ", logrus.Fields{"group_id": group_id})
	authResponse, _ := r.Usecase.DeleteGroup(ctx, group_id)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/**********************************************UploadProfilePhoto***********************************/
// @Tags group service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary SearchUsersInGroup
// @Description SearchUsersInGroup
// @Accept  mpfd
// @Produce  json
// @Param photo formData file true "photo"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/group/search_user_in_group [post]
func (r *GroupController) UploadGroupProfilePhoto(c echo.Context) error {
	maxSize := r.Config.ImageSize // allow only 1MB of file size
	fmt.Println(maxSize, "hbbhbd")
	file, handler, err := c.Request().FormFile("photo")
	if err != nil {
		logger.Logger.WithError(err).WithField("err", err).Errorf("File not found.", file)
		return c.JSON(400, &models.Response{Status: false, ResponseCode: http.StatusBadRequest, Msg: "Please select an image."})
	}
	logger.Logger.Info("File data", file)
	size := c.Request().ParseMultipartForm(maxSize)
	if size != nil {
		logger.Logger.Info("image size is to large")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: 400, Msg: "Image size is too large."})

	}
	defer file.Close()
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(r.Config.S3Region),
		Credentials: credentials.NewStaticCredentials(
			r.Config.S3AccessId,  // id
			r.Config.S3SecretKey, // secret
			""),              // token can be left blank for now
	})
	if err != nil {
		logger.Logger.WithError(err).WithField("err", err).Errorf("Could not upload file")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: http.StatusBadRequest, Msg: "Could not upload file."})
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	fileName, err := config.UploadFileToS3(s, fileBytes, "profile/"+handler.Filename, handler)
	if err != nil {
		logger.Logger.WithError(err).WithField("err", err).Errorf("Could not upload file.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: http.StatusBadRequest, Msg: "Could not upload file."})
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if ctx == nil {
		ctx = context.Background()
	}

	authResponse, _ := r.Usecase.UploadGroupProfilePhoto(ctx, file, handler, fileName)
	if authResponse == nil {

		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*******************************************Make or Remove Admin*******************************/
// @Tags group service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary MakeOrRemoveAdmin
// @Description MakeOrRemoveAdmin
// @Accept  json
// @Produce  json
// @Param group_id query int64 true "group_id"
// @Param new_admin_id query int64 true "new_admin_id"
// @Param user_id query int64 true "user_id"
// @Param method_type query string true "method_type"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/group/make_admin [post]
func (r *GroupController) MakeOrRemoveAdmin(c echo.Context) error {
	userid1, _ := strconv.Atoi(c.QueryParam("user_id"))
	user_id := int64(userid1)
	group_id1, _ := strconv.Atoi(c.QueryParam("group_id"))
	group_id := int64(group_id1)
	new_admin_id1, _ := strconv.Atoi(c.QueryParam("new_admin_id"))
	new_admin_id := int64(new_admin_id1)
	method_type := c.QueryParam("method_type")
	if method_type == "" {
		logger.Logger.Error("Type is missing")
		return c.JSON(400, models.Response{Status: false, ResponseCode: 400, Msg: "Please enter method type."})
	}
	if user_id == 0 || group_id == 0 || new_admin_id == 0 {
		logger.Logger.Error("User id, group id and admin id are missing")
		return c.JSON(400, models.Response{Status: false, ResponseCode: 400, Msg: "Please enter user id and group id and new admin id."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request from controller part of fetch status is ", logrus.Fields{"user_id": user_id})
	authResponse, _ := r.Usecase.MakeOrRemoveAdmin(ctx, user_id, group_id, new_admin_id, method_type)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/***************************************Get group details by id *****************************/
// @Tags group service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary GetGroupDetails
// @Description GetGroupDetails
// @Accept  json
// @Produce  json
// @Param group_id query int64 true "group_id"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/group/get_group_details [get]
func (r *GroupController) GetGroupDetails(c echo.Context) error {
	group_id1, _ := strconv.Atoi(c.QueryParam("group_id"))
	group_id := int64(group_id1)
	if group_id1 == 0 {
		logger.Logger.Error("Group id is missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: 400, Msg: "Please enter group id."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request in controller part of get group details ", logrus.Fields{"group_id": group_id})
	authResponse, _ := r.Usecase.GetGroupDetails(ctx, group_id)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/****************************************Accept/Decline invitation*******************************/
// @Tags group service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary AcceptAndDeclineGroupInvitation
// @Description AcceptAndDeclineGroupInvitation
// @Accept  json
// @Produce  json
// @Param AcceptAndDeclineGroupInvitation body models.AcceptGroupInvitationDto true "AcceptAndDeclineGroupInvitation"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/group/accept_and_decline_group_invitation [post]
func (r *GroupController) AcceptAndDeclineGroupInvitation(c echo.Context) error {
	var g models.AcceptGroupInvitationDto
	c.Bind(&g)
	group := models.AcceptGroupInvitationDto{
		GroupId: g.GroupId,
		UserId:  g.UserId,
		Type:    g.Type,
	}
	if group.GroupId == 0 || group.UserId == 0 {
		logger.Logger.Error("Group id and user_id is missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: 400, Msg: "Please enter group id and user id."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request in controller part of accept and decline group invitation", logrus.Fields{"group_id": g.GroupId, "user_id": g.UserId})
	authResponse, _ := r.Usecase.AcceptAndDeclineGroupInvitation(ctx, group)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/****************************************Get All Group details of a user*****************************/
// @Tags group service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary GetAllGroupDetailsOfUser
// @Description GetAllGroupDetailsOfUser
// @Accept  json
// @Produce  json
// @Param user_id query int64 true "user_id"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/group/get_all_group_details_of_user [get]
func (r *GroupController) GetAllGroupDetailsOfUser(c echo.Context) error {
	user_id1, _ := strconv.Atoi(c.QueryParam("user_id"))
	user_id := int64(user_id1)
	if user_id1 == 0 {
		logger.Logger.Error("User id is missing.")
		return c.JSON(400, &models.Response{Status: false, ResponseCode: http.StatusBadRequest, Msg: "Please enter user id."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request from controller part of get all group details of user is ", logrus.Fields{"user__id": user_id})
	authResponse, _ := r.Usecase.GetAllGroupDetailsOfUser(ctx, user_id)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}
