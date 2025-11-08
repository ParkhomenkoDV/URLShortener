package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

// PostgreSQLRepository реализует хранение данных в PostgreSQL
type PostgreSQLRepository struct {
	pool *pgxpool.Pool // Пул соединений с БД
}

// NewPostgreSQLRepository создает новый PostgreSQL репозиторий и выполняет миграции
func NewPostgreSQLRepository(dsn string) (URLRepository, error) {
	// Создаем пул соединений с базой данных
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	repo := &PostgreSQLRepository{
		pool: pool,
	}

	// Выполняем миграции схемы базы данных
	if err := repo.runMigrations(dsn); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to run migrations: %v", err)
	}

	return repo, nil
}

// runMigrations выполняет миграции базы данных из папки migrations
func (r *PostgreSQLRepository) runMigrations(dsn string) error {
	// Используем database/sql для совместимости с пакетом миграций
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// Создаем драйвер для миграций
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	// Инициализируем мигратор
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Путь к файлам миграций
		"postgres", driver)
	if err != nil {
		return err
	}

	// Применяем миграции
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

// GetLongValue возвращает оригинальный URL по короткому ключу
func (r *PostgreSQLRepository) GetLongValue(shortURL string) (string, error) {
	var originalURL string

	// Ищем только активные (не удаленные) URL
	err := r.pool.QueryRow(context.Background(),
		"SELECT original_url FROM urls WHERE short_url = $1 AND (is_deleted = FALSE OR is_deleted IS NULL)",
		shortURL).Scan(&originalURL)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New("not found key in database")
		}
		return "", fmt.Errorf("failed to get value: %v", err)
	}

	return originalURL, nil
}

// GetShortValue возвращает короткий ключ по оригинальному URL
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

// SetValue сохраняет пару URL в базе данных
// Использует INSERT с ON CONFLICT для обработки дубликатов
func (r *PostgreSQLRepository) SetValue(shortURL, originalURL, userID string) error {
	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	var result string
	// Пытаемся вставить запись, игнорируя конфликты по original_url
	err = tx.QueryRow(context.Background(),
		`INSERT INTO urls (short_url, original_url, user_id)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (original_url) DO NOTHING
		 RETURNING short_url`,
		shortURL, originalURL, userID).Scan(&result)

	// Если запись не вставлена (конфликт), возвращаем ошибку
	if errors.Is(err, sql.ErrNoRows) {
		return ErrRowExists
	}
	if err != nil {
		return fmt.Errorf("failed to insert url: %v", err)
	}

	// Фиксируем транзакцию
	if err = tx.Commit(context.Background()); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// SetValuesBatch сохраняет пакет URL пар в одной транзакции
func (r *PostgreSQLRepository) SetValuesBatch(pairs map[string]string, userID string) error {
	if len(pairs) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	// Вставляем каждую пару в транзакции
	for shortURL, originalURL := range pairs {
		// Удаляем возможные конфликтующие записи с тем же short_url
		_, err = tx.Exec(context.Background(),
			`DELETE FROM urls WHERE short_url = $1 AND original_url != $2`,
			shortURL, originalURL)
		if err != nil {
			return fmt.Errorf("failed to delete conflicting short_url: %v", err)
		}

		// Вставляем новую запись
		_, err = tx.Exec(context.Background(),
			`INSERT INTO urls (short_url, original_url, user_id) 
			 VALUES ($1, $2, $3)`,
			shortURL, originalURL, userID)
		if err != nil {
			return fmt.Errorf("failed to upsert url: %v", err)
		}
	}

	// Фиксируем всю транзакцию
	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// GetUserURLs возвращает все активные URL пользователя
func (r *PostgreSQLRepository) GetUserURLs(userID string) ([]map[string]string, error) {
	rows, err := r.pool.Query(context.Background(),
		"SELECT short_url, original_url FROM urls WHERE user_id = $1 AND (is_deleted = FALSE OR is_deleted IS NULL)",
		userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user urls: %v", err)
	}
	defer rows.Close()

	var urls []map[string]string
	for rows.Next() {
		var shortURL, originalURL string
		if err := rows.Scan(&shortURL, &originalURL); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		urls = append(urls, map[string]string{
			"short_url":    shortURL,
			"original_url": originalURL,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %v", err)
	}

	return urls, nil
}

// DeleteURLsBatch помечает URL как удаленные в пакетном режиме
func (r *PostgreSQLRepository) DeleteURLsBatch(shortURLs []string, userID string) error {
	if len(shortURLs) == 0 {
		return nil
	}

	tx, err := r.pool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	// Создаем динамический запрос с IN clause
	placeholders := make([]string, len(shortURLs))
	args := make([]interface{}, len(shortURLs)+1)
	args[0] = userID

	for i, shortURL := range shortURLs {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		args[i+1] = shortURL
	}

	query := fmt.Sprintf(`
		UPDATE urls 
		SET is_deleted = TRUE 
		WHERE user_id = $1 AND short_url IN (%s)`,
		strings.Join(placeholders, ","))

	_, err = tx.Exec(context.Background(), query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete urls: %v", err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// IsDeleted проверяет статус удаления URL
func (r *PostgreSQLRepository) IsDeleted(shortURL string) (bool, error) {
	var isDeleted bool
	err := r.pool.QueryRow(context.Background(),
		"SELECT COALESCE(is_deleted, FALSE) FROM urls WHERE short_url = $1", shortURL).Scan(&isDeleted)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("not found key in database")
		}
		return false, fmt.Errorf("failed to check deletion status: %v", err)
	}

	return isDeleted, nil
}

// Close закрывает пул соединений с базой данных
func (r *PostgreSQLRepository) Close() error {
	r.pool.Close()
	return nil
}
