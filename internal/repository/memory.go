package repository

import (
	"errors"
	"sync"
)

// MemoryRepository реализация репозитория для хранения в памяти
type MemoryRepository struct {
	data    map[string]string
	userMap map[string]string // map[short] = long
	mu      sync.RWMutex
}

// NewMemory создает новый репозиторий для работы с памятью
func NewMemory() URLRepository {
	return &MemoryRepository{
		data:    make(map[string]string),
		userMap: make(map[string]string),
	}
}

// GetValue получает оригинальный URL по короткому
func (mr *MemoryRepository) GetLongValue(shortURL string) (string, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	if value, ok := mr.data[shortURL]; ok {
		return value, nil
	}
	return "", errors.New("not found key in database")
}

// GetShortValue получает короткий URL по оригинальному
func (mr *MemoryRepository) GetShortValue(originalURL string) (string, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	for short, long := range mr.data {
		if long == originalURL {
			return short, nil
		}
	}
	return "", errors.New("not found key in database")
}

// SetValue сохраняет пару короткий URL - оригинальный URL
func (mr *MemoryRepository) SetValue(shortURL, originalURL, userID string) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	if _, ok := mr.data[shortURL]; ok {
		return ErrRowExists
	}
	mr.data[shortURL] = originalURL
	mr.userMap[shortURL] = userID
	return nil
}

// SetValuesBatch сохраняет пакет пар короткий URL - оригинальный URL
func (mr *MemoryRepository) SetValuesBatch(pairs map[string]string, userID string) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	for key, value := range pairs {
		if _, ok := mr.data[key]; ok {
			return ErrRowExists
		}
		mr.data[key] = value
		mr.userMap[key] = userID
	}
	return nil
}

// Close закрывает соединение с хранилищем (для памяти это заглушка)
func (mr *MemoryRepository) Close() error {
	return nil
}

// GetUserURLs получает все URL пользователя
func (mr *MemoryRepository) GetUserURLs(userID string) ([]map[string]string, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	var urls []map[string]string
	for shortURL, originalURL := range mr.data {
		if userID == mr.userMap[shortURL] {
			urls = append(urls, map[string]string{
				"short_url":    shortURL,
				"original_url": originalURL,
			})
		}
	}
	return urls, nil
}
