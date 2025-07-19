package site

import (
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

const TableNameUserAuth = "site_user_auth"

var (
	FieldUserAuthId       = field.NewInt(TableNameUserAuth, "id")
	FieldUserAuthUserId   = field.NewInt(TableNameUserAuth, "user_id")
	FieldUserAuthAreaCode = field.NewString(TableNameUserAuth, "area_code")
	FieldUserAuthMobile   = field.NewString(TableNameUserAuth, "mobile")
	FieldUserAuthEmail    = field.NewString(TableNameUserAuth, "email")
	FieldUserAuthPassword = field.NewString(TableNameUserAuth, "password")
)

type UserAuth struct {
	Id       int    `gorm:"column:id" json:"id"`               // ID
	UserId   int    `gorm:"column:user_id" json:"user_id"`     // 用户ID
	AreaCode string `gorm:"column:area_code" json:"area_code"` // 手机国家码
	Mobile   string `gorm:"column:mobile" json:"mobile"`       // 手机号
	Email    string `gorm:"column:email" json:"email"`         // 邮箱
	Password string `gorm:"column:password" json:"password"`   // 站内账号是密码、第三方登录是Token
}

func (*UserAuth) TableName() string {
	return TableNameUserAuth
}

func (a *UserAuth) IsEmpty() bool {
	if a == nil {
		return true
	}
	return a.Id == 0
}

// UserAuthWhereOption 条件项
type UserAuthWhereOption func(*gorm.DB) *gorm.DB

// WhereUserAuthId 查询ID
func WhereUserAuthId(value int) UserAuthWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserAuthId.Eq(value))
	}
}

// WhereUserAuthUserId 查询用户ID
func WhereUserAuthUserId(value int) UserAuthWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserAuthUserId.Eq(value))
	}
}

// WhereUserAuthAreaCode 查询手机国家码
func WhereUserAuthAreaCode(value string) UserAuthWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		if value == "" {
			return db
		}

		return db.Where(FieldUserAuthAreaCode.Eq(value))
	}
}

// WhereUserAuthMobile 查询手机号
func WhereUserAuthMobile(value string) UserAuthWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		if value == "" {
			return db
		}

		return db.Where(FieldUserAuthMobile.Eq(value))
	}
}

// WhereUserAuthEmail 查询邮箱
func WhereUserAuthEmail(value string) UserAuthWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		if value == "" {
			return db
		}

		return db.Where(FieldUserAuthEmail.Eq(value))
	}
}

// WhereUserAuthPassword 查询站内账号是密码、第三方登录是Token
func WhereUserAuthPassword(value string) UserAuthWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserAuthPassword.Eq(value))
	}
}

// UserAuthUpdateOption 修改项
type UserAuthUpdateOption func(map[string]interface{})

// SetUserAuthId 设置ID
func SetUserAuthId(value int) UserAuthUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserAuthId.ColumnName().String()] = value
	}
}

// SetUserAuthUserId 设置用户ID
func SetUserAuthUserId(value int) UserAuthUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserAuthUserId.ColumnName().String()] = value
	}
}

// SetUserAuthAreaCode 设置手机国家码
func SetUserAuthAreaCode(value string) UserAuthUpdateOption {
	return func(updates map[string]interface{}) {
		if value != "" {
			updates[FieldUserAuthAreaCode.ColumnName().String()] = value
		}
	}
}

// SetUserAuthMobile 设置手机号
func SetUserAuthMobile(value string) UserAuthUpdateOption {
	return func(updates map[string]interface{}) {
		if value != "" {
			updates[FieldUserAuthMobile.ColumnName().String()] = value
		}
	}
}

// SetUserAuthEmail 设置邮箱
func SetUserAuthEmail(value string) UserAuthUpdateOption {
	return func(updates map[string]interface{}) {
		if value != "" {
			updates[FieldUserAuthEmail.ColumnName().String()] = value
		}
	}
}

// SetUserAuthPassword 设置站内账号是密码、第三方登录是Token
func SetUserAuthPassword(value string) UserAuthUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserAuthPassword.ColumnName().String()] = value
	}
}
