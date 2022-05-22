package cache

import (
	"context"
	"errors"
	"time"
)

var (
	// DefaultExpireTime Default expiration time
	DefaultExpireTime = time.Hour * 24
	// DefaultNotFoundExpireTime The expiration time when the result is empty is 1 minute, commonly used for the
	// cache time when the data is empty (cache penetration)
	DefaultNotFoundExpireTime = time.Minute
	// NotFoundPlaceholder .
	NotFoundPlaceholder = "*"

	// DefaultClient Generate a cache client, where keyPrefix is generally a business prefix
	DefaultClient Cache

	// ErrPlaceholder .
	ErrPlaceholder = errors.New("cache: placeholder")
	// ErrSetMemoryWithNotFound .
	ErrSetMemoryWithNotFound = errors.New("cache: set memory cache err for not found")
)

// Cache Define the cache driver interface
type Cache interface {
	Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, val interface{}) error
	MultiSet(ctx context.Context, valMap map[string]interface{}, expiration time.Duration) error
	MultiGet(ctx context.Context, keys []string, valueMap interface{}) error
	Del(ctx context.Context, keys ...string) error
	SetCacheWithNotFound(ctx context.Context, key string) error
}

// Set data
func Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	return DefaultClient.Set(ctx, key, val, expiration)
}

// Get data
func Get(ctx context.Context, key string, val interface{}) error {
	return DefaultClient.Get(ctx, key, val)
}

// MultiSet batch set
func MultiSet(ctx context.Context, valMap map[string]interface{}, expiration time.Duration) error {
	return DefaultClient.MultiSet(ctx, valMap, expiration)
}

// MultiGet Bulk acquisition
func MultiGet(ctx context.Context, keys []string, valueMap interface{}) error {
	return DefaultClient.MultiGet(ctx, keys, valueMap)
}

// Del batch deletion
func Del(ctx context.Context, keys ...string) error {
	return DefaultClient.Del(ctx, keys...)
}

// SetCacheWithNotFound .
func SetCacheWithNotFound(ctx context.Context, key string) error {
	return DefaultClient.SetCacheWithNotFound(ctx, key)
}
