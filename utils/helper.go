package utils

import (
	"context"
	"math/rand"
	"regexp"
	"time"
)

type contextKey string
const CurrentAccountID contextKey = "currentAccountID"

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

func IsValidChannelType(channel string) bool {
	switch channel {
	case "email": //, "sms", "in_app":
		return true
	default:
		return false
	}
}

func GetCurrentAccountID(ctx context.Context) int {
	accountID, ok := ctx.Value(CurrentAccountID).(int)
	if !ok {
		return -1 // or handle the error as needed
	}
	return accountID
}

func SetCurrentAccountID(ctx context.Context, accountID int) context.Context {
	return context.WithValue(ctx, CurrentAccountID, accountID)
}

func ParseTime(timeStr string) time.Time {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{} // return zero time if parsing fails
	}
	return t
}

func GetCurrentTime() time.Time {
	return time.Now().UTC()
}