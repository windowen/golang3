package site

import "time"

type StartupImage struct {
	Id           int       `gorm:"column:id" json:"id"`                       // 启动图主键id
	BatchNumber  string    `gorm:"column:batch_number" json:"batch_number"`   // 批次号(因为每种语言一条记录,需要这个字段关联这几条记录,类似于parentId，列表只显示一条的话查询language_code是zh中文的就行)
	CountryCode  string    `gorm:"column:country_code" json:"country_code"`   // 国家编码
	LanguageCode string    `gorm:"column:language_code" json:"language_code"` // 语言编码
	Name         string    `gorm:"column:name" json:"name"`                   // 启动图名称
	ImageUrl     string    `gorm:"column:image_url" json:"image_url"`         // 启动图地址
	ThumbnailUrl string    `gorm:"column:thumbnail_url" json:"thumbnail_url"` // 启动图缩略图地址
	Status       int       `gorm:"column:status" json:"status"`               // 1-启用2-禁用
	Sort         int       `gorm:"column:sort" json:"sort"`                   // 排序
	Remark       string    `gorm:"column:remark" json:"remark"`               // 备注
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`       // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`       // 更新时间
}

func (StartupImage) TableName() string {
	return "site_startup_image"
}
