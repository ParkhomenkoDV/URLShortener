package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCreateDB тестирует создание и инициализацию новой базы данных
func TestCreateDB(t *testing.T) {
	t.Run("Create new database instance with correct initial state", func(t *testing.T) {
		db := New()

		assert.NotNil(t, db, "Созданный экземпляр DB не должен быть nil")
		assert.NotNil(t, db.data, "Мапа data должна быть инициализирована")
		assert.NotNil(t, db.userMap, "Мапа userMap должна быть инициализирована")
		assert.Equal(t, 0, db.counter, "Счетчик должен начинаться с 0")
		assert.NotNil(t, db.persistence, "Persistence слой должен быть инициализирован")
		assert.Empty(t, db.data, "Новая БД должна содержать пустую мапу data")
		assert.Empty(t, db.userMap, "Новая БД должна содержать пустую мапу userMap")
	})
}

// TestDB_Get тестирует функциональность получения данных из хранилища
func TestDB_Get(t *testing.T) {
	// Тест 1: Получение существующего ключа
	t.Run("Get existing key returns correct value and exists flag", func(t *testing.T) {
		db := New()
		db.data["test-key"] = "test-value"

		value, exists := db.Get("test-key")

		assert.True(t, exists, "Для существующего ключа должен возвращаться exists=true")
		assert.Equal(t, "test-value", value, "Должно возвращаться корректное значение")
	})

	// Тест 2: Получение несуществующего ключа
	t.Run("Get non-existing key returns empty value and false exists flag", func(t *testing.T) {
		db := New()

		value, exists := db.Get("non-existing")

		assert.False(t, exists, "Для несуществующего ключа должен возвращаться exists=false")
		assert.Equal(t, "", value, "Для несуществующего ключа должна возвращаться пустая строка")
	})
}

// TestDB_Set тестирует базовую функциональность сохранения данных
func TestDB_Set(t *testing.T) {
	t.Run("Set value without user increments counter and stores data correctly", func(t *testing.T) {
		db := New()

		db.Set("key1", "value1")

		// Проверяем что данные сохранились
		value, exists := db.Get("key1")
		assert.True(t, exists, "Ключ должен существовать после вызова Set")
		assert.Equal(t, "value1", value, "Должно сохраняться корректное значение")

		// Проверяем что счетчик увеличился
		assert.Equal(t, 1, db.counter, "Счетчик должен увеличиться после добавления записи")

		// Проверяем что userID пустой для анонимной записи
		userID, userExists := db.GetUserID("key1")
		assert.True(t, userExists, "Запись должна существовать в userMap")
		assert.Equal(t, "", userID, "UserID должен быть пустым для анонимной записи")
	})
}

// TestDB_SetWithUser тестирует сохранение данных с привязкой к пользователю
func TestDB_SetWithUser(t *testing.T) {
	// Тест 1: Сохранение одной записи с пользователем
	t.Run("Set value with user stores both data and user mapping", func(t *testing.T) {
		db := New()

		db.SetWithUser("key1", "value1", "user-123")

		// Проверяем данные
		value, exists := db.Get("key1")
		assert.True(t, exists, "Ключ должен существовать после вызова SetWithUser")
		assert.Equal(t, "value1", value, "Должно сохраняться корректное значение")

		// Проверяем счетчик
		assert.Equal(t, 1, db.counter, "Счетчик должен увеличиться после добавления записи")

		// Проверяем привязку пользователя
		userID, userExists := db.GetUserID("key1")
		assert.True(t, userExists, "Запись должна существовать в userMap")
		assert.Equal(t, "user-123", userID, "Должен сохраняться корректный userID")
	})

	// Тест 2: Сохранение нескольких записей с разными пользователями
	t.Run("Set multiple values with different users maintains correct mappings", func(t *testing.T) {
		db := New()

		db.SetWithUser("key1", "value1", "user-1")
		db.SetWithUser("key2", "value2", "user-2")
		db.SetWithUser("key3", "value3", "user-1") // Один пользователь, несколько URL

		assert.Equal(t, 3, db.counter, "Счетчик должен равняться количеству добавленных записей")

		// Проверяем значения
		value1, _ := db.Get("key1")
		assert.Equal(t, "value1", value1, "Первая запись должна содержать корректное значение")

		value2, _ := db.Get("key2")
		assert.Equal(t, "value2", value2, "Вторая запись должна содержать корректное значение")

		// Проверяем пользователей
		userID1, _ := db.GetUserID("key1")
		assert.Equal(t, "user-1", userID1, "Первая запись должна быть привязана к user-1")

		userID2, _ := db.GetUserID("key2")
		assert.Equal(t, "user-2", userID2, "Вторая запись должна быть привязана к user-2")
	})
}

// TestDB_GetUserID тестирует получение идентификатора пользователя по короткому URL
func TestDB_GetUserID(t *testing.T) {
	// Тест 1: Получение userID для существующей записи
	t.Run("Get existing user ID returns correct user and exists flag", func(t *testing.T) {
		db := New()
		db.SetWithUser("key1", "value1", "user-123")

		userID, exists := db.GetUserID("key1")

		assert.True(t, exists, "Для существующей записи должен возвращаться exists=true")
		assert.Equal(t, "user-123", userID, "Должен возвращаться корректный userID")
	})

	// Тест 2: Получение userID для несуществующей записи
	t.Run("Get non-existing user ID returns empty string and false exists flag", func(t *testing.T) {
		db := New()

		userID, exists := db.GetUserID("non-existing")

		assert.False(t, exists, "Для несуществующей записи должен возвращаться exists=false")
		assert.Equal(t, "", userID, "Для несуществующей записи должна возвращаться пустая строка")
	})
}

// TestDB_GetUserURLs тестирует получение всех URL пользователя
func TestDB_GetUserURLs(t *testing.T) {
	// Тест 1: Получение URL существующего пользователя с несколькими записями
	t.Run("Get URLs for existing user with multiple records", func(t *testing.T) {
		db := New()

		db.SetWithUser("short1", "https://example.com/1", "user-1")
		db.SetWithUser("short2", "https://example.com/2", "user-2")
		db.SetWithUser("short3", "https://example.com/3", "user-1") // Второй URL для user-1

		urls := db.GetUserURLs("user-1")

		assert.Len(t, urls, 2, "Должно вернуться 2 URL для пользователя user-1")

		// Проверяем что получили правильные URL
		foundShort1 := false
		foundShort3 := false

		for _, urlMap := range urls {
			shortURL := urlMap["short_url"]
			switch shortURL {
			case "short1":
				foundShort1 = true
				assert.Equal(t, "https://example.com/1", urlMap["original_url"],
					"Первый URL должен содержать корректный оригинальный URL")
			case "short3":
				foundShort3 = true
				assert.Equal(t, "https://example.com/3", urlMap["original_url"],
					"Третий URL должен содержать корректный оригинальный URL")
			}
		}

		assert.True(t, foundShort1, "Должен присутствовать short1 в результатах")
		assert.True(t, foundShort3, "Должен присутствовать short3 в результатах")
	})

	// Тест 2: Получение URL для несуществующего пользователя
	t.Run("Get URLs for non-existing user returns empty slice", func(t *testing.T) {
		db := New()
		db.SetWithUser("short1", "https://example.com/1", "user-1")

		urls := db.GetUserURLs("non-existing-user")

		assert.Empty(t, urls, "Для несуществующего пользователя должен возвращаться пустой слайс")
	})

	// Тест 3: Получение URL когда в БД нет данных
	t.Run("Get URLs when database is empty returns empty slice", func(t *testing.T) {
		db := New()

		urls := db.GetUserURLs("any-user")

		assert.Empty(t, urls, "Для пустой БД должен возвращаться пустой слайс")
	})
}

// TestDB_SaveToFile тестирует сохранение данных в файл
func TestDB_SaveToFile(t *testing.T) {
	t.Run("Save data to file creates file with correct content", func(t *testing.T) {
		db := New()
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test_save.json")

		db.SetWithUser("short1", "https://example.com/1", "user-1")
		db.SetWithUser("short2", "https://example.com/2", "user-2")

		err := db.SaveToFile(filePath)

		assert.NoError(t, err, "Сохранение в файл не должно возвращать ошибку")

		// Проверяем что файл был создан
		_, err = os.Stat(filePath)
		assert.NoError(t, err, "Файл должен существовать после сохранения")

		// Проверяем содержимое файла
		content, err := os.ReadFile(filePath)
		assert.NoError(t, err, "Чтение файла не должно возвращать ошибку")
		assert.Contains(t, string(content), "short1", "Файл должен содержать первый короткий URL")
		assert.Contains(t, string(content), "https://example.com/1", "Файл должен содержать первый оригинальный URL")
		assert.Contains(t, string(content), "user-1", "Файл должен содержать первого пользователя")
	})
}

// TestDB_LoadFromFile тестирует загрузку данных из файла
func TestDB_LoadFromFile(t *testing.T) {
	// Тест 1: Загрузка данных из существующего файла
	t.Run("Load data from existing file restores complete state", func(t *testing.T) {
		// Создаем первую БД и сохраняем данные
		db1 := New()
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test_load.json")

		db1.SetWithUser("short1", "https://example.com/1", "user-1")
		db1.SetWithUser("short2", "https://example.com/2", "user-2")

		err := db1.SaveToFile(filePath)
		assert.NoError(t, err, "Предварительное сохранение должно завершиться успешно")

		// Создаем новую БД и загружаем данные
		db2 := New()
		err = db2.LoadFromFile(filePath)
		assert.NoError(t, err, "Загрузка из файла не должна возвращать ошибку")

		// Проверяем что данные загрузились корректно
		value1, exists1 := db2.Get("short1")
		assert.True(t, exists1, "Первый ключ должен существовать после загрузки")
		assert.Equal(t, "https://example.com/1", value1, "Первое значение должно загрузиться корректно")

		value2, exists2 := db2.Get("short2")
		assert.True(t, exists2, "Второй ключ должен существовать после загрузки")
		assert.Equal(t, "https://example.com/2", value2, "Второе значение должно загрузиться корректно")

		// Проверяем пользователей
		userID1, _ := db2.GetUserID("short1")
		assert.Equal(t, "user-1", userID1, "Первый userID должен загрузиться корректно")

		userID2, _ := db2.GetUserID("short2")
		assert.Equal(t, "user-2", userID2, "Второй userID должен загрузиться корректно")

		// Проверяем счетчик
		assert.Equal(t, 2, db2.counter, "Счетчик должен восстановиться из файла")
	})

	// Тест 2: Загрузка из несуществующего файла
	t.Run("Load from non-existing file clears current data without error", func(t *testing.T) {
		db := New()
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "non_existing.json")

		err := db.LoadFromFile(filePath)

		assert.NoError(t, err, "Загрузка несуществующего файла не должна возвращать ошибку")
		assert.Empty(t, db.data, "Data должна быть пустой после загрузки несуществующего файла")
		assert.Empty(t, db.userMap, "UserMap должна быть пустой после загрузки несуществующего файла")
		assert.Equal(t, 0, db.counter, "Счетчик должен быть сброшен в 0")
	})
}

// TestDB_Integration тестирует полный рабочий цикл хранилища
func TestDB_Integration(t *testing.T) {
	t.Run("Complete storage workflow: add, query, save, load, verify", func(t *testing.T) {
		db := New()
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "integration.json")

		// Добавляем тестовые данные
		db.SetWithUser("abc123", "https://google.com", "user-1")
		db.SetWithUser("xyz789", "https://example.com", "user-2")
		db.SetWithUser("def456", "https://github.com", "user-1") // Второй URL для user-1

		// Проверяем получение URL пользователей
		user1URLs := db.GetUserURLs("user-1")
		assert.Len(t, user1URLs, 2, "User-1 должен иметь 2 URL")

		user2URLs := db.GetUserURLs("user-2")
		assert.Len(t, user2URLs, 1, "User-2 должен иметь 1 URL")

		// Сохраняем данные в файл
		err := db.SaveToFile(filePath)
		assert.NoError(t, err, "Сохранение должно завершиться успешно")

		// Создаем новую БД и загружаем данные
		newDB := New()
		err = newDB.LoadFromFile(filePath)
		assert.NoError(t, err, "Загрузка должна завершиться успешно")

		// Проверяем что все данные загрузились корректно
		newUser1URLs := newDB.GetUserURLs("user-1")
		assert.Len(t, newUser1URLs, 2, "После загрузки user-1 должен иметь 2 URL")

		newUser2URLs := newDB.GetUserURLs("user-2")
		assert.Len(t, newUser2URLs, 1, "После загрузки user-2 должен иметь 1 URL")

		// Проверяем конкретные значения
		value, exists := newDB.Get("abc123")
		assert.True(t, exists, "Ключ abc123 должен существовать после загрузки")
		assert.Equal(t, "https://google.com", value, "Значение abc123 должно загрузиться корректно")

		userID, userExists := newDB.GetUserID("abc123")
		assert.True(t, userExists, "UserID для abc123 должен существовать после загрузки")
		assert.Equal(t, "user-1", userID, "UserID для abc123 должен загрузиться корректно")
	})
}
