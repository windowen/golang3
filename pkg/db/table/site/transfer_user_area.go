package site

import "time"

type TransferUserArea struct {
	Id        int       `gorm:"column:id" json:"id"`                 // ID
	UserId    int       `gorm:"column:user_id" json:"user_id"`       // 用户ID
	AreaId    int       `gorm:"column:area_id" json:"area_id"`       // 地区id
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"` // 创建时间
	DeletedId time.Time `gorm:"column:deleted_id" json:"deleted_id"` // 删除时间
}

func (TransferUserArea) TableName() string {
	return "transfer_user_area"
}
