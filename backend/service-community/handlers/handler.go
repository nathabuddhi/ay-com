package handlers

import (
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func NewHandlers(db *gorm.DB) *Handler {
	return &Handler{
		DB: db,
	}
}
