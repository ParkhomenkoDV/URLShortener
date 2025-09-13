package config

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
		cfg, err := NewConfig()
		require.NoError(t, err)
		require.Equal(t, "localhost:8080", cfg.ServerAddress)
		require.Equal(t, "http://localhost:8080", cfg.BaseURL)
		require.Equal(t, "data/db.json", cfg.FileStorage)
	})

	t.Run("invalid base URL panics", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
		flag.String("b", "invalid-url", "")
		require.Panics(t, func() { NewConfig() })
	})

	t.Run("environment variables", func(t *testing.T) {
		os.Setenv("SERVER_ADDRESS", "127.0.0.1:9090")
		os.Setenv("BASE_URL", "https://example.com")
		os.Setenv("FILE_STORAGE_PATH", "/tmp/db.json")
		defer func() {
			os.Unsetenv("SERVER_ADDRESS")
			os.Unsetenv("BASE_URL")
			os.Unsetenv("FILE_STORAGE_PATH")
		}()

		flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
		cfg, err := NewConfig()
		require.NoError(t, err)
		require.Equal(t, "127.0.0.1:9090", cfg.ServerAddress)
		require.Equal(t, "https://example.com", cfg.BaseURL)
		require.Equal(t, "/tmp/db.json", cfg.FileStorage)
	})
}
