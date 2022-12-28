package config

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"pandoMessagingWalletService/com.pando.messaging/logger"
	models "pandoMessagingWalletService/com.pando.messaging/model"
	"strconv"
	"time"
)

func UploadFileToIPFS(files []byte, handler *multipart.FileHeader) (*models.IPFSResult, error) {
	path, _ := os.Getwd()
	conf, errs := GetConfig(path + "/com.pando.messaging/env/")
	if errs != nil {
		logger.Logger.Info("config data not found for email.")
	}
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	timestamp := time.Now().UnixNano()
	real_time := strconv.Itoa(int(timestamp))
	part, _ := writer.CreateFormFile("file", real_time+"_"+handler.Filename)
	part.Write(files)
	writer.Close()
	req, _ := http.NewRequest("POST", conf.IPFSURL, payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logger.Logger.WithError(err).WithField("error", err).Error("File is not uploaded on IPFS.")
		return nil, nil
	}
	var ipfs models.IPFSResult
	data, _ := ioutil.ReadAll(res.Body)
	_ = json.Unmarshal(data, &ipfs)
	return &ipfs, nil
}
