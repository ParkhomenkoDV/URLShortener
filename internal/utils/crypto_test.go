package utils

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateShortURL(t *testing.T) {
	t.Run("generates URL of correct length", func(t *testing.T) {
		short := GenerateShortURL(8)
		require.Len(t, short, 8)
	})

	t.Run("generates URL with valid characters", func(t *testing.T) {
		short := GenerateShortURL(8)
		match, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, short)
		require.True(t, match, "contains invalid characters")
	})

	t.Run("generates unique URLs", func(t *testing.T) {
		url1 := GenerateShortURL(8)
		url2 := GenerateShortURL(8)
		require.NotEqual(t, url1, url2)
	})
}
