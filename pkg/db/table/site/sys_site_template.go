package site

import (
	"time"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"serverApi/pkg/constant"
)

const TableNameSysTemplate = "sys_site_template"

var (
	FieldSysTemplateId              = field.NewInt(TableNameSysTemplate, "id")
	FieldSysTemplateCountryCode     = field.NewString(TableNameSysTemplate, "country_code")
	FieldSysTemplateLanguageCode    = field.NewString(TableNameSysTemplate, "language_code")
	FieldSysTemplateCategory        = field.NewInt(TableNameSysTemplate, "category")
	FieldSysTemplateTemplateCode    = field.NewString(TableNameSysTemplate, "template_code")
	FieldSysTemplateSubject         = field.NewString(TableNameSysTemplate, "subject")
	FieldSysTemplateTemplateContent = field.NewString(TableNameSysTemplate, "template_content")
	FieldSysTemplateStatus          = field.NewInt(TableNameSysTemplate, "status")
	FieldSysTemplateCreatedAt       = field.NewTime(TableNameSysTemplate, "created_at")
)

type SysTemplate struct {
	Id              int       `gorm:"column:id" json:"id"`                             // 站内信模板表Id
	CountryCode     string    `gorm:"column:country_code" json:"country_code"`         // 国家编码
	LanguageCode    string    `gorm:"column:language_code" json:"language_code"`       // 语言 短信、邮件模版使用
	Category        int       `gorm:"column:category" json:"category"`                 // 通知类型1-充值成功通知2-充值失败通知3-提现成功通知4-提现失败通知5-兑换钻石6-管理后台通知7-公告维护 8-短信 9-邮件
	TemplateCode    string    `gorm:"column:template_code" json:"template_code"`       // 模版code  如短信登陆、注册模版
	Subject         string    `gorm:"column:subject" json:"subject"`                   // 邮件主题
	TemplateContent string    `gorm:"column:template_content" json:"template_content"` // 模板内容
	Status          int       `gorm:"column:status" json:"status"`                     // 状态0-停用1-启用
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`             // 创建日期
}

func (*SysTemplate) TableName() string {
	return TableNameSysTemplate
}

func (a *SysTemplate) IsEmpty() bool {
	if a == nil {
		return true
	}
	return a.Id == 0
}

// DefaultScope 过滤已删除 status=3
func (a *SysTemplate) DefaultScope(db *gorm.DB) *gorm.DB {
	return db.Where(FieldSysTemplateStatus.Eq(constant.StatusNormal))
}

// SysTemplateWhereOption 条件项
type SysTemplateWhereOption func(*gorm.DB) *gorm.DB

// WhereSysTemplateId 查询站内信模板表Id
func WhereSysTemplateId(value int) SysTemplateWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldSysTemplateId.Eq(value))
	}
}

// WhereSysTemplateCountryCode 查询国家编码
func WhereSysTemplateCountryCode(value string) SysTemplateWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldSysTemplateCountryCode.Eq(value))
	}
}

// WhereSysTemplateLanguageCode 查询语言 短信、邮件模版使用
func WhereSysTemplateLanguageCode(value string) SysTemplateWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldSysTemplateLanguageCode.Eq(value))
	}
}

// WhereSysTemplateCategory 查询通知类型1-充值成功通知2-充值失败通知3-提现成功通知4-提现失败通知5-兑换钻石6-管理后台通知7-公告维护 8-短信 9-邮件
func WhereSysTemplateCategory(value int) SysTemplateWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldSysTemplateCategory.Eq(value))
	}
}

// WhereSysTemplateTemplateCode 查询模版code  如短信登陆、注册模版
func WhereSysTemplateTemplateCode(value string) SysTemplateWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldSysTemplateTemplateCode.Eq(value))
	}
}

// WhereSysTemplateSubject 查询邮件主题
func WhereSysTemplateSubject(value string) SysTemplateWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldSysTemplateSubject.Eq(value))
	}
}

// WhereSysTemplateTemplateContent 查询模板内容
func WhereSysTemplateTemplateContent(value string) SysTemplateWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldSysTemplateTemplateContent.Eq(value))
	}
}

// WhereSysTemplateStatus 查询状态0-停用1-启用
func WhereSysTemplateStatus(value int) SysTemplateWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldSysTemplateStatus.Eq(value))
	}
}

// WhereSysTemplateCreatedAt 查询创建日期
func WhereSysTemplateCreatedAt(value time.Time) SysTemplateWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldSysTemplateCreatedAt.Eq(value))
	}
}

// SysTemplateUpdateOption 修改项
type SysTemplateUpdateOption func(map[string]interface{})

// SetSysTemplateId 设置站内信模板表Id
func SetSysTemplateId(value int) SysTemplateUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldSysTemplateId.ColumnName().String()] = value
	}
}

// SetSysTemplateCountryCode 设置国家编码
func SetSysTemplateCountryCode(value string) SysTemplateUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldSysTemplateCountryCode.ColumnName().String()] = value
	}
}

// SetSysTemplateLanguageCode 设置语言 短信、邮件模版使用
func SetSysTemplateLanguageCode(value string) SysTemplateUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldSysTemplateLanguageCode.ColumnName().String()] = value
	}
}

// SetSysTemplateCategory 设置通知类型1-充值成功通知2-充值失败通知3-提现成功通知4-提现失败通知5-兑换钻石6-管理后台通知7-公告维护 8-短信 9-邮件
func SetSysTemplateCategory(value int) SysTemplateUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldSysTemplateCategory.ColumnName().String()] = value
	}
}

// SetSysTemplateTemplateCode 设置模版code  如短信登陆、注册模版
func SetSysTemplateTemplateCode(value string) SysTemplateUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldSysTemplateTemplateCode.ColumnName().String()] = value
	}
}

// SetSysTemplateSubject 设置邮件主题
func SetSysTemplateSubject(value string) SysTemplateUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldSysTemplateSubject.ColumnName().String()] = value
	}
}

// SetSysTemplateTemplateContent 设置模板内容
func SetSysTemplateTemplateContent(value string) SysTemplateUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldSysTemplateTemplateContent.ColumnName().String()] = value
	}
}

// SetSysTemplateStatus 设置状态0-停用1-启用
func SetSysTemplateStatus(value int) SysTemplateUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldSysTemplateStatus.ColumnName().String()] = value
	}
}

// SetSysTemplateCreatedAt 设置创建日期
func SetSysTemplateCreatedAt(value time.Time) SysTemplateUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldSysTemplateCreatedAt.ColumnName().String()] = value
	}
}
