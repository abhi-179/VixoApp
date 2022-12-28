package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	booking "pandoMessagingWalletService/com.pando.messaging/service"
	"strconv"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

type BookingController struct {
	Usecase booking.BookingUsecase
}

/**********************************************BookTicket************************************/
// @Tags booking service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Book tickets
// @Description book tickets
// @Accept  json
// @Produce  json
// @Param booking body models.BookingDetailDto true "Book ticket"
// @Success 200 {object} interface{}
// @Header 200 {string} Token "qwerty"
// @Failure 404  "Not Found"
// @Router /api/v1/booking/book_ticket [post]
func (r *BookingController) BookTicket(c echo.Context) error {
	var bookTicket models.BookingDetailDto
	c.Bind(&bookTicket)
	id1, _ := strconv.Atoi(c.Request().Header.Get("id"))
	logger.Logger.Info("id ", id1)
	id := int64(id1)
	if id != bookTicket.UserId {
		logger.Logger.Error("Wrong auth token.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter correct auth token."})
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
	err := v.Struct(bookTicket)
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
	logger.Logger.Info("Request received from controller part of book ticket is ", bookTicket)
	authResponse, _ := r.Usecase.BookTicket(ctx, bookTicket)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)

}

/**************************************RefundToken*******************************************/
func (r *BookingController) RefundToken(c echo.Context) error {
	concertId1, _ := strconv.Atoi(c.QueryParam("concert_id"))
	concertId := int64(concertId1)
	if concertId == 0 {
		logger.Logger.Error("Concert id is missing.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter concert id."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of refund token is ", logrus.Fields{"concert_id": concertId})
	authResponse, _ := r.Usecase.RefundToken(ctx, concertId)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*************************************ViewTickets*********************************************/
// @Tags booking service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary View tickets
// @Description view tickets
// @Accept  json
// @Produce  json
// @Param user_id query int64 true "user_id"
// @Param concert_id query int64 true "concert_id"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/booking/view_tickets [get]
func (r *BookingController) ViewTickets(c echo.Context) error {
	concert_id1, _ := strconv.Atoi(c.QueryParam("concert_id"))
	concertId := int64(concert_id1)
	user_id1, _ := strconv.Atoi(c.QueryParam("user_id"))
	userId := int64(user_id1)
	id1, _ := strconv.Atoi(c.Request().Header.Get("id"))
	logger.Logger.Info("id ", id1)
	id := int64(id1)
	if id != userId {
		logger.Logger.Error("Wrong auth token.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter correct auth token."})
	}
	if concertId == 0 || userId == 0 {
		logger.Logger.Error("concert id or user id is missing.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter conert_id or user_id."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of view ticket is ", logrus.Fields{"concert_id": concertId, "user_id": userId})
	authResponse, _ := r.Usecase.ViewTickets(ctx, concertId, userId)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/****************************************VerifyTicketCode*************************************/
// @Tags booking service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Verify ticket
// @Description verify ticket
// @Accept  json
// @Produce  json
// @Param verify ticket body models.VerifyTicketsDto true "Verify ticket"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/booking/verify_ticket [post]
func (r *BookingController) VerifyTicketCode(c echo.Context) error {
	var Req map[string]interface{}
	c.Bind(&Req)
	id1, _ := strconv.Atoi(c.Request().Header.Get("id"))
	id := int64(id1)
	user_id1 := fmt.Sprintf("%v", Req["user_id"])
	userId, _ := strconv.ParseInt(user_id1, 10, 64)
	if id != userId {
		logger.Logger.Error("Wrong auth token.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter correct auth token."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if Req["user_id"] == 0 || Req["ticket_code"] == "" || Req["concert_id"] == 0 {
		logger.Logger.Error("User id, ticket code or concert id may be missing.")
		return c.JSON(http.StatusBadRequest, &models.Response{Status: false, Msg: "Please enter user_id and ticket_code both."})
	}
	logger.Logger.Info("Request received from controller part of verify ticket code is ", Req)
	authResponse, _ := r.Usecase.VerifyTicketCode(ctx, Req)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*****************************************SendTicket***************************************/
// @Tags booking service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Send ticket
// @Description send ticket
// @Accept  json
// @Produce  json
// @Param send ticket body models.SendTicketDto true "send ticket"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/booking/send_ticket [post]
func (r *BookingController) SentTicket(c echo.Context) error {
	var req map[string]interface{}
	c.Bind(&req)
	id1, _ := strconv.Atoi(c.Request().Header.Get("id"))
	id := int64(id1)
	sender_user_id1 := fmt.Sprintf("%v", req["sender_user_id"])
	sender_user_id, _ := strconv.ParseInt(sender_user_id1, 10, 64)
	if id != sender_user_id {
		logger.Logger.Error("Wrong auth token.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter correct auth token."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if req["sender_user_id"] == 0 || req["ticket_code"] == "" || req["receiver_user_id"] == 0 {
		logger.Logger.Error("sender_user id, receiver_user_id and ticket code are missing.")
		return c.JSON(http.StatusBadRequest, &models.Response{Status: false, Msg: "Please enter sender_user_id, receiver_user_id and ticket_code."})
	}
	logger.Logger.Info("Request received from controller part of send ticket is ", req)
	authResponse, _ := r.Usecase.SendTicket(ctx, req)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/****************************************Cancel Ticket***************************************/
// @Tags booking service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Cancel ticket
// @Description cancel ticket
// @Accept  json
// @Produce  json
// @Param cancel ticket body models.CancelTicketDto true "cancel ticket"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/booking/cancel_ticket [post]
func (r *BookingController) CancelTicket(c echo.Context) error {
	var req map[string]interface{}
	c.Bind(&req)
	id1, _ := strconv.Atoi(c.Request().Header.Get("id"))
	id := int64(id1)
	userId1 := fmt.Sprintf("%v", req["user_id"])
	userId, _ := strconv.ParseInt(userId1, 10, 64)
	if id != userId {
		logger.Logger.Error("Wrong auth token.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter correct auth token."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if req["user_id"] == 0 || req["ticket_code"] == "" || req["concert_id"] == 0 {
		logger.Logger.Error("user id, concert_id and ticket code are missing.")
		return c.JSON(http.StatusBadGateway, &models.Response{Status: false, Msg: "Please enter user_id, concert_id and ticket_code."})
	}
	logger.Logger.Info("Request received from controller part of cancel ticket is ", logrus.Fields{"user_id": req["user_id"], "concert_id": req["concert_id"], "ticket_code": req["ticket_code"]})
	authResponse, _ := r.Usecase.CancelTicket(ctx, req)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}
