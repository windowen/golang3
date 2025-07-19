package site

import (
	"context"
	"errors"

	"gorm.io/gorm"

	model "serverApi/pkg/db/table/site"
	"serverApi/pkg/tools/errs"
)

func (s *Site) CreateStartupImage(ctx context.Context, model *model.StartupImage) error {
	return s.DB.WithContext(ctx).Create(model).Error
}

func (s *Site) DeleteStartupImage(ctx context.Context, id int) error {
	return s.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.StartupImage{}).Error
}

func (s *Site) UpdateStartupImage(ctx context.Context, id int, data map[string]interface{}) error {
	return s.DB.WithContext(ctx).Model(&model.StartupImage{}).Where("id = ?", id).Updates(data).Error
}

// FindStartupImageByCountryCode 通过国家编码获取启动图
func (s *Site) FindStartupImageByCountryCode(ctx context.Context, countryCode string) (*model.StartupImage, error) {
	var data model.StartupImage
	if err := s.DB.WithContext(ctx).Where("status=1 and country_code=?", countryCode).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(gorm.ErrRecordNotFound)
		}
		return nil, errs.Wrap(err)
	}
	return &data, nil
}

func (s *Site) ListStartupImage(ctx context.Context, offset, limit int) ([]*model.StartupImage, error) {
	var models []*model.StartupImage
	if err := s.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return models, nil
}

func (s *Site) CountStartupImage(ctx context.Context) (int64, error) {
	var count int64
	if err := s.DB.WithContext(ctx).Model(&model.StartupImage{}).Count(&count).Error; err != nil {
		return 0, errs.Wrap(err)
	}
	return count, nil
}
