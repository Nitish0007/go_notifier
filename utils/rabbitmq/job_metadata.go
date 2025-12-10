package rabbitmq_utils

import (
	"context"
	"fmt"
	"time"

	redis_utils "github.com/Nitish0007/go_notifier/utils/redis"
	"github.com/redis/go-redis/v9"
)

type JobMetadata struct {
	RetryCount int `json:"retry_count"`
	MaxRetries int `json:"max_retries"`
	RetryDelay time.Duration `json:"retry_delay"`
}

func NewJobMetadata(retryCount int, maxRetries int, retryDelay time.Duration) *JobMetadata {
	return &JobMetadata{
		RetryCount: retryCount,
		MaxRetries: maxRetries,
		RetryDelay: retryDelay,
	}
}
func StoreJobMetadata(ctx context.Context, jobID string, metadata JobMetadata) error {
	key := fmt.Sprintf("jmd:%s", jobID)
	err := redis_utils.SetRedisJSON(ctx, key, metadata, 0)
	if err != nil {
		return fmt.Errorf("failed to store job metadata: %w", err)
	}
	return nil
}

func GetJobMetadata(ctx context.Context, jobID string) (*JobMetadata, error) {
	key := fmt.Sprintf("jmd:%s", jobID)
	metadata, err := redis_utils.GetRedisJSON(ctx, key)
	if err == redis.Nil {
		// create new metadata
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get job metadata: %w", err)
	}
	return metadata.(*JobMetadata), nil
}

func DeleteJobMetadata(ctx context.Context, jobID string) error {
	key := fmt.Sprintf("jmd:%s", jobID)
	err := redis_utils.DeleteRedisJSON(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete job metadata: %w", err)
	}
	return nil
}