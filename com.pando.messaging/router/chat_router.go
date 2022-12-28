package router

import (
	controller "pandoMessagingWalletService/com.pando.messaging/controller"
	chat "pandoMessagingWalletService/com.pando.messaging/service"

	"github.com/labstack/echo/v4"
)

func NewChatController(e *echo.Echo, chatusecase chat.ChatUsecase) {
	handler := &controller.ChatController{
		Usecase: chatusecase,
	}
	// e.POST("api/v1/chat/save_wallpaper_details", handler.SaveWallpapersDetails)
	e.GET("api/v1/chat/fetch_wallpaper_details", handler.FetchWallpapersDetails)
	e.POST("api/v1/chat/post_status", handler.PostStatus)
	e.GET("api/v1/chat/fetch_status", handler.FetchStatus)
	e.GET("api/v1/chat/delete_status", handler.DeleteStatus)
	e.GET("api/v1/chat/search_status_by_username", handler.SearchStatusByUsername)
	e.POST("api/v1/chat/report_chat", handler.ReportChat)
	e.POST("api/v1/chat/save_blocked_contacts", handler.SaveBlockUserDetails)
	e.GET("api/v1/chat/fetch_blocked_contacts", handler.FetchBlockedContactDetails)
	e.GET("api/v1/chat/fetch_blocked_users", handler.FetchBlockedUserDetails)
	e.POST("api/v1/chat/unblock_user", handler.Unblock_user)
	e.POST("api/v1/chat/save_group_chat_setting", handler.SaveGroupChatSetting)
	e.GET("api/v1/chat/fetch_group_chat_setting", handler.FetchGroupChatSetting)

}
