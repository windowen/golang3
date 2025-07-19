package redis_lock

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"

	"serverApi/pkg/zlogger"
)

const (
	// 删除操作原子性
	delLockKeyScript = "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
	// 尝试获取锁定睡眠时间
	tryGetLockSleepTimes = 5 * time.Millisecond
	// 锁过期时间
	expireTime = 10 * time.Second
	// 获取锁超时时间
	acquireTimeout = 10 * time.Second
)

type RedisLock struct {
	rdb redis.UniversalClient
}

func NewRedisLock(redisClient redis.UniversalClient) *RedisLock {
	return &RedisLock{
		rdb: redisClient,
	}
}

// 获取锁
func (r *RedisLock) acquireLock(ctx context.Context, lockKey, value string, expireTime time.Duration) (bool, error) {
	result := r.rdb.SetNX(ctx, lockKey, value, expireTime)
	if result.Err() != nil {
		return false, result.Err()
	}

	return result.Val(), nil
}

// TryLock 尝试获取分布式锁
func (r *RedisLock) TryLock(ctx context.Context, lockKey, value string) bool {
	endTime := time.Now().Add(acquireTimeout)

	for time.Now().Before(endTime) {
		acquired, err := r.acquireLock(ctx, lockKey, value, expireTime)
		if err != nil {
			zlogger.Errorf("TryLock |lockKey:%v| err: %v", lockKey, err)
			return false
		}

		if acquired {
			return true
		}

		time.Sleep(tryGetLockSleepTimes)
	}

	return false
}

// ReleaseLock 释放分布式锁
func (r *RedisLock) ReleaseLock(ctx context.Context, lockKey, value string) bool {
	result := r.rdb.Eval(ctx, delLockKeyScript, []string{lockKey}, value)

	if result.Err() != nil {
		zlogger.Errorf("ReleaseLock |lockKey:%v| err: %v", lockKey, result.Err())
		return false
	}

	return true
}
