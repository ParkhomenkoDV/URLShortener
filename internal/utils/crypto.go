package utils

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateShortURL(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic("failed to generate random bytes: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}
