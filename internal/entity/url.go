package entity

// URL представляет доменную модель URL
type URL struct {
	ID          int
	ShortURL    string
	OriginalURL string
	UserID      string
}

// NewURL создаёт новую доменную сущность URL
func NewURL(id int, shortURL, originalURL, userID string) *URL {
	return &URL{
		ID:          id,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
		UserID:      userID,
	}
}

// GetShortURL возвращает сокращённый URL
func (u *URL) GetShortURL() string {
	return u.ShortURL
}

// GetOriginalURL возвращает оригинальный URL
func (u *URL) GetOriginalURL() string {
	return u.OriginalURL
}

// GetID возвращает идентификатор
func (u *URL) GetID() int {
	return u.ID
}

// GetUserID возвращает идентификатор пользователя
func (u *URL) GetUserID() string {
	return u.UserID
}
