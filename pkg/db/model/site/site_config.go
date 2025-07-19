package site

import (
	"context"

	"serverApi/pkg/constant"
	table "serverApi/pkg/db/table/site"
	"serverApi/pkg/tools/errs"
)

// ListSiteConfigs 获取系统站点配置
func (s *Site) ListSiteConfigs(ctx context.Context) ([]*table.SysConfig, error) {
	var (
		sc []*table.SysConfig
	)

	if err := s.DB.WithContext(ctx).
		Where(table.FieldSysConfigStatus.Eq(constant.StatusNormal)).
		Find(&sc).Error; err != nil {
		return nil, errs.Wrap(err)
	}

	return sc, nil
}
