package repository

// URLRepository определяет контракт для работы с хранилищем URL
type URLRepository interface {
	// GetLongValue возвращает оригинальный URL по короткому ключу
	GetLongValue(shortURL string) (string, error)

	// GetShortValue возвращает короткий ключ по оригинальному URL
	GetShortValue(originalURL string) (string, error)

	// SetValue сохраняет связь между коротким и оригинальным URL
	SetValue(shortURL, originalURL, userID string) error

	// SetValuesBatch сохраняет пакет URL пар атомарно
	SetValuesBatch(pairs map[string]string, userID string) error

	// GetUserURLs возвращает все URL, принадлежащие пользователю
	GetUserURLs(userID string) ([]map[string]string, error)

	// DeleteURLsBatch помечает URL как удаленные для указанного пользователя
	DeleteURLsBatch(shortURLs []string, userID string) error

	// IsDeleted проверяет, помечен ли URL как удаленный
	IsDeleted(shortURL string) (bool, error)

	// Close освобождает ресурсы хранилища
	Close() error
}
