package router

import (
	controller "pandoMessagingWalletService/com.pando.messaging/controller"
	booking "pandoMessagingWalletService/com.pando.messaging/service"

	"github.com/labstack/echo/v4"
)

func NewBookingController(e *echo.Echo, bookingUsecase booking.BookingUsecase) {
	handler := &controller.BookingController{
		Usecase: bookingUsecase,
	}
	e.POST("api/v1/booking/book_ticket", handler.BookTicket)
	e.GET("api/v1/booking/refund", handler.RefundToken)
	e.GET("api/v1/booking/view_tickets", handler.ViewTickets)
	e.POST("api/v1/booking/verify_ticket", handler.VerifyTicketCode)
	e.POST("api/v1/booking/send_ticket", handler.SentTicket)
	e.POST("api/v1/booking/cancel_ticket", handler.CancelTicket)

}
