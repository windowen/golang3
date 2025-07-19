package site

import (
	"context"
	"errors"

	"gorm.io/gorm"

	model "serverApi/pkg/db/table/site"
	"serverApi/pkg/tools/errs"
)

func (s *Site) CreatePlatformMsg(ctx context.Context, model *model.PlatformMsg) error {
	return s.DB.WithContext(ctx).Create(model).Error
}

func (s *Site) DeletePlatformMsg(ctx context.Context, id int) error {
	return s.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.PlatformMsg{}).Error
}

func (s *Site) UpdatePlatformMsg(ctx context.Context, id int, data map[string]interface{}) error {
	return s.DB.WithContext(ctx).Model(&model.PlatformMsg{}).Where("id = ?", id).Updates(data).Error
}

func (s *Site) FindPlatformMsgById(ctx context.Context, id int) (*model.PlatformMsg, error) {
	var data model.PlatformMsg
	if err := s.DB.WithContext(ctx).First(&data, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(gorm.ErrRecordNotFound)
		}
		return nil, errs.Wrap(err)
	}
	return &data, nil
}

func (s *Site) ListPlatformMsg(ctx context.Context, offset, limit int) ([]*model.PlatformMsg, error) {
	var models []*model.PlatformMsg
	if err := s.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return models, nil
}

func (s *Site) CountPlatformMsg(ctx context.Context) (int64, error) {
	var count int64
	if err := s.DB.WithContext(ctx).Model(&model.PlatformMsg{}).Count(&count).Error; err != nil {
		return 0, errs.Wrap(err)
	}
	return count, nil
}

// PlatformMsgList 平台消息列表
func (s *Site) PlatformMsgList(ctx context.Context, pageNumber int32, pageSize int32, langCode string) ([]*model.PlatformMsg, int64, error) {
	db := s.DB.WithContext(ctx).
		Table("site_platform_msg").
		Where("status=1 and language_code=?", langCode).
		Order("created_at desc")

	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, errs.Wrap(err)
	}
	if pageSize > 0 {
		db = db.Limit(int(pageSize)).Offset(int((pageNumber - 1) * pageSize))
	}

	resp := make([]*model.PlatformMsg, 0)
	err = db.Find(&resp).Error
	if err != nil {
		return nil, 0, err
	}

	return resp, count, nil
}

// RotatingList 跑马灯列表
func (s *Site) RotatingList(ctx context.Context, category int32, countryCode string, languageCode string) ([]*model.RotatingAdver, error) {
	db := s.DB.WithContext(ctx).
		Table("site_rotating_adver").
		Where("status=1 and country_code=? and category=? and language_code=?", countryCode, category, languageCode).
		Order("sort asc,created_at desc")

	resp := make([]*model.RotatingAdver, 0)
	err := db.Find(&resp).Error
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// LastPlatformMsg 最后一条平台消息
func (s *Site) LastPlatformMsg(ctx context.Context, id int) (*model.PlatformMsg, error) {
	var data model.PlatformMsg
	if err := s.DB.WithContext(ctx).Model(&model.PlatformMsg{}).
		Where("status=1").
		Order("id desc").
		First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(gorm.ErrRecordNotFound)
		}
		return nil, errs.Wrap(err)
	}
	return &data, nil
}
