package auth

import (
	"net/http/httptest"
	"strings"
	"testing"
)

// TestAuthService проверяет основную функциональность сервиса аутентификации
func TestAuthService(t *testing.T) {
	// Создаем сервис аутентификации
	authService := NewAuthService("test-secret-key")

	// Тест 1: Генерация пользователя
	userID := authService.GenerateUserID()
	if userID == "" {
		t.Error("Generated user ID should not be empty")
	}

	// Тест 2: Создание подписанной куки
	cookie := authService.CreateSignedCookie(userID)
	if cookie == nil {
		t.Error("Cookie should not be nil")
		return // Добавляем return чтобы избежать дальнейшего выполнения
	}
	if cookie.Name != "user_id" {
		t.Errorf("Cookie name should be 'user_id', got %s", cookie.Name)
	}

	// Тест 3: Валидация куки
	extractedUserID, valid := authService.ValidateCookie(cookie.Value)
	if !valid {
		t.Error("Cookie should be valid")
	}
	if extractedUserID != userID {
		t.Errorf("Extracted user ID should be %s, got %s", userID, extractedUserID)
	}

	// Тест 4: Невалидная кука
	_, valid = authService.ValidateCookie("invalid-cookie")
	if valid {
		t.Error("Invalid cookie should not be valid")
	}

	// Тест 5: Получение или создание пользователя с валидной кукой
	// Используем httptest для создания правильного запроса
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(cookie)

	extractedUserID2, newCookie := authService.GetOrCreateUserID(req)
	if extractedUserID2 != userID {
		t.Errorf("Should return existing user ID %s, got %s", userID, extractedUserID2)
	}
	if newCookie != nil {
		t.Error("Should not create new cookie for existing valid user")
	}

	// Тест 6: Создание нового пользователя при отсутствии куки
	reqWithoutCookie := httptest.NewRequest("GET", "/", nil)
	newUserID, newCookie2 := authService.GetOrCreateUserID(reqWithoutCookie)

	if newUserID == "" {
		t.Error("New user ID should not be empty")
	}
	if newCookie2 == nil {
		t.Error("New cookie should be created")
	}
	if newUserID == userID {
		t.Error("New user ID should be different from existing user ID")
	}
}

// TestSignValue проверяет подписание значений
func TestSignValue(t *testing.T) {
	authService := NewAuthService("test-secret-key")

	value := "test-value"
	signature1 := authService.SignValue(value)
	signature2 := authService.SignValue(value)

	if signature1 != signature2 {
		t.Error("Signatures for the same value should be identical")
	}

	if len(signature1) == 0 {
		t.Error("Signature should not be empty")
	}

	// Тест с другим значением
	differentSignature := authService.SignValue("different-value")
	if signature1 == differentSignature {
		t.Error("Signatures for different values should be different")
	}
}

// TestCookieFormat проверяет формат куки
func TestCookieFormat(t *testing.T) {
	authService := NewAuthService("test-secret-key")

	userID := "test-user-id"
	cookie := authService.CreateSignedCookie(userID)

	// Кука должна содержать userID и подпись, разделенные ":"
	parts := strings.Split(cookie.Value, ":")
	if len(parts) != 2 {
		t.Errorf("Cookie value should have 2 parts separated by ':', got %d parts", len(parts))
	}

	if parts[0] != userID {
		t.Errorf("First part should be user ID %s, got %s", userID, parts[0])
	}

	// Вторая часть должна быть подписью
	expectedSignature := authService.SignValue(userID)
	if parts[1] != expectedSignature {
		t.Errorf("Second part should be signature %s, got %s", expectedSignature, parts[1])
	}
}
