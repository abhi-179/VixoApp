package router

import (
	controller "pandoMessagingWalletService/com.pando.messaging/controller"
	call "pandoMessagingWalletService/com.pando.messaging/service"

	"github.com/labstack/echo/v4"
)

func NewCallController(e *echo.Echo, callusecase call.CallUsecase) {
	handler := &controller.CallController{
		Usecase: callusecase,
	}
	e.POST("api/v1/call/save_call_details", handler.SaveCallLogs)
	e.GET("api/v1/call/fetch_all_call_details", handler.FetchAllCallLogs)
	e.GET("api/v1/call/fetch_missed_call_details", handler.FetchMissedCallLogs)
	e.GET("api/v1/call/delete_call_details", handler.DeleteCallLogs)
}
