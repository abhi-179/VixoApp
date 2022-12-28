package service

import (
	"context"
	"fmt"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	repo "pandoMessagingWalletService/com.pando.messaging/repository"
	"strconv"
)

type bookingUsecase struct {
	repository repo.BookingRepository
}

func NewBookingUsecase(repo repo.BookingRepository) BookingUsecase {
	return &bookingUsecase{
		repository: repo,
	}
}

/*****************************************BookTicket********************************************/
func (r *bookingUsecase) BookTicket(ctx context.Context, bookTicket models.BookingDetailDto) (*models.Response, error) {
	var message string
	if bookTicket.TotalTickets > 1 {
		message = "tickets"
	} else {
		message = "ticket"
	}
	logger.Logger.Info("Request received from service part of book tickets api.")
	return r.repository.BookTicket(ctx, bookTicket, message)
}

/***************************************Refund Token********************************************/
func (r *bookingUsecase) RefundToken(ctx context.Context, concertId int64) (*models.Response, error) {
	logger.Logger.Info("Request received from service part of refund token api.")
	return r.repository.RefundToken(ctx, concertId)
}

/*****************************************ViewTickets*******************************************/
func (r *bookingUsecase) ViewTickets(ctx context.Context, concertId, userId int64) (*models.Response, error) {
	logger.Logger.Info("Request received from service part of view token api.")
	return r.repository.ViewTickets(ctx, concertId, userId)
}

/*****************************************VerifyTicketCode************************************/
func (r *bookingUsecase) VerifyTicketCode(ctx context.Context, flow map[string]interface{}) (*models.Response, error) {
	userId1 := fmt.Sprintf("%v", flow["user_id"])
	userId, _ := strconv.ParseInt(userId1, 10, 64)
	concertId1 := fmt.Sprintf("%v", flow["concert_id"])
	concertId, _ := strconv.ParseInt(concertId1, 10, 64)
	ticketCode := fmt.Sprintf("%v", flow["ticket_code"])
	logger.Logger.Info("Request received from service part of verify ticket code api.")
	return r.repository.VerifyTicketCode(ctx, userId, concertId, ticketCode)
}

/****************************************SendTicket****************************************/
func (r *bookingUsecase) SendTicket(ctx context.Context, flow map[string]interface{}) (*models.Response, error) {
	senderUserId1 := fmt.Sprintf("%v", flow["sender_user_id"])
	senderUserId, _ := strconv.ParseInt(senderUserId1, 10, 64)
	receiverUserId1 := fmt.Sprintf("%v", flow["receiver_user_id"])
	receiverUserId, _ := strconv.ParseInt(receiverUserId1, 10, 64)
	ticketCode := fmt.Sprintf("%v", flow["ticket_code"])
	logger.Logger.Info("Request received from service part of send ticket api.")
	return r.repository.SendTicket(ctx, senderUserId, receiverUserId, ticketCode)
}

/****************************************CancelTicket*************************************/
func (r *bookingUsecase) CancelTicket(ctx context.Context, flow map[string]interface{}) (*models.Response, error) {
	userId1 := fmt.Sprintf("%v", flow["user_id"])
	userId, _ := strconv.ParseInt(userId1, 10, 64)
	concertId1 := fmt.Sprintf("%v", flow["concert_id"])
	concertId, _ := strconv.ParseInt(concertId1, 10, 64)
	ticketCode := fmt.Sprintf("%v", flow["ticket_code"])
	logger.Logger.Info("Request received from service part of cancel ticket api.")
	return r.repository.CancelTicket(ctx, userId, concertId, ticketCode)
}
