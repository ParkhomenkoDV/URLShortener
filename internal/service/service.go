package service

import (
	"database/sql"
	"errors"

	"github.com/ParkhomenkoDV/URLShortener/internal/config"
	"github.com/ParkhomenkoDV/URLShortener/internal/repository"
	"github.com/ParkhomenkoDV/URLShortener/pkg/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Сервис сокращения ссылок
type Service struct {
	Repository    repository.URLRepository
	Configuration *config.Config
}

// Конструктор для сервиса
func New(repo repository.URLRepository, configuration *config.Config) *Service {
	return &Service{
		Repository:    repo,
		Configuration: configuration,
	}
}

// Создание сокращенного URL
func (s *Service) CreateShortURL(url, userID string) (string, error) {
	var shortURL string

	// Генерируем сокращенную уникальную ссылку
	for {
		shortURL = utils.GenerateShortKey(s.Configuration.LengthKey)
		if _, err := s.Repository.GetLongValue(shortURL); err == nil {
			continue
		}
		break
	}

	// Сохраняем в репозитории
	if err := s.Repository.SetValue(shortURL, url, userID); err != nil {
		if errors.Is(err, repository.ErrRowExists) {
			// Если ссылка уже существует
			if shortURL, err = s.Repository.GetShortValue(url); err == nil {
				return shortURL, repository.ErrRowExists
			}
		}
		return "", err
	}

	return shortURL, nil
}

// CreateShortURLsBatch создает сокращенные URL
func (s *Service) CreateShortURLsBatch(urls []string, userID string) (map[string]string, error) {
	if len(urls) == 0 {
		return make(map[string]string), nil
	}

	result, pairs := make(map[string]string), make(map[string]string)

	// Генерируем короткие URL для каждого исходного URL
	for _, originalURL := range urls {
		var shortURL string

		// Генерируем уникальную короткую ссылку
		for {
			shortURL = utils.GenerateShortKey(s.Configuration.LengthKey)
			if _, err := s.Repository.GetLongValue(shortURL); err == nil {
				continue
			}
			// Проверяем на коализии, по хорошему надо бы сделать рекурсию
			if _, exists := pairs[shortURL]; exists {
				continue
			}
			break
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

// Получение полного URL
func (s *Service) GetFullURL(shortURL string) (string, error) {
	// Ищем полный URL в репозитории, или выдаем ошибку
	if url, err := s.Repository.GetLongValue(shortURL); err == nil {
		return url, nil
	} else {
		return "", errors.New("not found")
	}
}

// GetUserURLs получает все URL пользователя
func (s *Service) GetUserURLs(userID string) ([]map[string]string, error) {
	return s.Repository.GetUserURLs(userID)
}

// Ping DB
func (s *Service) PingPostgreSQL() error {
	db, err := sql.Open("pgx", s.Configuration.AddressDB)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

// Close закрывает соединение с репозиторием
func (s *Service) Close() error {
	return s.Repository.Close()
}
