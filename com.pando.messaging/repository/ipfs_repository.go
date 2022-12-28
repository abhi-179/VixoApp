package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"sync"

	"pandoMessagingWalletService/com.pando.messaging/config"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"

	"gorm.io/gorm"
)

type ipfsRepository struct {
	DBConn *gorm.DB
}

var S3_REGION, S3_ACCESS_ID, S3_SECRET_KEY, AWS_URL, IPFSURL string

func NewIPFSRepository(conn *gorm.DB, conf *config.Config) IPFSRepository {
	S3_ACCESS_ID = conf.S3AccessId
	S3_REGION = conf.S3Region
	AWS_URL = conf.AwsURL
	S3_SECRET_KEY = conf.S3SecretKey
	IPFSURL = conf.IPFSURL
	return &ipfsRepository{
		DBConn: conn,
	}
}

/***********************************************Upload file to ipfs ***************************/
func (r *ipfsRepository) UploadFile(ctx context.Context, files multipart.File, handler *multipart.FileHeader) (*models.IPFSResult, error) {
	var wg sync.WaitGroup
	wg.Add(1)
	var filename string
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(S3_REGION),
		Credentials: credentials.NewStaticCredentials(
			S3_ACCESS_ID,  // id
			S3_SECRET_KEY, // secret
			""),           // token can be left blank for now
	})
	if err != nil {
		logger.Logger.WithError(err).WithField("err", err).Errorf("Could not upload file")
		return nil, err
	}
	fileBytes, err := ioutil.ReadAll(files)
	if err != nil {
		return nil, err
	}
	file, err := config.UploadFileToIPFS(fileBytes, handler)
	if err != nil {
		logger.Logger.WithError(err).WithField("err", err).Errorf("Could not upload file.")
		return nil, err
	}
	go func() {
		fileName, err := config.UploadFileToS3(s, fileBytes, "chat/"+handler.Filename, handler)
		if err != nil {
			logger.Logger.WithError(err).WithField("err", err).Errorf("Could not upload file.")
			return
		}
		filename = fileName
		wg.Done()
	}()

	wg.Wait()
	return &models.IPFSResult{Name: file.Name, Hash: file.Hash, Size: file.Size, AwsURL: AWS_URL + filename}, nil
}

/**********************************SaveBackupFilehash********************************/
func (r *ipfsRepository) SaveBackupFilehash(ctx context.Context, req models.Backups) (*models.Response, error) {
	ipfs := models.Backups{}
	db := r.DBConn.Where("user_id = ?", req.UserId).First(&ipfs)
	if db.RowsAffected != 0 {
		update := r.DBConn.Where("user_id = ?", req.UserId).Find(&ipfs).Updates(map[string]interface{}{"filehash": req.Filehash, "file_size": req.FileSize, "file_name": req.FileName, "backup_nature": req.BackupNature})
		if update.RowsAffected == 0 {
			logger.Logger.WithError(update.Error).WithField("error", update.Error).Error("Filehash is not updated.")
			return &models.Response{Status: false, ResponseCode: 400, Msg: "Filehash is not updated."}, nil
		}
		return &models.Response{Status: true, ResponseCode: 200, Msg: "Filehash updated."}, nil
	} else {
		data := r.DBConn.Create(&req)
		if data.Error != nil {
			logger.Logger.WithError(data.Error).WithField("error", data.Error).Error("Filehash is not saved.")
			return &models.Response{Status: false, ResponseCode: 400, Msg: "Filehash is not saved."}, nil
		}
		return &models.Response{Status: true, ResponseCode: 200, Msg: "Filehash is saved"}, nil
	}
}

/*******************************************GetBackupFilehash******************************/
func (r *ipfsRepository) GetBackupFilehash(ctx context.Context, userId int64) (*models.Response, error) {
	ipfs := models.Backups{}
	db := r.DBConn.Where("user_id=?", userId).First(&ipfs)
	if db.Error != nil {
		logger.Logger.WithError(db.Error).WithField("error", db.Error).Error("Filehash is not found.")
		return &models.Response{Status: true, ResponseCode: 200, Msg: "Filehash is not found."}, nil
	}
	return &models.Response{Status: true, ResponseCode: 200, Msg: "Filehash is found", Backups: &ipfs}, nil
}

/******************************************UploadFileToIPFS******************************/
func (r *ipfsRepository) UploadFileToIPFS(ctx context.Context, files multipart.File, handler *multipart.FileHeader, data string) (*models.IPFSResult, error) {
	if data != "" {
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("file", data)
		err := writer.Close()
		if err != nil {
			return nil, err
		}
		client := &http.Client{}
		req, err := http.NewRequest("POST", IPFSURL, payload)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()
		var ipfs models.IPFSResult
		body, _ := ioutil.ReadAll(res.Body)
		_ = json.Unmarshal(body, &ipfs)
		return &models.IPFSResult{Name: ipfs.Name, Hash: ipfs.Hash, Size: ipfs.Size}, nil

	} else {
		fileBytes, err := ioutil.ReadAll(files)
		if err != nil {
			return nil, err
		}
		file, err := config.UploadFileToIPFS(fileBytes, handler)
		if err != nil {
			logger.Logger.WithError(err).WithField("err", err).Errorf("Could not upload file.")
			return nil, err
		}
		defer files.Close()
		return &models.IPFSResult{Name: file.Name, Hash: file.Hash, Size: file.Size}, nil
	}
}
