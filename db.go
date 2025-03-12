package main

import (
	"log"
	"os"
	"time"

	"github.com/Bedrock-Technology/lambda/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

func loadDatabase(dsn string) error {
	d, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\n", log.LstdFlags|log.Lshortfile), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		}),
	})
	if err != nil {
		return err
	}

	db = d
	return db.AutoMigrate(&model.ClaimInfo{})
}
