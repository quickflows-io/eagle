package repository

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"

	"github.com/go-eagle/eagle/internal/model"
	"github.com/go-eagle/eagle/pkg/cache"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/redis"
)

var g singleflight.Group

// CreateUser create user
func (d *repository) CreateUser(ctx context.Context, user *model.UserBaseModel) (id uint64, err error) {
	err = d.orm.Create(&user).Error
	if err != nil {
		//prom.BusinessErrCount.Incr("mysql: CreateUser")
		return 0, errors.Wrap(err, "[repo.user_base] create user err")
	}

	return user.ID, nil
}

// UpdateUser Update user information
func (d *repository) UpdateUser(ctx context.Context, id uint64, userMap map[string]interface{}) error {
	user, err := d.GetUser(ctx, id)
	if err != nil {
		//prom.BusinessErrCount.Incr("mysql: getOneUser")
		return errors.Wrap(err, "[repo.user_base] update user data err")
	}

	// delete cache
	err = d.userCache.DelUserBaseCache(ctx, id)
	if err != nil {
		log.Warnf("[repo.user_base] delete user cache err: %v", err)
	}

	err = d.orm.Model(user).Updates(userMap).Error
	if err != nil {
		//prom.BusinessErrCount.Incr("mysql: UpdateUser")
	}
	return err
}

// GetUser get users
// The update strategy of the cache uses the Cache Aside Pattern
// see: https://coolshell.cn/articles/17416.html
func (d *repository) GetUser(ctx context.Context, uid uint64) (userBase *model.UserBaseModel, err error) {
	ctx, span := d.tracer.Start(ctx, "GetUser", oteltrace.WithAttributes(
		attribute.String("param.uid", cast.ToString(uid)),
	))
	defer span.End()

	var data *model.UserBaseModel

	userBase, err = d.userCache.GetUserBaseCache(ctx, uid)
	if errors.Is(err, cache.ErrPlaceholder) {
		span.SetName("eq ErrPlaceholder")
		span.RecordError(err)
		return nil, ErrNotFound
	} else if errors.Is(err, redis.ErrRedisNotFound) {
		// use sync/singleflight mode to get data
		// demo see: https://github.com/go-demo/singleflight-demo/blob/master/main.go
		// https://juejin.cn/post/6844904084445593613
		key := fmt.Sprintf("get_user_base_%d", uid)
		val, err, _ := g.Do(key, func() (interface{}, error) {
			data := new(model.UserBaseModel)
			// ?????????????????????
			err = d.orm.WithContext(ctx).First(data, uid).Error
			// if data is empty, set not found cache to prevent cache penetration(????????????)
			if errors.Is(err, ErrNotFound) {
				err = d.userCache.SetCacheWithNotFound(ctx, uid)
				if err != nil {
					span.SetName("SetCacheWithNotFound")
					span.RecordError(err)
					log.Warnf("[repo.user_base] SetCacheWithNotFound err, uid: %d", uid)
				}
				return nil, ErrNotFound
			} else if err != nil {
				span.SetName("find from db err")
				span.RecordError(err)
				//prom.BusinessErrCount.Incr("mysql: getOneUser")
				return nil, errors.Wrapf(err, "[repo.user_base] query db err")
			}

			// set cache
			err = d.userCache.SetUserBaseCache(ctx, uid, data, cache.DefaultExpireTime)
			if err != nil {
				span.SetName("SetUserBaseCache")
				span.RecordError(err)
				return nil, errors.Wrap(err, "[repo.user_base] SetUserBaseCache err")
			}
			return data, nil
		})

		if err != nil && err != ErrNotFound {
			span.SetName("sg.do")
			span.RecordError(err)
			return nil, errors.Wrap(err, "[repo.user_base] get user base err via single flight do")
		}
		if val != nil {
			data = val.(*model.UserBaseModel)
		}
	} else if err != nil {
		// fail fast, if cache error return, don't request to db
		return nil, err
	}

	// cache hit
	if userBase != nil {
		//prom.CacheHit.Incr("getOneUser")
		log.WithContext(ctx).Warnf("[repo.user_base] get user base data from cache, uid: %d", uid)
		return
	}

	// cache miss
	//prom.CacheMiss.Incr("getOneUser")

	return data, nil
}

// GetUsersByIds Get users in batches
func (d *repository) GetUsersByIds(ctx context.Context, userIDs []uint64) ([]*model.UserBaseModel, error) {
	users := make([]*model.UserBaseModel, 0)

	// Bulk fetch from cache
	userCacheMap, err := d.userCache.MultiGetUserBaseCache(ctx, userIDs)
	if err != nil {
		return users, errors.Wrap(err, "[repo.user_base] multi get user cache data err")
	}

	// query miss
	for _, userID := range userIDs {
		idx := d.userCache.GetUserBaseCacheKey(userID)
		userModel, ok := userCacheMap[idx]
		if !ok {
			userModel, err = d.GetUser(ctx, userID)
			if err != nil {
				log.Warnf("[repo.user_base] get user model err: %v", err)
				continue
			}
		}
		users = append(users, userModel)
	}
	return users, nil
}

// GetUserByPhone Get users by phone number
func (d *repository) GetUserByPhone(ctx context.Context, phone int64) (*model.UserBaseModel, error) {
	user := model.UserBaseModel{}
	err := d.orm.Where("phone = ?", phone).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "[repo.user_base] get user err by phone")
	}

	return &user, nil
}

// GetUserByEmail Get mobile number by email
func (d *repository) GetUserByEmail(ctx context.Context, email string) (*model.UserBaseModel, error) {
	userBase := model.UserBaseModel{}
	err := d.orm.Where("email = ?", email).First(&userBase).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.Wrap(err, "[repo.user_base] get user err by email")
	}

	return &userBase, nil
}

// UserIsExist To determine whether the user exists, the username and email address must be unique
func (d *repository) UserIsExist(user *model.UserBaseModel) (bool, error) {
	err := d.orm.Where("username = ? or email = ?", user.Username, user.Email).First(&model.UserBaseModel{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
