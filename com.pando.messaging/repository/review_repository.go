package repository

import (
	"context"
	"net/url"
	"pandoMessagingWalletService/com.pando.messaging/config"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"

	"gorm.io/gorm"
)

type reviewRepository struct {
	DBConn *gorm.DB
}

func NewReviewRepository(conn *gorm.DB, conf *config.Config) ReviewRepository {
	return &reviewRepository{
		DBConn: conn,
	}
}

/********************************************SaveReviewAndFeedback******************************/
func (r *reviewRepository) SaveReviewAndFeedback(ctx context.Context, review models.Reviews) (*models.Response, error) {
	check := r.DBConn.Where("concert_id = ? and user_id = ?", review.ConcertId, review.UserId).Find(&review)
	if check.RowsAffected != 0 {
		logger.Logger.WithError(check.Error).WithField("error", check.Error).Error("You have already reviewed.")
		return &models.Response{Status: true, Msg: "You have already reviewed."}, nil
	}
	db := r.DBConn.Create(&review)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Review is not saved.")
		return &models.Response{Status: false, Msg: "Review is not saved."}, nil
	}
	return &models.Response{Status: true, Msg: "Review is Saved successfully."}, nil
}

/*******************************************GetReviewAndFeedback********************************/
func (r *reviewRepository) GetReviewAndFeedback(ctx context.Context, concertId int64, query url.Values) (*models.Response, error) {
	reviews := []models.ReviewDto{}
	pagination := config.GeneratePaginationFromRequest(query)
	offset := (pagination.Page) * pagination.Limit
	queryBuider := r.DBConn.Limit(pagination.Limit).Offset(offset)
	db := queryBuider.Table("reviews as r").Select("user_id,concert_id,rating,feedback,username,profile_pic_url").Joins("left join users as u on r.user_id=u.id").Where("concert_id = ?", concertId).Find(&reviews)
	if db.RowsAffected == 0 {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Reviews not found.")
		return &models.Response{Status: true, Msg: "Reveiws not found."}, nil
	}
	return &models.Response{Status: true, Msg: "Reviews found.", Reviews: &reviews}, nil
}

/*****************************************GetOverallReview***********************************/
func (r *reviewRepository) GetOverallReview(ctx context.Context, concertId int64) (*models.Response, error) {
	var rating float64
	db := r.DBConn.Table("reviews").Select("AVG(rating)").Where("concert_id = ?", concertId).Find(&rating)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Review not found.")
		return &models.Response{Status: true, Msg: "Review not found."}, nil
	}
	return &models.Response{Status: true, Msg: "Review found", Rating: rating}, nil
}
