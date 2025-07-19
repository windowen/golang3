package site

import "time"

type UserLoginLogs struct {
	Id          int       `gorm:"column:id" json:"id"` // ID
	UserId      int       `gorm:"column:user_id" json:"user_id"`
	DeviceType  int       `gorm:"column:device_type" json:"device_type"`   // 设备类型 1.ios 2.android 3.web 4.爬虫 5.其它
	DeviceModel string    `gorm:"column:device_model" json:"device_model"` // 设备型号
	Ip          string    `gorm:"column:ip" json:"ip"`                     // 访问IP地址
	IsDel       int       `gorm:"column:is_del" json:"is_del"`             // 是否删除 0-正常 1-删除
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`     // 创建时间
	DeletedAt   time.Time `gorm:"column:deleted_at" json:"deleted_at"`     // 删除时间
}

func (UserLoginLogs) UserLoginLogs() string {
	return "site_user_login_logs"
}

func (ull *UserLoginLogs) IsEmpty() bool {
	if ull == nil {
		return true
	}
	return ull.Id == 0
}
