package site

import (
	"context"

	"serverApi/pkg/db/dbconn"
	table "serverApi/pkg/db/table/site"
	"serverApi/pkg/tools/errs"
)

func (s *Site) CreateFamilyAnchors(ctx context.Context, table *table.FamilyAnchors) error {
	return s.DB.WithContext(ctx).Create(table).Error
}

func (s *Site) DeleteFamilyAnchors(ctx context.Context, id int) error {
	return s.DB.WithContext(ctx).Where("id = ?", id).Delete(&table.FamilyAnchors{}).Error
}

func (s *Site) UpdateFamilyAnchors(ctx context.Context, familyId, userId int, opts ...table.FamilyAnchorsUpdateOption) error {
	updates := make(map[string]interface{})

	for _, opt := range opts {
		opt(updates)
	}

	return s.DB.WithContext(ctx).
		Model(&table.FamilyAnchors{}).
		Where(table.FieldFamilyAnchorsFamilyId.Eq(familyId)).
		Where(table.FieldFamilyAnchorsUserId.Eq(userId)).
		Updates(updates).Error
}

func (s *Site) FindFamilyAnchorsById(ctx context.Context, id int) (*table.FamilyAnchors, error) {
	var (
		family table.FamilyAnchors
		query  = s.DB.WithContext(ctx).Scopes(family.DefaultScope)
	)

	err := query.First(&family, id).Error
	if err := dbconn.CheckErr(err); err != nil {
		return nil, errs.Wrap(err)
	}

	return &family, nil
}

func (s *Site) FindFamilyAnchors(ctx context.Context, opts ...table.FamilyAnchorsWhereOption) (*table.FamilyAnchors, error) {
	var (
		family table.FamilyAnchors
		query  = s.DB.WithContext(ctx).Scopes(family.DefaultScope)
	)

	for _, opt := range opts {
		query = opt(query)
	}

	err := query.First(&family).Error
	if err = dbconn.CheckErr(err); err != nil {
		return nil, err
	}

	return &family, nil
}

func (s *Site) ListFamilyAnchors(ctx context.Context, familyId int) ([]*table.FamilyAnchors, error) {
	var (
		family        table.FamilyAnchors
		familyAnchors []*table.FamilyAnchors
		query         = s.DB.WithContext(ctx).Scopes(family.DefaultScope)
	)

	if err := query.Where(table.FieldFamilyAnchorsFamilyId.Eq(familyId)).Find(&familyAnchors).Error; err != nil {
		return nil, errs.Wrap(err)
	}

	return familyAnchors, nil
}

func (s *Site) CountFamilyAnchors(ctx context.Context) (int64, error) {
	var count int64
	if err := s.DB.WithContext(ctx).Model(&table.FamilyAnchors{}).Count(&count).Error; err != nil {
		return 0, errs.Wrap(err)
	}
	return count, nil
}
