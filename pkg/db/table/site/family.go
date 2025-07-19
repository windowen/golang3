package site

import (
	"time"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"serverApi/pkg/constant"
)

const TableNameFamily = "live_family"

var (
	FieldFamilyId        = field.NewInt(TableNameFamily, "id")
	FieldFamilyUserId    = field.NewInt(TableNameFamily, "user_id")
	FieldFamilyName      = field.NewString(TableNameFamily, "name")
	FieldFamilyRatio     = field.NewInt(TableNameFamily, "ratio")
	FieldFamilyRemark    = field.NewString(TableNameFamily, "remark")
	FieldFamilyStatus    = field.NewInt(TableNameFamily, "status")
	FieldFamilyCreatedAt = field.NewTime(TableNameFamily, "created_at")
	FieldFamilyUpdatedAt = field.NewTime(TableNameFamily, "updated_at")
)

type Family struct {
	Id        int       `gorm:"column:id" json:"id"`                 // ID
	UserId    int       `gorm:"column:user_id" json:"user_id"`       // 家族长id
	Name      string    `gorm:"column:name" json:"name"`             // 家族名称
	Ratio     int       `gorm:"column:ratio" json:"ratio"`           // 平台抽成比例  实际*100
	Remark    string    `gorm:"column:remark" json:"remark"`         // 备注
	Status    int       `gorm:"column:status" json:"status"`         // 状态 1-正常 2-禁用 3-删除
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"` // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"` // 更新时间
}

func (*Family) TableName() string {
	return TableNameFamily
}

func (a *Family) IsEmpty() bool {
	if a == nil {
		return true
	}
	return a.Id == 0
}

// DefaultScope 获取正常 status=1
func (a *Family) DefaultScope(db *gorm.DB) *gorm.DB {
	return db.Where(FieldFamilyStatus.Eq(constant.StatusNormal))
}

// FamilyWhereOption 条件项
type FamilyWhereOption func(*gorm.DB) *gorm.DB

// WhereFamilyId 查询ID
func WhereFamilyId(value int) FamilyWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFamilyId.Eq(value))
	}
}

// WhereFamilyUserId 查询家族长id
func WhereFamilyUserId(value int) FamilyWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFamilyUserId.Eq(value))
	}
}

// WhereFamilyName 查询家族名称
func WhereFamilyName(value string) FamilyWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFamilyName.Eq(value))
	}
}

// WhereFamilyRatio 查询平台抽成比例  实际*100
func WhereFamilyRatio(value int) FamilyWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFamilyRatio.Eq(value))
	}
}

// WhereFamilyRemark 查询备注
func WhereFamilyRemark(value string) FamilyWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFamilyRemark.Eq(value))
	}
}

// WhereFamilyStatus 查询状态 1-正常 2-禁用 3-删除
func WhereFamilyStatus(value int) FamilyWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFamilyStatus.Eq(value))
	}
}

// WhereFamilyCreatedAt 查询创建时间
func WhereFamilyCreatedAt(value time.Time) FamilyWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFamilyCreatedAt.Eq(value))
	}
}

// WhereFamilyUpdatedAt 查询更新时间
func WhereFamilyUpdatedAt(value time.Time) FamilyWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldFamilyUpdatedAt.Eq(value))
	}
}

// FamilyUpdateOption 修改项
type FamilyUpdateOption func(map[string]interface{})

// SetFamilyId 设置ID
func SetFamilyId(value int) FamilyUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFamilyId.ColumnName().String()] = value
	}
}

// SetFamilyUserId 设置家族长id
func SetFamilyUserId(value int) FamilyUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFamilyUserId.ColumnName().String()] = value
	}
}

// SetFamilyName 设置家族名称
func SetFamilyName(value string) FamilyUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFamilyName.ColumnName().String()] = value
	}
}

// SetFamilyRatio 设置平台抽成比例  实际*100
func SetFamilyRatio(value int) FamilyUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFamilyRatio.ColumnName().String()] = value
	}
}

// SetFamilyRemark 设置备注
func SetFamilyRemark(value string) FamilyUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFamilyRemark.ColumnName().String()] = value
	}
}

// SetFamilyStatus 设置状态 1-正常 2-禁用 3-删除
func SetFamilyStatus(value int) FamilyUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFamilyStatus.ColumnName().String()] = value
	}
}

// SetFamilyCreatedAt 设置创建时间
func SetFamilyCreatedAt(value time.Time) FamilyUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFamilyCreatedAt.ColumnName().String()] = value
	}
}

// SetFamilyUpdatedAt 设置更新时间
func SetFamilyUpdatedAt(value time.Time) FamilyUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldFamilyUpdatedAt.ColumnName().String()] = value
	}
}
