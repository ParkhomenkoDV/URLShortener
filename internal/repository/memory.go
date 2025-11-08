package repository

import (
	"errors"
	"sync"
)

// MemoryRepository реализует хранение данных в оперативной памяти
type MemoryRepository struct {
	data       map[string]string // shortURL -> originalURL
	userMap    map[string]string // shortURL -> userID
	deletedMap map[string]bool   // shortURL -> isDeleted
	mu         sync.RWMutex      // Мьютекс для безопасного доступа из горутин
}

// NewMemory создает новый репозиторий в памяти
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

// checkURLExists проверяет существование URL и его статус удаления
// Должен вызываться только под блокировкой чтения или записи
func (mr *MemoryRepository) checkURLExists(shortURL string) (string, error) {
	// Проверяем существование URL
	value, exists := mr.data[shortURL]
	if !exists {
		return "", errors.New("not found key in database")
	}

	// Проверяем статус удаления
	if deleted, exists := mr.deletedMap[shortURL]; exists && deleted {
		return "", errors.New("not found key in database")
	}

	return value, nil
}

// GetLongValue возвращает оригинальный URL по короткому ключу
func (mr *MemoryRepository) GetLongValue(shortURL string) (string, error) {
	var result string
	var err error

	mr.withReadLock(func() error {
		result, err = mr.checkURLExists(shortURL)
		return err
	})

	return result, err
}

// GetShortValue возвращает короткий ключ по оригинальному URL
func (mr *MemoryRepository) GetShortValue(originalURL string) (string, error) {
	var result string
	var err error

	mr.withReadLock(func() error {
		// Линейный поиск по значениям (оптимизировать при необходимости)
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

// SetValue сохраняет пару короткий-оригинальный URL
func (mr *MemoryRepository) SetValue(shortURL, originalURL, userID string) error {
	return mr.withWriteLock(func() error {
		// Проверяем, не существует ли уже такой ключ
		if _, ok := mr.data[shortURL]; ok {
			return ErrRowExists
		}
		mr.data[shortURL] = originalURL
		mr.userMap[shortURL] = userID
		return nil
	})
}

// SetValuesBatch сохраняет пакет URL пар атомарно
func (mr *MemoryRepository) SetValuesBatch(pairs map[string]string, userID string) error {
	return mr.withWriteLock(func() error {
		// Проверяем все ключи на конфликты
		for key := range pairs {
			if _, ok := mr.data[key]; ok {
				return ErrRowExists
			}
		}

		// Сохраняем все пары
		for key, value := range pairs {
			mr.data[key] = value
			mr.userMap[key] = userID
		}
		return nil
	})
}

// Close закрывает репозиторий (для in-memory реализации это no-op)
func (mr *MemoryRepository) Close() error {
	return nil
}

// GetUserURLs возвращает все активные URL пользователя
func (mr *MemoryRepository) GetUserURLs(userID string) ([]map[string]string, error) {
	var urls []map[string]string

	mr.withReadLock(func() error {
		for shortURL, originalURL := range mr.data {
			// Фильтруем по пользователю и статусу удаления
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

// DeleteURLsBatch помечает URL как удаленные для указанного пользователя
func (mr *MemoryRepository) DeleteURLsBatch(shortURLs []string, userID string) error {
	return mr.withWriteLock(func() error {
		for _, shortURL := range shortURLs {
			// Проверяем права доступа пользователя
			if ownerID, exists := mr.userMap[shortURL]; exists && ownerID == userID {
				mr.deletedMap[shortURL] = true
			}
		}
		return nil
	})
}

// IsDeleted проверяет статус удаления URL
func (mr *MemoryRepository) IsDeleted(shortURL string) (bool, error) {
	var result bool
	var err error

	mr.withReadLock(func() error {
		// Проверяем существование URL
		if _, exists := mr.data[shortURL]; !exists {
			err = errors.New("not found key in database")
			return err
		}

		// Возвращаем статус удаления
		if deleted, exists := mr.deletedMap[shortURL]; exists {
			result = deleted
		} else {
			result = false
		}
		return nil
	})

	return result, err
}
