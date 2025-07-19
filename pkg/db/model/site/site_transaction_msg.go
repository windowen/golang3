package site

import (
	"context"
	"errors"

	"gorm.io/gorm"

	model "serverApi/pkg/db/table/site"
	"serverApi/pkg/tools/errs"
)

func (s *Site) CreateTransactionMsg(ctx context.Context, model *model.TransactionMsg) error {
	return s.DB.WithContext(ctx).Create(model).Error
}

func (s *Site) DeleteTransactionMsg(ctx context.Context, id int) error {
	return s.DB.WithContext(ctx).Where("id = ?", id).Delete(&model.TransactionMsg{}).Error
}

// UpdateTransactionMsg 更新用户的所有交易消息为已读
func (s *Site) UpdateTransactionMsg(ctx context.Context, userId int) error {
	return s.DB.WithContext(ctx).Model(&model.TransactionMsg{}).
		Where("user_id = ? and read_status=0", userId).Update("read_status", 1).Error
}

// UpdatePartTransactionMsg 更新用户的部分交易消息为已读
func (s *Site) UpdatePartTransactionMsg(ctx context.Context, userId int, ids []int32) error {
	return s.DB.WithContext(ctx).Model(&model.TransactionMsg{}).
		Where("user_id = ? and id in(?)", userId, ids).Update("read_status", 1).Error
}

func (s *Site) FindTransactionMsgById(ctx context.Context, id int) (*model.TransactionMsg, error) {
	var data model.TransactionMsg
	if err := s.DB.WithContext(ctx).First(&data, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(gorm.ErrRecordNotFound)
		}
		return nil, errs.Wrap(err)
	}
	return &data, nil
}

func (s *Site) ListTransactionMsg(ctx context.Context, offset, limit int) ([]*model.TransactionMsg, error) {
	var models []*model.TransactionMsg
	if err := s.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return models, nil
}

func (s *Site) CountTransactionMsg(ctx context.Context, userId int) (int64, error) {
	var count int64
	if err := s.DB.WithContext(ctx).Model(&model.TransactionMsg{}).
		Where("user_id=? and read_status=0 and status=0", userId).
		Count(&count).Error; err != nil {
		return 0, errs.Wrap(err)
	}
	return count, nil
}

// TransactionMsgList 交易消息列表
func (s *Site) TransactionMsgList(ctx context.Context, pageNumber int32, pageSize int32, userId int) ([]*model.TransactionMsg, int64, error) {
	db := s.DB.WithContext(ctx).
		Table("site_transaction_msg").
		Where("user_id=? and status=0", userId).
		Order("created_at desc")

	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, 0, errs.Wrap(err)
	}
	if pageSize > 0 {
		db = db.Limit(int(pageSize)).Offset(int((pageNumber - 1) * pageSize))
	}

	resp := make([]*model.TransactionMsg, 0)
	err = db.Find(&resp).Error
	if err != nil {
		return nil, 0, err
	}

	return resp, count, nil
}

// ClearTransactionMsgReadAll 清除交易已读消息
func (s *Site) ClearTransactionMsgReadAll(ctx context.Context, userId int) error {
	return s.DB.WithContext(ctx).Model(&model.TransactionMsg{}).
		Where("user_id = ? and read_status=1", userId).Update("status", 1).Error
}

// LastTransactionMsg 最后一条交易消息
func (s *Site) LastTransactionMsg(ctx context.Context, id int) (*model.TransactionMsg, error) {
	var data model.TransactionMsg
	if err := s.DB.WithContext(ctx).Model(&model.TransactionMsg{}).
		Where("user_id=? and status=0", id).
		Order("id desc").
		First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(gorm.ErrRecordNotFound)
		}
		return nil, errs.Wrap(err)
	}
	return &data, nil

}
