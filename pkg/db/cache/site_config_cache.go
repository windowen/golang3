package cache

import (
	"context"

	"github.com/redis/go-redis/v9"

	constsR "serverApi/pkg/constant/redis"
	"serverApi/pkg/zlogger"
)

// FindSysSiteConfig 获取系统站点配置缓存
func FindSysSiteConfig(ctx context.Context, rdb redis.UniversalClient, configKey string) (string, error) {
	result, err := rdb.HGet(ctx, constsR.CommonSysSiteCfgHash, configKey).Result()
	if err != nil {
		zlogger.Errorf("FindSysSiteConfig HGet Error: |configKey:%v| err: %v", configKey, err)
		return "", err
	}

	return result, nil
}

// FetchSysSiteConfigs 获取所有站点配置
func FetchSysSiteConfigs(ctx context.Context, rdb redis.UniversalClient) (map[string]string, error) {
	result, err := rdb.HGetAll(ctx, constsR.CommonSysSiteCfgHash).Result()
	if err != nil {
		zlogger.Errorf("GetSysSiteConfigs HGetAll Error: | err: %v", err)
		return nil, err
	}

	return result, nil
}
