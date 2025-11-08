package service

import (
	"testing"

	"github.com/ParkhomenkoDV/URLShortener/internal/config"
	"github.com/ParkhomenkoDV/URLShortener/internal/repository"
	"github.com/stretchr/testify/assert"
)

// TestURLShortnerService содержит тесты для сервиса сокращения URL
func TestURLShortnerService(t *testing.T) {
	// Инициализируем репозиторий в памяти для тестов
	repo := repository.NewMemory()
	configuration := config.New()
	service := New(repo, configuration)
	defer service.Close()

	// Тест создания и получения короткой ссылки
	t.Run("Create and get short URL", func(t *testing.T) {
		originalURL := "https://example.com/very/long/url"

		// Создаем короткую ссылку
		shortURL, err := service.CreateShortURL(originalURL, "")
		assert.NoError(t, err)
		assert.NotEmpty(t, shortURL)
		assert.Len(t, shortURL, configuration.LengthKey)

		// Получаем оригинальную ссылку по короткой
		fullURL, err := service.GetFullURL(shortURL)
		assert.NoError(t, err)
		assert.Equal(t, originalURL, fullURL)
	})

	// Тест попытки получения несуществующей короткой ссылки
	t.Run("Get non-existent short URL returns error", func(t *testing.T) {
		nonExistentKey := "nonexist"

		// Пытаемся получить несуществующую ссылку
		_, err := service.GetFullURL(nonExistentKey)
		assert.Error(t, err)
		assert.Equal(t, "not found", err.Error())
	})

	// Тест генерации уникальных коротких URL для разных исходных URL
	t.Run("Generate unique short URLs", func(t *testing.T) {
		url1 := "https://first.com"
		url2 := "https://second.com"

		// Генерируем две короткие ссылки
		short1, err1 := service.CreateShortURL(url1, "")
		short2, err2 := service.CreateShortURL(url2, "")

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, short1, short2) // Убеждаемся, что ссылки разные

		// Проверяем, что они ведут на соответствующие оригинальные URL
		full1, err1 := service.GetFullURL(short1)
		full2, err2 := service.GetFullURL(short2)

		assert.NoError(t, err1)
		assert.Equal(t, url1, full1)
		assert.NoError(t, err2)
		assert.Equal(t, url2, full2)
	})

	// Тест обработки пустого URL
	t.Run("Empty URL handling", func(t *testing.T) {
		emptyURL := ""

		// Проверяем, что сервис корректно обрабатывает пустой URL
		shortURL, err := service.CreateShortURL(emptyURL, "")
		assert.NoError(t, err)
		assert.NotEmpty(t, shortURL)

		// Проверяем, что можно получить пустой URL по короткой ссылке
		fullURL, err := service.GetFullURL(shortURL)
		assert.NoError(t, err)
		assert.Equal(t, emptyURL, fullURL)
	})

	// Тест проверки соединения с PostgreSQL (ожидается ошибка при использовании in-memory репозитория)
	t.Run("DB error", func(t *testing.T) {
		err := service.PingPostgreSQL()
		assert.Error(t, err)
	})
}
