package cache

import (
	"context"
	"encoding/json"

	constsR "serverApi/pkg/constant/redis"
)

// CommonGlobalAreaGet 获取地区、语言
func (uc *UserCache) CommonGlobalAreaGet(ctx context.Context) (string, error) {
	return uc.rdb.Get(ctx, constsR.CommonGlobalArea).Result()
}

func (uc *UserCache) CommonGlobalAreaSet(ctx context.Context, value any) error {
	cacheData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return uc.rdb.Set(ctx, constsR.CommonGlobalArea, cacheData, uc.expires).Err()
}
