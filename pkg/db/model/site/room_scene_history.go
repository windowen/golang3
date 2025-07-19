package site

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"serverApi/pkg/constant"
	"serverApi/pkg/db/dbconn"
	table "serverApi/pkg/db/table/site"
	"serverApi/pkg/tools/errs"
)

func (s *Site) CreateRoomSceneHistory(ctx context.Context, table *table.RoomSceneHistory) error {
	return s.DB.WithContext(ctx).Create(table).Error
}

func (s *Site) DeleteRoomSceneHistory(ctx context.Context, id int) error {
	return s.DB.WithContext(ctx).Where("id = ?", id).Delete(&table.RoomSceneHistory{}).Error
}

func (s *Site) UpdateRoomSceneHistory(ctx context.Context, roomId int, opts ...table.RoomSceneHistoryUpdateOption) error {
	updates := make(map[string]interface{})

	for _, opt := range opts {
		opt(updates)
	}

	return s.DB.WithContext(ctx).Model(&table.RoomSceneHistory{}).Where(table.FieldRoomSceneHistoryRoomId.Eq(roomId)).Updates(updates).Error
}

func (s *Site) FindRoomSceneHistoryById(ctx context.Context, id int) (*table.RoomSceneHistory, error) {
	var rshModel table.RoomSceneHistory

	if err := s.DB.WithContext(ctx).First(&rshModel, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(gorm.ErrRecordNotFound)
		}
		return nil, errs.Wrap(err)
	}
	return &rshModel, nil
}

func (s *Site) FindRoomSceneHistory(ctx context.Context, opts ...table.RoomSceneHistoryWhereOption) (*table.RoomSceneHistory, error) {
	var (
		roomSceneHistory table.RoomSceneHistory
		query            = s.DB.WithContext(ctx).Scopes(roomSceneHistory.DefaultScope)
	)

	for _, opt := range opts {
		query = opt(query)
	}

	err := query.First(&roomSceneHistory).Error
	if err = dbconn.CheckErr(err); err != nil {
		return nil, err
	}

	return &roomSceneHistory, nil
}

func (s *Site) ListRoomSceneHistory(ctx context.Context, offset, limit int) ([]*table.RoomSceneHistory, error) {
	var tables []*table.RoomSceneHistory
	if err := s.DB.WithContext(ctx).Limit(limit).Offset(offset).Find(&tables).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return tables, nil
}

func (s *Site) CountRoomSceneHistory(ctx context.Context) (int64, error) {
	var count int64
	if err := s.DB.WithContext(ctx).Model(&table.RoomSceneHistory{}).Count(&count).Error; err != nil {
		return 0, errs.Wrap(err)
	}
	return count, nil
}

func (s *Site) ListRoomSceneHistory30Day(ctx context.Context, roomId int) ([]*table.RoomSceneHistory, error) {
	var (
		model []*table.RoomSceneHistory
		query = s.DB.WithContext(ctx).Where(table.FieldRoomSceneHistoryStatus.Eq(constant.RoomSceneStatusEnd))
		now   = time.Now()
		// 获取30天前的时间
		thirtyDaysAgo = now.AddDate(0, 0, -30)
		// 将时间设置为当天的0点
		thirtyDaysAgoMidnight = time.Date(thirtyDaysAgo.Year(), thirtyDaysAgo.Month(), thirtyDaysAgo.Day(), 0, 0, 0, 0, thirtyDaysAgo.Location())
	)

	if err := query.
		Where(table.FieldRoomSceneHistoryRoomId.Eq(roomId)).
		Where(table.FieldRoomSceneHistoryOpenAt.Gte(thirtyDaysAgoMidnight)).
		Find(&model).Error; err != nil {
		return nil, errs.Wrap(err)
	}

	return model, nil
}
