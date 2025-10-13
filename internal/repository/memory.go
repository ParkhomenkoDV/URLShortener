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

// NewMemoryRepository создает новый репозиторий для работы с памятью
func NewMemoryRepository() URLRepository {
	return &MemoryRepository{
		data: make(map[string]string),
	}
}

// GetValue получает оригинальный URL по короткому
func (r *MemoryRepository) GetFullValue(shortURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if value, ok := r.data[shortURL]; ok {
		return value, nil
	}
	return "", errors.New("not found key in database")
}

// GetShortValue получает короткий URL по оригинальному
func (r *MemoryRepository) GetShortValue(originalURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for short, long := range r.data {
		if long == originalURL {
			return short, nil
		}
	}
	return "", errors.New("not found key in database")
}

// SetValue сохраняет пару короткий URL - оригинальный URL
func (r *MemoryRepository) SetValue(shortURL, originalURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[shortURL]; ok {
		return ErrRowExists
	}
	r.data[shortURL] = originalURL
	return nil
}

// SetValuesBatch сохраняет пакет пар короткий URL - оригинальный URL
func (r *MemoryRepository) SetValuesBatch(pairs map[string]string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for key, value := range pairs {
		if _, ok := r.data[key]; ok {
			return ErrRowExists
		}
		r.data[key] = value
	}
	return nil
}

// Close закрывает соединение с хранилищем (для памяти это заглушка)
func (r *MemoryRepository) Close() error {
	return nil
}
