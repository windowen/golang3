package site

import (
	"time"

	"gorm.io/gen/field"
	"gorm.io/gorm"
)

const TableNameSysConfig = "sys_site_config"

var (
	FieldSysConfigId         = field.NewInt(TableNameSysConfig, "id")
	FieldSysConfigCategory   = field.NewInt(TableNameSysConfig, "category")
	FieldSysConfigName       = field.NewString(TableNameSysConfig, "name")
	FieldSysConfigConfigCode = field.NewString(TableNameSysConfig, "config_code")
	FieldSysConfigContent    = field.NewString(TableNameSysConfig, "content")
	FieldSysConfigStatus     = field.NewInt(TableNameSysConfig, "status")
	FieldSysConfigCreatedAt  = field.NewTime(TableNameSysConfig, "created_at")
)

type SysConfig struct {
	Id         int       `gorm:"column:id" json:"id"`
	Category   int       `gorm:"column:category" json:"category"`       // 配置类型 1- App配置 2-直播配置 3-财务配置
	Name       string    `gorm:"column:name" json:"name"`               // 配置名称
	ConfigCode string    `gorm:"column:config_code" json:"config_code"` // 配置code 唯一
	Content    string    `gorm:"column:content" json:"content"`         // 配置内容
	Status     int       `gorm:"column:status" json:"status"`           // 状态0-停用1-启用
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`   // 创建日期
}

func (*SysConfig) TableName() string {
	return TableNameSysConfig
}

func (a *SysConfig) IsEmpty() bool {
	if a == nil {
		return true
	}
	return a.Id == 0
}

// SysConfigWhereOption 条件项
type SysConfigWhereOption func(*gorm.DB) *gorm.DB

// WhereSysConfigId 查询
func WhereSysConfigId(value int) SysConfigWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldSysConfigId.Eq(value))
	}
}

// WhereSysConfigCategory 查询配置类型 1- App配置 2-直播配置 3-财务配置
func WhereSysConfigCategory(value int) SysConfigWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldSysConfigCategory.Eq(value))
	}
}

// WhereSysConfigName 查询配置名称
func WhereSysConfigName(value string) SysConfigWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldSysConfigName.Eq(value))
	}
}

// WhereSysConfigConfigCode 查询配置code 唯一
func WhereSysConfigConfigCode(value string) SysConfigWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldSysConfigConfigCode.Eq(value))
	}
}

// WhereSysConfigContent 查询配置内容
func WhereSysConfigContent(value string) SysConfigWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldSysConfigContent.Eq(value))
	}
}

// WhereSysConfigStatus 查询状态0-停用1-启用
func WhereSysConfigStatus(value int) SysConfigWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldSysConfigStatus.Eq(value))
	}
}

// WhereSysConfigCreatedAt 查询创建日期
func WhereSysConfigCreatedAt(value time.Time) SysConfigWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldSysConfigCreatedAt.Eq(value))
	}
}

// SysConfigUpdateOption 修改项
type SysConfigUpdateOption func(map[string]interface{})

// SetSysConfigId 设置
func SetSysConfigId(value int) SysConfigUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldSysConfigId.ColumnName().String()] = value
	}
}

// SetSysConfigCategory 设置配置类型 1- App配置 2-直播配置 3-财务配置
func SetSysConfigCategory(value int) SysConfigUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldSysConfigCategory.ColumnName().String()] = value
	}
}

// SetSysConfigName 设置配置名称
func SetSysConfigName(value string) SysConfigUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldSysConfigName.ColumnName().String()] = value
	}
}

// SetSysConfigConfigCode 设置配置code 唯一
func SetSysConfigConfigCode(value string) SysConfigUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldSysConfigConfigCode.ColumnName().String()] = value
	}
}

// SetSysConfigContent 设置配置内容
func SetSysConfigContent(value string) SysConfigUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldSysConfigContent.ColumnName().String()] = value
	}
}

// SetSysConfigStatus 设置状态0-停用1-启用
func SetSysConfigStatus(value int) SysConfigUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldSysConfigStatus.ColumnName().String()] = value
	}
}

// SetSysConfigCreatedAt 设置创建日期
func SetSysConfigCreatedAt(value time.Time) SysConfigUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldSysConfigCreatedAt.ColumnName().String()] = value
	}
}
