package repository

import (
	"errors"
	"log"
)

// ErrRowExists ошибка, которая возникает, когда запись уже существует
var ErrRowExists = errors.New("short URL already exists")

// CreateRepository создает репозиторий в зависимости от конфигурации по приоритету: PostgreSQL -> File -> Memory
func CreateRepository(databaseDSN, filePath string) URLRepository {
	// Если есть DATABASE_DSN и он не пустой, используем PostgreSQL
	if databaseDSN != "" {
		log.Printf("Use PostgreSQL with DSN: %s", databaseDSN)
		repo, err := NewPostgreSQLRepository(databaseDSN)
		if err != nil {
			log.Printf("Create PostgreSQL error: %v", err)
		} else {
			return repo
		}
	}

	// Если есть FILE_STORAGE_PATH и он не пустой, используем файловое хранилище
	if filePath != "" {
		log.Printf("Use File with path: %s", filePath)
		return NewFile(filePath)
	}

	// Иначе используем память
	log.Println("Use Memory")
	return NewMemory()
}
