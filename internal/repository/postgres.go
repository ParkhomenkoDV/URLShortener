package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

// PostgreSQLRepository реализация репозитория для работы с PostgreSQL
type PostgreSQLRepository struct {
	pool *pgxpool.Pool
}

// NewPostgreSQLRepository создает новый репозиторий для работы с PostgreSQL
func NewPostgreSQLRepository(dsn string) (URLRepository, error) {
	// Подключаемся к базе данных
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	repo := &PostgreSQLRepository{
		pool: pool,
	}

	// Выполняем миграции
	if err := repo.runMigrations(dsn); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	return repo, nil
}

// runMigrations выполняет миграции базы данных
func (r *PostgreSQLRepository) runMigrations(dsn string) error {
	// Открываем соединение через database/sql для миграций
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// GetValue получает оригинальный URL по короткому
func (r *PostgreSQLRepository) GetFullValue(shortURL string) (string, error) {
	var originalURL string
	err := r.pool.QueryRow(context.Background(),
		"SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("not found key in database")
		}
		return "", fmt.Errorf("failed to get value: %v", err)
	}

	return originalURL, nil
}

// GetShortValue получает оригинальный URL по короткому
func (r *PostgreSQLRepository) GetShortValue(originalURL string) (string, error) {
	var shortURL string
	err := r.pool.QueryRow(context.Background(),
		"SELECT short_url FROM urls WHERE original_url = $1", originalURL).Scan(&shortURL)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("not found key in database")
		}
		return "", fmt.Errorf("failed to get value: %v", err)
	}

	return shortURL, nil
}

// SetValue сохраняет пару короткий URL - оригинальный URL
func (r *PostgreSQLRepository) SetValue(shortURL, originalURL string) error {
	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	var result string
	err = tx.QueryRow(context.Background(),
		`INSERT INTO urls (short_url, original_url)
		 VALUES ($1, $2)
		 ON CONFLICT (original_url) DO NOTHING
		 RETURNING short_url`,
		shortURL, originalURL).Scan(&result)

	// Запись уже существует
	if errors.Is(err, sql.ErrNoRows) {
		return ErrRowExists
	}
	if err != nil {
		return fmt.Errorf("failed to insert url: %v", err)
	}

	if err = tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// SetValuesBatch сохраняет пакет пар короткий URL - оригинальный URL
func (r *PostgreSQLRepository) SetValuesBatch(pairs map[string]string) error {
	if len(pairs) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	// Подготавливаем пакетную вставку
	for shortURL, originalURL := range pairs {
		// Удаляем любую запись с таким short_url (если она не та, что будет обновлена)
		_, err = tx.Exec(context.Background(),
			`DELETE FROM urls WHERE short_url = $1 AND original_url != $2`,
			shortURL, originalURL)
		if err != nil {
			return fmt.Errorf("failed to delete conflicting short_url: %v", err)
		}

		// Вставляем или обновляем по original_url
		_, err = tx.Exec(context.Background(),
			`INSERT INTO urls (short_url, original_url) 
			 VALUES ($1, $2)`,
			shortURL, originalURL)
		if err != nil {
			return fmt.Errorf("failed to upsert url: %v", err)
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// Close закрывает соединение с базой данных
func (r *PostgreSQLRepository) Close() error {
	r.pool.Close()
	return nil
}
