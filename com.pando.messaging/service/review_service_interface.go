package service

import (
	"context"
	"net/url"
	models "pandoMessagingWalletService/com.pando.messaging/model"
)

type ReviewUsecase interface {
	SaveReviewAndFeedback(ctx context.Context, review models.Reviews) (*models.Response, error)
	GetReviewAndFeedback(ctx context.Context, concertId int64, query url.Values) (*models.Response, error)
	GetOverallReview(ctx context.Context, concertId int64) (*models.Response, error)
}
