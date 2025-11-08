package persistence

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ParkhomenkoDV/URLShortener/internal/model"
	"github.com/stretchr/testify/assert"
)

// TestNewFileJSONPersistence тестирует создание нового экземпляра FileJSONPersistence
func TestNewFileJSONPersistence(t *testing.T) {
	t.Run("Create new FileJSONPersistence instance", func(t *testing.T) {
		persistence := NewFileJSONPersistence()
		assert.NotNil(t, persistence, "Созданный экземпляр не должен быть nil")
		assert.IsType(t, &FileJSONPersistence{}, persistence,
			"Тип возвращаемого объекта должен быть *FileJSONPersistence")
	})
}

// TestFileJSONPersistence_Save тестирует функциональность сохранения данных
func TestFileJSONPersistence_Save(t *testing.T) {
	persistence := NewFileJSONPersistence()
	tempDir := t.TempDir() // Создаем временную директорию для тестов

	// Тест 1: Сохранение данных с несколькими URL
	t.Run("Save multiple URL records to file", func(t *testing.T) {
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
		assert.NoError(t, err, "Сохранение данных не должно возвращать ошибку")

		// Проверяем, что файл был создан
		_, err = os.Stat(filePath)
		assert.NoError(t, err, "Файл должен существовать после сохранения")

		// Проверяем содержимое файла
		content, err := os.ReadFile(filePath)
		assert.NoError(t, err, "Чтение файла не должно возвращать ошибку")
		assert.Contains(t, string(content), "abc123", "Файл должен содержать короткий URL")
		assert.Contains(t, string(content), "https://example.com", "Файл должен содержать оригинальный URL")
		assert.Contains(t, string(content), "user-1", "Файл должен содержать идентификатор пользователя")
	})

	// Тест 2: Сохранение в несуществующую директорию (должно создавать директорию автоматически)
	t.Run("Save to non-existent directory creates directory automatically", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "newdir", "test.json")

		data := map[string]string{"key": "value"}
		userMap := map[string]string{"key": "user"}

		err := persistence.Save(filePath, data, userMap)
		assert.NoError(t, err, "Сохранение должно создавать директории автоматически")

		// Проверяем, что файл и директория были созданы
		_, err = os.Stat(filePath)
		assert.NoError(t, err, "Файл должен быть создан в новой директории")
	})

	// Тест 3: Сохранение пустых данных
	t.Run("Save empty data creates valid JSON file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "empty.json")

		data := make(map[string]string)
		userMap := make(map[string]string)

		err := persistence.Save(filePath, data, userMap)
		assert.NoError(t, err, "Сохранение пустых данных не должно возвращать ошибку")

		// Проверяем что создается валидный JSON с пустым массивом
		content, err := os.ReadFile(filePath)
		assert.NoError(t, err, "Чтение файла не должно возвращать ошибку")
		assert.Equal(t, "[]", string(content),
			"Пустые данные должны сохраняться как пустой JSON массив")
	})
}

// TestFileJSONPersistence_Load тестирует функциональность загрузки данных
func TestFileJSONPersistence_Load(t *testing.T) {
	persistence := NewFileJSONPersistence()
	tempDir := t.TempDir()

	// Тест 1: Загрузка существующего файла с данными
	t.Run("Load existing file with data returns correct mappings", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test_load.json")

		// Сначала сохраняем тестовые данные
		originalData := map[string]string{
			"abc123": "https://example.com",
			"xyz789": "https://google.com",
		}
		originalUserMap := map[string]string{
			"abc123": "user-1",
			"xyz789": "user-2",
		}

		err := persistence.Save(filePath, originalData, originalUserMap)
		assert.NoError(t, err, "Предварительное сохранение должно завершиться успешно")

		// Загружаем данные и проверяем корректность
		data, userMap, maxID, err := persistence.Load(filePath)
		assert.NoError(t, err, "Загрузка данных не должна возвращать ошибку")
		assert.Equal(t, originalData, data, "Загруженные данные должны соответствовать сохраненным")
		assert.Equal(t, originalUserMap, userMap, "Загруженные userMap должны соответствовать сохраненным")
		assert.Equal(t, 2, maxID, "MaxID должен равняться количеству элементов (2)")
	})

	// Тест 2: Загрузка несуществующего файла (должна возвращать пустые данные без ошибки)
	t.Run("Load non-existent file returns empty data without error", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "non_existent.json")

		data, userMap, maxID, err := persistence.Load(filePath)
		assert.NoError(t, err, "Загрузка несуществующего файла не должна возвращать ошибку")
		assert.Empty(t, data, "Data должна быть пустой для несуществующего файла")
		assert.Empty(t, userMap, "UserMap должна быть пустой для несуществующего файла")
		assert.Equal(t, 0, maxID, "MaxID должен быть 0 для несуществующего файла")
	})

	// Тест 3: Загрузка пустого файла
	t.Run("Load empty file returns empty data without error", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "empty_file.json")

		// Создаем пустой файл
		err := os.WriteFile(filePath, []byte(""), 0644)
		assert.NoError(t, err, "Создание пустого файла не должно возвращать ошибку")

		data, userMap, maxID, err := persistence.Load(filePath)
		assert.NoError(t, err, "Загрузка пустого файла не должна возвращать ошибку")
		assert.Empty(t, data, "Data должна быть пустой для пустого файла")
		assert.Empty(t, userMap, "UserMap должна быть пустой для пустого файла")
		assert.Equal(t, 0, maxID, "MaxID должен быть 0 для пустого файла")
	})
}

// TestFileJSONPersistence_SaveAndLoad_Integration тестирует полный цикл сохранения и загрузки
func TestFileJSONPersistence_SaveAndLoad_Integration(t *testing.T) {
	persistence := NewFileJSONPersistence()
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "integration_test.json")

	t.Run("Complete save and load cycle preserves data integrity", func(t *testing.T) {
		// Подготовка тестовых данных
		originalData := map[string]string{
			"short1": "https://example.com/1",
			"short2": "https://example.com/2",
			"short3": "https://example.com/3",
		}
		originalUserMap := map[string]string{
			"short1": "user-a",
			"short2": "user-b",
			"short3": "user-a", // Один пользователь может иметь несколько URL
		}

		// Сохранение данных
		err := persistence.Save(filePath, originalData, originalUserMap)
		assert.NoError(t, err, "Сохранение должно завершиться успешно")

		// Загрузка данных
		loadedData, loadedUserMap, maxID, err := persistence.Load(filePath)
		assert.NoError(t, err, "Загрузка должна завершиться успешно")

		// Проверка целостности данных
		assert.Equal(t, originalData, loadedData,
			"Загруженные данные должны полностью соответствовать сохраненным")
		assert.Equal(t, originalUserMap, loadedUserMap,
			"Загруженные userMap должны полностью соответствовать сохраненным")
		assert.Equal(t, 3, maxID,
			"MaxID должен равняться количеству сохраненных элементов")
	})
}

// TestFileJSONPersistence_saveRecordsToFile тестирует внутренний метод сохранения записей
func TestFileJSONPersistence_saveRecordsToFile(t *testing.T) {
	persistence := NewFileJSONPersistence()
	tempDir := t.TempDir()

	t.Run("Save URL records to file creates valid JSON", func(t *testing.T) {
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
		assert.NoError(t, err, "Сохранение записей не должно возвращать ошибку")

		// Проверяем содержимое файла
		content, err := os.ReadFile(filePath)
		assert.NoError(t, err, "Чтение файла не должно возвращать ошибку")
		assert.Contains(t, string(content), "abc123", "Файл должен содержать первый короткий URL")
		assert.Contains(t, string(content), "https://example.com", "Файл должен содержать первый оригинальный URL")
		assert.Contains(t, string(content), "user-1", "Файл должен содержать первого пользователя")
	})
}

// TestFileJSONPersistence_loadRecordsFromFile тестирует внутренний метод загрузки записей
func TestFileJSONPersistence_loadRecordsFromFile(t *testing.T) {
	persistence := NewFileJSONPersistence()
	tempDir := t.TempDir()

	// Тест 1: Загрузка записей из существующего файла
	t.Run("Load records from existing file returns correct data", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "load_records.json")

		// Создаем тестовый JSON файл с данными
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
		assert.NoError(t, err, "Создание тестового файла не должно возвращать ошибку")

		records, err := persistence.loadRecordsFromFile(filePath)
		assert.NoError(t, err, "Загрузка записей не должна возвращать ошибку")
		assert.Len(t, records, 2, "Должно быть загружено 2 записи")
		assert.Equal(t, "abc123", records[0].ShortURL, "Первая запись должна иметь правильный ShortURL")
		assert.Equal(t, "https://example.com", records[0].OriginalURL, "Первая запись должна иметь правильный OriginalURL")
		assert.Equal(t, "user-1", records[0].UserID, "Первая запись должна иметь правильный UserID")
	})

	// Тест 2: Загрузка из несуществующего файла
	t.Run("Load from non-existent file returns empty slice without error", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "non_existent_records.json")

		records, err := persistence.loadRecordsFromFile(filePath)
		assert.NoError(t, err, "Загрузка несуществующего файла не должна возвращать ошибку")
		assert.Empty(t, records, "Для несуществующего файла должен возвращаться пустой слайс")
	})
}
