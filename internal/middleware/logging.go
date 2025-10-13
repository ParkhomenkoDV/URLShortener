package middleware

import (
	"go.uber.org/zap"
)

var Sugar zap.SugaredLogger // Singleton для логгера

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
