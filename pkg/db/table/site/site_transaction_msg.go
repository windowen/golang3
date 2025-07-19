package site

import "time"

type TransactionMsg struct {
	Id          int       `gorm:"column:id" json:"id"`                     // 交易消息表Id
	Category    int       `gorm:"column:category" json:"category"`         // 通知类型1-充值成功通知2-充值失败通知3-提现成功通知4-提现失败通知5-兑换钻石
	UserId      int       `gorm:"column:user_id" json:"user_id"`           // 接收通知的用户Id
	CountryCode string    `gorm:"column:country_code" json:"country_code"` // 国家编码
	Content     string    `gorm:"column:content" json:"content"`           // 通知内容
	ReadStatus  int       `gorm:"column:read_status" json:"read_status"`   // 0-未读1-已读
	Status      int       `gorm:"column:status" json:"status"`             // 0-未删除1-已删除
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`     // 创建时间
}

func (TransactionMsg) TableName() string {
	return "site_transaction_msg"
}
