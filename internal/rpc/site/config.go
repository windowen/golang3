package site

import (
	"context"

	"github.com/redis/go-redis/v9"

	"serverApi/pkg/constant"
	constsR "serverApi/pkg/constant/redis"
	"serverApi/pkg/db/cache"
	commonPb "serverApi/pkg/protobuf/common"
	pb "serverApi/pkg/protobuf/site"
	"serverApi/pkg/tools/cast"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/zlogger"
)

// SiteConfigs 系统站点配置
func (s *siteSrv) SiteConfigs(ctx context.Context, _ *commonPb.Empty) (*pb.SiteConfigsResp, error) {
	// 获取系统配置
	configs, err := s.fetchSysSiteConfigs(ctx, s.redis)
	if err != nil {
		zlogger.Errorf("SiteConfigs fetchSysSiteConfigs | err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	resp := &pb.SiteConfigsResp{}

	// S3地址
	if value, ok := configs[constant.SysSiteCfgS3Code]; ok {
		resp.ResourceUri = value
	}

	// 钻石兑换比列
	if value, ok := configs[constant.SysSiteCfgDiamondExchangeUsdRatio]; ok {
		resp.DiamondExchangeUsdRatio = cast.ToInt32(value)
	}

	// 是否开启换绑
	if value, ok := configs[constant.SysSiteCfgSecurityChangeBind]; ok {
		resp.IsChangeBind = cast.ToInt32(value)
	}

	return resp, nil
}

// SysDefaultAvatar 系统默认头像
func (s *siteSrv) SysDefaultAvatar(ctx context.Context) string {
	// 获取系统配置
	configs, err := s.fetchSysSiteConfigs(ctx, s.redis)
	if err != nil {
		zlogger.Errorf("SysDefaultAvatar fetchSysSiteConfigs | err: %v", err)
		return ""
	}

	// 系统默认头像
	if value, ok := configs[constant.SysSiteCfgDefaultAvatarCode]; ok {
		return value
	}

	return ""
}

func (s *siteSrv) fetchSysSiteConfigs(ctx context.Context, rdb redis.UniversalClient) (map[string]string, error) {
	configs, err := cache.FetchSysSiteConfigs(ctx, s.redis)
	if err != nil || len(configs) == 0 {
		if err != nil {
			zlogger.Errorf("fetchSysSiteConfigs | err: %v", err)
		}

		siteConfigs, err := s.siteDB.ListSiteConfigs(ctx)
		if err != nil {
			zlogger.Errorf("fetchSysSiteConfigs ListSiteConfigs | err: %v", err)
			return nil, err
		}

		var configs = make(map[string]string, len(siteConfigs))
		for _, config := range siteConfigs {
			configs[config.ConfigCode] = config.Content
		}

		err = rdb.HSet(ctx, constsR.CommonSysSiteCfgHash, configs).Err()
		if err != nil {
			zlogger.Errorf("fetchSysSiteConfigs hSet | err: %v", err)
			return nil, nil
		}

		return configs, nil
	}

	return configs, nil
}
