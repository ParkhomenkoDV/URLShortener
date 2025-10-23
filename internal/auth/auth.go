package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// AuthService предоставляет функциональность для аутентификации пользователей
type AuthService struct {
	secretKey []byte
}

// NewAuthService создает новый экземпляр AuthService
func NewAuthService(secretKey string) *AuthService {
	return &AuthService{
		secretKey: []byte(secretKey),
	}
}

// GenerateUserID создает новый уникальный идентификатор пользователя
func (a *AuthService) GenerateUserID() string {
	return uuid.New().String()
}

// SignValue создает подпись для значения
func (a *AuthService) SignValue(value string) string {
	h := hmac.New(sha256.New, a.secretKey)
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

// CreateSignedCookie создает подписанную куку с ID пользователя
func (a *AuthService) CreateSignedCookie(userID string) *http.Cookie {
	signature := a.SignValue(userID)
	cookieValue := fmt.Sprintf("%s:%s", userID, signature)

	return &http.Cookie{
		Name:     "user_id",
		Value:    cookieValue,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(30 * 24 * time.Hour.Seconds()),
	}
}

// ValidateCookie проверяет подлинность куки и извлекает ID пользователя
func (a *AuthService) ValidateCookie(cookieValue string) (string, bool) {
	if cookieValue == "" {
		return "", false
	}

	// Разделяем значение куки на userID и подпись
	parts := []string{}
	colonIndex := -1
	for i := len(cookieValue) - 1; i >= 0; i-- {
		if cookieValue[i] == ':' {
			colonIndex = i
			break
		}
	}

	if colonIndex == -1 {
		return "", false
	}

	parts = append(parts, cookieValue[:colonIndex])
	parts = append(parts, cookieValue[colonIndex+1:])

	if len(parts) != 2 {
		return "", false
	}

	userID := parts[0]
	providedSignature := parts[1]

	// Проверяем подпись
	expectedSignature := a.SignValue(userID)
	if !hmac.Equal([]byte(providedSignature), []byte(expectedSignature)) {
		return "", false
	}

	return userID, true
}

// GetOrCreateUserID извлекает ID пользователя из куки или создает нового пользователя
func (a *AuthService) GetOrCreateUserID(r *http.Request) (string, *http.Cookie) {
	// Пытаемся получить куку
	cookie, err := r.Cookie("user_id")
	if err == nil {
		// Проверяем валидность куки
		if userID, valid := a.ValidateCookie(cookie.Value); valid {
			return userID, nil
		}
	}

	// Создаем нового пользователя
	newUserID := a.GenerateUserID()
	newCookie := a.CreateSignedCookie(newUserID)
	return newUserID, newCookie
}
