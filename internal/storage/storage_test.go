package storage

import (
	"os"
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

		// Test Delete
		err := db.Delete("key1")
		require.NoError(t, err)
		_, exists = db.Get("key1")
		require.False(t, exists)

		// Test Delete non-existent key
		err = db.Delete("nonexistent")
		require.Error(t, err)
	})

	t.Run("save and load from file", func(t *testing.T) {
		tempFile := "/tmp/test_db.json"
		defer os.Remove(tempFile)

		db1 := New()
		db1.Set("key1", "https://example.com")
		db1.Set("key2", "https://google.com")

		// Test SaveToFile
		err := db1.SaveToFile(tempFile)
		require.NoError(t, err)

		// Test LoadFromFile
		db2 := New()
		err = db2.LoadFromFile(tempFile)
		require.NoError(t, err)

		value, exists := db2.Get("key1")
		require.True(t, exists)
		require.Equal(t, "https://example.com", value)

		value, exists = db2.Get("key2")
		require.True(t, exists)
		require.Equal(t, "https://google.com", value)
	})

	t.Run("load from non-existent file", func(t *testing.T) {
		db := New()
		err := db.LoadFromFile("/nonexistent/file.json")
		require.Error(t, err)
	})

	t.Run("load from empty file", func(t *testing.T) {
		tempFile := "/tmp/empty_test.json"
		defer os.Remove(tempFile)

		// Create empty file
		err := os.WriteFile(tempFile, []byte{}, 0644)
		require.NoError(t, err)

		db := New()
		err = db.LoadFromFile(tempFile)
		require.Error(t, err)
	})
}
