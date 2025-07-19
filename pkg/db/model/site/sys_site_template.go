package site

import (
	"context"

	"serverApi/pkg/constant"
	table "serverApi/pkg/db/table/site"
)

func (s *Site) SysTmpl(ctx context.Context) ([]*table.SysTemplate, error) {
	var (
		stpl  table.SysTemplate
		stpls []*table.SysTemplate
		query = s.DB.WithContext(ctx).Scopes(stpl.DefaultScope)
	)

	err := query.Where(table.FieldSysTemplateCategory.In(constant.TmplTypeSms, constant.TmplTypeEmail)).Find(&stpls).Error
	if err != nil {
		return nil, err
	}

	return stpls, nil
}
