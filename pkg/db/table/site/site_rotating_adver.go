package site

import "time"

type RotatingAdver struct {
	Id           int       `gorm:"column:id" json:"id"`                       // 走马灯主键id
	CountryCode  string    `gorm:"column:country_code" json:"country_code"`   // 国家编码
	LanguageCode string    `gorm:"column:language_code" json:"language_code"` // 语言编码
	BatchNumber  string    `gorm:"column:batch_number" json:"batch_number"`   // 批次号
	Name         string    `gorm:"column:name" json:"name"`                   // 走马灯名称
	Category     int       `gorm:"column:category" json:"category"`           // 走马灯类别1-热门2-场馆游戏3-充值提现
	Content      string    `gorm:"column:content" json:"content"`             // 走马灯内容
	Status       int       `gorm:"column:status" json:"status"`               // 1-启用2-禁用
	Sort         int       `gorm:"column:sort" json:"sort"`                   // 排序
	Remark       string    `gorm:"column:remark" json:"remark"`               // 备注
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`       // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`       // 更新时间
}

func (RotatingAdver) TableName() string {
	return "site_rotating_adver"
}
