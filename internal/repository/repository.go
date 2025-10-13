package repository

// URLRepository интерфейс для работы с URL
type URLRepository interface {
	// GetFullValue получает оригинальный URL по короткому
	GetFullValue(shortURL string) (string, error)
	// GetShortValue получает короткий URL по оригинальному
	GetShortValue(shortURL string) (string, error)
	// SetValue сохраняет пару короткий URL - оригинальный URL
	SetValue(shortURL, originalURL string) error
	// SetValuesBatch сохраняет пакет пар короткий URL - оригинальный URL
	SetValuesBatch(pairs map[string]string) error
	// Close закрывает соединение с хранилищем
	Close() error
}
