package app

import (
	"os"

	"github.com/joho/godotenv"
	"kerdo.dev/taavi/logger"
)

func InitEnv() {
	// Only load .env file in development mode
	if os.Getenv("TAAVI_ENV") != "dev" {
		return
	}

	if err := godotenv.Load(); err != nil {
		logger.Errorw("error loading environment variables", logger.M{
			"err": err.Error(),
		})
	}
}
