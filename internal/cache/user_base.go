package cache

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/pkg/cache"
	"github.com/go-eagle/eagle/pkg/encoding"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/redis"
)

const (
	// PrefixUserBaseCacheKey cache prefix, rule: business+module+{ID}
	PrefixUserBaseCacheKey = "eagle:user:base:%d"
)

// Cache cache
type Cache struct {
	cache  cache.Cache
	tracer trace.Tracer
	//localCache cache.Cache
}

// NewUserCache new a user cache
func NewUserCache() *Cache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	return &Cache{
		cache: cache.NewRedisCache(redis.RedisClient, cachePrefix, jsonEncoding, func() interface{} {
			return &model.UserBaseModel{}
		},
		),
		tracer: otel.Tracer("user cache"),
	}
}

// GetUserBaseCacheKey get cache key
func (c *Cache) GetUserBaseCacheKey(userID uint64) string {
	return fmt.Sprintf(PrefixUserBaseCacheKey, userID)
}

// SetUserBaseCache write to user cache
func (c *Cache) SetUserBaseCache(ctx context.Context, userID uint64, user *model.UserBaseModel, duration time.Duration) error {
	ctx, span := c.tracer.Start(ctx, "SetUserBaseCache")
	defer span.End()

	if user == nil || user.ID == 0 {
		return nil
	}
	cacheKey := c.GetUserBaseCacheKey(userID)
	err := c.cache.Set(ctx, cacheKey, user, duration)
	if err != nil {
		return err
	}
	return nil
}

// GetUserBaseCache Get user cache
func (c *Cache) GetUserBaseCache(ctx context.Context, userID uint64) (data *model.UserBaseModel, err error) {
	ctx, span := c.tracer.Start(ctx, "GetUserBaseCache")
	defer span.End()

	client := getCacheClient(ctx)

	cacheKey := c.GetUserBaseCacheKey(userID)
	//err = u.cache.Get(cacheKey, &data)
	err = client.Get(ctx, cacheKey, &data)
	if err != nil {
		log.WithContext(ctx).Warnf("get err from redis, err: %+v", err)
		return nil, err
	}
	return data, nil
}

// MultiGetUserBaseCache Get user caches in batches
func (c *Cache) MultiGetUserBaseCache(ctx context.Context, userIDs []uint64) (map[string]*model.UserBaseModel, error) {
	ctx, span := c.tracer.Start(ctx, "MultiGetUserBaseCache")
	defer span.End()
	var keys []string
	for _, v := range userIDs {
		cacheKey := c.GetUserBaseCacheKey(v)
		keys = append(keys, cacheKey)
	}

	// Need to instantiate make here, if it is defined directly in the return parameter, it will report nil map
	userMap := make(map[string]*model.UserBaseModel)
	err := c.cache.MultiGet(ctx, keys, userMap)
	if err != nil {
		return nil, err
	}
	return userMap, nil
}

// DelUserBaseCache delete user cache
func (c *Cache) DelUserBaseCache(ctx context.Context, userID uint64) error {
	ctx, span := c.tracer.Start(ctx, "DelUserBaseCache")
	defer span.End()
	cacheKey := c.GetUserBaseCacheKey(userID)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound .
func (c *Cache) SetCacheWithNotFound(ctx context.Context, userID uint64) error {
	ctx, span := c.tracer.Start(ctx, "SetCacheWithNotFound")
	defer span.End()
	cacheKey := c.GetUserBaseCacheKey(userID)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
