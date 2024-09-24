package controllers

import (
	"backendtku/config"

	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
	Config *config.Config
}

func NewHandler(db *gorm.DB, cfg *config.Config) *Handler {
	return &Handler{
		DB: db,
		Config: cfg,
	}
}

