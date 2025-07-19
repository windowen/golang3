package site

import (
	"context"
	"encoding/json"
	"errors"
	"sort"

	"serverApi/pkg/common/mctx"
	constsR "serverApi/pkg/constant/redis"
	"serverApi/pkg/db/cache"
	commonPb "serverApi/pkg/protobuf/common"
	pb "serverApi/pkg/protobuf/site"
	"serverApi/pkg/queue"
	"serverApi/pkg/rocketmq"
	"serverApi/pkg/tools/cast"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/tools/strhelper"
	"serverApi/pkg/tools/utils"
	"serverApi/pkg/zlogger"
)

// SiteBanner 站点轮播
func (s *siteSrv) SiteBanner(ctx context.Context, req *pb.SiteBannerReq) (*pb.SiteBannerResp, error) {
	// 获取轮播
	banner, err := s.fetchSiteBanner(ctx, cast.ToInt(req.GetCategory()))
	if err != nil {
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 排序
	sort.Sort(cache.ByBannerSort(banner))

	var res []*pb.SiteBannerResp_SiteBannerItem
	for _, val := range banner {
		res = append(res, &pb.SiteBannerResp_SiteBannerItem{
			Uri:      val.Uri,
			ShowType: cast.ToInt32(val.ShowType),
			ExtInfo:  val.ExtInfo,
		})
	}

	return &pb.SiteBannerResp{List: res}, nil
}

func (s *siteSrv) fetchSiteBanner(ctx context.Context, category int) ([]*cache.SiteBannerCache, error) {
	var (
		bannerData  []*cache.SiteBannerCache
		countryCode = mctx.GetCountryCode(ctx)
	)

	banner, err := cache.FetchSiteBanner(ctx, s.redis)
	if err != nil || len(banner) == 0 {
		if err != nil {
			zlogger.Errorf("fetchSiteBanner | err:%v", err)
		}

		siteBanner, err := s.siteDB.ListSiteBanner(ctx)
		if err != nil {
			zlogger.Errorf("fetchSiteBanner ListSiteBanner | err:%v", err)
			return nil, err
		}

		// 开启管道
		pipe := s.redis.Pipeline()
		for _, val := range siteBanner {
			bannerCache := &cache.SiteBannerCache{
				Sort:        val.Sort,
				Category:    val.Category,
				Uri:         val.Uri,
				ShowType:    val.ShowType,
				ExtInfo:     val.ExtInfo,
				CountryJson: val.CountryJson,
			}

			marshal, err := json.Marshal(bannerCache)
			if err != nil {
				zlogger.Errorf("fetchSiteBanner Marshal |err:%v", err)
				continue
			}
			pipe.HSet(ctx, constsR.CommonSiteBannerHash, val.Id, marshal)

			var countryInfo []string
			err = strhelper.Json2Struct(val.CountryJson, &countryInfo)
			if err != nil {
				zlogger.Errorf("fetchSiteBanner db Json2Struct val.CountryJson Unmarshal |countryJson:%v|err: %v", val.CountryJson, err)
				continue
			}

			if val.Category == category && utils.SliceHas(countryInfo, countryCode) {
				bannerData = append(bannerData, bannerCache)
			}
		}

		// 提交
		_, err = pipe.Exec(ctx)
		if err != nil {
			zlogger.Errorf("fetchSiteBanner pipe.Exec |err: %v", err)
		}

		return bannerData, err
	}

	// 格式化缓存
	for _, val := range banner {
		var bannerCache *cache.SiteBannerCache
		err := json.Unmarshal([]byte(val), &bannerCache)
		if err != nil {
			zlogger.Errorf("fetchSiteBanner json.Unmarshal |banner:%v| err: %v", val, err)
			continue
		}

		var countryInfo []string
		err = strhelper.Json2Struct(bannerCache.CountryJson, &countryInfo)
		if err != nil {
			zlogger.Errorf("fetchSiteBanner redis Json2Struct val.CountryJson Unmarshal |countryJson:%v|err: %v", bannerCache.CountryJson, err)
			continue
		}

		if bannerCache.Category != category || !utils.SliceHas(countryInfo, countryCode) {
			continue
		}

		bannerData = append(bannerData, &cache.SiteBannerCache{
			Sort:     bannerCache.Sort,
			Category: bannerCache.Category,
			Uri:      bannerCache.Uri,
			ShowType: bannerCache.ShowType,
			ExtInfo:  bannerCache.ExtInfo,
		})
	}

	return bannerData, nil
}

// StartupImage 站点启动图
func (s *siteSrv) StartupImage(ctx context.Context, _ *commonPb.Empty) (*pb.StartupImageResp, error) {
	countryCode := mctx.GetCountryCode(ctx)

	startupImage, err := s.siteDB.FindStartupImageByCountryCode(ctx, countryCode)
	if err != nil {
		return nil, err
	}

	return &pb.StartupImageResp{
		CountryCode:  startupImage.CountryCode,
		LanguageCode: startupImage.LanguageCode,
		Name:         startupImage.Name,
		ImageUrl:     startupImage.ImageUrl,
		ThumbnailUrl: startupImage.ThumbnailUrl,
	}, nil
}

// Stay 站点启动图
func (s *siteSrv) Stay(ctx context.Context, req *pb.PageStayReq) (*commonPb.Empty, error) {
	if req.StayTime <= 0 {
		return nil, errors.New("stay time is zero")
	}

	// 用户停留时长
	rocketmq.PublishJson(rocketmq.StatsEvent, &queue.EventStats{
		EventType: queue.EventPageStay,
		PageName:  req.PageName,
		Timestamp: int64(req.StayTime),
	})

	return &commonPb.Empty{}, nil
}

// BannerClick 页面点击
func (s *siteSrv) BannerClick(ctx context.Context, req *pb.BannerClickReq) (*commonPb.Empty, error) {
	// 用户点击banner
	switch req.GetBannerType() {
	case 1:
		rocketmq.PublishJson(rocketmq.StatsEvent, &queue.EventStats{
			EventType: queue.EventHomepageBannerClicks,
		})
	case 2:
		rocketmq.PublishJson(rocketmq.StatsEvent, &queue.EventStats{
			EventType: queue.EventRecommendedBannerClicks,
		})
	case 3:
		rocketmq.PublishJson(rocketmq.StatsEvent, &queue.EventStats{
			EventType: queue.EventPopularBannerClicks,
		})
	}

	return &commonPb.Empty{}, nil
}
