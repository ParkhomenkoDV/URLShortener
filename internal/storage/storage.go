package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var mutex sync.Mutex

type DB struct {
	data  map[string]string
	count int
}

// New - создание нового объекта БД.
func New() *DB {
	return &DB{
		data:  make(map[string]string),
		count: 0,
	}
}

// Get - получение значения по ключу.
func (db *DB) Get(key string) (string, bool) {
	mutex.Lock()
	value, exists := db.data[key]
	mutex.Unlock()

	return value, exists
}

// Set - установка значения по ключу.
func (db *DB) Set(key, value string) {
	mutex.Lock()
	db.data[key] = value
	db.count++
	mutex.Unlock()
}

func (db *DB) Delete(key string) error {
	_, exists := db.data[key]
	if !exists {
		return fmt.Errorf("key not found")
	}
	delete(db.data, key)
	db.count--
	return nil
}

// record - запись об URLs.
type record struct {
	ID          int
	ShortURL    string
	OriginalURL string
}

// SaveToFile - сохранение данных в JSON файл.
func (db *DB) SaveToFile(filePath string) error {
	var records []record

	counter := 0
	for shortURL, originalURL := range db.data {
		counter++
		record := record{counter, shortURL, originalURL}
		records = append(records, record)
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	// Создаем директорию, если её нет
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// LoadFromFile - загрузка данных из JSON файла.
func (db *DB) LoadFromFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) { // файла не существует
		return err
	}

	bytes, err := os.ReadFile(filePath) // ошибка чтения файла
	if err != nil {
		return err
	}

	if len(bytes) == 0 { // пустой файл
		return fmt.Errorf("empty file")
	}

	var records []record

	if err := json.Unmarshal(bytes, &records); err != nil { // ошибка десериализации JSON
		return err
	}

	data, counter := make(map[string]string), 0
	for _, record := range records {
		data[record.ShortURL] = record.OriginalURL
		if record.ID > counter {
			counter = record.ID
		}
	}

	db.data = data
	db.count = counter

	return nil
}
