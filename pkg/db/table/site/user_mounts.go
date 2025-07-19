package site

import (
	"time"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"serverApi/pkg/constant"
)

const TableNameUserMounts = "site_user_mounts"

var (
	FieldUserMountsId          = field.NewInt(TableNameUserMounts, "id")
	FieldUserMountsUserId      = field.NewInt(TableNameUserMounts, "user_id")
	FieldUserMountsMountsId    = field.NewInt(TableNameUserMounts, "mounts_id")
	FieldUserMountsIsSelected  = field.NewInt(TableNameUserMounts, "is_selected")
	FieldUserMountsStatus      = field.NewInt(TableNameUserMounts, "status")
	FieldUserMountsExpiredTime = field.NewTime(TableNameUserMounts, "expired_time")
	FieldUserMountsCreatedAt   = field.NewTime(TableNameUserMounts, "created_at")
	FieldUserMountsUpdatedAt   = field.NewTime(TableNameUserMounts, "updated_at")
)

type UserMounts struct {
	Id          int       `gorm:"column:id" json:"id"`                     // ID
	UserId      int       `gorm:"column:user_id" json:"user_id"`           // 用户ID
	MountsId    int       `gorm:"column:mounts_id" json:"mounts_id"`       // 坐骑id
	IsSelected  int       `gorm:"column:is_selected" json:"is_selected"`   // 是否使用
	Status      int       `gorm:"column:status" json:"status"`             // 状态 1-正常 2-禁用 3-删除
	ExpiredTime time.Time `gorm:"column:expired_time" json:"expired_time"` // 过期时间
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`     // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`     // 更新时间
}

func (*UserMounts) TableName() string {
	return TableNameUserMounts
}

func (a *UserMounts) IsEmpty() bool {
	if a == nil {
		return true
	}
	return a.Id == 0
}

// DefaultScope 过滤已删除 status=3
func (a *UserMounts) DefaultScope(db *gorm.DB) *gorm.DB {
	return db.Where(FieldUserMountsStatus.Eq(constant.StatusNormal))
}

// UserMountsWhereOption 条件项
type UserMountsWhereOption func(*gorm.DB) *gorm.DB

// WhereUserMountsId 查询ID
func WhereUserMountsId(value int) UserMountsWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserMountsId.Eq(value))
	}
}

// WhereUserMountsUserId 查询用户ID
func WhereUserMountsUserId(value int) UserMountsWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserMountsUserId.Eq(value))
	}
}

// WhereUserMountsMountsId 查询坐骑id
func WhereUserMountsMountsId(value int) UserMountsWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserMountsMountsId.Eq(value))
	}
}

// WhereUserMountsIsSelected 查询是否使用
func WhereUserMountsIsSelected(value int) UserMountsWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserMountsIsSelected.Eq(value))
	}
}

// WhereUserMountsStatus 查询状态 1-正常 2-禁用 3-删除
func WhereUserMountsStatus(value int) UserMountsWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserMountsStatus.Eq(value))
	}
}

// WhereUserMountsExpiredTime 查询过期时间
func WhereUserMountsExpiredTime(value time.Time) UserMountsWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserMountsExpiredTime.Gte(value))
	}
}

// WhereUserMountsCreatedAt 查询创建时间
func WhereUserMountsCreatedAt(value time.Time) UserMountsWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserMountsCreatedAt.Eq(value))
	}
}

// WhereUserMountsUpdatedAt 查询更新时间
func WhereUserMountsUpdatedAt(value time.Time) UserMountsWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserMountsUpdatedAt.Eq(value))
	}
}

// UserMountsUpdateOption 修改项
type UserMountsUpdateOption func(map[string]interface{})

// SetUserMountsId 设置ID
func SetUserMountsId(value int) UserMountsUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserMountsId.ColumnName().String()] = value
	}
}

// SetUserMountsUserId 设置用户ID
func SetUserMountsUserId(value int) UserMountsUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserMountsUserId.ColumnName().String()] = value
	}
}

// SetUserMountsMountsId 设置坐骑id
func SetUserMountsMountsId(value int) UserMountsUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserMountsMountsId.ColumnName().String()] = value
	}
}

// SetUserMountsIsSelected 设置是否使用
func SetUserMountsIsSelected(value int) UserMountsUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserMountsIsSelected.ColumnName().String()] = value
	}
}

// SetUserMountsStatus 设置状态 1-正常 2-禁用 3-删除
func SetUserMountsStatus(value int) UserMountsUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserMountsStatus.ColumnName().String()] = value
	}
}

// SetUserMountsExpiredTime 设置过期时间
func SetUserMountsExpiredTime(value time.Time) UserMountsUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserMountsExpiredTime.ColumnName().String()] = value
	}
}

// SetUserMountsCreatedAt 设置创建时间
func SetUserMountsCreatedAt(value time.Time) UserMountsUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserMountsCreatedAt.ColumnName().String()] = value
	}
}

// SetUserMountsUpdatedAt 设置更新时间
func SetUserMountsUpdatedAt(value time.Time) UserMountsUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserMountsUpdatedAt.ColumnName().String()] = value
	}
}
