package site

import "time"

type Vip struct {
	Id        int       `gorm:"column:id" json:"id"`                 // ID
	Name      string    `gorm:"column:name" json:"name"`             // 等级名称
	Level     int       `gorm:"column:level" json:"level"`           // 等级
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"` // 创建时间
	DeletedAt time.Time `gorm:"column:deleted_at" json:"deleted_at"` // 删除时间
}

func (Vip) TableName() string {
	return "site_vip"
}
