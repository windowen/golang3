package cache

import (
	"context"

	"github.com/redis/go-redis/v9"

	constsR "serverApi/pkg/constant/redis"
	"serverApi/pkg/zlogger"
)

type SiteBannerCache struct {
	Sort        int    `json:"sort"`
	Category    int    `json:"category"`
	Uri         string `json:"uri"`
	ShowType    int    `json:"showType"`
	ExtInfo     string `json:"extInfo"`
	CountryJson string `json:"countryJson"`
}

// FindSiteBanner 获取系统站点配置缓存
func FindSiteBanner(ctx context.Context, rdb redis.UniversalClient, configKey string) (string, error) {
	result, err := rdb.HGet(ctx, constsR.CommonSiteBannerHash, configKey).Result()
	if err != nil {
		zlogger.Errorf("FindSiteBanner HGet Error: |configKey:%v| err: %v", configKey, err)
		return "", err
	}

	return result, nil
}

// FetchSiteBanner 获取所有站点配置
func FetchSiteBanner(ctx context.Context, rdb redis.UniversalClient) (map[string]string, error) {
	result, err := rdb.HGetAll(ctx, constsR.CommonSiteBannerHash).Result()
	if err != nil {
		zlogger.Errorf("FindSiteBanner HGetAll Error: | err: %v", err)
		return nil, err
	}

	return result, nil
}

type ByBannerSort []*SiteBannerCache

func (a ByBannerSort) Len() int           { return len(a) }
func (a ByBannerSort) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByBannerSort) Less(i, j int) bool { return a[i].Sort < a[j].Sort }
