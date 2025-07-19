package site

import (
	"time"

	"gorm.io/gen/field"
	"gorm.io/gorm"
)

const TableNameUserWallet = "site_user_wallet"

var (
	FieldUserWalletId             = field.NewInt(TableNameUserWallet, "id")
	FieldUserWalletUserId         = field.NewInt(TableNameUserWallet, "user_id")
	FieldUserWalletBalance        = field.NewFloat64(TableNameUserWallet, "balance")
	FieldUserWalletFreeze         = field.NewFloat64(TableNameUserWallet, "freeze")
	FieldUserWalletDiamond        = field.NewInt(TableNameUserWallet, "diamond")
	FieldUserWalletRechargeAmount = field.NewFloat64(TableNameUserWallet, "recharge_amount")
	FieldUserWalletRebateAmount   = field.NewFloat64(TableNameUserWallet, "rebate_amount")
	FieldUserWalletCreatedAt      = field.NewTime(TableNameUserWallet, "created_at")
	FieldUserWalletUpdatedAt      = field.NewTime(TableNameUserWallet, "updated_at")
)

type UserWallet struct {
	Id              int       `gorm:"column:id" json:"id"` // ID
	UserId          int       `gorm:"column:user_id" json:"user_id"`
	Balance         float64   `gorm:"column:balance" json:"balance"`                     // 总余额(美金)
	Freeze          float64   `gorm:"column:freeze" json:"freeze"`                       // 冻结金额
	Diamond         int       `gorm:"column:diamond" json:"diamond"`                     // 钻石(只有整数)
	RechargeAmount  float64   `gorm:"column:recharge_amount" json:"recharge_amount"`     // 充值金额
	RebateAmount    float64   `gorm:"column:rebate_amount" json:"rebate_amount"`         // 返点金额
	Trc20Address    string    `gorm:"column:trc20_address" json:"trc20_address"`         // trc20支付地址
	Trc20PrivateKey string    `gorm:"column:trc20_private_key" json:"trc20_private_key"` // trc20支付私钥
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`               // 创建时间
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`               // 更新时间
}

func (*UserWallet) TableName() string {
	return TableNameUserWallet
}

func (a *UserWallet) IsEmpty() bool {
	if a == nil {
		return true
	}
	return a.Id == 0
}

// UserWalletWhereOption 条件项
type UserWalletWhereOption func(*gorm.DB) *gorm.DB

// WhereUserWalletId 查询ID
func WhereUserWalletId(value int) UserWalletWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserWalletId.Eq(value))
	}
}

// WhereUserWalletUserId 查询
func WhereUserWalletUserId(value int) UserWalletWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserWalletUserId.Eq(value))
	}
}

// WhereUserWalletBalance 查询总余额(美金)
func WhereUserWalletBalance(value float64) UserWalletWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserWalletBalance.Eq(value))
	}
}

// WhereUserWalletFreeze 查询冻结金额
func WhereUserWalletFreeze(value float64) UserWalletWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserWalletFreeze.Eq(value))
	}
}

// WhereUserWalletDiamond 查询钻石(只有整数)
func WhereUserWalletDiamond(value int) UserWalletWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserWalletDiamond.Eq(value))
	}
}

// WhereUserWalletRechargeAmount 查询充值金额
func WhereUserWalletRechargeAmount(value float64) UserWalletWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserWalletRechargeAmount.Eq(value))
	}
}

// WhereUserWalletRebateAmount 查询返点金额
func WhereUserWalletRebateAmount(value float64) UserWalletWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserWalletRebateAmount.Eq(value))
	}
}

// WhereUserWalletCreatedAt 查询创建时间
func WhereUserWalletCreatedAt(value time.Time) UserWalletWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserWalletCreatedAt.Eq(value))
	}
}

// WhereUserWalletUpdatedAt 查询更新时间
func WhereUserWalletUpdatedAt(value time.Time) UserWalletWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserWalletUpdatedAt.Eq(value))
	}
}

// UserWalletUpdateOption 修改项
type UserWalletUpdateOption func(map[string]interface{})

// SetUserWalletId 设置ID
func SetUserWalletId(value int) UserWalletUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserWalletId.ColumnName().String()] = value
	}
}

// SetUserWalletUserId 设置
func SetUserWalletUserId(value int) UserWalletUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserWalletUserId.ColumnName().String()] = value
	}
}

// SetUserWalletBalance 设置总余额(美金)
func SetUserWalletBalance(value float64) UserWalletUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserWalletBalance.ColumnName().String()] = value
	}
}

// SetUserWalletFreeze 设置冻结金额
func SetUserWalletFreeze(value float64) UserWalletUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserWalletFreeze.ColumnName().String()] = value
	}
}

// SetUserWalletDiamond 设置钻石(只有整数)
func SetUserWalletDiamond(value int) UserWalletUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserWalletDiamond.ColumnName().String()] = value
	}
}

// SetUserWalletRechargeAmount 设置充值金额
func SetUserWalletRechargeAmount(value float64) UserWalletUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserWalletRechargeAmount.ColumnName().String()] = value
	}
}

// SetUserWalletRebateAmount 设置返点金额
func SetUserWalletRebateAmount(value float64) UserWalletUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserWalletRebateAmount.ColumnName().String()] = value
	}
}

// SetUserWalletCreatedAt 设置创建时间
func SetUserWalletCreatedAt(value time.Time) UserWalletUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserWalletCreatedAt.ColumnName().String()] = value
	}
}

// SetUserWalletUpdatedAt 设置更新时间
func SetUserWalletUpdatedAt(value time.Time) UserWalletUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserWalletUpdatedAt.ColumnName().String()] = value
	}
}
