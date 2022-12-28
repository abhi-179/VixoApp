package controller

import (
	"context"
	"net/http"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	review "pandoMessagingWalletService/com.pando.messaging/service"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ReviewController struct {
	Usecase review.ReviewUsecase
}

// @Tags review service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Save review
// @Description save review
// @Accept  json
// @Produce  json
// @Param save review body models.SaveReviewDto true "save review"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/review/save_review_and_feedback [post]
func (r *ReviewController) SaveReviewAndFeedback(c echo.Context) error {
	var req models.Reviews
	c.Bind(&req)
	if req.ConcertId == 0 || req.UserId == 0 {
		logger.Logger.Error("Concert or user id is null.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter concert id and user id."})
	}
	id1, _ := strconv.Atoi(c.Request().Header.Get("id"))
	id := int64(id1)
	if id != req.UserId {
		logger.Logger.Error("Wrong auth token.")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter correct auth token."})
	}
	if req.Rating == 0 || req.Rating > 5 {
		logger.Logger.Error("Please enter review between 1 to 5")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter review between 1 to 5."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of save review and feedback is ", logrus.Fields{"concert_id": req.ConcertId, "user_id": req.UserId, "rating": req.Rating, "feedback": req.Feedback})
	authResponse, _ := r.Usecase.SaveReviewAndFeedback(ctx, req)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*************************************GetReviewAndFeedback************************************/
// @Tags review service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Get reviews and feedback
// @Description get reviews and feedback
// @Accept  json
// @Produce  json
// @Param concert_id query int64 true "concert_id"
// @Param limit query int false "limit"
// @Param page query int false "page"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/review/get_review_and_feedback [get]
func (r *ReviewController) GetReviewAndFeedback(c echo.Context) error {
	concert_id1, _ := strconv.Atoi(c.QueryParam("concert_id"))
	concertId := int64(concert_id1)
	query := c.Request().URL.Query()
	if concertId == 0 {
		logger.Logger.Error("Please enter concert id")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter concert."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of get review and feedback is ", logrus.Fields{"concert_id": concertId, "query": query})
	authResponse, _ := r.Usecase.GetReviewAndFeedback(ctx, concertId, query)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}

/*****************************************GetOverallReview**********************************/
// @Tags review service
// @Security ApiKeyAuth
// @securityDefinitions.apikey ApiKeyAuth Authorization
// @in header
// @name Authorization
// @Summary Get overall reviews
// @Description get overall reviews
// @Accept  json
// @Produce  json
// @Param concert_id query int64 true "concert_id"
// @Success 200 {object} interface{}
// @Header 200 {string} Token string
// @Failure 404  "Not Found"
// @Router /api/v1/review/get_overall_rating [get]
func (r *ReviewController) GetOverallReview(c echo.Context) error {
	concert_id1, _ := strconv.Atoi(c.QueryParam("concert_id"))
	concertId := int64(concert_id1)
	if concertId == 0 {
		logger.Logger.Error("Please enter concert id")
		return c.JSON(400, &models.Response{Status: false, Msg: "Please enter concert."})
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	logger.Logger.Info("Request received from controller part of get overall review and feedback is ", logrus.Fields{"concert_id": concertId})
	authResponse, _ := r.Usecase.GetOverallReview(ctx, concertId)
	if authResponse == nil {
		return c.JSON(http.StatusUnauthorized, authResponse)
	}
	if !authResponse.Status {
		return c.JSON(http.StatusBadRequest, authResponse)
	}
	return c.JSON(http.StatusOK, authResponse)
}
