package site

import (
	"time"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"serverApi/pkg/constant"
)

const TableNameFollow = "site_follow"

var (
	FieldFollowId        = field.NewInt(TableNameFollow, "id")
	FieldFollowFollowId  = field.NewInt(TableNameFollow, "follow_id")
	FieldFollowUserId    = field.NewInt(TableNameFollow, "user_id")
	FieldFollowStatus    = field.NewInt(TableNameFollow, "status")
	FieldFollowCreatedAt = field.NewTime(TableNameFollow, "created_at")
	FieldFollowUpdatedAt = field.NewTime(TableNameFollow, "updated_at")
)

type Follow struct {
	Id        int       `gorm:"column:id" json:"id"`
	FollowId  int       `gorm:"column:follow_id" json:"follow_id"`   // 关注用户id
	UserId    int       `gorm:"column:user_id" json:"user_id"`       // 粉丝id
	Status    int       `gorm:"column:status" json:"status"`         // 状态 1-关注 2-取关
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"` // 更新时间
}

type FollowUserItem struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Avatar    string    `json:"avatar"`
	Sex       int       `json:"sex"`
	LevelId   int       `json:"level_id"`
	Sign      string    `json:"sign"`
	CreatedAt time.Time `json:"updated_at"`
}

func (*Follow) TableName() string {
	return TableNameFollow
}

func (a *Follow) IsEmpty() bool {
	if a == nil {
		return true
	}
	return a.Id == 0
}

// DefaultScope 过滤已取关 status=2
func (a *Follow) DefaultScope(db *gorm.DB) *gorm.DB {
	return db.Where(FieldFollowStatus.Neq(constant.UserFollowUnlock))
}

// FollowWhereOption 条件项
type FollowWhereOption func(*gorm.DB) *gorm.DB

// WhereFollowId 查询
func WhereFollowId(value int) FollowWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFollowId.Eq(value))
	}
}

// WhereFollowFollowId 查询关注用户id
func WhereFollowFollowId(value int) FollowWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFollowFollowId.Eq(value))
	}
}

// WhereFollowUserId 查询粉丝id
func WhereFollowUserId(value int) FollowWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFollowUserId.Eq(value))
	}
}

// WhereFollowStatus 查询状态 1-关注 2-取关
func WhereFollowStatus(value int) FollowWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFollowStatus.Eq(value))
	}
}

// WhereFollowCreatedAt 查询创建时间
func WhereFollowCreatedAt(value time.Time) FollowWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFollowCreatedAt.Eq(value))
	}
}

// WhereFollowUpdatedAt 查询更新时间
func WhereFollowUpdatedAt(value time.Time) FollowWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFollowUpdatedAt.Eq(value))
	}
}

// FollowUpdateOption 修改项
type FollowUpdateOption func(map[string]interface{})

// SetFollowId 设置
func SetFollowId(value int) FollowUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFollowId.ColumnName().String()] = value
	}
}

// SetFollowFollowId 设置关注用户id
func SetFollowFollowId(value int) FollowUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFollowFollowId.ColumnName().String()] = value
	}
}

// SetFollowUserId 设置粉丝id
func SetFollowUserId(value int) FollowUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFollowUserId.ColumnName().String()] = value
	}
}

// SetFollowStatus 设置状态 1-关注 2-取关
func SetFollowStatus(value int) FollowUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFollowStatus.ColumnName().String()] = value
	}
}

// SetFollowCreatedAt 设置创建时间
func SetFollowCreatedAt(value time.Time) FollowUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFollowCreatedAt.ColumnName().String()] = value
	}
}

// SetFollowUpdatedAt 设置更新时间
func SetFollowUpdatedAt(value time.Time) FollowUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFollowUpdatedAt.ColumnName().String()] = value
	}
}
