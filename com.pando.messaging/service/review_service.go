package service

import (
	"context"
	"net/url"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	repo "pandoMessagingWalletService/com.pando.messaging/repository"
)

type reviewUsecase struct {
	repository repo.ReviewRepository
}

func NewReviewUsecase(repo repo.ReviewRepository) ReviewUsecase {
	return &reviewUsecase{
		repository: repo,
	}
}

/***********************************SaveReviewAndFeedback*************************************/
func (r *reviewUsecase) SaveReviewAndFeedback(ctx context.Context, review models.Reviews) (*models.Response, error) {
	logger.Logger.Info("Request received from service part of save review and feedback api.")
	return r.repository.SaveReviewAndFeedback(ctx, review)
}

/************************************GetReviewAndFeedback***************************************/
func (r *reviewUsecase) GetReviewAndFeedback(ctx context.Context, concertId int64, query url.Values) (*models.Response, error) {
	logger.Logger.Info("Request received from service part of get review and feedback api.")
	return r.repository.GetReviewAndFeedback(ctx, concertId, query)
}

/************************************GetOverallReview*******************************************/
func (r *reviewUsecase) GetOverallReview(ctx context.Context, concertId int64) (*models.Response, error) {
	logger.Logger.Info("Request received from service part of get overall review and feedback api.")
	return r.repository.GetOverallReview(ctx, concertId)
}
