package database

import (
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=Pratama05 dbname=ay-user-service port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.Logger.LogMode(logger.Info)
	if err != nil {
		zap.L().Panic("Failed to connect to USER SERBICE DATABASE: " + err.Error())
	}
	zap.L().Info("Connected to USER SERVICE DATABASE")
	return db
}
