package storage

import (
	"github.com/ParkhomenkoDV/URLShortener/internal/persistence"
)

// Структура для БД
type DB struct {
	data        map[string]string
	counter     int
	persistence persistence.JSONPersistence
}

// Создание БД
func CreateDB() *DB {
	return &DB{
		data:        make(map[string]string),
		counter:     0,
		persistence: persistence.NewFileJSONPersistence(),
	}
}

// Получение значения по ключу
func (db *DB) Get(key string) (string, bool) {
	value, exists := db.data[key]
	return value, exists
}

// Установка значения
func (db *DB) Set(key, value string) {
	db.data[key] = value
	db.counter++
}

// Сохраняет данные в файл в формате JSON
func (db *DB) SaveToFile(filePath string) error {
	return db.persistence.Save(filePath, db.data)
}

// Загружает данные из файла в формате JSON
func (db *DB) LoadFromFile(filePath string) error {
	data, maxCounter, err := db.persistence.Load(filePath)
	if err != nil {
		return err
	}

	db.data = data
	db.counter = maxCounter

	return nil
}
