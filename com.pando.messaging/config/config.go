package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost                 string `mapstructure:"DB_HOST"`
	DBPort                 string `mapstructure:"DB_PORT"`
	DBUser                 string `mapstructure:"DB_USER"`
	DBPass                 string `mapstructure:"DB_PASS"`
	DBName                 string `mapstructure:"DB_NAME"`
	DBType                 string `mapstructure:"DB_TYPE"`
	LogFile                string `mapstructure:"LOG_FILE"`
	LogLevel               string `mapstructure:"LOG_LEVEL"`
	HostPort               string `mapstructure:"HOST_PORT"`
	KafkaTopic             string `mapstructure:"KAFKA_TOPIC"`
	KafkaURl               string `mapstructure:"KAFKA_URL"`
	Registry_Type          string `mapstructure:"EUREKA_REGISTRY_TYPE"`
	ServiceRegistry_URL    string `mapstructure:"EUREKA_SERVICE_REGISTRY_URL"`
	Appname                string `mapstructure:"APPNAME"`
	WalletStatementFileUrl string `mapstructure:"WALLET_STATEMENT_FILE_URL"`
	LogoUrl                string `mapstructure:"LOGO_URL"`
	SmtpHost               string `mapstructure:"SMTP_HOST"`
	SmtpUser               string `mapstructure:"SMTP_USER"`
	SmtpPass               string `mapstructure:"SMTP_PASS"`
	SmtpSenderEmail        string `mapstructure:"SMTP_SENDER_EMAIL"`
	TicketTemplatePath     string `mapstructure:"TICKET_TEMPLATE_PATH"`
	StoragePath            string `mapstructure:"STORAGE_PATH"`
	AdminWalletId          string `mapstructure:"ADMIN_WALLET_ID"`
	WalletApiUrl           string `mapstructure:"WALLET_API_URL"`
	S3Bucket               string `mapstructure:"S3_BUCKET"`
	S3Region               string `mapstructure:"S3_REGION"`
	S3AccessId             string `mapstructure:"S3_ACCESS_ID"`
	S3SecretKey            string `mapstructure:"S3_SECRET_KEY"`
	IPFSURL                string `mapstructure:"IPFS_URL"`
	AwsURL                 string `mapstructure:"AWS_URL"`
	ImageSize              int64  `mapstructure:"IMAGE_SIZE"`
	Batch_size             int    `mapstructure:"BATCH_SIZE"`
	ReportCount1           int64  `mapstructure:"REPORT_COUNT1"`
	ReportCount2           int64  `mapstructure:"REPORT_COUNT2"`
	ReportCount3           int64  `mapstructure:"REPORT_COUNT3"`
	ReportCount4           int64  `mapstructure:"REPORT_COUNT4"`
}

// GetConfig - Function to get Config
func GetConfig(path string) (config *Config, err error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("dev")
	v.SetConfigType("env")
	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err != nil {
		fmt.Println("error in reading config file.")
		return
	}
	err = v.Unmarshal(&config)
	return
}
