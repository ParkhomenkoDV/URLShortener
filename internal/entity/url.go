package entity

type URL struct {
	ID          int
	ShortURL    string
	OriginalURL string
	UserID      string
	IsDeleted   bool
}

// NewURL создаёт новый объект URL
func NewURL(id int, shortURL, originalURL, userID string) *URL {
	return &URL{
		ID:          id,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		UserID:      userID,
		IsDeleted:   false,
	}
}

// GetShortURL возвращает сокращённый URL
func (url *URL) GetShortURL() string {
	return url.ShortURL
}

// GetOriginalURL возвращает оригинальный URL
func (url *URL) GetOriginalURL() string {
	return url.OriginalURL
}

// GetID возвращает ID
func (url *URL) GetID() int {
	return url.ID
}

// GetUserID возвращает идентификатор пользователя
func (url *URL) GetUserID() string {
	return url.UserID
}
