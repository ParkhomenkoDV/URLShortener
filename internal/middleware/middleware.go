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

// UserIDKey является ключом для хранения ID пользователя в контексте
const UserIDKey = "userID"

// Middleware для логирования запросов
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Обслуживаем запрос
		c.Next()

		duration := time.Since(start)

		Sugar.Infoln(
			"uri", c.Request.RequestURI,
			"method", c.Request.Method,
			"status", c.Writer.Status(),
			"duration", duration,
			"client_ip", c.ClientIP(),
		)
	}
}

// GzipMiddleware middleware для обработки gzip сжатия
func GzipMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Обрабатываем входящие сжатые данные
		if c.GetHeader("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid gzip data"})
				return
			}
			defer reader.Close()

			// Заменяем тело запроса на распакованное
			body, err := io.ReadAll(reader)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Error reading gzip data"})
				return
			}

			c.Request.Body = io.NopCloser(bytes.NewReader(body))
			c.Request.ContentLength = int64(len(body))
			c.Request.Header.Del("Content-Encoding")

			// Восстанавливаем правильный Content-Type для разных случаев
			contentType := c.GetHeader("Content-Type")
			if contentType == "application/x-gzip" {
				// Определяем тип содержимого по тому, что делаем
				if strings.Contains(c.Request.URL.Path, "/api/") {
					c.Request.Header.Set("Content-Type", "application/json")
				} else {
					c.Request.Header.Set("Content-Type", "text/plain")
				}
			}
		}

		// Проверяем, поддерживает ли клиент gzip для ответа
		if !strings.Contains(c.GetHeader("Accept-Encoding"), "gzip") {
			c.Next()
			return
		}

		// Устанавливаем заголовки для сжатого ответа
		c.Header("Content-Encoding", "gzip")
		c.Header("Vary", "Accept-Encoding")

		// Создаем gzip writer для ответа
		gz := gzip.NewWriter(c.Writer)
		defer func() {
			gz.Close()
		}()

		// Оборачиваем ResponseWriter
		c.Writer = &gzipWriter{
			ResponseWriter: c.Writer,
			writer:         gz,
		}

		c.Next()
	}
}

// AuthMiddleware middleware для аутентификации пользователей
// Интерфейс для сервиса аутентификации
type AuthServiceInterface interface {
	GetOrCreateUserID(r *http.Request) (string, *http.Cookie)
	ValidateCookie(cookieValue string) (string, bool)
}

func AuthMiddleware(authService AuthServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, cookie := authService.GetOrCreateUserID(c.Request)
		
		// Сохраняем userID в контексте для использования в хендлерах
		c.Set(UserIDKey, userID)
		
		// Если была создана новая кука, устанавливаем её
		if cookie != nil {
			http.SetCookie(c.Writer, cookie)
		}
		
		c.Next()
	}
}

// RequireAuthMiddleware middleware, который требует валидную существующую куку
func RequireAuthMiddleware(authService AuthServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Проверяем наличие куки в запросе
		cookie, err := c.Request.Cookie("user_id")
		if err != nil {
			// Кука отсутствует
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		
		// Проверяем валидность куки
		userID, valid := authService.ValidateCookie(cookie.Value)
		if !valid || userID == "" {
			// Кука невалидна
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		
		// Сохраняем userID в контексте для использования в хендлерах
		c.Set(UserIDKey, userID)
		
		c.Next()
	}
}