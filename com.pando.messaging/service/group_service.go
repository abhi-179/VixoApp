package service

import (
	"context"
	"mime/multipart"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	repo "pandoMessagingWalletService/com.pando.messaging/repository"
)

//	"encoding/json"

type groupUsecase struct {
	repository repo.GroupRepository
}

func NewgroupUsecase(repo repo.GroupRepository) GroupUsecase {
	return &groupUsecase{
		repository: repo,
	}
}

/********************************************Create Group******************************************/
func (r *groupUsecase) Create_Group(ctx context.Context, flow models.Group) (*models.Response, error) {
	logger.Logger.Info("Request received from create group service part.")
	return r.repository.Create_Group(ctx, flow)
}

/*******************************************Add Users to group*************************************/
func (r *groupUsecase) AddUserToGroup(ctx context.Context, flow models.Groups) (*models.Response, error) {
	logger.Logger.Info("Request received from add users to group service part.")
	return r.repository.AddUserToGroup(ctx, flow)
}

/*****************************************Remove Users from group***********************************/
func (r *groupUsecase) RemoveUsersFromGroup(ctx context.Context, group_id int64, user_id int64, admin_id int64) (*models.Response, error) {
	logger.Logger.Info("Request received from remove users from group service part.")
	return r.repository.RemoveUsersFromGroup(ctx, group_id, user_id, admin_id)
}

/******************************************LeaveGroup************************************************/
func (r *groupUsecase) LeaveGroup(ctx context.Context, group_id int64, user_id int64) (*models.Response, error) {
	logger.Logger.Info("Request received from leave group service part.")
	return r.repository.LeaveGroup(ctx, group_id, user_id)
}

/*****************************************EditGroupInfo********************************************/
func (r *groupUsecase) EditGroupInfo(ctx context.Context, group_id int64, new_group_name string, profile_pic_url string) (*models.Response, error) {
	logger.Logger.Info("Request received from edit group name service part.")
	return r.repository.EditGroupInfo(ctx, group_id, new_group_name, profile_pic_url)
}

/*****************************************SearchUsersInGroup****************************************/
func (r *groupUsecase) SearchUsersInGroup(ctx context.Context, group_id string, username string) (*models.Response, error) {
	logger.Logger.Info("Request received from search users in group service part.")
	return r.repository.SearchUsersInGroup(ctx, group_id, username)
}

/******************************************DeleteGroup*************************************/
func (r *groupUsecase) DeleteGroup(ctx context.Context, group_id int64) (*models.Response, error) {
	logger.Logger.Info("Request received from delete group service part.")
	return r.repository.DeleteGroup(ctx, group_id)
}

/******************************************UploadProfilePhoto***************************************/
func (r *groupUsecase) UploadGroupProfilePhoto(ctx context.Context, file multipart.File, handler *multipart.FileHeader, filename string) (*models.Response, error) {
	logger.Logger.Info("Request received from Upload Profile photo service part.")
	return r.repository.UploadGroupProfilePhoto(ctx, file, handler, filename)
}

/******************************************Make or Remove Admin*******************************/
func (r *groupUsecase) MakeOrRemoveAdmin(ctx context.Context, user_id int64, group_id int64, new_admin_id int64, method_type string) (*models.Response, error) {
	logger.Logger.Info("Request received from make or remove admin service part.")
	return r.repository.MakeOrRemoveAdmin(ctx, user_id, group_id, new_admin_id, method_type)
}

/******************************************Get group details*******************************/
func (r *groupUsecase) GetGroupDetails(ctx context.Context, group_id int64) (*models.Response, error) {
	logger.Logger.Info("Request received from get group details service part.")
	return r.repository.GetGroupDetails(ctx, group_id)
}

/************************************AcceptAndDeclineGroupInvitation****************************/
func (r *groupUsecase) AcceptAndDeclineGroupInvitation(ctx context.Context, group models.AcceptGroupInvitationDto) (*models.Response, error) {
	logger.Logger.Info("Request received from accept and decline group chat invitation service part.")
	return r.repository.AcceptAndDeclineGroupInvitation(ctx, group)
}

/****************************************Get all group details of user******************************************/
func (r *groupUsecase) GetAllGroupDetailsOfUser(ctx context.Context, user_id int64) (*models.Response, error) {
	logger.Logger.Info("Request received from get all group details of user service part.")
	return r.repository.GetAllGroupDetailsOfUser(ctx, user_id)
}
