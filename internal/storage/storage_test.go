package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDB(t *testing.T) {
	t.Run("Create new database", func(t *testing.T) {
		db := New()

		assert.NotNil(t, db)
		assert.NotNil(t, db.data)
		assert.NotNil(t, db.userMap)
		assert.Equal(t, 0, db.counter)
		assert.NotNil(t, db.persistence)
		assert.Empty(t, db.data)
		assert.Empty(t, db.userMap)
	})
}

func TestDB_Get(t *testing.T) {
	t.Run("Get existing key", func(t *testing.T) {
		db := New()
		db.data["test-key"] = "test-value"

		value, exists := db.Get("test-key")

		assert.True(t, exists)
		assert.Equal(t, "test-value", value)
	})

	t.Run("Get non-existing key", func(t *testing.T) {
		db := New()

		value, exists := db.Get("non-existing")

		assert.False(t, exists)
		assert.Equal(t, "", value)
	})
}

func TestDB_Set(t *testing.T) {
	t.Run("Set value without user", func(t *testing.T) {
		db := New()

		db.Set("key1", "value1")

		value, exists := db.Get("key1")
		assert.True(t, exists)
		assert.Equal(t, "value1", value)
		assert.Equal(t, 1, db.counter)

		// Проверяем, что userID пустой
		userID, userExists := db.GetUserID("key1")
		assert.True(t, userExists)
		assert.Equal(t, "", userID)
	})
}

func TestDB_SetWithUser(t *testing.T) {
	t.Run("Set value with user", func(t *testing.T) {
		db := New()

		db.SetWithUser("key1", "value1", "user-123")

		value, exists := db.Get("key1")
		assert.True(t, exists)
		assert.Equal(t, "value1", value)
		assert.Equal(t, 1, db.counter)

		userID, userExists := db.GetUserID("key1")
		assert.True(t, userExists)
		assert.Equal(t, "user-123", userID)
	})

	t.Run("Set multiple values with different users", func(t *testing.T) {
		db := New()

		db.SetWithUser("key1", "value1", "user-1")
		db.SetWithUser("key2", "value2", "user-2")
		db.SetWithUser("key3", "value3", "user-1")

		assert.Equal(t, 3, db.counter)

		// Проверяем значения
		value1, _ := db.Get("key1")
		assert.Equal(t, "value1", value1)

		value2, _ := db.Get("key2")
		assert.Equal(t, "value2", value2)

		// Проверяем пользователей
		userID1, _ := db.GetUserID("key1")
		assert.Equal(t, "user-1", userID1)

		userID2, _ := db.GetUserID("key2")
		assert.Equal(t, "user-2", userID2)
	})
}

func TestDB_GetUserID(t *testing.T) {
	t.Run("Get existing user ID", func(t *testing.T) {
		db := New()
		db.SetWithUser("key1", "value1", "user-123")

		userID, exists := db.GetUserID("key1")

		assert.True(t, exists)
		assert.Equal(t, "user-123", userID)
	})

	t.Run("Get non-existing user ID", func(t *testing.T) {
		db := New()

		userID, exists := db.GetUserID("non-existing")

		assert.False(t, exists)
		assert.Equal(t, "", userID)
	})
}

func TestDB_GetUserURLs(t *testing.T) {
	t.Run("Get URLs for existing user", func(t *testing.T) {
		db := New()

		db.SetWithUser("short1", "https://example.com/1", "user-1")
		db.SetWithUser("short2", "https://example.com/2", "user-2")
		db.SetWithUser("short3", "https://example.com/3", "user-1")

		urls := db.GetUserURLs("user-1")

		assert.Len(t, urls, 2)

		// Проверяем, что получили правильные URL
		foundShort1 := false
		foundShort3 := false

		for _, urlMap := range urls {
			shortURL := urlMap["short_url"]
			switch shortURL {
			case "short1":
				foundShort1 = true
				assert.Equal(t, "https://example.com/1", urlMap["original_url"])
			case "short3":
				foundShort3 = true
				assert.Equal(t, "https://example.com/3", urlMap["original_url"])
			}
		}

		assert.True(t, foundShort1)
		assert.True(t, foundShort3)
	})

	t.Run("Get URLs for non-existing user", func(t *testing.T) {
		db := New()

		db.SetWithUser("short1", "https://example.com/1", "user-1")

		urls := db.GetUserURLs("non-existing-user")

		assert.Empty(t, urls)
	})

	t.Run("Get URLs when no data", func(t *testing.T) {
		db := New()

		urls := db.GetUserURLs("any-user")

		assert.Empty(t, urls)
	})
}

func TestDB_SaveToFile(t *testing.T) {
	t.Run("Save data to file", func(t *testing.T) {
		db := New()
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test_save.json")

		db.SetWithUser("short1", "https://example.com/1", "user-1")
		db.SetWithUser("short2", "https://example.com/2", "user-2")

		err := db.SaveToFile(filePath)

		assert.NoError(t, err)

		// Проверяем, что файл создался
		_, err = os.Stat(filePath)
		assert.NoError(t, err)

		// Проверяем содержимое
		content, err := os.ReadFile(filePath)
		assert.NoError(t, err)
		assert.Contains(t, string(content), "short1")
		assert.Contains(t, string(content), "https://example.com/1")
		assert.Contains(t, string(content), "user-1")
	})
}

func TestDB_LoadFromFile(t *testing.T) {
	t.Run("Load data from file", func(t *testing.T) {
		// Создаем первую БД и сохраняем данные
		db1 := New()
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "test_load.json")

		db1.SetWithUser("short1", "https://example.com/1", "user-1")
		db1.SetWithUser("short2", "https://example.com/2", "user-2")

		err := db1.SaveToFile(filePath)
		assert.NoError(t, err)

		// Создаем новую БД и загружаем данные
		db2 := New()
		err = db2.LoadFromFile(filePath)
		assert.NoError(t, err)

		// Проверяем, что данные загрузились
		value1, exists1 := db2.Get("short1")
		assert.True(t, exists1)
		assert.Equal(t, "https://example.com/1", value1)

		value2, exists2 := db2.Get("short2")
		assert.True(t, exists2)
		assert.Equal(t, "https://example.com/2", value2)

		// Проверяем пользователей
		userID1, _ := db2.GetUserID("short1")
		assert.Equal(t, "user-1", userID1)

		userID2, _ := db2.GetUserID("short2")
		assert.Equal(t, "user-2", userID2)

		// Проверяем счетчик
		assert.Equal(t, 2, db2.counter)
	})

	t.Run("Load from non-existing file", func(t *testing.T) {
		db := New()
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "non_existing.json")

		err := db.LoadFromFile(filePath)

		assert.NoError(t, err)
		assert.Empty(t, db.data)
		assert.Empty(t, db.userMap)
		assert.Equal(t, 0, db.counter)
	})
}

func TestDB_Integration(t *testing.T) {
	t.Run("Full workflow integration test", func(t *testing.T) {
		db := New()
		tempDir := t.TempDir()
		filePath := filepath.Join(tempDir, "integration.json")

		// Добавляем данные
		db.SetWithUser("abc123", "https://google.com", "user-1")
		db.SetWithUser("xyz789", "https://example.com", "user-2")
		db.SetWithUser("def456", "https://github.com", "user-1")

		// Проверяем получение URL пользователя
		user1URLs := db.GetUserURLs("user-1")
		assert.Len(t, user1URLs, 2)

		user2URLs := db.GetUserURLs("user-2")
		assert.Len(t, user2URLs, 1)

		// Сохраняем в файл
		err := db.SaveToFile(filePath)
		assert.NoError(t, err)

		// Создаем новую БД и загружаем
		newDB := New()
		err = newDB.LoadFromFile(filePath)
		assert.NoError(t, err)

		// Проверяем, что все данные загрузились корректно
		newUser1URLs := newDB.GetUserURLs("user-1")
		assert.Len(t, newUser1URLs, 2)

		newUser2URLs := newDB.GetUserURLs("user-2")
		assert.Len(t, newUser2URLs, 1)

		// Проверяем конкретные значения
		value, exists := newDB.Get("abc123")
		assert.True(t, exists)
		assert.Equal(t, "https://google.com", value)

		userID, userExists := newDB.GetUserID("abc123")
		assert.True(t, userExists)
		assert.Equal(t, "user-1", userID)
	})
}
