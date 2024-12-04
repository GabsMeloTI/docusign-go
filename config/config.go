package config

import (
	"github.com/joho/godotenv"
	"os"
	"strings"
)

type Config struct {
	ServerName         string
	ServerPort         string
	Environment        string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBDatabase         string
	DBSSLMode          string
	DBDriver           string
	DocuSignHost       string
	DocuSignAccountId  string
	DocuSignApiKey     string
	DocuSignUsername   string
	DocuSignRSAKey     string
	NgrokAuthToken     string
	AwsRegion          string
	AwsAccessKey       string
	AwsSecretAccessKey string
}

func NewConfig() Config {
	if os.Getenv("ENVIRONMENT") == "" {
		if err := godotenv.Load(".env"); err != nil {
			panic("Error loading env file")
		}
	}
	return Config{
		ServerName:         os.Getenv("SERVER_NAME"),
		ServerPort:         os.Getenv("SERVER_PORT"),
		Environment:        os.Getenv("ENVIRONMENT"),
		DBHost:             os.Getenv("POSTGRES_HOST"),
		DBPort:             os.Getenv("POSTGRES_PORT"),
		DBUser:             os.Getenv("POSTGRES_USER"),
		DBPassword:         os.Getenv("POSTGRES_PASSWORD"),
		DBDatabase:         os.Getenv("POSTGRES_DB"),
		DBSSLMode:          os.Getenv("DB_SSL_MODE"),
		DBDriver:           os.Getenv("DB_DRIVER"),
		DocuSignHost:       os.Getenv("DOCUSIGN_HOST"),
		DocuSignAccountId:  os.Getenv("DOCUSIGN_ACCTID"),
		DocuSignApiKey:     os.Getenv("DOCUSIGN_APIKEY"),
		DocuSignUsername:   os.Getenv("DOCUSIGN_USERNAME"),
		DocuSignRSAKey:     strings.ReplaceAll(os.Getenv("DOCUSIGN_RSA_PRIVATE_KEY"), "\\n", "\n"),
		NgrokAuthToken:     os.Getenv("NGROK_AUTHTOKEN"),
		AwsRegion:          os.Getenv("AWS_REGION"),
		AwsAccessKey:       os.Getenv("AWS_ACCESS_KEY_ID"),
		AwsSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}
}
