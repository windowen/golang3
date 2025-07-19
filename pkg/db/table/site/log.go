package site

import "time"

type Log struct {
	Id            int       `gorm:"column:id" json:"id"`                         // ID
	UserId        int       `gorm:"column:user_id" json:"user_id"`               // 用户ID
	DeviceType    int       `gorm:"column:device_type" json:"device_type"`       // 设备类型 1.ios 2.android 3.web 4.爬虫 5.其它
	DeviceModel   string    `gorm:"column:device_model" json:"device_model"`     // 设备型号
	Resolution    string    `gorm:"column:resolution" json:"resolution"`         // 分辨率
	Uri           string    `gorm:"column:uri" json:"uri"`                       // URI
	RequestParam  string    `gorm:"column:request_param" json:"request_param"`   // 请求参数
	ResponseParam string    `gorm:"column:response_param" json:"response_param"` // 响应参数
	Ip            string    `gorm:"column:ip" json:"ip"`                         // 访问IP地址
	Area          string    `gorm:"column:area" json:"area"`                     // 请求者所在区域
	Elapsed       int       `gorm:"column:elapsed" json:"elapsed"`               // 请求耗时
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`         // 创建时间
}

func (Log) TableName() string {
	return "sys_log"
}
