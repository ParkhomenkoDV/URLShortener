package handler

import (
	"testing"

	"github.com/ParkhomenkoDV/URLShortener/internal/service/server"
	"github.com/stretchr/testify/require"
)

func TestNewHandler(t *testing.T) {
	cfg := server.Config{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080",
	}

	h := New(cfg)
	require.NotNil(t, h)
	require.Equal(t, cfg, h.config)
	require.NotNil(t, h.data)
}
