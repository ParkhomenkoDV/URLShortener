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
type AuthService struct {
	secretKey []byte
}

// New создает новый экземпляр AuthService
func New(secretKey string) *AuthService {
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
	if value == "" {
		return ""
	}
	h := hmac.New(sha256.New, a.secretKey)
	h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

// CreateSignedCookie создает подписанную куку с ID пользователя
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
		HttpOnly: true,
		Secure:   false, // установить в true для production
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
	parts := strings.Split(cookieValue, ":")
	if len(parts) != 2 {
		return "", false
	}

	userID := parts[0]
	providedSignature := parts[1]

	// Проверяем подпись
	expectedSignature := a.SignValue(userID)

	return userID, hmac.Equal([]byte(providedSignature), []byte(expectedSignature))
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
