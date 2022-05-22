package cache

import (
	"context"
	"reflect"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/go-eagle/eagle/pkg/encoding"
	"github.com/go-eagle/eagle/pkg/log"
)

// redisCache redis cache structure
type redisCache struct {
	client            *redis.Client
	KeyPrefix         string
	encoding          encoding.Encoding
	DefaultExpireTime time.Duration
	newObject         func() interface{}
}

// NewRedisCache A new cache, the client parameter can be passed in, which is convenient for unit testing
func NewRedisCache(client *redis.Client, keyPrefix string, encoding encoding.Encoding, newObject func() interface{}) Cache {
	return &redisCache{
		client:    client,
		KeyPrefix: keyPrefix,
		encoding:  encoding,
		newObject: newObject,
	}
}

func (c *redisCache) Set(ctx context.Context, key string, val interface{}, expiration time.Duration) error {
	buf, err := encoding.Marshal(c.encoding, val)
	if err != nil {
		return errors.Wrapf(err, "marshal data err, value is %+v", val)
	}

	cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
	if err != nil {
		return errors.Wrapf(err, "build cache key err, key is %+v", key)
	}
	if expiration == 0 {
		expiration = DefaultExpireTime
	}
	err = c.client.Set(ctx, cacheKey, buf, expiration).Err()
	if err != nil {
		return errors.Wrapf(err, "redis set err: %+v", err)
	}
	return nil
}

func (c *redisCache) Get(ctx context.Context, key string, val interface{}) error {
	cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
	if err != nil {
		return errors.Wrapf(err, "build cache key err, key is %+v", key)
	}

	bytes, err := c.client.Get(ctx, cacheKey).Bytes()
	if err != nil {
		return err
	}

	// Prevent Unmarshal from reporting an error when data is empty
	if string(bytes) == "" {
		return nil
	}
	if string(bytes) == NotFoundPlaceholder {
		return ErrPlaceholder
	}
	err = encoding.Unmarshal(c.encoding, bytes, val)
	if err != nil {
		return errors.Wrapf(err, "unmarshal data error, key=%s, cacheKey=%s type=%v, json is %+v ",
			key, cacheKey, reflect.TypeOf(val), string(bytes))
	}
	return nil
}

func (c *redisCache) MultiSet(ctx context.Context, valueMap map[string]interface{}, expiration time.Duration) error {
	if len(valueMap) == 0 {
		return nil
	}
	if expiration == 0 {
		expiration = DefaultExpireTime
	}
	// The key-value is paired, so the capacity here is twice that of the map
	paris := make([]interface{}, 0, 2*len(valueMap))
	for key, value := range valueMap {
		buf, err := encoding.Marshal(c.encoding, value)
		if err != nil {
			log.Warnf("marshal data err: %+v, value is %+v", err, value)
			continue
		}
		cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
		if err != nil {
			log.Warnf("build cache key err: %+v, key is %+v", err, key)
			continue
		}
		paris = append(paris, []byte(cacheKey))
		paris = append(paris, buf)
	}
	pipeline := c.client.Pipeline()
	err := pipeline.MSet(ctx, paris...).Err()
	if err != nil {
		return errors.Wrapf(err, "redis multi set error")
	}
	for i := 0; i < len(paris); i = i + 2 {
		switch paris[i].(type) {
		case []byte:
			pipeline.Expire(ctx, string(paris[i].([]byte)), expiration)
		default:
			log.Warnf("redis expire is unsupported key type: %+v", reflect.TypeOf(paris[i]))
		}
	}
	_, err = pipeline.Exec(ctx)
	if err != nil {
		return errors.Wrapf(err, "redis multi set pipeline exec error")
	}
	return nil
}

func (c *redisCache) MultiGet(ctx context.Context, keys []string, value interface{}) error {
	if len(keys) == 0 {
		return nil
	}
	cacheKeys := make([]string, len(keys))
	for index, key := range keys {
		cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
		if err != nil {
			return errors.Wrapf(err, "build cache key err, key is %+v", key)
		}
		cacheKeys[index] = cacheKey
	}
	values, err := c.client.MGet(ctx, cacheKeys...).Result()
	if err != nil {
		return errors.Wrapf(err, "redis MGet error, keys is %+v", keys)
	}

	// Injected into map via reflection
	valueMap := reflect.ValueOf(value)
	for i, value := range values {
		if value == nil {
			continue
		}
		object := c.newObject()
		err = encoding.Unmarshal(c.encoding, []byte(value.(string)), object)
		if err != nil {
			log.Warnf("unmarshal data error: %+v, key=%s, cacheKey=%s type=%v", err,
				keys[i], cacheKeys[i], reflect.TypeOf(value))
			continue
		}
		valueMap.SetMapIndex(reflect.ValueOf(cacheKeys[i]), reflect.ValueOf(object))
	}
	return nil
}

func (c *redisCache) Del(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	// Batch build cacheKey
	cacheKeys := make([]string, len(keys))
	for index, key := range keys {
		cacheKey, err := BuildCacheKey(c.KeyPrefix, key)
		if err != nil {
			log.Warnf("build cache key err: %+v, key is %+v", err, key)
			continue
		}
		cacheKeys[index] = cacheKey
	}
	err := c.client.Del(ctx, cacheKeys...).Err()
	if err != nil {
		return errors.Wrapf(err, "redis delete error, keys is %+v", keys)
	}
	return nil
}

func (c *redisCache) SetCacheWithNotFound(ctx context.Context, key string) error {
	return c.client.Set(ctx, key, NotFoundPlaceholder, DefaultNotFoundExpireTime).Err()
}
