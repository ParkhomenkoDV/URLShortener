package persistence

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ParkhomenkoDV/URLShortener/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestNewFileJSONPersistence(t *testing.T) {
	t.Run("Create new FileJSONPersistence", func(t *testing.T) {
		persistence := NewFileJSONPersistence()
		assert.NotNil(t, persistence)
		assert.IsType(t, &FileJSONPersistence{}, persistence)
	})
}

func TestFileJSONPersistence_Save(t *testing.T) {
	persistence := NewFileJSONPersistence()
	tempDir := t.TempDir()

	t.Run("Save data to file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test_save.json")

		data := map[string]string{
			"abc123": "https://example.com",
			"xyz789": "https://google.com",
		}
		userMap := map[string]string{
			"abc123": "user-1",
			"xyz789": "user-2",
		}

		err := persistence.Save(filePath, data, userMap)
		assert.NoError(t, err)

		// Проверяем, что файл создался
		_, err = os.Stat(filePath)
		assert.NoError(t, err)

		// Проверяем содержимое файла
		content, err := os.ReadFile(filePath)
		assert.NoError(t, err)
		assert.Contains(t, string(content), "abc123")
		assert.Contains(t, string(content), "https://example.com")
		assert.Contains(t, string(content), "user-1")
	})

	t.Run("Save to non-existent directory", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "newdir", "test.json")

		data := map[string]string{"key": "value"}
		userMap := map[string]string{"key": "user"}

		err := persistence.Save(filePath, data, userMap)
		assert.NoError(t, err)

		// Проверяем, что директория и файл созданы
		_, err = os.Stat(filePath)
		assert.NoError(t, err)
	})

	t.Run("Save empty data", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "empty.json")

		data := make(map[string]string)
		userMap := make(map[string]string)

		err := persistence.Save(filePath, data, userMap)
		assert.NoError(t, err)

		content, err := os.ReadFile(filePath)
		assert.NoError(t, err)
		assert.Equal(t, "[]", string(content))
	})
}

func TestFileJSONPersistence_Load(t *testing.T) {
	persistence := NewFileJSONPersistence()
	tempDir := t.TempDir()

	t.Run("Load existing file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test_load.json")

		// Сначала сохраняем данные
		originalData := map[string]string{
			"abc123": "https://example.com",
			"xyz789": "https://google.com",
		}
		originalUserMap := map[string]string{
			"abc123": "user-1",
			"xyz789": "user-2",
		}

		err := persistence.Save(filePath, originalData, originalUserMap)
		assert.NoError(t, err)

		// Загружаем данные
		data, userMap, maxID, err := persistence.Load(filePath)
		assert.NoError(t, err)
		assert.Equal(t, originalData, data)
		assert.Equal(t, originalUserMap, userMap)
		assert.Equal(t, 2, maxID) // Два элемента, максимальный ID = 2
	})

	t.Run("Load non-existent file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "non_existent.json")

		data, userMap, maxID, err := persistence.Load(filePath)
		assert.NoError(t, err)
		assert.Empty(t, data)
		assert.Empty(t, userMap)
		assert.Equal(t, 0, maxID)
	})

	t.Run("Load empty file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "empty_file.json")

		// Создаем пустой файл
		err := os.WriteFile(filePath, []byte(""), 0644)
		assert.NoError(t, err)

		data, userMap, maxID, err := persistence.Load(filePath)
		assert.NoError(t, err)
		assert.Empty(t, data)
		assert.Empty(t, userMap)
		assert.Equal(t, 0, maxID)
	})
}

func TestFileJSONPersistence_SaveAndLoad_Integration(t *testing.T) {
	persistence := NewFileJSONPersistence()
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "integration_test.json")

	t.Run("Save and load cycle", func(t *testing.T) {
		// Подготовка данных
		originalData := map[string]string{
			"short1": "https://example.com/1",
			"short2": "https://example.com/2",
			"short3": "https://example.com/3",
		}
		originalUserMap := map[string]string{
			"short1": "user-a",
			"short2": "user-b",
			"short3": "user-a",
		}

		// Сохранение
		err := persistence.Save(filePath, originalData, originalUserMap)
		assert.NoError(t, err)

		// Загрузка
		loadedData, loadedUserMap, maxID, err := persistence.Load(filePath)
		assert.NoError(t, err)

		// Проверка
		assert.Equal(t, originalData, loadedData)
		assert.Equal(t, originalUserMap, loadedUserMap)
		assert.Equal(t, 3, maxID)
	})
}

func TestFileJSONPersistence_saveRecordsToFile(t *testing.T) {
	persistence := NewFileJSONPersistence()
	tempDir := t.TempDir()

	t.Run("Save records to file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "records.json")

		records := []model.URLRecord{
			{
				ID:          1,
				ShortURL:    "abc123",
				OriginalURL: "https://example.com",
				UserID:      "user-1",
			},
			{
				ID:          2,
				ShortURL:    "xyz789",
				OriginalURL: "https://google.com",
				UserID:      "user-2",
			},
		}

		err := persistence.saveRecordsToFile(filePath, records)
		assert.NoError(t, err)

		// Проверяем содержимое
		content, err := os.ReadFile(filePath)
		assert.NoError(t, err)
		assert.Contains(t, string(content), "abc123")
		assert.Contains(t, string(content), "https://example.com")
		assert.Contains(t, string(content), "user-1")
	})
}

func TestFileJSONPersistence_loadRecordsFromFile(t *testing.T) {
	persistence := NewFileJSONPersistence()
	tempDir := t.TempDir()

	t.Run("Load records from file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "load_records.json")

		// Создаем тестовый файл
		testData := `[
			{
				"id": 1,
				"short_url": "abc123",
				"original_url": "https://example.com",
				"user_id": "user-1"
			},
			{
				"id": 2,
				"short_url": "xyz789",
				"original_url": "https://google.com",
				"user_id": "user-2"
			}
		]`

		err := os.WriteFile(filePath, []byte(testData), 0644)
		assert.NoError(t, err)

		records, err := persistence.loadRecordsFromFile(filePath)
		assert.NoError(t, err)
		assert.Len(t, records, 2)
		assert.Equal(t, "abc123", records[0].ShortURL)
		assert.Equal(t, "https://example.com", records[0].OriginalURL)
		assert.Equal(t, "user-1", records[0].UserID)
	})

	t.Run("Load from non-existent file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "non_existent_records.json")

		records, err := persistence.loadRecordsFromFile(filePath)
		assert.NoError(t, err)
		assert.Empty(t, records)
	})
}
