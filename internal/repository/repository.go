package repository

// URLRepository интерфейс для работы с URL
type URLRepository interface {
	// GetLongValue получает оригинальный URL по короткому
	GetLongValue(shortURL string) (string, error)
	// GetShortValue получает короткий URL по оригинальному
	GetShortValue(shortURL string) (string, error)
	// SetValue сохраняет пару короткий URL - оригинальный URL
	SetValue(shortURL, originalURL, userID string) error
	// SetValuesBatch сохраняет хэш-таблицу короткий URL - оригинальный URL
	SetValuesBatch(pairs map[string]string, userID string) error
	// GetUserURLs получает все URL пользователя
	GetUserURLs(userID string) ([]map[string]string, error)
	// Close закрывает соединение с хранилищем
	Close() error
}
