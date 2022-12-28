package repository

import (
	"context"
	models "pandoMessagingWalletService/com.pando.messaging/model"
)

type BookingRepository interface {
	BookTicket(ctx context.Context, bookTicket models.BookingDetailDto, message string) (*models.Response, error)
	RefundToken(ctx context.Context, concertId int64) (*models.Response, error)
	ViewTickets(ctx context.Context, concertId, userId int64) (*models.Response, error)
	VerifyTicketCode(ctx context.Context, userId, concertId int64, ticketCode string) (*models.Response, error)
	SendTicket(ctx context.Context, senderUserId, receiverUserId int64, ticketCode string) (*models.Response, error)
	CancelTicket(ctx context.Context, userId, concertId int64, ticketCode string) (*models.Response, error)
}
