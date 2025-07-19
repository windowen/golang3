package site

import (
	"context"

	"serverApi/pkg/db/dbconn"
	table "serverApi/pkg/db/table/site"
)

func (s *Site) GetUserAuth(ctx context.Context, opts ...table.UserAuthWhereOption) (*table.UserAuth, error) {
	var (
		ua    table.UserAuth
		query = s.DB.WithContext(ctx)
	)

	for _, opt := range opts {
		query = opt(query)
	}

	err := query.First(&ua).Error
	if err = dbconn.CheckErr(err); err != nil {
		return nil, err
	}

	return &ua, nil
}

func (s *Site) UpdateUserAuth(ctx context.Context, userID int, opts ...table.UserAuthUpdateOption) error {
	updates := make(map[string]interface{})

	for _, opt := range opts {
		opt(updates)
	}

	err := s.DB.WithContext(ctx).
		Model(&table.UserAuth{}).
		Where(table.FieldUserAuthUserId.Eq(userID)).
		Updates(updates).Error

	if err != nil {
		return err
	}

	return nil
}
