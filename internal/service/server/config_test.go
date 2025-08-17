package server

import (
	"flag"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
		cfg := New()
		require.Equal(t, "localhost:8080", cfg.ServerAddress)
		require.Equal(t, "http://localhost:8080", cfg.BaseURL)
	})

	t.Run("invalid base URL panics", func(t *testing.T) {
		flag.CommandLine = flag.NewFlagSet("test", flag.ContinueOnError)
		flag.String("b", "invalid-url", "")
		require.Panics(t, func() { New() })
	})
}
