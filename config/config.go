package config

import (
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
	SessionSecret string
	EnableSwagger bool
}

func GetAppConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Logger.Warn().Msg(".env file not found. proceeding with environment variables.")
	}
	appPort := lookupEnvWithDefault("PORT", "8000")
	appPortNumeric, err := strconv.Atoi(appPort)
	if err != nil {
		log.Fatal().Msgf("failed to convert %s to numeric value", appPort)
	}
	dbHost := lookupMandatoryEnv("DB_HOST")
	dbPort := lookupMandatoryEnv("DB_PORT")
	dbPortNumeric, err := strconv.Atoi(dbPort)
	if err != nil {
		log.Fatal().Msgf("failed to convert %s to numeric value", dbPort)
	}
	dbUser := lookupMandatoryEnv("DB_USER")
	dbPassword := lookupMandatoryEnv("DB_PASSWORD")
	dbName := lookupMandatoryEnv("DB_NAME")
	sessionSecret := lookupMandatoryEnv("SESSION_SECRET")
	enableSwagger := lookupEnvWithDefault("ENABLE_SWAGGER", "false")
	enableSwaggerBool, err := strconv.ParseBool(enableSwagger)
	if err != nil {
		log.Fatal().Msgf("failed to convert %s to boolean value", enableSwagger)
	}
	return &AppConfig{
		AppPort:       appPortNumeric,
		DbHost:        dbHost,
		DbPort:        dbPortNumeric,
		DbUser:        dbUser,
		DbPassword:    dbPassword,
		DbName:        dbName,
		SessionSecret: sessionSecret,
		EnableSwagger: enableSwaggerBool,
	}

}

func lookupEnvWithDefault(env, defaultValue string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		return defaultValue
	}
	return value
}

func lookupMandatoryEnv(env string) string {
	value, ok := os.LookupEnv(env)
	if !ok {
		log.Fatal().Msgf("failed to read a mandatory %s env. variable", env)
	}
	return value
}
