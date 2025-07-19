package site

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"serverApi/pkg/db/dbconn"
	table "serverApi/pkg/db/table/site"
	"serverApi/pkg/tools/errs"
)

func (s *Site) CreateFamily(ctx context.Context, table *table.Family) error {
	return s.DB.WithContext(ctx).Create(table).Error
}

func (s *Site) DeleteFamily(ctx context.Context, id int) error {
	return s.DB.WithContext(ctx).Where("id = ?", id).Delete(&table.Family{}).Error
}

func (s *Site) UpdateFamily(ctx context.Context, id int, data map[string]interface{}) error {
	return s.DB.WithContext(ctx).Model(&table.Family{}).Where("id = ?", id).Updates(data).Error
}

func (s *Site) FindFamilyById(ctx context.Context, id int) (*table.Family, error) {
	var family table.Family
	if err := s.DB.WithContext(ctx).First(&family, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(gorm.ErrRecordNotFound)
		}
		return nil, errs.Wrap(err)
	}
	return &family, nil
}

func (s *Site) FindFamily(ctx context.Context, opts ...table.FamilyWhereOption) (*table.Family, error) {
	var (
		family table.Family
		query  = s.DB.WithContext(ctx).Scopes(family.DefaultScope)
	)

	for _, opt := range opts {
		query = opt(query)
	}

	err := query.First(&family).Error
	if err = dbconn.CheckErr(err); err != nil {
		return nil, err
	}

	return &family, nil
}

func (s *Site) ListFamily(ctx context.Context, offset, limit int) ([]*table.Family, error) {
	var tables []*table.Family
	if err := s.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&tables).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return tables, nil
}

func (s *Site) CountFamily(ctx context.Context) (int64, error) {
	var count int64
	if err := s.DB.WithContext(ctx).Model(&table.Family{}).Count(&count).Error; err != nil {
		return 0, errs.Wrap(err)
	}
	return count, nil
}
