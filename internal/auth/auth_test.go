package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestAuthService проверяет основную функциональность сервиса аутентификации
func TestAuthService(t *testing.T) {
	// Создаем сервис аутентификации с тестовым ключом
	authService := New("test-secret-key")

	t.Run("Test GenerateUserID", func(t *testing.T) {
		userID := authService.GenerateUserID()
		if userID == "" {
			t.Error("Generated user ID should not be empty")
		}
	})

	t.Run("Test CreateSignedCookie", func(t *testing.T) {
		userID := authService.GenerateUserID()
		cookie := authService.CreateSignedCookie(userID)

		if cookie == nil {
			t.Error("Cookie should not be nil")
			return
		}

		// Проверяем основные атрибуты безопасности куки
		if cookie.Name != "user_id" {
			t.Errorf("Cookie name should be 'user_id', got %s", cookie.Name)
		}
		if !cookie.HttpOnly {
			t.Error("HttpOnly should be true for XSS protection")
		}
		if cookie.SameSite != http.SameSiteLaxMode {
			t.Error("SameSite should be Lax for CSRF protection")
		}
		if cookie.MaxAge != int(30*24*time.Hour.Seconds()) {
			t.Errorf("MaxAge should be 30 days, got %d", cookie.MaxAge)
		}
	})

	t.Run("Test ValidateCookie", func(t *testing.T) {
		userID := authService.GenerateUserID()
		cookie := authService.CreateSignedCookie(userID)
		extractedUserID, valid := authService.ValidateCookie(cookie.Value)

		if !valid {
			t.Error("Cookie should be valid")
		}
		if extractedUserID != userID {
			t.Errorf("Extracted user ID should be %s, got %s", userID, extractedUserID)
		}
	})

	t.Run("Test Invalid Cookie", func(t *testing.T) {
		_, valid := authService.ValidateCookie("invalid-cookie")
		if valid {
			t.Error("Invalid cookie should not be valid")
		}
	})

	t.Run("Test GetOrCreateUserID with valid cookie", func(t *testing.T) {
		userID := authService.GenerateUserID()
		cookie := authService.CreateSignedCookie(userID)
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(cookie)

		extractedUserID, newCookie := authService.GetOrCreateUserID(req)
		if extractedUserID != userID {
			t.Errorf("Should return existing user ID %s, got %s", userID, extractedUserID)
		}
		if newCookie != nil {
			t.Error("Should not create new cookie for existing valid user")
		}
	})

	t.Run("Test GetOrCreateUserID without cookie", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		newUserID, newCookie := authService.GetOrCreateUserID(req)

		if newUserID == "" {
			t.Error("New user ID should not be empty")
		}
		if newCookie == nil {
			t.Error("New cookie should be created when no valid cookie exists")
		}
	})
}

// TestSignValue проверяет корректность работы подписи значений
func TestSignValue(t *testing.T) {
	authService := New("test-secret-key")

	// Тест 1: Проверка идентичности подписей для одинаковых значений
	t.Run("Test signature consistency", func(t *testing.T) {
		value := "test-value"
		signature1 := authService.SignValue(value)
		signature2 := authService.SignValue(value)

		if signature1 != signature2 {
			t.Error("Signatures for the same value should be identical")
		}

		if len(signature1) == 0 {
			t.Error("Signature should not be empty")
		}
	})

	// Тест 2: Проверка уникальности подписей для разных значений
	t.Run("Test unique signatures", func(t *testing.T) {
		value1 := "test-value"
		value2 := "another-test-value"
		sig1 := authService.SignValue(value1)
		sig2 := authService.SignValue(value2)

		if sig1 == sig2 {
			t.Error("Signatures for different values should be different")
		}
	})

	// Тест 3: Проверка подписи для пустого значения
	t.Run("Test empty value", func(t *testing.T) {
		emptySig := authService.SignValue("")
		if emptySig != "" {
			t.Error("Signature for empty value should be empty")
		}
	})
}

// TestCookieFormat проверяет формат и структуру создаваемых кук
func TestCookieFormat(t *testing.T) {
	authService := New("test-secret-key")

	// Тест 1: Проверка корректного формата куки
	t.Run("Test valid cookie format", func(t *testing.T) {
		userID := "test-user-id"
		cookie := authService.CreateSignedCookie(userID)

		if cookie == nil {
			t.Error("Cookie should not be nil")
			return
		}

		// Проверяем структуру куки
		if cookie.Name != "user_id" {
			t.Errorf("Cookie name should be 'user_id', got %s", cookie.Name)
		}
		if !cookie.HttpOnly {
			t.Error("HttpOnly should be true")
		}
		if cookie.SameSite != http.SameSiteLaxMode {
			t.Error("SameSite should be Lax")
		}
		if cookie.MaxAge != int(30*24*time.Hour.Seconds()) {
			t.Errorf("MaxAge should be 30 days, got %d", cookie.MaxAge)
		}

		// Проверяем формат значения куки: "user_id:signature"
		parts := strings.Split(cookie.Value, ":")
		if len(parts) != 2 {
			t.Errorf("Cookie value should have 2 parts separated by ':', got %d parts", len(parts))
		}
		if parts[0] != userID {
			t.Errorf("First part should be user ID %s, got %s", userID, parts[0])
		}
		expectedSignature := authService.SignValue(userID)
		if parts[1] != expectedSignature {
			t.Errorf("Second part should be signature %s, got %s", expectedSignature, parts[1])
		}
	})

	// Тест 2: Проверка создания куки с пустым userID
	t.Run("Test empty userID", func(t *testing.T) {
		cookie := authService.CreateSignedCookie("")
		if cookie != nil {
			t.Error("Cookie should be nil for empty userID")
		}
	})
}

// TestValidateCookie проверяет различные сценарии валидации кук
func TestValidateCookie(t *testing.T) {
	authService := New("test-secret-key")

	// Тест 1: Валидная кука с правильным форматом
	t.Run("Test valid cookie", func(t *testing.T) {
		userID := "test-user-id"
		cookieValue := authService.CreateSignedCookie(userID).Value

		extractedID, valid := authService.ValidateCookie(cookieValue)
		if !valid {
			t.Error("Валидная кука должна быть признана действительной")
		}
		if extractedID != userID {
			t.Errorf("Ожидался ID: %s, получен: %s", userID, extractedID)
		}
	})

	// Тест 2: Пустая строка
	t.Run("Test empty string", func(t *testing.T) {
		_, valid := authService.ValidateCookie("")
		if valid {
			t.Error("Пустая строка должна быть признана недействительной")
		}
	})

	// Тест 3: Некорректный формат (нет разделителя)
	t.Run("Test invalid format no delimiter", func(t *testing.T) {
		_, valid := authService.ValidateCookie("test-user-id")
		if valid {
			t.Error("Кука без разделителя должна быть признана недействительной")
		}
	})

	// Тест 4: Некорректный формат (лишние разделители)
	t.Run("Test invalid format extra delimiter", func(t *testing.T) {
		_, valid := authService.ValidateCookie("test-user-id:signature:extra")
		if valid {
			t.Error("Кука с лишними разделителями должна быть признана недействительной")
		}
	})

	// Тест 5: Неверная подпись
	t.Run("Test invalid signature", func(t *testing.T) {
		userID := "test-user-id"
		// Создаем корректную куку
		cookieValue := authService.CreateSignedCookie(userID).Value
		// Изменяем подпись
		parts := strings.Split(cookieValue, ":")
		parts[1] = "invalid-signature"
		modifiedCookie := strings.Join(parts, ":")

		_, valid := authService.ValidateCookie(modifiedCookie)
		if valid {
			t.Error("Кука с неверной подписью должна быть признана недействительной")
		}
	})

	// Тест 6: Невалидный hex-формат подписи
	t.Run("Test invalid hex signature", func(t *testing.T) {
		userID := "test-user-id"
		// Создаем корректную куку
		cookieValue := authService.CreateSignedCookie(userID).Value
		// Заменяем подпись на невалидный hex
		parts := strings.Split(cookieValue, ":")
		parts[1] = "invalid-hex-characters"
		modifiedCookie := strings.Join(parts, ":")

		_, valid := authService.ValidateCookie(modifiedCookie)
		if valid {
			t.Error("Кука с невалидным hex должен быть признана недействительной")
		}
	})

	// Тест 7: Слишком короткий userID (граничный случай)
	t.Run("Test short userID", func(t *testing.T) {
		userID := "a" // минимально возможный userID
		cookieValue := authService.CreateSignedCookie(userID).Value
		extractedID, valid := authService.ValidateCookie(cookieValue)
		if !valid {
			t.Error("Кука с коротким userID должна быть валидной")
		}
		if extractedID != userID {
			t.Errorf("Ожидался ID: %s, получен: %s", userID, extractedID)
		}
	})
}
