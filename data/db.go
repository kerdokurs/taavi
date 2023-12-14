package data

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"kerdo.dev/taavi/logger"
)

var DB *gorm.DB

func Init() {
	var err error

	dsn := getDsn()
	if DB, err = gorm.Open(postgres.Open(dsn)); err != nil {
		logger.Fatalw("error connecting to database", logger.M{
			"err": err.Error(),
		})
	}

	if err = DB.AutoMigrate(
		&Job{},
		&ExecutedJob{},
		&JobMeta{},
	); err != nil {
		logger.Fatalw("error migrating database", logger.M{
			"err": err.Error(),
		})
	}
}

func getDsn() string {
	dsn := os.Getenv("DB_DSN")
	if dsn != "" {
		return dsn
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Tallinn",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
}
