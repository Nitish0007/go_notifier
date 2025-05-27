package utils

import (
	"math/rand"
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