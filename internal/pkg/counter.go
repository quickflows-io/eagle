// Package pkg counter, which can be used for statistical use of various models of the business
// Scenario: Commonly used for repeated strategies, or anti-cheating processing control package pkg
package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	redis2 "github.com/go-eagle/eagle/pkg/redis"
)

const (
	// PrefixCounter counter key
	PrefixCounter = "eagle:counter:%s"
	// DefaultStep default step key
	DefaultStep = 1
	// DefaultExpirationTime .
	DefaultExpirationTime = 600 * time.Second
)

// Counter define struct
type Counter struct {
	client *redis.Client
}

// NewCounter create a counter
func NewCounter() *Counter {
	return &Counter{
		client: redis2.RedisClient,
	}
}

// GetKey 获取key
func (c *Counter) GetKey(key string) string {
	return fmt.Sprintf(PrefixCounter, key)
}

// SetCounter set counter
func (c *Counter) SetCounter(ctx context.Context, idStr string, expiration time.Duration) (int64, error) {
	key := c.GetKey(idStr)
	ret, err := c.client.IncrBy(ctx, key, DefaultStep).Result()
	if err != nil {
		return 0, err
	}
	_, _ = c.client.Expire(ctx, key, expiration).Result()
	return ret, nil
}

// GetCounter get total count
func (c *Counter) GetCounter(ctx context.Context, idStr string) (int64, error) {
	key := c.GetKey(idStr)
	return c.client.Get(ctx, key).Int64()
}

// DelCounter del count
func (c *Counter) DelCounter(ctx context.Context, idStr string) int64 {
	key := c.GetKey(idStr)
	var keys []string
	keys = append(keys, key)
	return c.client.Del(ctx, keys...).Val()
}
