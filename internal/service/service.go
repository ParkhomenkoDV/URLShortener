package service

import (
	"database/sql"
	"errors"

	"github.com/ParkhomenkoDV/URLShortener/internal/config"
	"github.com/ParkhomenkoDV/URLShortener/internal/repository"
	"github.com/ParkhomenkoDV/URLShortener/pkg/utils"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Структура для сервиса сокращения ссылок
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
func (u *Service) CreateShortURL(url string) (string, error) {
	// Инициализация результата
	var shortURL string

	// Генерируем сокращенную уникальную ссылку
	for {
		shortURL = utils.GenerateShortKey(6)
		if _, err := u.Repository.GetFullValue(shortURL); err == nil {
			continue
		}
		break
	}

	// Сохраняем в репозитории
	if err := u.Repository.SetValue(shortURL, url); err != nil {
		if errors.Is(err, repository.ErrRowExists) {
			// Если ссылка уже существует
			if shortURL, err = u.Repository.GetShortValue(url); err == nil {
				return shortURL, repository.ErrRowExists
			}
		}
		return "", err
	}

	return shortURL, nil
}

// CreateShortURLsBatch создает сокращенные URL для пакета URL
func (u *Service) CreateShortURLsBatch(urls []string) (map[string]string, error) {
	if len(urls) == 0 {
		return make(map[string]string), nil
	}

	result := make(map[string]string)
	pairs := make(map[string]string)

	// Генерируем короткие URL для каждого исходного URL
	for _, originalURL := range urls {
		var shortURL string

		// Генерируем уникальную короткую ссылку
		for {
			shortURL = utils.GenerateShortKey(6)
			if _, err := u.Repository.GetFullValue(shortURL); err == nil {
				continue
			}
			// Проверяем также, что этот ключ не используется в текущем пакете
			if _, exists := pairs[shortURL]; exists {
				continue
			}
			break
		}

		pairs[shortURL] = originalURL
		result[originalURL] = shortURL
	}

	// Сохраняем пакет в репозитории
	if err := u.Repository.SetValuesBatch(pairs); err != nil {
		return nil, err
	}

	return result, nil
}

// Получение полного URL
func (u *Service) GetFullURL(shortURL string) (string, error) {
	// Ищем полный URL в репозитории, или выдаем ошибку
	if url, err := u.Repository.GetFullValue(shortURL); err == nil {
		return url, nil
	} else {
		return "", errors.New("not found")
	}
}

// Ping DB
func (u *Service) PingPostgreSQL() error {
	db, err := sql.Open("pgx", u.Configuration.AddressDB)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Ping()
}

// Close закрывает соединение с репозиторием
func (u *Service) Close() error {
	return u.Repository.Close()
}
