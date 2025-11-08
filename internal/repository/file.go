package repository

import (
	"errors"
	"sync"

	"github.com/ParkhomenkoDV/URLShortener/internal/persistence"
)

// FileRepository реализует хранение данных в файле с использованием JSON
type FileRepository struct {
	data         map[string]string           // shortURL -> originalURL
	reversedData map[string]string           // originalURL -> shortURL (обратное отображение)
	userMap      map[string]string           // shortURL -> userID
	deletedMap   map[string]bool             // shortURL -> isDeleted
	mu           sync.RWMutex                // Мьютекс для безопасного доступа из горутин
	filePath     string                      // Путь к файлу данных
	persistence  persistence.JSONPersistence // Сервис для работы с файловой системой
}

// NewFile создает новый файловый репозиторий и загружает данные из файла
func NewFile(filePath string) URLRepository {
	repo := &FileRepository{
		data:         make(map[string]string),
		reversedData: make(map[string]string),
		userMap:      make(map[string]string),
		deletedMap:   make(map[string]bool),
		filePath:     filePath,
		persistence:  persistence.NewFileJSONPersistence(),
	}

	// Загружаем существующие данные из файла при инициализации
	data, userMap, _, err := repo.persistence.Load(filePath)
	if err == nil {
		repo.data = data
		repo.userMap = userMap
	}

	// Создаем обратное отображение для быстрого поиска по оригинальному URL
	for originalURL, shortURL := range data {
		repo.reversedData[originalURL] = shortURL
	}

	return repo
}

// GetLongValue возвращает оригинальный URL по его короткой версии
func (r *FileRepository) GetLongValue(shortURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Проверяем, не удален ли URL
	if deleted, exists := r.deletedMap[shortURL]; exists && deleted {
		return "", errors.New("not found key in database")
	}

	if value, ok := r.data[shortURL]; ok {
		return value, nil
	}
	return "", errors.New("not found key in database")
}

// GetShortValue возвращает короткий URL по оригинальному URL
func (r *FileRepository) GetShortValue(originalURL string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if value, ok := r.reversedData[originalURL]; ok {
		return value, nil
	}
	return "", errors.New("not found key in database")
}

// SetValue сохраняет связь между коротким и оригинальным URL
func (r *FileRepository) SetValue(shortURL, originalURL, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверяем, не существует ли уже такой короткий URL
	if _, ok := r.data[shortURL]; ok {
		return ErrRowExists
	}

	r.data[shortURL] = originalURL
	r.userMap[shortURL] = userID
	r.reversedData[originalURL] = shortURL // Обновляем обратное отображение

	// Сохраняем изменения в файл
	return r.persistence.Save(r.filePath, r.data, r.userMap)
}

// SetValuesBatch сохраняет пакет URL пар в одном вызове
func (r *FileRepository) SetValuesBatch(pairs map[string]string, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Проверяем все пары на конфликты перед сохранением
	for key := range pairs {
		if _, ok := r.data[key]; ok {
			return ErrRowExists
		}
	}

	// Сохраняем все пары
	for key, value := range pairs {
		r.data[key] = value
		r.userMap[key] = userID
		r.reversedData[value] = key // Обновляем обратное отображение
	}

	// Сохраняем изменения в файл
	return r.persistence.Save(r.filePath, r.data, r.userMap)
}

// Close сохраняет данные в файл и закрывает репозиторий
func (r *FileRepository) Close() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Финальное сохранение данных в файл
	return r.persistence.Save(r.filePath, r.data, r.userMap)
}

// GetUserURLs возвращает все URL, принадлежащие указанному пользователю
func (r *FileRepository) GetUserURLs(userID string) ([]map[string]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var urls []map[string]string
	for shortURL, originalURL := range r.data {
		// Проверяем принадлежность пользователю и статус удаления
		if userID == r.userMap[shortURL] {
			if deleted, exists := r.deletedMap[shortURL]; !exists || !deleted {
				urls = append(urls, map[string]string{
					"short_url":    shortURL,
					"original_url": originalURL,
				})
			}
		}
	}
	return urls, nil
}

// DeleteURLsBatch помечает указанные URL как удаленные для заданного пользователя
// Примечание: в текущей реализации статус удаления не сохраняется в файл
func (r *FileRepository) DeleteURLsBatch(shortURLs []string, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, shortURL := range shortURLs {
		// Проверяем, принадлежит ли URL пользователю
		if ownerID, exists := r.userMap[shortURL]; exists && ownerID == userID {
			r.deletedMap[shortURL] = true
		}
	}

	// ВАЖНО: Для файлового хранилища deleteMap не сохраняется в файл
	// В реальной реализации необходимо добавить сохранение статуса удаления
	return nil
}

// IsDeleted проверяет, помечен ли URL как удаленный
func (r *FileRepository) IsDeleted(shortURL string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Проверяем существование URL
	if _, exists := r.data[shortURL]; !exists {
		return false, errors.New("not found key in database")
	}

	// Возвращаем статус удаления (по умолчанию false)
	if deleted, exists := r.deletedMap[shortURL]; exists {
		return deleted, nil
	}
	return false, nil
}
