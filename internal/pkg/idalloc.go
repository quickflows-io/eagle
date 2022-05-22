// Package pkg ID distributor, mainly using redis for distribution
package pkg

import "github.com/go-eagle/eagle/pkg/redis"

// IDAlloc define struct
type IDAlloc struct {
	idGenerator *redis.IDAlloc
}

// NewIDAlloc create a id alloc
func NewIDAlloc() *IDAlloc {
	return &IDAlloc{
		idGenerator: redis.NewIDAlloc(redis.RedisClient),
	}
}

// GetUserID generate user id from redis
func (i *IDAlloc) GetUserID() (int64, error) {
	return i.idGenerator.GetNewID("user_id", 1)
}
