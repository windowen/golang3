package site

import (
	"context"

	"serverApi/pkg/db/dbconn"
	table "serverApi/pkg/db/table/site"
)

func (s *Site) FindUserUseMount(ctx context.Context, opts ...table.UserMountsWhereOption) (*table.UserMounts, error) {
	var (
		um    table.UserMounts
		query = s.DB.WithContext(ctx).Scopes(um.DefaultScope)
	)

	for _, opt := range opts {
		query = opt(query)
	}

	err := query.First(&um).Error
	if err = dbconn.CheckErr(err); err != nil {
		return nil, err
	}

	return &um, nil
}
