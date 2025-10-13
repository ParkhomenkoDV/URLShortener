package repository

import (
	"errors"
	"log"
)

// ErrRowExists ошибка, которая возникает, когда запись уже существует
var ErrRowExists = errors.New("short URL already exists")

// CreateRepository создает репозиторий в зависимости от конфигурации
// Приоритет: PostgreSQL -> File -> Memory
func CreateRepository(databaseDSN, filePath string) URLRepository {
	// Если есть DATABASE_DSN и он не пустой, используем PostgreSQL
	if databaseDSN != "" && !isDefaultPostgresValue(databaseDSN) {
		log.Printf("Используем PostgreSQL репозиторий с DSN: %s", databaseDSN)
		repo, err := NewPostgreSQLRepository(databaseDSN)
		if err != nil {
			log.Printf("Ошибка создания PostgreSQL репозитория: %v. Переходим к файловому хранилищу", err)
		} else {
			return repo
		}
	}

	// Если есть FILE_STORAGE_PATH и он не пустой, используем файловое хранилище
	if filePath != "" {
		log.Printf("Используем файловый репозиторий с путем: %s", filePath)
		return NewFileRepository(filePath)
	}

	// Иначе используем память
	log.Println("Используем репозиторий в памяти")
	return NewMemoryRepository()
}

// isDefaultPostgresValue проверяет, является ли значение DSN дефолтным значением из флагов
func isDefaultPostgresValue(dsn string) bool {
	defaultDSN := ""
	return dsn == defaultDSN
}
