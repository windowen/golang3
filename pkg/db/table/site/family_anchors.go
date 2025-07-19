package site

import (
	"time"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"serverApi/pkg/constant"
)

const TableNameFamilyAnchors = "live_family_anchors"

var (
	FieldFamilyAnchorsId          = field.NewInt(TableNameFamilyAnchors, "id")
	FieldFamilyAnchorsFamilyId    = field.NewInt(TableNameFamilyAnchors, "family_id")
	FieldFamilyAnchorsPatriarchId = field.NewInt(TableNameFamilyAnchors, "patriarch_id")
	FieldFamilyAnchorsUserId      = field.NewInt(TableNameFamilyAnchors, "user_id")
	FieldFamilyAnchorsStatus      = field.NewInt(TableNameFamilyAnchors, "status")
	FieldFamilyAnchorsCreatedAt   = field.NewTime(TableNameFamilyAnchors, "created_at")
	FieldFamilyAnchorsUpdatedAt   = field.NewTime(TableNameFamilyAnchors, "updated_at")
)

type FamilyAnchors struct {
	Id          int       `gorm:"column:id" json:"id"`                     // ID
	FamilyId    int       `gorm:"column:family_id" json:"family_id"`       // 家族id
	PatriarchId int       `gorm:"column:patriarch_id" json:"patriarch_id"` // 族长id
	UserId      int       `gorm:"column:user_id" json:"user_id"`           // 主播id
	Status      int       `gorm:"column:status" json:"status"`             // 状态 1-正常 3-踢出
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`     // 创建时间
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`     // 更新时间
}

func (*FamilyAnchors) TableName() string {
	return TableNameFamilyAnchors
}

func (a *FamilyAnchors) IsEmpty() bool {
	if a == nil {
		return true
	}
	return a.Id == 0
}

// DefaultScope 过滤不可用
func (a *FamilyAnchors) DefaultScope(db *gorm.DB) *gorm.DB {
	return db.Where(FieldFamilyAnchorsStatus.Eq(constant.StatusNormal))
}

// FamilyAnchorsWhereOption 条件项
type FamilyAnchorsWhereOption func(*gorm.DB) *gorm.DB

// WhereFamilyAnchorsId 查询ID
func WhereFamilyAnchorsId(value int) FamilyAnchorsWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFamilyAnchorsId.Eq(value))
	}
}

// WhereFamilyAnchorsFamilyId 查询家族id
func WhereFamilyAnchorsFamilyId(value int) FamilyAnchorsWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFamilyAnchorsFamilyId.Eq(value))
	}
}

// WhereFamilyAnchorsPatriarchId 查询族长id
func WhereFamilyAnchorsPatriarchId(value int) FamilyAnchorsWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFamilyAnchorsPatriarchId.Eq(value))
	}
}

// WhereFamilyAnchorsUserId 查询主播id
func WhereFamilyAnchorsUserId(value int) FamilyAnchorsWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFamilyAnchorsUserId.Eq(value))
	}
}

// WhereFamilyAnchorsStatus 查询状态 1-正常 3-踢出
func WhereFamilyAnchorsStatus(value int) FamilyAnchorsWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFamilyAnchorsStatus.Eq(value))
	}
}

// WhereFamilyAnchorsCreatedAt 查询创建时间
func WhereFamilyAnchorsCreatedAt(value time.Time) FamilyAnchorsWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFamilyAnchorsCreatedAt.Eq(value))
	}
}

// WhereFamilyAnchorsUpdatedAt 查询更新时间
func WhereFamilyAnchorsUpdatedAt(value time.Time) FamilyAnchorsWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFamilyAnchorsUpdatedAt.Eq(value))
	}
}

// FamilyAnchorsUpdateOption 修改项
type FamilyAnchorsUpdateOption func(map[string]interface{})

// SetFamilyAnchorsId 设置ID
func SetFamilyAnchorsId(value int) FamilyAnchorsUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFamilyAnchorsId.ColumnName().String()] = value
	}
}

// SetFamilyAnchorsFamilyId 设置家族id
func SetFamilyAnchorsFamilyId(value int) FamilyAnchorsUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFamilyAnchorsFamilyId.ColumnName().String()] = value
	}
}

// SetFamilyAnchorsPatriarchId 设置族长id
func SetFamilyAnchorsPatriarchId(value int) FamilyAnchorsUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFamilyAnchorsPatriarchId.ColumnName().String()] = value
	}
}

// SetFamilyAnchorsUserId 设置主播id
func SetFamilyAnchorsUserId(value int) FamilyAnchorsUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFamilyAnchorsUserId.ColumnName().String()] = value
	}
}

// SetFamilyAnchorsStatus 设置状态 1-正常 3-踢出
func SetFamilyAnchorsStatus(value int) FamilyAnchorsUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFamilyAnchorsStatus.ColumnName().String()] = value
	}
}

// SetFamilyAnchorsCreatedAt 设置创建时间
func SetFamilyAnchorsCreatedAt(value time.Time) FamilyAnchorsUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFamilyAnchorsCreatedAt.ColumnName().String()] = value
	}
}

// SetFamilyAnchorsUpdatedAt 设置更新时间
func SetFamilyAnchorsUpdatedAt(value time.Time) FamilyAnchorsUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFamilyAnchorsUpdatedAt.ColumnName().String()] = value
	}
}
