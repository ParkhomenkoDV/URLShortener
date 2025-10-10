package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// Генерация короткого ключа заданной длины
func GenerateShortKey(length int) string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate random bytes: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}
