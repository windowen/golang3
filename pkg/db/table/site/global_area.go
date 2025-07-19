package site

import (
	"time"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"serverApi/pkg/constant"
)

const TableNameGlobalArea = "sys_global_area"

var (
	FieldGlobalAreaId          = field.NewInt(TableNameGlobalArea, "id")
	FieldGlobalAreaMobileCode  = field.NewString(TableNameGlobalArea, "mobile_code")
	FieldGlobalAreaCountryCode = field.NewString(TableNameGlobalArea, "country_code")
	FieldGlobalAreaName        = field.NewString(TableNameGlobalArea, "name")
	FieldGlobalAreaStatus      = field.NewInt(TableNameGlobalArea, "status")
	FieldGlobalAreaCreatedAt   = field.NewTime(TableNameGlobalArea, "created_at")
	FieldGlobalAreaDeletedAt   = field.NewTime(TableNameGlobalArea, "deleted_at")
)

type GlobalArea struct {
	Id          int       `gorm:"column:id" json:"id"`                     // ID
	MobileCode  string    `gorm:"column:mobile_code" json:"mobile_code"`   // 手机编码(如855)
	CountryCode string    `gorm:"column:country_code" json:"country_code"` // 国家代码（如BR巴西）
	Name        string    `gorm:"column:name" json:"name"`                 // 国家名称
	Status      int       `gorm:"column:status" json:"status"`             // 1-启用2-禁用 3-删除
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`     // 创建时间
	DeletedAt   time.Time `gorm:"column:deleted_at" json:"deleted_at"`     // 删除时间
}

func (*GlobalArea) TableName() string {
	return TableNameGlobalArea
}

func (a *GlobalArea) IsEmpty() bool {
	if a == nil {
		return true
	}
	return a.Id == 0
}

// DefaultScope 过滤已删除 status=3
func (a *GlobalArea) DefaultScope(db *gorm.DB) *gorm.DB {
	return db.Where(FieldGlobalAreaStatus.Eq(constant.StatusNormal))
}

// GlobalAreaWhereOption 条件项
type GlobalAreaWhereOption func(*gorm.DB) *gorm.DB

// WhereGlobalAreaId 查询ID
func WhereGlobalAreaId(value int) GlobalAreaWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldGlobalAreaId.Eq(value))
	}
}

// WhereGlobalAreaMobileCode 查询手机编码(如855)
func WhereGlobalAreaMobileCode(value string) GlobalAreaWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldGlobalAreaMobileCode.Eq(value))
	}
}

// WhereGlobalAreaCountryCode 查询国家代码（如BR巴西）
func WhereGlobalAreaCountryCode(value string) GlobalAreaWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldGlobalAreaCountryCode.Eq(value))
	}
}

// WhereGlobalAreaName 查询国家名称
func WhereGlobalAreaName(value string) GlobalAreaWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldGlobalAreaName.Eq(value))
	}
}

// WhereGlobalAreaStatus 查询1-启用2-禁用 3-删除
func WhereGlobalAreaStatus(value int) GlobalAreaWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldGlobalAreaStatus.Eq(value))
	}
}

// WhereGlobalAreaCreatedAt 查询创建时间
func WhereGlobalAreaCreatedAt(value time.Time) GlobalAreaWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldGlobalAreaCreatedAt.Eq(value))
	}
}

// WhereGlobalAreaDeletedAt 查询删除时间
func WhereGlobalAreaDeletedAt(value time.Time) GlobalAreaWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldGlobalAreaDeletedAt.Eq(value))
	}
}

// GlobalAreaUpdateOption 修改项
type GlobalAreaUpdateOption func(map[string]interface{})

// SetGlobalAreaId 设置ID
func SetGlobalAreaId(value int) GlobalAreaUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldGlobalAreaId.ColumnName().String()] = value
	}
}

// SetGlobalAreaMobileCode 设置手机编码(如855)
func SetGlobalAreaMobileCode(value string) GlobalAreaUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldGlobalAreaMobileCode.ColumnName().String()] = value
	}
}

// SetGlobalAreaCountryCode 设置国家代码（如BR巴西）
func SetGlobalAreaCountryCode(value string) GlobalAreaUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldGlobalAreaCountryCode.ColumnName().String()] = value
	}
}

// SetGlobalAreaName 设置国家名称
func SetGlobalAreaName(value string) GlobalAreaUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldGlobalAreaName.ColumnName().String()] = value
	}
}

// SetGlobalAreaStatus 设置1-启用2-禁用 3-删除
func SetGlobalAreaStatus(value int) GlobalAreaUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldGlobalAreaStatus.ColumnName().String()] = value
	}
}

// SetGlobalAreaCreatedAt 设置创建时间
func SetGlobalAreaCreatedAt(value time.Time) GlobalAreaUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldGlobalAreaCreatedAt.ColumnName().String()] = value
	}
}

// SetGlobalAreaDeletedAt 设置删除时间
func SetGlobalAreaDeletedAt(value time.Time) GlobalAreaUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldGlobalAreaDeletedAt.ColumnName().String()] = value
	}
}
