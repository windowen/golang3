package site

import (
	"context"

	"serverApi/pkg/tools/cast"

	"serverApi/pkg/constant"
	"serverApi/pkg/db/dbconn"
	table "serverApi/pkg/db/table/site"
	pb "serverApi/pkg/protobuf/site"
	"serverApi/pkg/tools/errs"
)

func (s *Site) CreateFollow(ctx context.Context, table *table.Follow) error {
	return s.DB.WithContext(ctx).Create(table).Error
}

func (s *Site) DeleteFollow(ctx context.Context, id int) error {
	return s.DB.WithContext(ctx).Where("id = ?", id).Delete(&table.Follow{}).Error
}

func (s *Site) UpdateFollow(ctx context.Context, id int, opts ...table.FollowUpdateOption) error {
	updates := make(map[string]interface{})

	for _, opt := range opts {
		opt(updates)
	}

	err := s.DB.WithContext(ctx).
		Model(&table.Follow{}).
		Where(table.FieldFollowId.Eq(id)).
		Updates(updates).Error

	if err != nil {
		return err
	}

	return nil
}

func (s *Site) FindFollowById(ctx context.Context, id int) (*table.Follow, error) {
	var (
		follow table.Follow
		query  = s.DB.WithContext(ctx).Scopes(follow.DefaultScope)
	)

	err := query.First(&follow, id).Error
	if err = dbconn.CheckErr(err); err != nil {
		return nil, err
	}

	return &follow, nil
}

func (s *Site) FindFollow(ctx context.Context, opts ...table.FollowWhereOption) (*table.Follow, error) {
	var (
		follow table.Follow
		query  = s.DB.WithContext(ctx).Scopes(follow.DefaultScope)
	)

	for _, opt := range opts {
		query = opt(query)
	}

	err := query.First(&follow).Error
	if err = dbconn.CheckErr(err); err != nil {
		return nil, err
	}

	return &follow, nil
}

func (s *Site) ListFollow(ctx context.Context, limit int) ([]*table.Follow, error) {
	var tables []*table.Follow
	if err := s.DB.WithContext(ctx).Limit(limit).Find(&tables).Error; err != nil {
		return nil, errs.Wrap(err)
	}
	return tables, nil
}

func (s *Site) FansCount(ctx context.Context, followId int) (int64, error) {
	var (
		count  int64
		follow table.Follow
		query  = s.DB.WithContext(ctx).Scopes(follow.DefaultScope)
	)

	if err := query.Model(&follow).Where(table.FieldFollowFollowId.Eq(followId)).Count(&count).Error; err != nil {
		return 0, errs.Wrap(err)
	}
	return count, nil
}

func (s *Site) FollowingCount(ctx context.Context, userId int) (int64, error) {
	var (
		count  int64
		follow table.Follow
		query  = s.DB.WithContext(ctx).Scopes(follow.DefaultScope)
	)
	if err := query.Model(&follow).Where(table.FieldFollowUserId.Eq(userId)).Count(&count).Error; err != nil {
		return 0, errs.Wrap(err)
	}
	return count, nil
}

func (s *Site) FollowFans(ctx context.Context, userId int, req *pb.FollowsReq) ([]*table.FollowUserItem, error) {
	var (
		res   []*table.FollowUserItem
		query = s.DB.WithContext(ctx)
		limit = cast.ToInt(req.GetPageSize())
	)

	query = query.Where("f.follow_id = ?", userId)
	query = query.Where("f.status = ?", constant.UserFollowFocusOn)
	if req.GetLastId() > 0 {
		query = query.Where("f.id < ?", req.GetLastId())
	}

	if limit <= 0 {
		limit = 20
	}

	err := query.
		Table("site_follow as f").
		Joins("join site_user as su on su.id = f.user_id").
		Select("f.id, f.user_id, f.created_at, su.nickname, su.avatar, su.level_id, su.sex, su.sign").
		Limit(limit).
		Order("f.id desc").
		Find(&res).Error
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return res, nil
}

func (s *Site) FollowFollowing(ctx context.Context, userId int, req *pb.FollowsReq) ([]*table.FollowUserItem, error) {
	var (
		res   []*table.FollowUserItem
		query = s.DB.WithContext(ctx)
		limit = cast.ToInt(req.GetPageSize())
	)

	query = query.Where("f.user_id = ?", userId)
	query = query.Where("f.status = ?", constant.UserFollowFocusOn)
	if req.GetLastId() > 0 {
		query = query.Where("f.id < ?", req.GetLastId())
	}

	if limit <= 0 {
		limit = 20
	}

	err := query.
		Table("site_follow as f").
		Joins("join site_user as su on su.id = f.follow_id").
		Select("f.id, f.follow_id as user_id, f.created_at, su.nickname, su.avatar, su.level_id, su.sex, su.sign").
		Limit(limit).
		Order("f.id desc").
		Find(&res).Error
	if err != nil {
		return nil, errs.Wrap(err)
	}

	return res, nil
}
