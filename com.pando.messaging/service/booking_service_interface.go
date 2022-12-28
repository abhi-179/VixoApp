package service

import (
	"context"
	models "pandoMessagingWalletService/com.pando.messaging/model"
)

type BookingUsecase interface {
	BookTicket(ctx context.Context, bookTicket models.BookingDetailDto) (*models.Response, error)
	RefundToken(ctx context.Context, concertId int64) (*models.Response, error)
	ViewTickets(ctx context.Context, concertId, userId int64) (*models.Response, error)
	VerifyTicketCode(ctx context.Context, flow map[string]interface{}) (*models.Response, error)
	SendTicket(ctx context.Context, flow map[string]interface{}) (*models.Response, error)
	CancelTicket(ctx context.Context, flow map[string]interface{}) (*models.Response, error)
}
