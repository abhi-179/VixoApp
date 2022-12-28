package main

import (
	"log"
	"pandoMessagingWalletService/com.pando.messaging/config/kafka"
	logging "pandoMessagingWalletService/com.pando.messaging/logger"
	"pandoMessagingWalletService/com.pando.messaging/repository"
	"pandoMessagingWalletService/com.pando.messaging/router"
	"pandoMessagingWalletService/com.pando.messaging/service"
	"pandoMessagingWalletService/docs"

	config "pandoMessagingWalletService/com.pando.messaging/config"
	database "pandoMessagingWalletService/com.pando.messaging/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	eureka "github.com/xuanbo/eureka-client"

	"sync"
)

var onceRest sync.Once
var (
	//listenAddrApi string

	// kafka
	kafkaBrokerUrl string

	kafkaTopic string
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	onceRest.Do(func() {
		e := echo.New()
		//Setting up the config
		config, err := config.GetConfig("./com.pando.messaging/env/")
		if err != nil {
			log.Print("config file not found")
		}
		//Setting up the Logger
		logger := logging.NewLogger(config.LogFile, config.LogLevel)
		logger.SetReportCaller(true)
		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		kafkaBrokerUrl = config.KafkaURl
		kafkaTopic = config.KafkaTopic
		kafkaProducer, err := kafka.Configure([]string{kafkaBrokerUrl}, "", kafkaTopic)
		if err != nil {
			logger.Error("unable to configure kafka")
			return
		}
		defer kafkaProducer.Close()
		client := eureka.NewClient(&eureka.Config{
			DefaultZone:                  config.ServiceRegistry_URL,
			App:                          config.Appname,
			Port:                         20707,
			RegistryFetchIntervalSeconds: 90,
			RenewalIntervalInSecs:        30,  // Renewal time, default 30S
			DurationInSecs:               300, // service time-effective time, the registration center exceeds that the time period does not receive heartbeat, will be removed from the service list

			Metadata: map[string]interface{}{
				"VERSION":              "0.1.0",
				"NODE_GROUP_ID":        0,
				"PRODUCT_CODE":         "DEFAULT",
				"PRODUCT_VERSION_CODE": "DEFAULT",
				"PRODUCT_ENV_CODE":     "DEFAULT",
				"SERVICE_VERSION_CODE": "DEFAULT",
			},
		})
		// start client, register、heartbeat、refresh
		client.Start()
		docs.SwaggerInfo.Title = "Chat & Wallet Service API"
		docs.SwaggerInfo.Description = "Documentation Wallet API v1.0"
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = "backend.message.pandotest.com/"
		docs.SwaggerInfo.Schemes = []string{"https"}
		e.GET("/swagger-ui/*", echoSwagger.WrapHandler)

		db := database.DB(config)
		walletRepo := repository.NewWalletRepository(db, config)
		walletUc := service.NewWalletUsecase(walletRepo)
		router.NewWalletController(e, walletUc)

		bookingRepo := repository.NewBookingRepository(db, config)
		bookingUc := service.NewBookingUsecase(bookingRepo)
		router.NewBookingController(e, bookingUc)

		reviewRepo := repository.NewReviewRepository(db, config)
		reviewUc := service.NewReviewUsecase(reviewRepo)
		router.NewReviewController(e, reviewUc)

		ipfsRepo := repository.NewIPFSRepository(db, config)
		ipfsUc := service.NewIPFSUsecase(ipfsRepo)
		router.NewIPFSController(e, ipfsUc)

		callRepo := repository.NewcallRepository(db)
		callUc := service.NewcallUsecase(callRepo)
		router.NewCallController(e, callUc)

		chatRepo := repository.NewchatRepository(db, config)
		chatUc := service.NewchatUsecase(chatRepo)
		router.NewChatController(e, chatUc)

		groupRepo := repository.NewgroupRepository(db, config)
		groupUc := service.NewgroupUsecase(groupRepo)
		router.NewGroupController(e, groupUc, config)

		if err := e.Start(config.HostPort); err != nil {
			logger.WithError(err).Fatal("not connected")
		}
	})
}
