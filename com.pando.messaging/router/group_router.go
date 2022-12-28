package router

import (
	"pandoMessagingWalletService/com.pando.messaging/config"
	controller "pandoMessagingWalletService/com.pando.messaging/controller"
	group "pandoMessagingWalletService/com.pando.messaging/service"

	"github.com/labstack/echo/v4"
)

func NewGroupController(e *echo.Echo, groupusecase group.GroupUsecase, conf *config.Config) {
	handler := &controller.GroupController{
		Usecase: groupusecase,
		Config:  conf,
	}

	e.POST("api/v1/group/create_group", handler.Create_Group)
	e.POST("api/v1/group/add_users_to_group", handler.AddUserToGroup)
	e.GET("api/v1/group/remove_user_from_group", handler.RemoveUsersFromGroup)
	e.GET("api/v1/group/leave_group", handler.LeaveGroup)
	e.POST("api/v1/group/edit_group_info", handler.EditGroupInfo)
	e.GET("api/v1/group/delete_group", handler.DeleteGroup)
	e.POST("api/v1/group/upload_profile_photo", handler.UploadGroupProfilePhoto)
	e.POST("api/v1/group/make_admin", handler.MakeOrRemoveAdmin)
	e.GET("api/v1/group/get_group_details", handler.GetGroupDetails)
	e.POST("api/v1/group/accept_and_decline_group_invitation", handler.AcceptAndDeclineGroupInvitation)
	e.GET("api/v1/group/get_all_group_details_of_user", handler.GetAllGroupDetailsOfUser)
}
