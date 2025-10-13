package entity

type URL struct {
	ID          int
	ShortURL    string
	OriginalURL string
}

// NewURL создаёт новый объект URL
func NewURL(id int, shortURL, originalURL string) *URL {
	return &URL{
		ID:          id,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
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

// GetID возвращает ID
func (u *URL) GetID() int {
	return u.ID
}
