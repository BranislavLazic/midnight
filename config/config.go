package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

type AppConfig struct {
	AppPort       int
	DbHost        string
	DbPort        int
	DbUser        string
	DbPassword    string
	DbName        string
	EnableSwagger bool
}

func GetAppConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Logger.Warn().Msg(".env file not found. proceeding with environment variables.")
	}
	appPort := lookupEnvWithDefault("PORT", "8000")
	appPortNumeric, err := strconv.Atoi(appPort)
	if err != nil {
		return nil, err
	}
	dbHost, err := lookupMandatoryEnv("DB_HOST")
	if err != nil {
		return nil, err
	}
	dbPort, err := lookupMandatoryEnv("DB_PORT")
	if err != nil {
		return nil, err
	}
	dbPortNumeric, err := strconv.Atoi(dbPort)
	if err != nil {
		return nil, err
	}
	dbUser, err := lookupMandatoryEnv("DB_USER")
	if err != nil {
		return nil, err
	}
	dbPassword, err := lookupMandatoryEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}
	dbName, err := lookupMandatoryEnv("DB_NAME")
	if err != nil {
		return nil, err
	}
	enableSwagger := lookupEnvWithDefault("ENABLE_SWAGGER", "false")
	enableSwaggerBool, err := strconv.ParseBool(enableSwagger)
	if err != nil {
		return nil, err
	}
	return &AppConfig{
		AppPort:       appPortNumeric,
		DbHost:        dbHost,
		DbPort:        dbPortNumeric,
		DbUser:        dbUser,
		DbPassword:    dbPassword,
		DbName:        dbName,
		EnableSwagger: enableSwaggerBool,
	}, nil

}

func lookupEnvWithDefault(env, defaultValue string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		return defaultValue
	}
	return value
}

func lookupMandatoryEnv(env string) (string, error) {
	value, ok := os.LookupEnv(env)
	if !ok {
		return "", errors.New(fmt.Sprintf("%s environment variable must be provided.", env))
	}
	return value, nil
}
