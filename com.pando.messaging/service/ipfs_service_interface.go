package service

import (
	"context"
	"mime/multipart"
	models "pandoMessagingWalletService/com.pando.messaging/model"
)

type IPFSUsecase interface {
	UploadFile(ctx context.Context, files multipart.File, handler *multipart.FileHeader) (*models.IPFSResult, error)
	SaveBackupFilehash(ctx context.Context, req models.Backups) (*models.Response, error)
	GetBackupFilehash(ctx context.Context, userId int64) (*models.Response, error)
	UploadFileToIPFS(ctx context.Context, files multipart.File, handler *multipart.FileHeader, data string) (*models.IPFSResult, error)
}
