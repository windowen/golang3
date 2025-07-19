package site

import (
	"context"

	"serverApi/pkg/constant"
	"serverApi/pkg/db/dbconn"
	table "serverApi/pkg/db/table/site"
)

func (s *Site) GetGlobalAreaByMobileCode(mobileCode string) (*table.GlobalArea, error) {
	var (
		globalArea table.GlobalArea
		query      = s.DB.Scopes(globalArea.DefaultScope)
	)

	err := query.Where(table.FieldGlobalAreaMobileCode.Eq(mobileCode)).
		Where(table.FieldGlobalAreaStatus.Eq(constant.StatusNormal)).
		First(&globalArea).Error
	if err = dbconn.CheckErr(err); err != nil {
		return nil, err
	}

	return &globalArea, nil
}

func (s *Site) GlobalAreas(ctx context.Context) ([]*table.GlobalArea, error) {
	var globalAreas []*table.GlobalArea

	err := s.DB.WithContext(ctx).Where(table.FieldGlobalAreaStatus.Eq(constant.StatusNormal)).Find(&globalAreas).Error
	if err != nil {
		return nil, err
	}

	return globalAreas, nil
}
