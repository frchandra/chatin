package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

type AppConfig struct {
	AppName      string
	IsProduction string
	AppUrl       string
	AppPort      string

	DBHost                    string
	DBUser                    string
	DBPassword                string
	DBName                    string
	DBPort                    string
	DBMaxIdleConnection       int
	DBMaxOpenConnection       int
	DBConnectionMaxLifeMinute time.Duration

	RedisPassword string
	RedisHost     string
	RedisPort     string

	APISecret     string
	TokenDuration string

	AccessSecret  string
	RefreshSecret string
	AccessMinute  time.Duration
	RefreshMinute time.Duration

	AdminName     string
	AdminEmail    string
	AdminPassword string

	DialogflowCredential string
	DialogflowProjectId  string
	DialogflowLanguage   string
	DialogflowTimeZone   string
	DialogflowSessionId  string
}

func NewAppConfig() *AppConfig {
	accessMinute, _ := time.ParseDuration(getEnv("ACCESS_MINUTE", "15m"))
	refreshMinute, _ := time.ParseDuration(getEnv("ACCESS_MINUTE", "120m"))
	dbMaxIdleConnection, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNECTION", "10"))
	dbMaxOpenConnection, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNECTION", "10"))
	dbConnectionMaxLifeMinute, _ := time.ParseDuration(getEnv("DB_CONNECTION_MAX_LIFE_MINUTE", "60m"))

	var appConfig = AppConfig{
		AppName:      getEnv("APP_NAME", "gmcgo"),
		IsProduction: getEnv("IS_PRODUCTION", "false"),
		AppUrl:       getEnv("APP_URL", "127.0.0.1"),
		AppPort:      getEnv("APP_PORT", "8080"),

		DBHost:                    getEnv("DB_HOST", "localhost"),
		DBUser:                    getEnv("DB_USER", "root"),
		DBPassword:                getEnv("DB_PASSWORD", "root"),
		DBName:                    getEnv("DB_NAME", "chatin_db"),
		DBPort:                    getEnv("DB_PORT", "5432"),
		DBMaxIdleConnection:       dbMaxIdleConnection,
		DBMaxOpenConnection:       dbMaxOpenConnection,
		DBConnectionMaxLifeMinute: dbConnectionMaxLifeMinute,

		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisHost:     getEnv("REDIS_HOST", "127.0.0.1"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),

		AccessSecret:  getEnv("ACCESS_SECRET", ""),
		RefreshSecret: getEnv("REFRESH_SECRET", ""),
		AccessMinute:  accessMinute,
		RefreshMinute: refreshMinute,

		AdminName:     getEnv("ADMIN_NAME", ""),
		AdminEmail:    getEnv("ADMIN_EMAIL", ""),
		AdminPassword: getEnv("ADMIN_PASSWORD", ""),

		DialogflowCredential: getEnv("DIALOGFLOW_CREDENTIAL", ""),
		DialogflowProjectId:  getEnv("DIALOGFLOW_PROJECT_ID", ""),
		DialogflowLanguage:   getEnv("DIALOGFLOW_LANGUAGE", ""),
		DialogflowTimeZone:   getEnv("DIALOGFLOW_TIME_ZONE", ""),
		DialogflowSessionId:  getEnv("DIALOGFLOW_SESSION_ID", ""),
	}
	return &appConfig

}

func getEnv(key, fallback string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error on loading .env file")
	}
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
