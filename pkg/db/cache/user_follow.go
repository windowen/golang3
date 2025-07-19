package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"serverApi/pkg/tools/cast"

	constsR "serverApi/pkg/constant/redis"
)

type FollowStatistics struct {
	FansCount   int32 `json:"fansCount"`
	FollowCount int32 `json:"followCount"`
}

// FollowUser 关注
func (uc *UserCache) FollowUser(ctx context.Context, userID, followID int) error {
	// 关注时间戳作为排序依据
	timestamp := time.Now().Unix()

	// 将 followID 添加到关注列表和排序集合
	uc.rdb.SAdd(ctx, fmt.Sprintf(constsR.UserFollowCacheKey, userID), followID)
	uc.rdb.ZAdd(ctx, fmt.Sprintf(constsR.UserFollowSortedCacheKey, userID), redis.Z{
		Score:  float64(timestamp),
		Member: followID,
	})

	// 将 userID 添加到对方的粉丝列表和排序集合
	uc.rdb.SAdd(ctx, fmt.Sprintf(constsR.UserFansCacheKey, followID), userID)
	uc.rdb.ZAdd(ctx, fmt.Sprintf(constsR.UserFansSortedCacheKey, followID), redis.Z{
		Score:  float64(timestamp),
		Member: userID,
	})

	return nil
}

// UnfollowUser 取关
func (uc *UserCache) UnfollowUser(ctx context.Context, userID, followID int) error {
	// 从用户的关注列表和排序集合中移除 followID
	err := uc.rdb.SRem(ctx, fmt.Sprintf(constsR.UserFollowCacheKey, userID), followID).Err()
	if err != nil {
		return fmt.Errorf("failed to remove followID from user follow list: %v", err)
	}

	err = uc.rdb.ZRem(ctx, fmt.Sprintf(constsR.UserFollowSortedCacheKey, userID), followID).Err()
	if err != nil {
		return fmt.Errorf("failed to remove followID from user follow sorted set: %v", err)
	}

	// 从被关注用户的粉丝列表和排序集合中移除 userID
	err = uc.rdb.SRem(ctx, fmt.Sprintf(constsR.UserFansCacheKey, followID), userID).Err()
	if err != nil {
		return fmt.Errorf("failed to remove userID from follower's fans list: %v", err)
	}

	err = uc.rdb.ZRem(ctx, fmt.Sprintf(constsR.UserFansSortedCacheKey, followID), userID).Err()
	if err != nil {
		return fmt.Errorf("failed to remove userID from follower's fans sorted set: %v", err)
	}

	return nil
}

// IsFollow 是否关注、粉丝
func (uc *UserCache) IsFollow(ctx context.Context, key string, userID, followID int) (bool, error) {
	isMember, err := uc.rdb.SIsMember(ctx, fmt.Sprintf(key, userID), followID).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check if user is following: %v", err)
	}

	return isMember, nil
}

// FollowStatistics 关注/粉丝统计
func (uc *UserCache) FollowStatistics(ctx context.Context, userID int) (*FollowStatistics, error) {
	fansCount, err := uc.rdb.ZCard(ctx, fmt.Sprintf(constsR.UserFansSortedCacheKey, userID)).Result()
	if err != nil {
		return nil, err
	}

	followCount, err := uc.rdb.ZCard(ctx, fmt.Sprintf(constsR.UserFollowSortedCacheKey, userID)).Result()
	if err != nil {
		return nil, err
	}

	return &FollowStatistics{
		FansCount:   cast.ToInt32(fansCount),
		FollowCount: cast.ToInt32(followCount),
	}, nil
}

// FollowList 获取列表
func (uc *UserCache) FollowList(ctx context.Context, key string, userID, lastId, pageSize int) ([]string, error) {
	var (
		err      error
		result   []string
		redisKey = fmt.Sprintf(key, userID)
	)

	if lastId == 0 {
		// 如果 lastId 为空，获取最新的关注列表
		result, err = uc.rdb.ZRevRange(ctx, redisKey, 0, int64(pageSize)-1).Result()
	} else {
		// 从 lastId 之后获取 pageSize 条记录
		lastScore, err := uc.rdb.ZScore(ctx, redisKey, cast.ToString(lastId)).Result()
		if err != nil {
			return nil, err
		}

		result, err = uc.rdb.ZRevRangeByScore(ctx, redisKey, &redis.ZRangeBy{
			Max:   fmt.Sprintf("%f", lastScore-1), // 获取上次返回的记录之前的数据
			Min:   "-inf",
			Count: int64(pageSize),
		}).Result()
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}
