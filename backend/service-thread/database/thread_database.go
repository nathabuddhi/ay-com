package database

import (
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.Logger.LogMode(logger.Info)
	if err != nil {
		zap.L().Panic("Failed to connect to THREAD SERBICE DATABASE: " + err.Error())
	}
	zap.L().Info("Connected to THREAD SERVICE DATABASE")
	return db
}
