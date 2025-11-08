package repository

import (
	"errors"
	"log"
)

// ErrRowExists возникает при попытке создать запись, которая уже существует
var ErrRowExists = errors.New("short URL already exists")

// CreateRepository создает репозиторий в зависимости от конфигурации
// Приоритет выбора хранилища: PostgreSQL -> File -> Memory
func CreateRepository(databaseDSN, filePath string) URLRepository {
	// Если указан DSN для PostgreSQL, пытаемся использовать его
	if databaseDSN != "" {
		log.Printf("Use PostgreSQL with DSN: %s", databaseDSN)
		repo, err := NewPostgreSQLRepository(databaseDSN)
		if err != nil {
			log.Printf("Create PostgreSQL error: %v", err)
		} else {
			return repo
		}
	}

	// Если указан путь к файлу, используем файловое хранилище
	if filePath != "" {
		log.Printf("Use File with path: %s", filePath)
		return NewFile(filePath)
	}

	// По умолчанию используем хранилище в памяти
	log.Println("Use Memory")
	return NewMemory()
}
