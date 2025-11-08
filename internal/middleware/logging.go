package middleware

import (
	"go.uber.org/zap"
)

// Sugar является глобальным экземпляром логгера для использования во всем приложении.
// Использует sugared logger из zap для удобного структурированного логирования.
var Sugar zap.SugaredLogger

// InitLogger инициализирует глобальный логгер для middleware и других компонентов приложения.
// Создает логгер в режиме development с удобным для отладки выводом.
// В случае ошибки инициализации вызывает panic, так как логирование критично для приложения.
func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // Обеспечивает сброс буферизованных логов перед выходом

	// Инициализируем глобальный sugared logger для удобного использования
	Sugar = *logger.Sugar()
}
