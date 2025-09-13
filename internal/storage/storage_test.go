package storage

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		db := New()

		// Test Set and Get
		db.Set("key1", "value1")
		value, exists := db.Get("key1")
		require.True(t, exists)
		require.Equal(t, "value1", value)

		// Test non-existent key
		_, exists = db.Get("nonexistent")
		require.False(t, exists)
	})

	t.Run("delete operations", func(t *testing.T) {
		db := New()
		db.Set("key1", "value1")

		// Test Delete existing key
		err := db.Delete("key1")
		require.NoError(t, err)
		_, exists := db.Get("key1")
		require.False(t, exists)

		// Test Delete non-existent key
		err = db.Delete("nonexistent")
		require.Error(t, err)
		require.Equal(t, "key not found", err.Error())
	})

	t.Run("count tracking", func(t *testing.T) {
		db := New()

		// Initial count should be 0
		require.Equal(t, 0, db.Count())

		// Count should increase on Set
		db.Set("key1", "value1")
		require.Equal(t, 1, db.Count())

		db.Set("key2", "value2")
		require.Equal(t, 2, db.Count())

		// Count should decrease on Delete
		db.Delete("key1")
		require.Equal(t, 1, db.Count())

		db.Delete("key2")
		require.Equal(t, 0, db.Count())
	})
}

func TestDBConcurrent(t *testing.T) {
	t.Run("concurrent set operations", func(t *testing.T) {
		db := New()
		var wg sync.WaitGroup
		iterations := 100

		wg.Add(iterations)
		for i := 0; i < iterations; i++ {
			go func(index int) {
				defer wg.Done()
				key := formatKey(index)
				db.Set(key, formatValue(index))
			}(i)
		}
		wg.Wait()

		// Verify all values were set correctly
		for i := 0; i < iterations; i++ {
			value, exists := db.Get(formatKey(i))
			require.True(t, exists)
			require.Equal(t, formatValue(i), value)
		}
	})

	t.Run("concurrent get and set operations", func(t *testing.T) {
		db := New()
		var wg sync.WaitGroup
		iterations := 50

		// Start with some initial data
		for i := 0; i < iterations; i++ {
			db.Set(formatKey(i), formatValue(i))
		}

		wg.Add(iterations * 2)
		for i := 0; i < iterations; i++ {
			// Concurrent gets
			go func(index int) {
				defer wg.Done()
				db.Get(formatKey(index))
			}(i)

			// Concurrent sets (updating values)
			go func(index int) {
				defer wg.Done()
				db.Set(formatKey(index), formatValue(index*2))
			}(i)
		}
		wg.Wait()

		// Verify final values
		for i := 0; i < iterations; i++ {
			value, exists := db.Get(formatKey(i))
			require.True(t, exists)
			require.Equal(t, formatValue(i*2), value)
		}
	})

	t.Run("concurrent delete operations", func(t *testing.T) {
		db := New()
		var wg sync.WaitGroup
		iterations := 100

		// Set up initial data
		for i := 0; i < iterations; i++ {
			db.Set(formatKey(i), formatValue(i))
		}

		wg.Add(iterations)
		for i := 0; i < iterations; i++ {
			go func(index int) {
				defer wg.Done()
				db.Delete(formatKey(index))
			}(i)
		}
		wg.Wait()

		// All keys should be deleted
		for i := 0; i < iterations; i++ {
			_, exists := db.Get(formatKey(i))
			require.False(t, exists)
		}
		require.Equal(t, 0, db.Count())
	})

	t.Run("concurrent set and delete same key", func(t *testing.T) {
		db := New()
		key := "key"
		iterations := 100

		var wg sync.WaitGroup
		wg.Add(iterations * 2)

		for i := 0; i < iterations; i++ {
			go func() {
				defer wg.Done()
				db.Set(key, "value")
			}()

			go func() {
				defer wg.Done()
				db.Delete(key)
			}()
		}

		wg.Wait()

		_, exists := db.Get(key)
		require.False(t, exists)
	})
}

func TestFileOperations(t *testing.T) {
	t.Run("save and load from file", func(t *testing.T) {
		tempFile := "/tmp/test_db.json"
		defer os.Remove(tempFile)

		db1 := New()
		db1.Set("key1", "https://example.com")
		db1.Set("key2", "https://google.com")
		db1.Set("key3", "https://github.com")

		err := db1.SaveToFile(tempFile)
		require.NoError(t, err)

		db2 := New()
		err = db2.LoadFromFile(tempFile)
		require.NoError(t, err)

		value, exists := db2.Get("key1")
		require.True(t, exists)
		require.Equal(t, "https://example.com", value)

		value, exists = db2.Get("key2")
		require.True(t, exists)
		require.Equal(t, "https://google.com", value)

		value, exists = db2.Get("key3")
		require.True(t, exists)
		require.Equal(t, "https://github.com", value)
	})

	t.Run("save empty db to file", func(t *testing.T) {
		tempFile := "/tmp/empty_db.json"
		defer os.Remove(tempFile)

		db := New()
		err := db.SaveToFile(tempFile)
		require.NoError(t, err)

		// Verify file was created
		_, err = os.Stat(tempFile)
		require.NoError(t, err)
	})

	t.Run("load from non-existent file", func(t *testing.T) {
		db := New()
		err := db.LoadFromFile("/nonexistent/file.json")
		require.Error(t, err)
	})

	t.Run("load from empty file", func(t *testing.T) {
		tempFile := "/tmp/empty_file.json"
		defer os.Remove(tempFile)

		// Create empty file
		err := os.WriteFile(tempFile, []byte{}, 0644)
		require.NoError(t, err)

		db := New()
		err = db.LoadFromFile(tempFile)
		require.Error(t, err)
		require.Equal(t, "empty file", err.Error())
	})

	t.Run("load from invalid json file", func(t *testing.T) {
		tempFile := "/tmp/invalid_json.json"
		defer os.Remove(tempFile)

		// Create invalid JSON file
		err := os.WriteFile(tempFile, []byte("{invalid json}"), 0644)
		require.NoError(t, err)

		db := New()
		err = db.LoadFromFile(tempFile)
		require.Error(t, err)
	})
}

func TestFileConcurrent(t *testing.T) {
	t.Run("concurrent save operations", func(t *testing.T) {
		tempFile1 := "/tmp/concurrent1.json"
		tempFile2 := "/tmp/concurrent2.json"
		defer os.Remove(tempFile1)
		defer os.Remove(tempFile2)

		db := New()
		db.Set("key1", "value1")
		db.Set("key2", "value2")

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			err := db.SaveToFile(tempFile1)
			require.NoError(t, err)
		}()

		go func() {
			defer wg.Done()
			err := db.SaveToFile(tempFile2)
			require.NoError(t, err)
		}()

		wg.Wait()

		// Both files should be created
		_, err := os.Stat(tempFile1)
		require.NoError(t, err)
		_, err = os.Stat(tempFile2)
		require.NoError(t, err)
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("overwrite existing key", func(t *testing.T) {
		db := New()
		db.Set("key1", "value1")
		db.Set("key1", "value2") // Overwrite!

		value, exists := db.Get("key1")
		require.True(t, exists)
		require.Equal(t, "value2", value)
		require.Equal(t, 1, db.Count()) // Count should not increase
	})

	t.Run("empty key", func(t *testing.T) {
		db := New()
		db.Set("", "empty key value")

		value, exists := db.Get("")
		require.True(t, exists)
		require.Equal(t, "empty key value", value)
	})

	t.Run("special characters in key", func(t *testing.T) {
		db := New()
		specialKey := "key-with-special-chars!@#$%^&*()"
		db.Set(specialKey, "special value")

		value, exists := db.Get(specialKey)
		require.True(t, exists)
		require.Equal(t, "special value", value)
	})
}

func formatKey(index int) string {
	return string(rune('a' + index%26))
}

func formatValue(index int) string {
	return string(rune('A' + index%26))
}
