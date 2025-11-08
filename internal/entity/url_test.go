package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewURL тестирует фабричный метод создания новой сущности URL
func TestNewURL(t *testing.T) {
	// Тест 1: Создание URL с корректными данными
	t.Run("Create new URL entity with valid data", func(t *testing.T) {
		id := 1
		shortURL := "abc123"
		originalURL := "https://example.com/very/long/url"
		userID := "user-123"

		url := NewURL(id, shortURL, originalURL, userID)

		assert.NotNil(t, url, "Созданный URL не должен быть nil")
		assert.Equal(t, id, url.ID, "ID должен соответствовать переданному значению")
		assert.Equal(t, shortURL, url.ShortURL, "ShortURL должен соответствовать переданному значению")
		assert.Equal(t, originalURL, url.OriginalURL, "OriginalURL должен соответствовать переданному значению")
		assert.Equal(t, userID, url.UserID, "UserID должен соответствовать переданному значению")
		assert.False(t, url.IsDeleted, "Новая запись не должна быть помечена как удаленная")
	})

	// Тест 2: Создание URL с нулевыми/пустыми значениями (граничный случай)
	t.Run("Create URL with empty and zero values", func(t *testing.T) {
		url := NewURL(0, "", "", "")

		assert.NotNil(t, url, "URL с нулевыми значениями все равно должен создаваться")
		assert.Equal(t, 0, url.ID, "ID должен быть 0")
		assert.Equal(t, "", url.ShortURL, "ShortURL должен быть пустой строкой")
		assert.Equal(t, "", url.OriginalURL, "OriginalURL должен быть пустой строкой")
		assert.Equal(t, "", url.UserID, "UserID должен быть пустой строкой")
		assert.False(t, url.IsDeleted, "Запись не должна быть помечена как удаленная")
	})
}

// TestURL_GetShortURL тестирует геттер для получения сокращенного URL
func TestURL_GetShortURL(t *testing.T) {
	t.Run("Get short URL from entity", func(t *testing.T) {
		expectedShortURL := "abc123"
		url := NewURL(1, expectedShortURL, "https://example.com", "user-123")

		actualShortURL := url.GetShortURL()

		assert.Equal(t, expectedShortURL, actualShortURL,
			"GetShortURL должен возвращать корректное значение ShortURL")
	})
}

// TestURL_GetOriginalURL тестирует геттер для получения оригинального URL
func TestURL_GetOriginalURL(t *testing.T) {
	t.Run("Get original URL from entity", func(t *testing.T) {
		expectedOriginalURL := "https://example.com/very/long/url"
		url := NewURL(1, "abc123", expectedOriginalURL, "user-123")

		actualOriginalURL := url.GetOriginalURL()

		assert.Equal(t, expectedOriginalURL, actualOriginalURL,
			"GetOriginalURL должен возвращать корректное значение OriginalURL")
	})
}

// TestURL_GetID тестирует геттер для получения идентификатора записи
func TestURL_GetID(t *testing.T) {
	t.Run("Get ID from URL entity", func(t *testing.T) {
		expectedID := 42
		url := NewURL(expectedID, "abc123", "https://example.com", "user-123")

		actualID := url.GetID()

		assert.Equal(t, expectedID, actualID,
			"GetID должен возвращать корректное значение ID")
	})
}

// TestURL_GetUserID тестирует геттер для получения идентификатора пользователя
func TestURL_GetUserID(t *testing.T) {
	t.Run("Get user ID from URL entity", func(t *testing.T) {
		expectedUserID := "user-456"
		url := NewURL(1, "abc123", "https://example.com", expectedUserID)

		actualUserID := url.GetUserID()

		assert.Equal(t, expectedUserID, actualUserID,
			"GetUserID должен возвращать корректное значение UserID")
	})
}

// TestURL_AllMethods тестирует все геттеры вместе для комплексной проверки
func TestURL_AllMethods(t *testing.T) {
	t.Run("Test all getter methods together with comprehensive data", func(t *testing.T) {
		// Подготавливаем тестовые данные
		expectedID := 100
		expectedShortURL := "xyz789"
		expectedOriginalURL := "https://google.com/search?q=golang"
		expectedUserID := "user-789"

		// Создаем сущность URL
		url := NewURL(expectedID, expectedShortURL, expectedOriginalURL, expectedUserID)

		// Проверяем что все геттеры возвращают ожидаемые значения
		assert.Equal(t, expectedID, url.GetID(),
			"GetID должен возвращать корректный идентификатор")
		assert.Equal(t, expectedShortURL, url.GetShortURL(),
			"GetShortURL должен возвращать корректный сокращенный URL")
		assert.Equal(t, expectedOriginalURL, url.GetOriginalURL(),
			"GetOriginalURL должен возвращать корректный оригинальный URL")
		assert.Equal(t, expectedUserID, url.GetUserID(),
			"GetUserID должен возвращать корректный идентификатор пользователя")
	})

	// Дополнительный тест: проверка целостности данных после создания
	t.Run("Data integrity after creation", func(t *testing.T) {
		testCases := []struct {
			name        string
			id          int
			shortURL    string
			originalURL string
			userID      string
		}{
			{"Normal case", 1, "short1", "https://example.com/1", "user1"},
			{"Special characters", 2, "short2", "https://example.com/path?query=param", "user2"},
			{"Long URLs", 3, "short3", "https://very-long-domain-name.com/very/long/path/with/many/segments", "user3"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				url := NewURL(tc.id, tc.shortURL, tc.originalURL, tc.userID)

				// Проверяем что все данные сохраняются корректно
				assert.Equal(t, tc.id, url.GetID())
				assert.Equal(t, tc.shortURL, url.GetShortURL())
				assert.Equal(t, tc.originalURL, url.GetOriginalURL())
				assert.Equal(t, tc.userID, url.GetUserID())
			})
		}
	})
}
