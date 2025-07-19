package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"

	constsR "serverApi/pkg/constant/redis"
	"serverApi/pkg/db/cache"
)

func TestFollow(t *testing.T) {
	var (
		err error
		ctx = context.Background()
	)

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis 在本地运行
		DB:   1,
	})

	userCache := cache.NewUserCache(rdb)

	// 是否关注
	isFollow, err := userCache.IsFollow(ctx, constsR.UserFollowCacheKey, 18, 132)
	if err != nil {
		return
	}

	if !isFollow {
		err = userCache.FollowUser(ctx, 18, 132)
		if err != nil {
			fmt.Println(err, "FollowUser")
			return
		}
	}

	statistics, err := userCache.FollowStatistics(ctx, 18)
	if err != nil {
		fmt.Println(err, "FollowStatistics")
		return
	}

	fmt.Println(statistics)

	pagination, err := userCache.FollowList(ctx, constsR.UserFansSortedCacheKey, 132, 0, 20)
	if err != nil {
		fmt.Println(err, "GetFollowWithPagination")
		return
	}

	fmt.Println(pagination)
}
