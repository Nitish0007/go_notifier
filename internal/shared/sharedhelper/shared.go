package sharedhelper

import (
	"context"
	"math/rand"
	"regexp"
	"time"
	"github.com/google/uuid"
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

func IsValidUUID(id string) bool {
	if id == "" {
		return false
	}
	if _, err := uuid.Parse(id); err != nil {
		return false
	}
	return true
}

func GetCurrentAccountID(ctx context.Context) int64 {
	accountID, ok := ctx.Value(CurrentAccountID).(int64)
	if !ok {
		return -1 // or handle the error as needed
	}
	return accountID
}

func SetCurrentAccountID(ctx context.Context, accountID int64) context.Context {
	return context.WithValue(ctx, CurrentAccountID, accountID)
}

func ParseTime(timeStr string, layout string) *time.Time {
	if layout == "" {
		layout = time.RFC3339
	}
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return nil
	}
	return &t
}

func GetCurrentTime() *time.Time {
	t := time.Now().UTC()
	return &t
}

func SetValueToContext(ctx context.Context, key string, value any) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetValueFromContext(ctx context.Context, key string) any {
	return ctx.Value(key)
}

// NOTE: comparable only support equality operators not ordering operators
func GetUniques[T comparable](list []T) []T {
	seen := make(map[T]struct{}, len(list))
	out := make([]T, 0, len(list))
	for _, v := range list {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}