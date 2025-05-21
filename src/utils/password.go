package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

func GenerateRandomPassword(length int) (string, error) {
	if length < 8 {
		length = 8
	}

	bytes := make([]byte, length+8)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	password := base64.StdEncoding.EncodeToString(bytes)
	password = strings.ReplaceAll(password, "+", "")
	password = strings.ReplaceAll(password, "/", "")
	password = strings.ReplaceAll(password, "=", "")

	// Обмежуємо довжину
	if len(password) > length {
		password = password[:length]
	}

	return password, nil
}
