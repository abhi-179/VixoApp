package service

import (
	"context"
	models "pandoMessagingWalletService/com.pando.messaging/model"
)

type ChatUsecase interface {
	PostStatus(ctx context.Context, status models.Status) (*models.Response, error)
	FetchStatus(ctx context.Context, user_id int64) (*models.Response, error)
	DeleteStatus(ctx context.Context, user_id int64, status_id []int64) (*models.Response, error)
	SearchStatusByUsername(ctx context.Context, username string) (*models.Response, error)
	ReportChat(ctx context.Context, flow models.Reports) (*models.Response, error)
	SaveBlockUserDetails(ctx context.Context, flow models.BlockedContacts) (*models.Response, error)
	FetchBlockedUserDetails(ctx context.Context, user_id int64) (*models.Response, error)
	FetchBlockedContactDetails(ctx context.Context, user_id int64) (*models.Response, error)
	Unblock_user(ctx context.Context, flow models.BlockedContacts) (*models.Response, error)
	FetchWallpapersDetails(ctx context.Context) (*models.Response, error)
	SaveGroupChatSetting(ctx context.Context, flow models.ChatSettings) (*models.Response, error)
	FetchGroupChatSetting(ctx context.Context, user_id int64) (*models.Response, error)
}
