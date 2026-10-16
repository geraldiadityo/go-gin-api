package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type BaseCache struct {
	rdb       *redis.Client
	namespace string
	ttl       time.Duration
}

func NewBaseCache(rdb *redis.Client, namespace string, defaultTTL time.Duration) *BaseCache {
	return &BaseCache{
		rdb:       rdb,
		namespace: namespace,
		ttl:       defaultTTL,
	}
}

func (c *BaseCache) FormatKey(subKey string) string {
	return fmt.Sprintf("%s:%s", c.namespace, subKey)
}

func (c *BaseCache) SetData(ctx context.Context, subKey string, value any, customTTL ...time.Duration) error {
	ttl := c.ttl
	if len(customTTL) > 0 && customTTL[0] > 0 {
		ttl = customTTL[0]
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return c.rdb.Set(ctx, c.FormatKey(subKey), bytes, ttl).Err()
}

func (c *BaseCache) InvalidNamespace(ctx context.Context) error {
	pattern := fmt.Sprintf("%s:*", c.namespace)

	var cursor uint64
	for {
		keys, nextCursor, err := c.rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

func (c *BaseCache) GetData(ctx context.Context, subKey string, dest any) (bool, error) {
	val, err := c.rdb.Get(ctx, c.FormatKey(subKey)).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(val), dest)
	return true, err
}

func FetchWithCache[T any](ctx context.Context, c *BaseCache, subKey string, fetchDB func() (T, error), customeTTL ...time.Duration) (T, error) {
	fullKey := c.FormatKey(subKey)
	var result T
	hit, err := c.GetData(ctx, subKey, &result)
	if err == nil && hit {
		log.Printf("\033[32m[CACHE HIT]\033[0m Mengambil data dari cache redis (key: %s)", fullKey)
		return result, nil
	}

	log.Printf("\033[33m[CACHE MISS]\033[0 Mengambil data dari database (key: %s)", fullKey)

	data, err := fetchDB()
	if err != nil {
		return result, err
	}

	_ = c.SetData(ctx, subKey, data, customeTTL...)

	return data, nil
}
