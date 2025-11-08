package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AuthService предоставляет функциональность для аутентификации пользователей
// через подписанные cookies. Использует HMAC-SHA256 для подписи идентификаторов.
type AuthService struct {
	secretKey []byte
}

// New создает новый экземпляр AuthService с указанным секретным ключом.
// Секретный ключ должен быть достаточно длинным и храниться в безопасности.
func New(secretKey string) *AuthService {
	return &AuthService{
		secretKey: []byte(secretKey),
	}
}

// GenerateUserID создает новый уникальный идентификатор пользователя
// используя UUID v4. Возвращает строку в формате UUID.
func (a *AuthService) GenerateUserID() string {
	return uuid.New().String()
}

// SignValue создает HMAC-SHA256 подпись для переданного значения.
// Возвращает hex-encoded строку подписи.
// Для пустого значения возвращает пустую строку.
func (a *AuthService) SignValue(value string) string {
	if value == "" {
		return ""
	}
	h := hmac.New(sha256.New, a.secretKey)
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

// CreateSignedCookie создает подписанную куку с ID пользователя.
// Формат значения куки: "user_id:signature".
// Возвращает nil если userID пустой.
func (a *AuthService) CreateSignedCookie(userID string) *http.Cookie {
	if userID == "" {
		return nil
	}
	signature := a.SignValue(userID)
	cookieValue := fmt.Sprintf("%s:%s", userID, signature)

	return &http.Cookie{
		Name:     "user_id",
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,                               // Защита от XSS атак
		Secure:   false,                              // TODO: установить в true для production (HTTPS)
		SameSite: http.SameSiteLaxMode,               // Защита от CSRF атак
		MaxAge:   int(30 * 24 * time.Hour.Seconds()), // 30 дней
	}
}

// ValidateCookie проверяет подлинность куки и извлекает ID пользователя.
// Возвращает userID и флаг валидности подписи.
// Проверяет формат значения и соответствие HMAC подписи.
func (a *AuthService) ValidateCookie(cookieValue string) (string, bool) {
	if cookieValue == "" {
		return "", false
	}

	// Ожидаемый формат: "user_id:signature"
	parts := strings.Split(cookieValue, ":")
	if len(parts) != 2 {
		return "", false
	}

	userID := parts[0]
	providedSignature := parts[1]

	// Проверяем подпись с использованием constant-time сравнения
	expectedSignature := a.SignValue(userID)

	return userID, hmac.Equal([]byte(providedSignature), []byte(expectedSignature))
}

// GetOrCreateUserID извлекает ID пользователя из куки или создает нового пользователя.
// Возвращает userID и новую куку (если была создана).
// Если кука существует и валидна - возвращает существующий userID и nil куку.
// Если кука отсутствует или невалидна - создает нового пользователя и возвращает новую куку.
func (a *AuthService) GetOrCreateUserID(r *http.Request) (string, *http.Cookie) {
	// Пытаемся получить существующую куку
	cookie, err := r.Cookie("user_id")
	if err == nil {
		// Проверяем валидность куки
		if userID, valid := a.ValidateCookie(cookie.Value); valid {
			return userID, nil
		}
		// Если кука невалидна, продолжаем создавать нового пользователя
	}

	// Создаем нового пользователя и куку
	newUserID := a.GenerateUserID()
	newCookie := a.CreateSignedCookie(newUserID)
	return newUserID, newCookie
}
