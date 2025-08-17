package handler

import (
	"sync"

	"github.com/ParkhomenkoDV/URLShortener/internal/service/server"
)

type Handler struct {
	//repo   storage.Repository // not used yet
	config server.Config
	mutex  sync.Mutex
	data   map[string]string
}

func New(config server.Config) *Handler {
	return &Handler{
		config: config,
		data:   make(map[string]string),
	}
}
