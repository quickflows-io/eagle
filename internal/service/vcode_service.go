package service

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-eagle/eagle/internal/repository"

	"github.com/pkg/errors"

	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/redis"
)

// Verification code service, mainly provides verification code generation and verification code acquisition
const (
	verifyCodeRedisKey = "app:login:vcode:%d" // verification code key
	maxDurationTime    = 10 * time.Minute     // Verification code validity period
)

// VCodeService define interface func
type VCodeService interface {
	GenLoginVCode(phone string) (int, error)
	CheckLoginVCode(phone int64, vCode int) bool
	GetLoginVCode(phone int64) (int, error)
}

type vcodeService struct {
	repo repository.Repository
}

var _ VCodeService = (*vcodeService)(nil)

func newVCode(svc *service) *vcodeService {
	return &vcodeService{repo: svc.repo}
}

// GenLoginVCode Generate check code
func (s *vcodeService) GenLoginVCode(phone string) (int, error) {
	// step1: generate random numbers
	vCodeStr := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))

	// step2: write to redis
	// use set, key use prefix + mobile phone number, cache for 10 minutes)
	key := fmt.Sprintf("app:login:vcode:%s", phone)
	err := redis.RedisClient.Set(context.Background(), key, vCodeStr, maxDurationTime).Err()
	if err != nil {
		return 0, errors.Wrap(err, "gen login code from redis set err")
	}

	vCode, err := strconv.Atoi(vCodeStr)
	if err != nil {
		return 0, errors.Wrap(err, "string convert int err")
	}

	return vCode, nil
}

// Mobile whitelist
var phoneWhiteLit = []int64{
	13010102020,
}

// isTestPhone Here you can add the test number, directly through
func isTestPhone(phone int64) bool {
	for _, val := range phoneWhiteLit {
		if val == phone {
			return true
		}
	}
	return false
}

// CheckLoginVCode Verify that the check code is correct
func (s *vcodeService) CheckLoginVCode(phone int64, vCode int) bool {
	if isTestPhone(phone) {
		return true
	}

	oldVCode, err := s.GetLoginVCode(phone)
	if err != nil {
		log.Warnf("[vcode_service] get verify code err, %v", err)
		return false
	}

	if vCode != oldVCode {
		return false
	}

	return true
}

// GetLoginVCode get check code
func (s *vcodeService) GetLoginVCode(phone int64) (int, error) {
	// Get it directly from redis
	key := fmt.Sprintf(verifyCodeRedisKey, phone)
	vcode, err := redis.RedisClient.Get(context.Background(), key).Result()
	if err == redis.ErrRedisNotFound {
		return 0, nil
	} else if err != nil {
		return 0, errors.Wrap(err, "redis get login vcode err")
	}

	verifyCode, err := strconv.Atoi(vcode)
	if err != nil {
		return 0, errors.Wrap(err, "strconv err")
	}

	return verifyCode, nil
}
