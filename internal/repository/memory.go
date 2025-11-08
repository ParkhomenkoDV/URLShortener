package repository

import (
	"errors"
	"sync"
)

// MemoryRepository реализация репозитория для хранения в памяти
type MemoryRepository struct {
	data       map[string]string // map[short] = long
	userMap    map[string]string // map[short] = ID
	deletedMap map[string]bool   // shortURL -> isDeleted
	mu         sync.RWMutex
}

// NewMemory создает новый репозиторий для работы с памятью
func NewMemory() URLRepository {
	return &MemoryRepository{
		data:       make(map[string]string),
		userMap:    make(map[string]string),
		deletedMap: make(map[string]bool),
	}
}

// withReadLock выполняет функцию под блокировкой чтения
func (mr *MemoryRepository) withReadLock(fn func() error) error {
	mr.mu.RLock()
	defer mr.mu.RUnlock()
	return fn()
}

// withWriteLock выполняет функцию под блокировкой записи
func (mr *MemoryRepository) withWriteLock(fn func() error) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	return fn()
}

// checkURLExists проверяет существование и статус удаления URL (должен вызываться под блокировкой)
func (mr *MemoryRepository) checkURLExists(shortURL string) (string, error) {
	// Проверяем, существует ли URL
	value, exists := mr.data[shortURL]
	if !exists {
		return "", errors.New("not found key in database")
	}

	// Проверяем, не удален ли URL
	if deleted, exists := mr.deletedMap[shortURL]; exists && deleted {
		return "", errors.New("not found key in database")
	}

	return value, nil
}

// GetValue получает оригинальный URL по короткому
func (mr *MemoryRepository) GetLongValue(shortURL string) (string, error) {
	var result string
	var err error

	mr.withReadLock(func() error {
		result, err = mr.checkURLExists(shortURL)
		return err
	})

	return result, err
}

// GetShortValue получает короткий URL по оригинальному
func (mr *MemoryRepository) GetShortValue(originalURL string) (string, error) {
	var result string
	var err error

	mr.withReadLock(func() error {
		for short, long := range mr.data {
			if long == originalURL {
				result = short
				return nil
			}
		}
		err = errors.New("not found key in database")
		return err
	})

	return result, err
}

// SetValue сохраняет пару короткий URL - оригинальный URL
func (mr *MemoryRepository) SetValue(shortURL, originalURL, userID string) error {
	return mr.withWriteLock(func() error {
		if _, ok := mr.data[shortURL]; ok {
			return ErrRowExists
		}
		mr.data[shortURL] = originalURL
		mr.userMap[shortURL] = userID
		return nil
	})
}

// SetValuesBatch сохраняет пакет пар короткий URL - оригинальный URL
func (mr *MemoryRepository) SetValuesBatch(pairs map[string]string, userID string) error {
	return mr.withWriteLock(func() error {
		for key, value := range pairs {
			if _, ok := mr.data[key]; ok {
				return ErrRowExists
			}
			mr.data[key] = value
			mr.userMap[key] = userID
		}
		return nil
	})
}

// Close закрывает соединение с хранилищем (для памяти это заглушка)
func (mr *MemoryRepository) Close() error {
	return nil
}

// GetUserURLs получает все не удаленные URL пользователя
func (mr *MemoryRepository) GetUserURLs(userID string) ([]map[string]string, error) {
	var urls []map[string]string

	mr.withReadLock(func() error {
		for shortURL, originalURL := range mr.data {
			// Проверяем, что URL принадлежит пользователю и не удален
			if userID == mr.userMap[shortURL] {
				if deleted, exists := mr.deletedMap[shortURL]; !exists || !deleted {
					urls = append(urls, map[string]string{
						"short_url":    shortURL,
						"original_url": originalURL,
					})
				}
			}
		}
		return nil
	})

	return urls, nil
}

// DeleteURLsBatch помечает множественные URL как удаленные для указанного пользователя
func (mr *MemoryRepository) DeleteURLsBatch(shortURLs []string, userID string) error {
	return mr.withWriteLock(func() error {
		for _, shortURL := range shortURLs {
			// Проверяем, принадлежит ли URL пользователю
			if ownerID, exists := mr.userMap[shortURL]; exists && ownerID == userID {
				mr.deletedMap[shortURL] = true
			}
		}
		return nil
	})
}

// IsDeleted проверяет, помечен ли URL как удаленный
func (mr *MemoryRepository) IsDeleted(shortURL string) (bool, error) {
	var result bool
	var err error

	mr.withReadLock(func() error {
		// Проверяем, существует ли URL вообще
		if _, exists := mr.data[shortURL]; !exists {
			err = errors.New("not found key in database")
			return err
		}

		// Проверяем статус удаления
		if deleted, exists := mr.deletedMap[shortURL]; exists {
			result = deleted
		} else {
			result = false
		}
		return nil
	})

	return result, err
}
