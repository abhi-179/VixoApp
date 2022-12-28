package service

import (
	"context"
	"mime/multipart"
	models "pandoMessagingWalletService/com.pando.messaging/model"
)

type GroupUsecase interface {
	Create_Group(ctx context.Context, flow models.Group) (*models.Response, error)
	AddUserToGroup(ctx context.Context, flow models.Groups) (*models.Response, error)
	RemoveUsersFromGroup(ctx context.Context, group_id int64, user_id int64, admin_id int64) (*models.Response, error)
	LeaveGroup(ctx context.Context, group_id int64, user_id int64) (*models.Response, error)
	EditGroupInfo(ctx context.Context, group_id int64, new_group_name string, profile_pic_url string) (*models.Response, error)
	SearchUsersInGroup(ctx context.Context, group_id string, username string) (*models.Response, error)
	DeleteGroup(ctx context.Context, group_id int64) (*models.Response, error)
	UploadGroupProfilePhoto(ctx context.Context, file multipart.File, handler *multipart.FileHeader, filename string) (*models.Response, error)
	MakeOrRemoveAdmin(ctx context.Context, user_id int64, group_id int64, new_admin_id int64, method_type string) (*models.Response, error)
	GetGroupDetails(ctx context.Context, group_id int64) (*models.Response, error)
	AcceptAndDeclineGroupInvitation(ctx context.Context, group models.AcceptGroupInvitationDto) (*models.Response, error)
	GetAllGroupDetailsOfUser(ctx context.Context, user_id int64) (*models.Response, error)
}
