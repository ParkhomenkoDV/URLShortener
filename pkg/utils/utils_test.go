package utils

import (
	"regexp"
	"testing"
)

// TestGenerateShortKey тестирует функцию генерации коротких ключей.
// Проверяет корректность длины, формат символов и уникальность генерируемых ключей.
func TestGenerateShortKey(t *testing.T) {
	// Фиксированная длина для тестирования
	const length = 6

	// Тест 1: Проверка корректности длины сгенерированного ключа
	t.Run("Generated key has correct length", func(t *testing.T) {
		key := GenerateShortKey(length)

		if len(key) != length {
			t.Fatalf("Generated key length incorrect: len(%s)=%d, expected: %d",
				key, len(key), length)
		}
	})

	// Тест 2: Проверка что ключ содержит только разрешенные URL-безопасные символы
	t.Run("Generated key contains only URL-safe characters", func(t *testing.T) {
		key := GenerateShortKey(length)

		// Регулярное выражение для URL-безопасных символов base64 encoding:
		// A-Z, a-z, 0-9, дефис (-), подчеркивание (_)
		validChars := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
		if !validChars.MatchString(key) {
			t.Errorf("Generated key contains prohibited symbols: %s. "+
				"Only A-Z, a-z, 0-9, -, _ are allowed.", key)
		}
	})

	// Тест 3: Проверка уникальности сгенерированных ключей (коллизии)
	t.Run("Consecutively generated keys are unique", func(t *testing.T) {
		key1 := GenerateShortKey(length)
		key2 := GenerateShortKey(length)

		// Проверяем что два последовательно сгенерированных ключа различаются
		// Это вероятностный тест, но вероятность коллизии для 6 символов крайне мала
		if key1 == key2 {
			t.Errorf("Key collision detected: two consecutive calls returned identical keys: %s",
				key1)
		}
	})

	// Тест 4: Проверка генерации ключей разной длины
	t.Run("Generate keys of different lengths", func(t *testing.T) {
		lengths := []int{4, 6, 8}

		for _, l := range lengths {
			t.Run("Length test", func(t *testing.T) {
				key := GenerateShortKey(l)
				if len(key) != l {
					t.Errorf("Key length mismatch for length %d: got %d, expected %d", l, len(key), l)
				}
			})
		}
	})

	// Тест 5: Проверка множественной генерации на уникальность (статистический тест)
	t.Run("Multiple generated keys are unique", func(t *testing.T) {
		const numKeys = 1000
		keys := make(map[string]bool, numKeys)

		// Генерируем 1000 ключей и проверяем на уникальность
		for i := 0; i < numKeys; i++ {
			key := GenerateShortKey(length)

			// Проверяем что ключ еще не встречался
			if keys[key] {
				t.Errorf("Duplicate key found after %d generations: %s", i, key)
				break
			}
			keys[key] = true
		}

		// Проверяем что сгенерировано нужное количество уникальных ключей
		if len(keys) != numKeys {
			t.Errorf("Expected %d unique keys, got %d", numKeys, len(keys))
		}
	})

	// Тест 6: Проверка что ключи состоят из печатных ASCII символов
	t.Run("Generated keys consist of printable ASCII characters", func(t *testing.T) {
		key := GenerateShortKey(length)

		// Проверяем что все символы в диапазоне печатных ASCII (от 33 до 126)
		for i, char := range key {
			if char < 33 || char > 126 {
				t.Errorf("Key contains non-printable ASCII character at position %d: %q (code %d)",
					i, char, char)
			}
		}
	})
}

// BenchmarkGenerateShortKey измеряет производительность функции генерации ключей
func BenchmarkGenerateShortKey(b *testing.B) {
	lengths := []int{4, 6, 8}

	for _, length := range lengths {
		b.Run("Length test", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				GenerateShortKey(length)
			}
		})
	}
}
