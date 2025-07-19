package site

import "time"

type ActivityMsg struct {
	Id           int       `gorm:"column:id" json:"id"`                       // 交易消息表Id
	Category     int       `gorm:"column:category" json:"category"`           // 活动类型(不用管默认是0)
	Content      string    `gorm:"column:content" json:"content"`             // 通知内容
	CountryCode  string    `gorm:"column:country_code" json:"country_code"`   // 国家编码(每个国家语言不一样)
	BatchNumber  string    `gorm:"column:batch_number" json:"batch_number"`   // 批次号
	LanguageCode string    `gorm:"column:language_code" json:"language_code"` // 语言编码
	Status       int       `gorm:"column:status" json:"status"`               // 0-禁用1-启用
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`       // 创建时间
}

func (ActivityMsg) TableName() string {
	return "site_activity_msg"
}
