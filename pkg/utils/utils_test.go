package utils

import (
	"regexp"
	"testing"
)

// Проверка длины и формата сгенерированной ссылки
func TestGenerateShortKey(t *testing.T) {
	length := 6
	key1 := GenerateShortKey(length)
	key2 := GenerateShortKey(length)

	if len(key1) != length {
		t.Fatalf("ожидалась длина %v, получено %d: %s", length, len(key1), key1)
	}

	// Ключ содержит только разрешённые символы
	validChars := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validChars.MatchString(key1) {
		t.Errorf("ключ содержит недопустимые символы: %s", key1)
	}

	// coalises
	if key1 == key2 {
		t.Errorf("coalese %v == %v", key1, key2)
	}
}
