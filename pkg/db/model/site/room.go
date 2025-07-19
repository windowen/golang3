package site

import (
	"context"

	"serverApi/pkg/db/dbconn"
	table "serverApi/pkg/db/table/site"
)

// FindRoomInfo 获取单个直播间
func (s *Site) FindRoomInfo(ctx context.Context, opts ...table.RoomWhereOption) (*table.Room, error) {
	var (
		room  table.Room
		query = s.DB.WithContext(ctx).Scopes(room.DefaultScope)
	)

	for _, opt := range opts {
		query = opt(query)
	}

	err := query.First(&room).Error
	if err = dbconn.CheckErr(err); err != nil {
		return nil, err
	}

	return &room, nil
}

// UpdateRoomInfo 设置直播间信息
func (s *Site) UpdateRoomInfo(ctx context.Context, userID int, opts ...table.RoomUpdateOption) error {
	updates := make(map[string]interface{})

	for _, opt := range opts {
		opt(updates)
	}

	err := s.DB.WithContext(ctx).
		Model(&table.Room{}).
		Where(table.FieldRoomUserId.Eq(userID)).
		Updates(updates).Error

	if err != nil {
		return err
	}

	return nil
}
