package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewURL(t *testing.T) {
	t.Run("Create new URL entity", func(t *testing.T) {
		id := 1
		shortURL := "abc123"
		originalURL := "https://example.com/very/long/url"
		userID := "user-123"

		url := NewURL(id, shortURL, originalURL, userID)

		assert.NotNil(t, url)
		assert.Equal(t, id, url.ID)
		assert.Equal(t, shortURL, url.ShortURL)
		assert.Equal(t, originalURL, url.OriginalURL)
		assert.Equal(t, userID, url.UserID)
	})

	t.Run("Create URL with empty values", func(t *testing.T) {
		url := NewURL(0, "", "", "")

		assert.NotNil(t, url)
		assert.Equal(t, 0, url.ID)
		assert.Equal(t, "", url.ShortURL)
		assert.Equal(t, "", url.OriginalURL)
		assert.Equal(t, "", url.UserID)
	})
}

func TestURL_GetShortURL(t *testing.T) {
	t.Run("Get short URL", func(t *testing.T) {
		shortURL := "abc123"
		url := NewURL(1, shortURL, "https://example.com", "user-123")

		result := url.GetShortURL()

		assert.Equal(t, shortURL, result)
	})
}

func TestURL_GetOriginalURL(t *testing.T) {
	t.Run("Get original URL", func(t *testing.T) {
		originalURL := "https://example.com/very/long/url"
		url := NewURL(1, "abc123", originalURL, "user-123")

		result := url.GetOriginalURL()

		assert.Equal(t, originalURL, result)
	})
}

func TestURL_GetID(t *testing.T) {
	t.Run("Get ID", func(t *testing.T) {
		id := 42
		url := NewURL(id, "abc123", "https://example.com", "user-123")

		result := url.GetID()

		assert.Equal(t, id, result)
	})
}

func TestURL_GetUserID(t *testing.T) {
	t.Run("Get user ID", func(t *testing.T) {
		userID := "user-456"
		url := NewURL(1, "abc123", "https://example.com", userID)

		result := url.GetUserID()

		assert.Equal(t, userID, result)
	})
}

func TestURL_AllMethods(t *testing.T) {
	t.Run("Test all getter methods together", func(t *testing.T) {
		id := 100
		shortURL := "xyz789"
		originalURL := "https://google.com/search?q=golang"
		userID := "user-789"

		url := NewURL(id, shortURL, originalURL, userID)

		assert.Equal(t, id, url.GetID())
		assert.Equal(t, shortURL, url.GetShortURL())
		assert.Equal(t, originalURL, url.GetOriginalURL())
		assert.Equal(t, userID, url.GetUserID())
	})
}
