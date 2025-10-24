package storage

import (
	"github.com/ParkhomenkoDV/URLShortener/internal/persistence"
)

// Структура для БД
type DB struct {
	data        map[string]string
	userMap     map[string]string // map[short] = long
	counter     int
	persistence persistence.JSONPersistence
}

// Конструктор БД
func New() *DB {
	return &DB{
		data:        make(map[string]string),
		userMap:     make(map[string]string),
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
	db.SetWithUser(key, value, "")
}

// Установка значения с указанием пользователя
func (db *DB) SetWithUser(key, value, userID string) {
	db.data[key] = value
	db.userMap[key] = userID
	db.counter++
}

// Сохраняет данные в файл в формате JSON
func (db *DB) SaveToFile(filePath string) error {
	return db.persistence.Save(filePath, db.data, db.userMap)
}

// Загружает данные из файла в формате JSON
func (db *DB) LoadFromFile(filePath string) error {
	data, userMap, maxCounter, err := db.persistence.Load(filePath)
	if err != nil {
		return err
	}

	db.data = data
	db.userMap = userMap
	db.counter = maxCounter

	return nil
}

// GetUserURLs возвращает все URL конкретного пользователя
func (db *DB) GetUserURLs(userID string) []map[string]string {
	var urls []map[string]string
	for shortURL, originalURL := range db.data {
		if db.userMap[shortURL] == userID {
			urls = append(urls, map[string]string{
				"short_url":    shortURL,
				"original_url": originalURL,
			})
		}
	}
	return urls
}

// GetUserID возвращает ID пользователя по короткому URL
func (db *DB) GetUserID(shortURL string) (string, bool) {
	userID, exists := db.userMap[shortURL]
	return userID, exists
}
