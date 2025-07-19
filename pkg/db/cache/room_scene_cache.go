package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	constsR "serverApi/pkg/constant/redis"
	"serverApi/pkg/tools/cast"
	"serverApi/pkg/zlogger"
)

// AddRoomScene 添加
func (rc *RoomCache) AddRoomScene(ctx context.Context, key string, roomId int, userId string) error {
	return rc.rdb.SAdd(ctx, fmt.Sprintf(key, roomId), userId).Err()
}

// CountRoomScene 数量
func (rc *RoomCache) CountRoomScene(ctx context.Context, key string, roomId int) (int64, error) {
	return rc.rdb.SCard(ctx, fmt.Sprintf(key, roomId)).Result()
}

// RemRoomScene 移除
func (rc *RoomCache) RemRoomScene(ctx context.Context, key string, roomId int, userId string) error {
	return rc.rdb.SRem(ctx, fmt.Sprintf(key, roomId), userId).Err()
}

// RoomSceneDataIsExist 是否存在
func (rc *RoomCache) RoomSceneDataIsExist(ctx context.Context, key string, roomId int, userId string) (bool, error) {
	return rc.rdb.SIsMember(ctx, fmt.Sprintf(key, roomId), userId).Result()
}

// ClearRoomScene 删除本场数据
func (rc *RoomCache) ClearRoomScene(ctx context.Context, key string, roomId int) error {
	return rc.rdb.Del(ctx, fmt.Sprintf(key, roomId)).Err()
}

// AddRoomScenePayUsers 新增直播间付费用户
func (rc *RoomCache) AddRoomScenePayUsers(ctx context.Context, roomId, userId int) error {
	pipe := rc.rdb.Pipeline()

	pipe.ZAdd(ctx, fmt.Sprintf(constsR.ScenePayUsers, roomId), redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: userId,
	})

	pipe.SAdd(ctx, fmt.Sprintf(constsR.ScenePayUsersSet, roomId), userId)

	// 提交
	_, err := pipe.Exec(ctx)
	if err != nil {
		zlogger.Errorf("AddRoomScenePayUsers Exec | err: %v", err)
		return err
	}

	return nil
}

// RoomSceneUserIsPay 是否付费
func (rc *RoomCache) RoomSceneUserIsPay(ctx context.Context, roomId, userId int) (bool, error) {
	exists, err := rc.rdb.SIsMember(ctx, fmt.Sprintf(constsR.ScenePayUsersSet, roomId), cast.ToString(userId)).Result()
	if err != nil {
		zlogger.Errorf("RoomSceneUserIsPay SIsMember | err: %v", err)
		return false, err
	}

	return exists, nil
}

// ClearRoomSceneUserPay 清理本场付费数据
func (rc *RoomCache) ClearRoomSceneUserPay(ctx context.Context, roomId int) error {
	pipe := rc.rdb.Pipeline()

	pipe.Del(ctx, fmt.Sprintf(constsR.ScenePayUsers, roomId))
	pipe.Del(ctx, fmt.Sprintf(constsR.ScenePayUsersSet, roomId))

	// 提交
	_, err := pipe.Exec(ctx)
	if err != nil {
		zlogger.Errorf("ClearRoomSceneUserPay Exec | err: %v", err)
		return err
	}

	return nil
}
