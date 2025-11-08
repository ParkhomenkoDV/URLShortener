package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// UserIDKey является ключом для хранения и извлечения ID пользователя из контекста Gin.
// Используется для передачи аутентифицированного userID между middleware и обработчиками.
const UserIDKey = "userID"

// AuthServiceInterface определяет контракт для сервиса аутентификации.
// Используется для внедрения зависимости в auth middleware без жесткой привязки к конкретной реализации.
type AuthServiceInterface interface {
	// GetOrCreateUserID извлекает userID из куки или создает нового пользователя
	GetOrCreateUserID(r *http.Request) (string, *http.Cookie)
	// ValidateCookie проверяет валидность куки и извлекает userID
	ValidateCookie(cookieValue string) (string, bool)
}

// LoggingMiddleware создает middleware для логирования HTTP-запросов.
// Логирует URI, метод, статус ответа, длительность обработки и IP клиента.
// Вызывается для каждого входящего запроса для мониторинга и отладки.
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Передаем управление следующему обработчику в цепочке
		c.Next()

		duration := time.Since(start)

		// Логируем детали запроса с использованием структурированного логгера
		Sugar.Infoln(
			"uri", c.Request.RequestURI,
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"duration", duration,
			"client_ip", c.ClientIP(),
		)
	}
}

// GzipMiddleware создает middleware для прозрачной обработки gzip сжатия HTTP-запросов и ответов.
// Поддерживает:
//   - Распаковку входящих сжатых запросов
//   - Сжатие исходящих ответов для клиентов, поддерживающих gzip
//   - Автоматическое определение Content-Type для распакованных данных
func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ОБРАБОТКА ВХОДЯЩИХ СЖАТЫХ ДАННЫХ
		if c.GetHeader("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid gzip data"})
				return
			}
			defer reader.Close()

			// Читаем и распаковываем тело запроса
			body, err := io.ReadAll(reader)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error reading gzip data"})
				return
			}

			// Заменяем тело запроса на распакованные данные
			c.Request.Body = io.NopCloser(bytes.NewReader(body))
			c.Request.ContentLength = int64(len(body))
			c.Request.Header.Del("Content-Encoding")

			// Корректируем Content-Type для распакованных данных
			contentType := c.GetHeader("Content-Type")
			if contentType == "application/x-gzip" {
				// Определяем тип содержимого на основе пути запроса
				if strings.Contains(c.Request.URL.Path, "/api/") {
					c.Request.Header.Set("Content-Type", "application/json")
				} else {
					c.Request.Header.Set("Content-Type", "text/plain")
				}
			}
		}

		// ОБРАБОТКА ИСХОДЯЩИХ СЖАТЫХ ДАННЫХ
		// Проверяем, поддерживает ли клиент gzip для ответа
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		// Устанавливаем заголовки для сжатого ответа
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding") // Указываем кэширующим прокси, что ответ зависит от Accept-Encoding

		// Создаем gzip writer для сжатия ответа
		gz := gzip.NewWriter(c.Writer)
		defer func() {
			gz.Close() // Гарантируем закрытие writer'а при выходе из функции
		}()

		// Оборачиваем ResponseWriter для прозрачного сжатия исходящих данных
		c.Writer = &gzipWriter{
			ResponseWriter: c.Writer,
			writer:         gz,
		}

		// Передаем управление следующему обработчику
		c.Next()
	}
}

// AuthMiddleware создает middleware для автоматической аутентификации пользователей.
// Для каждого запроса:
//   - Извлекает userID из существующей валидной куки ИЛИ
//   - Создает нового пользователя и устанавливает новую куку
//   - Сохраняет userID в контексте для использования в обработчиках
//
// Этот middleware не требует обязательной аутентификации - он всегда создает пользователя если нужно.
func AuthMiddleware(authService AuthServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем userID из куки или создаем нового пользователя
		userID, cookie := authService.GetOrCreateUserID(c.Request)

		// Сохраняем userID в контексте для использования в хендлерах
		c.Set(UserIDKey, userID)

		// Если была создана новая кука (новый пользователь), устанавливаем её в ответ
		if cookie != nil {
			http.SetCookie(c.Writer, cookie)
		}

		c.Next()
	}
}

// RequireAuthMiddleware создает middleware, который требует валидную существующую куку аутентификации.
// В отличие от AuthMiddleware, этот middleware прерывает запрос с 401 статусом если:
//   - Кука отсутствует ИЛИ
//   - Кука невалидна
//
// Используется для защиты endpoint'ов, требующих обязательной аутентификации.
func RequireAuthMiddleware(authService AuthServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем наличие куки аутентификации в запросе
		cookie, err := c.Request.Cookie("user_id")
		if err != nil {
			// Кука отсутствует - возвращаем ошибку аутентификации
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Проверяем валидность куки и извлекаем userID
		userID, valid := authService.ValidateCookie(cookie.Value)
		if !valid || userID == "" {
			// Кука невалидна - возвращаем ошибку аутентификации
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Сохраняем userID в контексте для использования в хендлерах
		c.Set(UserIDKey, userID)

		c.Next()
	}
}
