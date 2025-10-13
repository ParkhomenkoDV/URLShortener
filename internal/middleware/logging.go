package middleware

import (
	"go.uber.org/zap"
)

// Singleton для логгера
var Sugar zap.SugaredLogger

// Инициализатор для логгера
func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		// вызываем панику, если ошибка
		panic(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	Sugar = *logger.Sugar()
}