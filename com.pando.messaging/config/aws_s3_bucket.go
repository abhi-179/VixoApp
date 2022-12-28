package config

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"

	"pandoMessagingWalletService/com.pando.messaging/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadFileToS3(s *session.Session, file []byte, TempFileName string, fileHeader *multipart.FileHeader) (string, error) {
	path, _ := os.Getwd()
	conf, errs := GetConfig(path + "/com.pando.messaging/env/")
	if errs != nil {
		logger.Logger.Info("config data not found for email.")
	}
	// get the file size and read
	size := fileHeader.Size
	tempFileName := TempFileName
	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(conf.S3Bucket),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:                 bytes.NewReader(file),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(file)),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return "", err
	}
	return tempFileName, err
}
