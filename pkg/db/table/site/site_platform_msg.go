package site

import "time"

type PlatformMsg struct {
	Id           int       `gorm:"column:id" json:"id"`                       // 平台消息表Id
	Content      string    `gorm:"column:content" json:"content"`             // 平台消息内容
	CountryCode  string    `gorm:"column:country_code" json:"country_code"`   // 国家编码(每个国家语言不一样)
	LanguageCode string    `gorm:"column:language_code" json:"language_code"` // 语言编码
	BatchNumber  string    `gorm:"column:batch_number" json:"batch_number"`   // 批次号
	Status       int       `gorm:"column:status" json:"status"`               // 0-禁用1-启用
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

func (PlatformMsg) TableName() string {
	return "site_platform_msg"
}
