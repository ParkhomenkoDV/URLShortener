package logger

import (
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Singleton для логгера
var sugar zap.SugaredLogger

// Инициализатор для логгера
func New() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	sugar = *logger.Sugar()
}

// loggingResponseWriter оборачивает http.ResponseWriter для отслеживания статуса
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}

// Middleware для логирования запросов
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Создаем кастомный ResponseWriter для отслеживания статуса
		lw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Обслуживаем запрос
		next.ServeHTTP(lw, r)

		duration := time.Since(start)

		// Получаем IP клиента
		clientIP := getClientIP(r)

		// Логируем информацию о запросе
		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", lw.statusCode,
			"duration", duration,
			"client_ip", clientIP,
		)
	})
}

// getClientIP возвращает IP адрес клиента
func getClientIP(r *http.Request) string {
	// Пробуем получить IP из заголовков (для прокси)
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}

	// Возвращаем IP из RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
