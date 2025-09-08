package handler

import (
	"sync"

	"github.com/ParkhomenkoDV/URLShortener/internal/config"
	"github.com/ParkhomenkoDV/URLShortener/internal/storage"
)

type Handler struct {
	//repo   storage.Repository // not used yet
	config config.Config
	mutex  sync.Mutex
	db     *storage.DB
}

func New(config config.Config, db *storage.DB) *Handler {
	return &Handler{
		config: config,
		db:     db,
	}
}
