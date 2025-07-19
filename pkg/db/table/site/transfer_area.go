package site

import "time"

type TransferArea struct {
	Id        int       `gorm:"column:id" json:"id"`                 // ID
	Name      string    `gorm:"column:name" json:"name"`             // 地区名称
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"` // 创建时间
	DeletedAt time.Time `gorm:"column:deleted_at" json:"deleted_at"` // 删除时间
}

func (TransferArea) TableName() string {
	return "transfer_area"
}
