package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGenerateConfig тестирует создание конфигурации с различными источниками параметров
func TestGenerateConfig(t *testing.T) {
	// Сохраняем оригинальные переменные окружения для восстановления после тестов
	originalServerAddr := os.Getenv("SERVER_ADDRESS")
	originalBaseURL := os.Getenv("BASE_URL")
	originalFilePath := os.Getenv("FILE_STORAGE_PATH")
	originalDatabaseDSN := os.Getenv("DATABASE_DSN")

	// Гарантируем восстановление оригинальных значений переменных окружения после тестов
	defer func() {
		restoreEnvVar("SERVER_ADDRESS", originalServerAddr)
		restoreEnvVar("BASE_URL", originalBaseURL)
		restoreEnvVar("FILE_STORAGE_PATH", originalFilePath)
		restoreEnvVar("DATABASE_DSN", originalDatabaseDSN)
	}()

	// Тест 1: Конфигурация со значениями по умолчанию
	t.Run("Default config generation", func(t *testing.T) {
		// Очищаем переменные окружения для тестирования значений по умолчанию
		os.Unsetenv("SERVER_ADDRESS")
		os.Unsetenv("BASE_URL")
		os.Unsetenv("FILE_STORAGE_PATH")
		os.Unsetenv("DATABASE_DSN")

		config := New()

		assert.NotNil(t, config, "Config should not be nil")
		assert.Equal(t, "http://", config.Protocol, "Default protocol should be http://")
		assert.Contains(t, config.Port, ":", "Port should contain colon prefix")
		assert.Contains(t, config.ShortAddress, "http://", "Short address should contain protocol")
		assert.NotEmpty(t, config.FilePath, "File path should not be empty")
		assert.Equal(t, "your-secret-key-change-in-production", config.AuthSecretKey,
			"Auth secret key should match default value")
	})

	// Тест 2: Конфигурация с переменными окружения (высший приоритет)
	t.Run("Config with environment variables", func(t *testing.T) {
		// Устанавливаем тестовые переменные окружения
		os.Setenv("SERVER_ADDRESS", "localhost:9090")
		os.Setenv("BASE_URL", "https://example.com")
		os.Setenv("FILE_STORAGE_PATH", "/tmp/test.json")
		os.Setenv("DATABASE_DSN", "postgres://user:pass@localhost/db")

		config := New()

		assert.NotNil(t, config, "Config should not be nil")
		assert.Equal(t, "http://", config.Protocol, "Protocol should remain http://")
		assert.Equal(t, ":9090", config.Port, "Port should be extracted from SERVER_ADDRESS")
		assert.Equal(t, "https://example.com", config.ShortAddress, "Short address should use BASE_URL")
		assert.Equal(t, "/tmp/test.json", config.FilePath, "File path should use FILE_STORAGE_PATH")
		assert.Equal(t, "postgres://user:pass@localhost/db", config.AddressDB, "DB address should use DATABASE_DSN")
		assert.Equal(t, "your-secret-key-change-in-production", config.AuthSecretKey,
			"Auth secret key should remain default")
	})
}

// TestConfigStruct тестирует создание и заполнение структуры Config напрямую
func TestConfigStruct(t *testing.T) {
	t.Run("ConfigStruct creation and field assignment", func(t *testing.T) {
		// Создаем конфигурацию напрямую для тестирования структуры
		config := &Config{
			Protocol:      "https://",
			Port:          ":3000",
			ShortAddress:  "https://short.ly",
			FilePath:      "/path/to/file.json",
			AddressDB:     "postgres://localhost/test",
			AuthSecretKey: "test-secret",
		}

		// Проверяем корректность присвоения всех полей
		assert.Equal(t, "https://", config.Protocol, "Protocol should be https://")
		assert.Equal(t, ":3000", config.Port, "Port should be :3000")
		assert.Equal(t, "https://short.ly", config.ShortAddress, "Short address should be https://short.ly")
		assert.Equal(t, "/path/to/file.json", config.FilePath, "File path should match")
		assert.Equal(t, "postgres://localhost/test", config.AddressDB, "DB address should match")
		assert.Equal(t, "test-secret", config.AuthSecretKey, "Auth secret key should match")
	})
}

// restoreEnvVar восстанавливает значение переменной окружения или удаляет ее если значение было пустым
func restoreEnvVar(key, originalValue string) {
	if originalValue != "" {
		os.Setenv(key, originalValue)
	} else {
		os.Unsetenv(key)
	}
}
