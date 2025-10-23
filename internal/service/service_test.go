package service

import (
	"testing"

	"github.com/ParkhomenkoDV/URLShortener/internal/config"
	"github.com/ParkhomenkoDV/URLShortener/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestURLShortnerService(t *testing.T) {
	// Инициализируем репозиторий в памяти для тестов
	repo := repository.NewMemory()
	configuration := config.Config{}
	service := New(repo, &configuration)
	defer service.Close()

	t.Run("Create and get short URL", func(t *testing.T) {
		originalURL := "https://example.com/very/long/url"

		// Создаем короткую ссылку
		shortURL, err := service.CreateShortURL(originalURL, "")
		assert.NoError(t, err)
		assert.NotEmpty(t, shortURL)
		assert.Len(t, shortURL, 6)

		// Получаем оригинальную ссылку
		fullURL, err := service.GetFullURL(shortURL)
		assert.NoError(t, err)
		assert.Equal(t, originalURL, fullURL)
	})

	t.Run("Get non-existent short URL returns error", func(t *testing.T) {
		nonExistentKey := "nonexist"

		// Пытаемся получить несуществующую ссылку
		_, err := service.GetFullURL(nonExistentKey)
		assert.Error(t, err)
		assert.Equal(t, "not found", err.Error())
	})

	t.Run("Generate unique short URLs", func(t *testing.T) {
		url1 := "https://first.com"
		url2 := "https://second.com"

		// Генерируем две короткие ссылки
		short1, err1 := service.CreateShortURL(url1, "")
		short2, err2 := service.CreateShortURL(url2, "")

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, short1, short2)

		// Проверяем, что они ведут на разные URL
		full1, err1 := service.GetFullURL(short1)
		full2, err2 := service.GetFullURL(short2)

		assert.NoError(t, err1)
		assert.Equal(t, url1, full1)
		assert.NoError(t, err2)
		assert.Equal(t, url2, full2)
	})

	t.Run("Empty URL handling", func(t *testing.T) {
		emptyURL := ""

		// Не должно паниковать при пустом URL
		shortURL, err := service.CreateShortURL(emptyURL, "")
		assert.NoError(t, err)
		assert.NotEmpty(t, shortURL)

		fullURL, err := service.GetFullURL(shortURL)
		assert.NoError(t, err)
		assert.Equal(t, emptyURL, fullURL)
	})

	t.Run("DB error", func(t *testing.T) {
		err := service.PingPostgreSQL()
		assert.Error(t, err)
	})
}
