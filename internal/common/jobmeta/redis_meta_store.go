package jobmeta

import (
	"fmt"
	"context"
	"github.com/Nitish0007/go_notifier/internal/shared/sharedhelper"
	redis_client "github.com/Nitish0007/go_notifier/internal/common/redis"
)

type RedisMetaStore struct {
	redisCache *redis_client.RedisCache
}

func NewRedisMetaStore() *RedisMetaStore {
	return &RedisMetaStore{
		redisCache: redis_client.NewRedisCache(),
	}
}

func (s *RedisMetaStore) Put(ctx context.Context, jobID string, meta sharedhelper.JobMetadata) error {
	key := fmt.Sprintf("jmd:%s", jobID)
	return s.redisCache.Set(ctx, key, meta, 0)
}