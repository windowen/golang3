package site

import (
	"time"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"serverApi/pkg/constant"
)

const TableNameGlobalLang = "sys_global_language"

var (
	FieldGlobalLangId           = field.NewInt(TableNameGlobalLang, "id")
	FieldGlobalLangLanguageCode = field.NewString(TableNameGlobalLang, "language_code")
	FieldGlobalLangLanguageName = field.NewString(TableNameGlobalLang, "language_name")
	FieldGlobalLangStatus       = field.NewInt(TableNameGlobalLang, "status")
	FieldGlobalLangCreatedAt    = field.NewTime(TableNameGlobalLang, "created_at")
	FieldGlobalLangDeletedAt    = field.NewTime(TableNameGlobalLang, "deleted_at")
)

type GlobalLang struct {
	Id           int       `gorm:"column:id" json:"id"`                       // ID
	LanguageCode string    `gorm:"column:language_code" json:"language_code"` // 语言代码( 如en表示英语)
	LanguageName string    `gorm:"column:language_name" json:"language_name"` // 语言名称
	Status       int       `gorm:"column:status" json:"status"`               // 1-启用2-禁用 3-删除
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`       // 创建时间
	DeletedAt    time.Time `gorm:"column:deleted_at" json:"deleted_at"`       // 删除时间
}

func (*GlobalLang) TableName() string {
	return TableNameGlobalLang
}

func (a *GlobalLang) IsEmpty() bool {
	if a == nil {
		return true
	}
	return a.Id == 0
}

// DefaultScope 过滤已删除 status=3
func (a *GlobalLang) DefaultScope(db *gorm.DB) *gorm.DB {
	return db.Where(FieldGlobalLangStatus.Eq(constant.StatusNormal))
}

// GlobalLangWhereOption 条件项
type GlobalLangWhereOption func(*gorm.DB) *gorm.DB

// WhereGlobalLangId 查询ID
func WhereGlobalLangId(value int) GlobalLangWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldGlobalLangId.Eq(value))
	}
}

// WhereGlobalLangLanguageCode 查询语言代码( 如en表示英语)
func WhereGlobalLangLanguageCode(value string) GlobalLangWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldGlobalLangLanguageCode.Eq(value))
	}
}

// WhereGlobalLangLanguageName 查询语言名称
func WhereGlobalLangLanguageName(value string) GlobalLangWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldGlobalLangLanguageName.Eq(value))
	}
}

// WhereGlobalLangStatus 查询1-启用2-禁用 3-删除
func WhereGlobalLangStatus(value int) GlobalLangWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldGlobalLangStatus.Eq(value))
	}
}

// WhereGlobalLangCreatedAt 查询创建时间
func WhereGlobalLangCreatedAt(value time.Time) GlobalLangWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldGlobalLangCreatedAt.Eq(value))
	}
}

// WhereGlobalLangDeletedAt 查询删除时间
func WhereGlobalLangDeletedAt(value time.Time) GlobalLangWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldGlobalLangDeletedAt.Eq(value))
	}
}

// GlobalLangUpdateOption 修改项
type GlobalLangUpdateOption func(map[string]interface{})

// SetGlobalLangId 设置ID
func SetGlobalLangId(value int) GlobalLangUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldGlobalLangId.ColumnName().String()] = value
	}
}

// SetGlobalLangLanguageCode 设置语言代码( 如en表示英语)
func SetGlobalLangLanguageCode(value string) GlobalLangUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldGlobalLangLanguageCode.ColumnName().String()] = value
	}
}

// SetGlobalLangLanguageName 设置语言名称
func SetGlobalLangLanguageName(value string) GlobalLangUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldGlobalLangLanguageName.ColumnName().String()] = value
	}
}

// SetGlobalLangStatus 设置1-启用2-禁用 3-删除
func SetGlobalLangStatus(value int) GlobalLangUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldGlobalLangStatus.ColumnName().String()] = value
	}
}

// SetGlobalLangCreatedAt 设置创建时间
func SetGlobalLangCreatedAt(value time.Time) GlobalLangUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldGlobalLangCreatedAt.ColumnName().String()] = value
	}
}

// SetGlobalLangDeletedAt 设置删除时间
func SetGlobalLangDeletedAt(value time.Time) GlobalLangUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldGlobalLangDeletedAt.ColumnName().String()] = value
	}
}
