package router

import (
	controller "pandoMessagingWalletService/com.pando.messaging/controller"
	review "pandoMessagingWalletService/com.pando.messaging/service"

	"github.com/labstack/echo/v4"
)

func NewReviewController(e *echo.Echo, reviewUsecase review.ReviewUsecase) {
	handler := &controller.ReviewController{
		Usecase: reviewUsecase,
	}
	e.POST("api/v1/review/save_review_and_feedback", handler.SaveReviewAndFeedback)
	e.GET("api/v1/review/get_review_and_feedback", handler.GetReviewAndFeedback)
	e.GET("api/v1/review/get_overall_rating", handler.GetOverallReview)
}
