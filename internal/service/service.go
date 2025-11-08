package service

import (
	"database/sql"
	"errors"

	"github.com/ParkhomenkoDV/URLShortener/internal/config"
	"github.com/ParkhomenkoDV/URLShortener/internal/repository"
	"github.com/ParkhomenkoDV/URLShortener/pkg/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Service предоставляет функциональность для сокращения URL
type Service struct {
	Repository    repository.URLRepository // Репозиторий для хранения данных
	Configuration *config.Config           // Конфигурация приложения
}

// New создает новый экземпляр сервиса
func New(repo repository.URLRepository, configuration *config.Config) *Service {
	return &Service{
		Repository:    repo,
		Configuration: configuration,
	}
}

// CreateShortURL создает сокращенную версию URL
// Если URL уже существует, возвращает существующую сокращенную версию
func (s *Service) CreateShortURL(url, userID string) (string, error) {
	var shortURL string

	// Генерируем уникальный короткий ключ, избегая коллизий
	for {
		shortURL = utils.GenerateShortKey(s.Configuration.LengthKey)
		// Проверяем, не существует ли уже такой ключ
		if _, err := s.Repository.GetLongValue(shortURL); err == nil {
			continue // Ключ существует, генерируем новый
		}
		break // Уникальный ключ найден
	}

	// Сохраняем связь в репозитории
	if err := s.Repository.SetValue(shortURL, url, userID); err != nil {
		if errors.Is(err, repository.ErrRowExists) {
			// Если URL уже существует, возвращаем существующую сокращенную версию
			if shortURL, err = s.Repository.GetShortValue(url); err == nil {
				return shortURL, repository.ErrRowExists
			}
		}
		return "", err
	}

	return shortURL, nil
}

// CreateShortURLsBatch создает сокращенные URL для пакета исходных URL
// Возвращает карту соответствия оригинальных URL и их сокращенных версий
func (s *Service) CreateShortURLsBatch(urls []string, userID string) (map[string]string, error) {
	if len(urls) == 0 {
		return make(map[string]string), nil
	}

	result := make(map[string]string) // оригинальный URL -> короткий URL
	pairs := make(map[string]string)  // короткий URL -> оригинальный URL

	// Генерируем уникальные короткие URL для каждого исходного URL
	for _, originalURL := range urls {
		var shortURL string

		// Генерируем уникальную короткую ссылку
		for {
			shortURL = utils.GenerateShortKey(s.Configuration.LengthKey)
			// Проверяем на существование в репозитории
			if _, err := s.Repository.GetLongValue(shortURL); err == nil {
				continue
			}
			// Проверяем на коллизии внутри текущего пакета
			if _, exists := pairs[shortURL]; exists {
				continue
			}
			break // Уникальный ключ найден
		}

		pairs[shortURL] = originalURL
		result[originalURL] = shortURL
	}

	// Сохраняем пакет в репозитории
	if err := s.Repository.SetValuesBatch(pairs, userID); err != nil {
		return nil, err
	}

	return result, nil
}

// GetFullURL возвращает оригинальный URL по его сокращенной версии
// Если URL помечен как удаленный, возвращает ошибку
func (s *Service) GetFullURL(shortURL string) (string, error) {
	// Проверяем, не удален ли URL
	if deleted, err := s.Repository.IsDeleted(shortURL); err == nil && deleted {
		return "", errors.New("url is deleted")
	}

	// Ищем полный URL в репозитории
	if url, err := s.Repository.GetLongValue(shortURL); err == nil {
		return url, nil
	} else {
		return "", errors.New("not found")
	}
}

// GetUserURLs возвращает все URL, созданные указанным пользователем
func (s *Service) GetUserURLs(userID string) ([]map[string]string, error) {
	return s.Repository.GetUserURLs(userID)
}

// PingPostgreSQL проверяет соединение с базой данных PostgreSQL
func (s *Service) PingPostgreSQL() error {
	db, err := sql.Open("pgx", s.Configuration.AddressDB)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

// DeleteURLsBatch помечает указанные URL как удаленные для заданного пользователя
func (s *Service) DeleteURLsBatch(shortURLs []string, userID string) error {
	return s.Repository.DeleteURLsBatch(shortURLs, userID)
}

// IsDeleted проверяет, помечен ли URL как удаленный
func (s *Service) IsDeleted(shortURL string) (bool, error) {
	return s.Repository.IsDeleted(shortURL)
}

// Close закрывает соединение с репозиторием
func (s *Service) Close() error {
	return s.Repository.Close()
}
