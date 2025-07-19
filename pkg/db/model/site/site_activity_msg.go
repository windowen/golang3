package site

import (
	"context"
	"errors"

	"gorm.io/gorm"

	model "serverApi/pkg/db/table/site"
	"serverApi/pkg/tools/errs"
)

func (s *Site) CreateActivityMsg(ctx context.Context, model *model.ActivityMsg) error {
	return s.DB.WithContext(ctx).Create(model).Error
}

func (s *Site) DeleteActivityMsg(ctx context.Context, id int) error {
	return s.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.ActivityMsg{}).Error
}

func (s *Site) UpdateActivityMsg(ctx context.Context, id int, data map[string]interface{}) error {
	return s.DB.WithContext(ctx).Model(&model.ActivityMsg{}).Where("id = ?", id).Updates(data).Error
}

func (s *Site) FindActivityMsgById(ctx context.Context, id int) (*model.ActivityMsg, error) {
	var data model.ActivityMsg
	if err := s.DB.WithContext(ctx).First(&data, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(gorm.ErrRecordNotFound)
		}
		return nil, errs.Wrap(err)
	}
	return &data, nil
}

func (s *Site) ListActivityMsg(ctx context.Context, offset, limit int) ([]*model.ActivityMsg, error) {
	var models []*model.ActivityMsg
	if err := s.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return models, nil
}

func (s *Site) CountActivityMsg(ctx context.Context) (int64, error) {
	var count int64
	if err := s.DB.WithContext(ctx).Model(&model.ActivityMsg{}).Count(&count).Error; err != nil {
		return 0, errs.Wrap(err)
	}
	return count, nil
}

// ActivityMsgList 活动消息列表
func (s *Site) ActivityMsgList(ctx context.Context, pageNumber int32, pageSize int32, langCode string) ([]*model.ActivityMsg, int64, error) {
	db := s.DB.WithContext(ctx).
		Table("site_activity_msg").
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

	resp := make([]*model.ActivityMsg, 0)
	err = db.Find(&resp).Error
	if err != nil {
		return nil, 0, err
	}

	return resp, count, nil
}

// LastActivityMsg 最后一条活动消息
func (s *Site) LastActivityMsg(ctx context.Context, id int) (*model.ActivityMsg, error) {
	var data model.ActivityMsg
	if err := s.DB.WithContext(ctx).Model(&model.ActivityMsg{}).
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
