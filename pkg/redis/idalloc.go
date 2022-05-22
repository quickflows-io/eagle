package redis

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"

	"github.com/go-eagle/eagle/pkg/log"
)

// IDAlloc id generator
// key is the business key, which consists of business prefix + function prefix + specific scene id
// The key is the business key, which is composed of business prefix + function prefix + specific scene id.
//For example, to generate user id, you can pass in user_id. Complete example: eagle:idalloc:user_id
type IDAlloc struct {
	// Redis instance, it is best to use a business-independent instance,
	//it is best to deploy a cluster to make id alloc highly available
	redisClient *redis.Client
}

// NewIDAlloc create a id alloc instance
func NewIDAlloc(conn *redis.Client) *IDAlloc {
	return &IDAlloc{
		redisClient: conn,
	}
}

// GetNewID generate id
func (ia *IDAlloc) GetNewID(key string, step int64) (int64, error) {
	key = ia.GetKey(key)
	id, err := ia.redisClient.IncrBy(context.Background(), key, step).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "redis incr err, key: %s", key)
	}

	if id == 0 {
		log.Warnf("[redis.idalloc] %s GetNewID failed", key)
		return 0, errors.Wrapf(err, "[redis.idalloc] %s GetNewID failed", key)
	}
	return id, nil
}

// GetCurrentID get current id
func (ia *IDAlloc) GetCurrentID(key string) (int64, error) {
	key = ia.GetKey(key)
	ret, err := ia.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return 0, errors.Wrapf(err, "redis get err, key: %s", key)
	}
	id, err := strconv.Atoi(ret)
	if err != nil {
		return 0, errors.Wrap(err, "str convert err")
	}
	return int64(id), nil
}

// GetKey get key
func (ia *IDAlloc) GetKey(key string) string {
	lockKey := "idalloc"
	return strings.Join([]string{lockKey, key}, ":")
}
