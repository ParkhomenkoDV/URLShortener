package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMemoryRepository содержит тесты для in-memory репозитория
func TestMemoryRepository(t *testing.T) {
	// Инициализация репозитория в памяти
	repo := NewMemory()
	defer repo.Close()

	// Тест успешного сохранения и получения значения
	t.Run("Set and Get value successfully", func(t *testing.T) {
		key := "testKey"
		value := "https://example.com"
		userID := "user123"

		// Записываем значение
		err := repo.SetValue(key, value, userID)
		assert.NoError(t, err)

		// Получаем значение и проверяем корректность
		result, err := repo.GetLongValue(key)
		assert.NoError(t, err)
		assert.Equal(t, value, result)
	})

	// Тест попытки получения несуществующего ключа
	t.Run("Get non-existent key returns error", func(t *testing.T) {
		nonExistentKey := "nonExistentKey"

		// Пытаемся получить несуществующий ключ
		result, err := repo.GetLongValue(nonExistentKey)
		assert.Error(t, err)
		assert.Equal(t, "", result)
		assert.Equal(t, "not found key in database", err.Error())
	})

	// Тест обработки конфликта при перезаписи существующего ключа
	t.Run("Overwrite existing key", func(t *testing.T) {
		key := "existingKey"
		firstValue := "https://first.com"
		secondValue := "https://second.com"
		userID := "user456"

		// Первая запись
		err := repo.SetValue(key, firstValue, userID)
		assert.NoError(t, err)
		firstResult, err := repo.GetLongValue(key)
		assert.NoError(t, err)
		assert.Equal(t, firstValue, firstResult)

		// Попытка перезаписи (должна вернуть ошибку)
		err = repo.SetValue(key, secondValue, userID)
		assert.Error(t, err)

		// Проверяем, что оригинальное значение не изменилось
		currentValue, err := repo.GetLongValue(key)
		assert.NoError(t, err)
		assert.Equal(t, firstValue, currentValue)
	})

	// Тест получения URL конкретного пользователя
	t.Run("Get user URLs", func(t *testing.T) {
		userID := "user789"

		// Создаем несколько URL для пользователя
		_ = repo.SetValue("key1", "https://example1.com", userID)
		_ = repo.SetValue("key2", "https://example2.com", userID)
		_ = repo.SetValue("key3", "https://example3.com", "anotherUser")

		// Получаем URL пользователя
		userURLs, err := repo.GetUserURLs(userID)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(userURLs)) // Должны получить только 2 URL
	})
}
