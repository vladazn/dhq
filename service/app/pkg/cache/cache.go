package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, bool, error)
	Set(ctx context.Context, key string, value string) error
	Clear(ctx context.Context, key string) error
}

type CacheStruct struct {
	pfx string
	r   *redis.Client
}

func NewCacheService(r *redis.Client, prefix string) Cache {
	return CacheStruct{prefix, r}
}

func (c CacheStruct) key(key string) string {
	return fmt.Sprintf("%v:%v", c.pfx, key)
}

func (c CacheStruct) Get(ctx context.Context, key string) (string, bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	result, err := c.r.Get(ctx, c.key(key)).Result()
	if err == redis.Nil {
		return "", false, nil
	}

	return result, true, nil
}

func (c CacheStruct) Set(ctx context.Context, key string, value string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return c.r.Set(ctx, c.key(key), value, 5*time.Hour).Err()
}

func (c CacheStruct) Clear(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return c.r.Del(ctx, c.key(key)).Err()
}
