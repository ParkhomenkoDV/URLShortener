package repository

import (
	"errors"
	"sync"

	"github.com/ParkhomenkoDV/URLShortener/internal/persistence"
)

// FileRepository реализация репозитория для хранения в файле
type FileRepository struct {
	data         map[string]string
	reversedData map[string]string
	userMap      map[string]string
	mu           sync.RWMutex
	filePath     string
	persistence  persistence.JSONPersistence
}

// NewFile создает новый репозиторий для работы с файлом
func NewFile(filePath string) URLRepository {
	repo := &FileRepository{
		data:         make(map[string]string),
		reversedData: make(map[string]string),
		userMap:      make(map[string]string),
		filePath:     filePath,
		persistence:  persistence.NewFileJSONPersistence(),
	}

	// Загружаем данные из файла при инициализации
	data, userMap, _, err := repo.persistence.Load(filePath)
	if err == nil {
		repo.data = data
		repo.userMap = userMap
	}

	// Формирование обратной мапы TODO: вынести в pkg/functions
	for originalURL, shortURL := range data {
		repo.reversedData[originalURL] = shortURL
	}

	return repo
}

// GetLongValue получает оригинальный URL по короткому
func (r *FileRepository) GetLongValue(shortURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if value, ok := r.data[shortURL]; ok {
		return value, nil
	}
	return "", errors.New("not found key in database")
}

// GetShortValue получает короткий URL по оригинальному
func (r *FileRepository) GetShortValue(originalURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if value, ok := r.reversedData[originalURL]; ok {
		return value, nil
	}
	return "", errors.New("not found key in database")
}

// SetValue сохраняет пару короткий URL - оригинальный URL
func (r *FileRepository) SetValue(shortURL, originalURL, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[shortURL]; ok {
		return ErrRowExists
	}
	r.data[shortURL] = originalURL
	r.userMap[shortURL] = userID

	// Сохраняем в файл
	return r.persistence.Save(r.filePath, r.data, r.userMap)
}

// SetValuesBatch сохраняет пакет пар короткий URL - оригинальный URL
func (r *FileRepository) SetValuesBatch(pairs map[string]string, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for key, value := range pairs {
		if _, ok := r.data[key]; ok {
			return ErrRowExists
		}
		r.data[key] = value
		r.userMap[key] = userID
	}

	// Сохраняем в файл
	return r.persistence.Save(r.filePath, r.data, r.userMap)
}

// Close закрывает соединение с хранилищем
func (r *FileRepository) Close() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Финальное сохранение в файл
	return r.persistence.Save(r.filePath, r.data, r.userMap)
}

// GetUserURLs получает все URL пользователя
func (r *FileRepository) GetUserURLs(userID string) ([]map[string]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var urls []map[string]string
	for shortURL, originalURL := range r.data {
		if userID == r.userMap[shortURL] {
			urls = append(urls, map[string]string{
				"short_url":    shortURL,
				"original_url": originalURL,
			})
		}
	}
	return urls, nil
}
