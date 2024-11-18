package cache

import (
	"context"
	"time"
)

// RedisCache interface defines methods for redis cache operations
type RedisCache interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{}) error
	Get(ctx context.Context, key string) (interface{}, error)
	Expire(ctx context.Context, key string, expirationTime time.Duration) error
	HashSet(ctx context.Context, key string, value interface{}) error
	HGetAll(ctx context.Context, key string) ([]interface{}, error)
}
