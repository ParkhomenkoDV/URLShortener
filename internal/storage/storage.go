package storage

import (
	"github.com/ParkhomenkoDV/URLShortener/internal/persistence"
)

// DB представляет in-memory хранилище для системы сокращения URL.
// Хранит маппинги коротких URL на оригинальные URL и связывает их с пользователями.
// Поддерживает сохранение и загрузку данных в/из JSON файлов через persistence слой.
type DB struct {
	data        map[string]string           // Маппинг коротких URL на оригинальные URL: map[shortURL]originalURL
	userMap     map[string]string           // Маппинг коротких URL на идентификаторы пользователей: map[shortURL]userID
	counter     int                         // Счетчик для генерации ID записей (используется при сохранении в JSON)
	persistence persistence.JSONPersistence // Интерфейс для сериализации/десериализации данных
}

// New создает и инициализирует новый экземпляр in-memory базы данных.
// Возвращает готовую к использованию БД с инициализированными мапами и FileJSONPersistence.
func New() *DB {
	return &DB{
		data:        make(map[string]string),
		userMap:     make(map[string]string),
		counter:     0,
		persistence: persistence.NewFileJSONPersistence(), // Используем файловую реализацию persistence
	}
}

// Get возвращает оригинальный URL по короткому идентификатору.
// Параметры:
//   - key: короткий идентификатор URL
//
// Возвращает:
//   - string: оригинальный URL (пустая строка если ключ не найден)
//   - bool: флаг существования ключа в хранилище
func (db *DB) Get(key string) (string, bool) {
	value, exists := db.data[key]
	return value, exists
}

// Set сохраняет пару короткий-оригинальный URL без привязки к пользователю.
// Увеличивает внутренний счетчик записей. Используется для анонимных сокращений.
// Параметры:
//   - key: короткий идентификатор URL
//   - value: оригинальный URL
func (db *DB) Set(key, value string) {
	db.SetWithUser(key, value, "")
}

// SetWithUser сохраняет пару короткий-оригинальный URL с привязкой к пользователю.
// Увеличивает внутренний счетчик записей. Используется для авторизованных пользователей.
// Параметры:
//   - key: короткий идентификатор URL
//   - value: оригинальный URL
//   - userID: идентификатор пользователя-владельца ссылки
func (db *DB) SetWithUser(key, value, userID string) {
	db.data[key] = value
	db.userMap[key] = userID
	db.counter++
}

// SaveToFile сохраняет все данные хранилища в JSON файл через persistence слой.
// Автоматически создает необходимые директории если они не существуют.
// Параметры:
//   - filePath: путь к файлу для сохранения
//
// Возвращает:
//   - error: ошибка если операция сохранения не удалась
func (db *DB) SaveToFile(filePath string) error {
	return db.persistence.Save(filePath, db.data, db.userMap)
}

// LoadFromFile загружает данные из JSON файла в хранилище через persistence слой.
// Заменяет текущие данные в памяти на загруженные из файла.
// Если файл не существует, очищает текущие данные.
// Параметры:
//   - filePath: путь к файлу для загрузки
//
// Возвращает:
//   - error: ошибка если операция загрузки не удалась (кроме случая несуществующего файла)
func (db *DB) LoadFromFile(filePath string) error {
	data, userMap, maxCounter, err := db.persistence.Load(filePath)
	if err != nil {
		return err
	}

	db.data = data
	db.userMap = userMap
	db.counter = maxCounter // Восстанавливаем счетчик для продолжения нумерации

	return nil
}

// GetUserURLs возвращает все сокращенные URL указанного пользователя.
// Возвращает мапу с информацией о каждом URL: короткий URL и оригинальный URL.
// Параметры:
//   - userID: идентификатор пользователя
//
// Возвращает:
//   - []map[string]string: слайс мап, где каждая мапа содержит:
//   - "short_url": короткий идентификатор
//   - "original_url": оригинальный URL
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

// GetUserID возвращает идентификатор пользователя по короткому URL.
// Параметры:
//   - shortURL: короткий идентификатор URL
//
// Возвращает:
//   - string: идентификатор пользователя (пустая строка если ключ не найден)
//   - bool: флаг существования записи
func (db *DB) GetUserID(shortURL string) (string, bool) {
	userID, exists := db.userMap[shortURL]
	return userID, exists
}
