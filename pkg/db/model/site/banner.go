package site

import (
	"context"

	"serverApi/pkg/constant"
	table "serverApi/pkg/db/table/site"
	"serverApi/pkg/tools/errs"
)

// ListSiteBanner 获取站点轮播
func (s *Site) ListSiteBanner(ctx context.Context) ([]*table.Banner, error) {
	var (
		sb []*table.Banner
	)

	if err := s.DB.WithContext(ctx).
		Where(table.FieldBannerStatus.Eq(constant.StatusNormal)).
		Find(&sb).Error; err != nil {
		return nil, errs.Wrap(err)
	}

	return sb, nil
}
