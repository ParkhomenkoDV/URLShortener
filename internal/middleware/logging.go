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
		panic(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	Sugar = *logger.Sugar()
}
