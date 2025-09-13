package server

import (
	"net/http"
	"os"
	"testing"

	"github.com/ParkhomenkoDV/URLShortener/internal/config"
	"github.com/ParkhomenkoDV/URLShortener/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	t.Run("server creation", func(t *testing.T) {
		cfg := &config.Config{
			ServerAddress: "localhost:8080",
			FileStorage:   "/tmp/test_db.json",
		}
		handler := http.NewServeMux()
		db := storage.New()

		server := New(cfg, handler, db)
		require.NotNil(t, server)
		require.Equal(t, cfg, server.config)
		require.Equal(t, db, server.db)
	})

	t.Run("save data", func(t *testing.T) {
		tempFile := "/tmp/test_save.json"
		defer os.Remove(tempFile)

		cfg := &config.Config{FileStorage: tempFile}
		db := storage.New()
		db.Set("testKey", "https://example.com")

		server := &Server{config: cfg, db: db}
		err := server.saveData()
		require.NoError(t, err)

		// Verify file was created and contains data
		_, err = os.Stat(tempFile)
		require.NoError(t, err)
	})

}
