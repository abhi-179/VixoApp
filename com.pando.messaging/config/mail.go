package config

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"

	"pandoMessagingWalletService/com.pando.messaging/logger"

	"github.com/scorredoira/email"
)

func SendMail(filename, to, message, subject string) {
	path, _ := os.Getwd()
	conf, errs := GetConfig(path + "/com.pando.messaging/env/")
	if errs != nil {
		logger.Logger.Info("config data not found for email.")
	}
	smtpHost := conf.SmtpHost // change to your SMTP provider address
	smtpPort := 587           // change to your SMPT provider port number
	smtpPass := conf.SmtpPass // change here
	smtpUser := conf.SmtpUser // change here

	// emailConf := &EmailConfig{smtpUser, smtpPass, smtpHost, smtpPort}

	emailauth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	sender := conf.SmtpSenderEmail // change here

	receivers := []string{
		to,
	} // change here

	emailContent := email.NewMessage(subject, message)
	emailContent.From.Address = sender
	emailContent.To = receivers

	err := emailContent.Attach(filename)

	if err != nil {
		fmt.Println(err)
	}

	// send out the email
	err = email.Send(smtpHost+":"+strconv.Itoa(smtpPort), //convert port number from int to string
		emailauth,
		emailContent)
	if err != nil {
		fmt.Println(err)
	}

}
