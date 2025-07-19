package site

import (
	"context"

	table "serverApi/pkg/db/table/site"
)

func (s *Site) GlobalLang(ctx context.Context) ([]*table.GlobalLang, error) {
	var (
		gl         table.GlobalLang
		globalLang []*table.GlobalLang
		query      = s.DB.WithContext(ctx).Scopes(gl.DefaultScope)
	)

	err := query.Find(&globalLang).Error
	if err != nil {
		return nil, err
	}

	return globalLang, nil
}
