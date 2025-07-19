package site

import "time"

type UserState struct {
	Id        int        `gorm:"column:id" json:"id"`                 // ID
	UserId    int        `gorm:"column:user_id" json:"user_id"`       // 用户ID
	Flag      string     `gorm:"column:flag" json:"flag"`             // 用户状态标识
	Value     int        `gorm:"column:value" json:"value"`           // 参数值
	Remark    string     `gorm:"column:remark" json:"remark"`         // 备注
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"` // 创建时间
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"` // 删除时间
}

func (UserState) TableName() string {
	return "site_user_state"
}
