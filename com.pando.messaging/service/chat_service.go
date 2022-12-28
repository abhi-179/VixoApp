package service

import (
	"context"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	repo "pandoMessagingWalletService/com.pando.messaging/repository"
)

type chatUsecase struct {
	repository repo.ChatRepository
}

func NewchatUsecase(repo repo.ChatRepository) ChatUsecase {
	return &chatUsecase{
		repository: repo,
	}
}

/******************************************Post Status******************************************/
func (r *chatUsecase) PostStatus(ctx context.Context, status models.Status) (*models.Response, error) {
	logger.Logger.Info("Request received from post status service part.")
	return r.repository.PostStatus(ctx, status)
}

/******************************************Fetch Status*****************************************/
func (r *chatUsecase) FetchStatus(ctx context.Context, user_id int64) (*models.Response, error) {
	logger.Logger.Info("Request received from fetch status service part.")
	return r.repository.FetchStatus(ctx, user_id)
}

/******************************************Delete Status***************************************/
func (r *chatUsecase) DeleteStatus(ctx context.Context, user_id int64, status_id []int64) (*models.Response, error) {
	logger.Logger.Info("Request received from delete status service part.")
	return r.repository.DeleteStatus(ctx, user_id, status_id)
}

/*******************************************SearchStatusByUsername*******************************/
func (r *chatUsecase) SearchStatusByUsername(ctx context.Context, username string) (*models.Response, error) {
	logger.Logger.Info("Request received from SearchStatusByUsername servive part.")
	return r.repository.SearchStatusByUsername(ctx, username)
}

/*********************************************Remove File From Ipfs********************************/
func (r *chatUsecase) ReportChat(ctx context.Context, flow models.Reports) (*models.Response, error) {
	logger.Logger.Info("Request received from report chat service part.")
	return r.repository.ReportChat(ctx, flow)
}

/********************************************Save Block User Details**********************************/
func (r *chatUsecase) SaveBlockUserDetails(ctx context.Context, flow models.BlockedContacts) (*models.Response, error) {
	logger.Logger.Info("Request received from save blocked user details service part.")
	return r.repository.SaveBlockUserDetails(ctx, flow)
}

/**********************************************Fetch Blocked users Details*********************************/
func (r *chatUsecase) FetchBlockedUserDetails(ctx context.Context, user_id int64) (*models.Response, error) {
	logger.Logger.Info("Request received from fetch blocked user details service part.")
	return r.repository.FetchBlockedUserDetails(ctx, user_id)
}

/**********************************************Fetch Blocked users Details*********************************/
func (r *chatUsecase) FetchBlockedContactDetails(ctx context.Context, user_id int64) (*models.Response, error) {
	logger.Logger.Info("Request received from fetch blocked user contacts service part.")
	return r.repository.FetchBlockedContactDetails(ctx, user_id)
}

/********************************************Unblock Blocked Users**********************************/
func (r *chatUsecase) Unblock_user(ctx context.Context, flow models.BlockedContacts) (*models.Response, error) {
	logger.Logger.Info("Request received from unblock blocked user service part.")
	return r.repository.Unblock_user(ctx, flow)
}

/**********************************************Fetch Wallpaper Details*********************************/
func (r *chatUsecase) FetchWallpapersDetails(ctx context.Context) (*models.Response, error) {
	logger.Logger.Info("Request received from fetch wallpapers details service part.")
	return r.repository.FetchWallpapersDetails(ctx)
}

/************************************************Save Group chat setting*************************/
func (r *chatUsecase) SaveGroupChatSetting(ctx context.Context, flow models.ChatSettings) (*models.Response, error) {
	logger.Logger.Info("Request received from save group chat setting service part.")
	return r.repository.SaveGroupChatSetting(ctx, flow)
}

/***********************************************Fetch Group Chat setting**************************/
func (r *chatUsecase) FetchGroupChatSetting(ctx context.Context, user_id int64) (*models.Response, error) {
	logger.Logger.Info("Request received from fetch group chat setting service part.")
	return r.repository.FetchGroupChatSetting(ctx, user_id)
}
