package repository

import (
	"errors"
	"sync"
)

// MemoryRepository реализация репозитория для хранения в памяти
type MemoryRepository struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewMemory создает новый репозиторий для работы с памятью
func NewMemory() URLRepository {
	return &MemoryRepository{
		data: make(map[string]string),
	}
}

// GetValue получает оригинальный URL по короткому
func (mr *MemoryRepository) GetFullValue(shortURL string) (string, error) {
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
func (mr *MemoryRepository) SetValue(shortURL, originalURL string) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	if _, ok := mr.data[shortURL]; ok {
		return ErrRowExists
	}
	mr.data[shortURL] = originalURL
	return nil
}

// SetValuesBatch сохраняет пакет пар короткий URL - оригинальный URL
func (mr *MemoryRepository) SetValuesBatch(pairs map[string]string) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	for key, value := range pairs {
		if _, ok := mr.data[key]; ok {
			return ErrRowExists
		}
		mr.data[key] = value
	}
	return nil
}

// Close закрывает соединение с хранилищем (для памяти это заглушка)
func (mr *MemoryRepository) Close() error {
	return nil
}
