package db

import (
	"fmt"
	config "pandoMessagingWalletService/com.pando.messaging/config"
	"pandoMessagingWalletService/com.pando.messaging/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DB(config *config.Config) *gorm.DB {
	psqlInfo1 := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPass, config.DBName)
	db, err := gorm.Open(postgres.Open(psqlInfo1), &gorm.Config{})
	if err != nil {
		logger.Logger.WithError(err).WithField("err", err).Errorf("Database not connected")
		panic(err)
	}
	logger.Logger.Info("Database connected successfully.")
	return db
}
