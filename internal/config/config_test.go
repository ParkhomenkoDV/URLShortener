package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateConfig(t *testing.T) {
	// Сохраняем оригинальные переменные окружения
	originalServerAddr := os.Getenv("SERVER_ADDRESS")
	originalBaseURL := os.Getenv("BASE_URL")
	originalFilePath := os.Getenv("FILE_STORAGE_PATH")
	originalDatabaseDSN := os.Getenv("DATABASE_DSN")

	// Восстанавливаем после теста
	defer func() {
		if originalServerAddr != "" {
			os.Setenv("SERVER_ADDRESS", originalServerAddr)
		} else {
			os.Unsetenv("SERVER_ADDRESS")
		}
		if originalBaseURL != "" {
			os.Setenv("BASE_URL", originalBaseURL)
		} else {
			os.Unsetenv("BASE_URL")
		}
		if originalFilePath != "" {
			os.Setenv("FILE_STORAGE_PATH", originalFilePath)
		} else {
			os.Unsetenv("FILE_STORAGE_PATH")
		}
		if originalDatabaseDSN != "" {
			os.Setenv("DATABASE_DSN", originalDatabaseDSN)
		} else {
			os.Unsetenv("DATABASE_DSN")
		}
	}()

	t.Run("Default config generation", func(t *testing.T) {
		// Очищаем переменные окружения для чистого теста
		os.Unsetenv("SERVER_ADDRESS")
		os.Unsetenv("BASE_URL")
		os.Unsetenv("FILE_STORAGE_PATH")
		os.Unsetenv("DATABASE_DSN")

		config := New()

		assert.NotNil(t, config)
		assert.Equal(t, "http://", config.Protocol)
		assert.Contains(t, config.Port, ":")               // Порт должен содержать двоеточие
		assert.Contains(t, config.ShortAddress, "http://") // Адрес должен содержать протокол
		assert.NotEmpty(t, config.FilePath)                // Путь к файлу не должен быть пустым
		assert.Equal(t, "your-secret-key-change-in-production", config.AuthSecretKey)
	})

	t.Run("Config with environment variables", func(t *testing.T) {
		// Устанавливаем переменные окружения
		os.Setenv("SERVER_ADDRESS", "localhost:9090")
		os.Setenv("BASE_URL", "https://example.com")
		os.Setenv("FILE_STORAGE_PATH", "/tmp/test.json")
		os.Setenv("DATABASE_DSN", "postgres://user:pass@localhost/db")

		config := New()

		assert.NotNil(t, config)
		assert.Equal(t, "http://", config.Protocol)
		assert.Equal(t, ":9090", config.Port)
		assert.Equal(t, "https://example.com", config.ShortAddress)
		assert.Equal(t, "/tmp/test.json", config.FilePath)
		assert.Equal(t, "postgres://user:pass@localhost/db", config.AddressDB)
		assert.Equal(t, "your-secret-key-change-in-production", config.AuthSecretKey)
	})
}

func TestConfigStruct(t *testing.T) {
	t.Run("ConfigStruct creation", func(t *testing.T) {
		config := &Config{
			Protocol:      "https://",
			Port:          ":3000",
			ShortAddress:  "https://short.ly",
			FilePath:      "/path/to/file.json",
			AddressDB:     "postgres://localhost/test",
			AuthSecretKey: "test-secret",
		}

		assert.Equal(t, "https://", config.Protocol)
		assert.Equal(t, ":3000", config.Port)
		assert.Equal(t, "https://short.ly", config.ShortAddress)
		assert.Equal(t, "/path/to/file.json", config.FilePath)
		assert.Equal(t, "postgres://localhost/test", config.AddressDB)
		assert.Equal(t, "test-secret", config.AuthSecretKey)
	})
}
