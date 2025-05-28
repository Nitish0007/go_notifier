package utils

import (
	"math/rand"
	"regexp"
	"time"
)

func GenerateAlphaNumericKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 32
	key := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range key {
		key[i] = charset[r.Intn(len(charset))]
	}
	return string(key)
}

func ValidateEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegex, email)
	return matched
}