package handler

import (
	"testing"

	"github.com/ParkhomenkoDV/URLShortener/internal/config"
	"github.com/ParkhomenkoDV/URLShortener/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	cfg := config.Config{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080",
	}
	db := storage.New()

	h := New(&cfg, db)
	require.NotNil(t, h)
	require.Equal(t, cfg, h.config)
	require.NotNil(t, h.db)
}
