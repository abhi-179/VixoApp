package service

import (
	"context"
	"mime/multipart"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	repo "pandoMessagingWalletService/com.pando.messaging/repository"
)

//	"encoding/json"

type ipfsUsecase struct {
	repository repo.IPFSRepository
}

func NewIPFSUsecase(repo repo.IPFSRepository) IPFSUsecase {
	return &ipfsUsecase{
		repository: repo,
	}
}

/***********************************Upload file to ipfs and aws********************************/
func (r *ipfsUsecase) UploadFile(ctx context.Context, files multipart.File, handler *multipart.FileHeader) (*models.IPFSResult, error) {
	logger.Logger.Info("Request received from upload file to ipfs and aws service part.")
	return r.repository.UploadFile(ctx, files, handler)
}

/**************************************Save Backup Filehash*************************************/
func (r *ipfsUsecase) SaveBackupFilehash(ctx context.Context, req models.Backups) (*models.Response, error) {
	logger.Logger.Info("Request received from save backup filehash service part.")
	return r.repository.SaveBackupFilehash(ctx, req)
}

/***************************************GetBackupFilehash***************************************/
func (r *ipfsUsecase) GetBackupFilehash(ctx context.Context, userId int64) (*models.Response, error) {
	logger.Logger.Info("Request received from get backup filehash service part.")
	return r.repository.GetBackupFilehash(ctx, userId)
}

/***********************************UploadFileToIPFS************************************/
func (r *ipfsUsecase) UploadFileToIPFS(ctx context.Context, files multipart.File, handler *multipart.FileHeader, data string) (*models.IPFSResult, error) {
	logger.Logger.Info("Request received from upload file to ipfs service part.")
	return r.repository.UploadFileToIPFS(ctx, files, handler, data)
}
