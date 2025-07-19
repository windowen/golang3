package site

import (
	"time"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"serverApi/pkg/constant"
)

const TableNameRoomSceneHistory = "live_room_scene_history"

var (
	FieldRoomSceneHistoryId             = field.NewInt(TableNameRoomSceneHistory, "id")
	FieldRoomSceneHistoryRoomId         = field.NewInt(TableNameRoomSceneHistory, "room_id")
	FieldRoomSceneHistoryUserId         = field.NewInt(TableNameRoomSceneHistory, "user_id")
	FieldRoomSceneHistoryGameId         = field.NewInt(TableNameRoomSceneHistory, "game_id")
	FieldRoomSceneHistoryPayRules       = field.NewInt(TableNameRoomSceneHistory, "pay_rules")
	FieldRoomSceneHistoryTrialDuration  = field.NewInt(TableNameRoomSceneHistory, "trial_duration")
	FieldRoomSceneHistoryUnitPrice      = field.NewInt(TableNameRoomSceneHistory, "unit_price")
	FieldRoomSceneHistoryIsOpenVibrator = field.NewInt(TableNameRoomSceneHistory, "is_open_vibrator")
	FieldRoomSceneHistoryViewsNum       = field.NewInt(TableNameRoomSceneHistory, "views_num")
	FieldRoomSceneHistoryGiftAmount     = field.NewInt(TableNameRoomSceneHistory, "gift_amount")
	FieldRoomSceneHistoryGiftNum        = field.NewInt(TableNameRoomSceneHistory, "gift_num")
	FieldRoomSceneHistoryFollowersNum   = field.NewInt(TableNameRoomSceneHistory, "followers_num")
	FieldRoomSceneHistoryLikeNum        = field.NewInt(TableNameRoomSceneHistory, "like_num")
	FieldRoomSceneHistoryOpenAt         = field.NewTime(TableNameRoomSceneHistory, "open_at")
	FieldRoomSceneHistoryEndAt          = field.NewTime(TableNameRoomSceneHistory, "end_at")
	FieldRoomSceneHistoryStatus         = field.NewInt(TableNameRoomSceneHistory, "status")
	FieldRoomSceneHistoryCreatedAt      = field.NewTime(TableNameRoomSceneHistory, "created_at")
	FieldRoomSceneHistoryUpdatedAt      = field.NewTime(TableNameRoomSceneHistory, "updated_at")
)

type RoomSceneHistory struct {
	Id             int       `gorm:"column:id" json:"id"`                             // 直播场次id
	RoomId         int       `gorm:"column:room_id" json:"room_id"`                   // 直播间id
	UserId         int       `gorm:"column:user_id" json:"user_id"`                   // 主播id
	GameId         int       `gorm:"column:game_id" json:"game_id"`                   // 直播游戏id
	PayRules       int       `gorm:"column:pay_rules" json:"pay_rules"`               // 付费规则 1- 免费 2-分钟付费 3-场次付费
	TrialDuration  int       `gorm:"column:trial_duration" json:"trial_duration"`     // 试看时长/秒
	UnitPrice      int       `gorm:"column:unit_price" json:"unit_price"`             // 单价 每分钟单价/每场单价，根据付费规则
	IsOpenVibrator int       `gorm:"column:is_open_vibrator" json:"is_open_vibrator"` // 是否开启跳蛋
	ViewsNum       int       `gorm:"column:views_num" json:"views_num"`               // 本次观看人数
	GiftAmount     int       `gorm:"column:gift_amount" json:"gift_amount"`           // 本次收到礼物金额/钻石
	GiftNum        int       `gorm:"column:gift_num" json:"gift_num"`                 // 本场送礼物人数
	FollowersNum   int       `gorm:"column:followers_num" json:"followers_num"`       // 本次关注数量
	LikeNum        int       `gorm:"column:like_num" json:"like_num"`                 // 本次点赞数量
	OpenAt         time.Time `gorm:"column:open_at" json:"open_at"`                   // 开播时间
	EndAt          time.Time `gorm:"column:end_at" json:"end_at"`                     // 结束时间
	Status         int       `gorm:"column:status" json:"status"`                     // 状态 1- 直播中 2-结束
	CreatedAt      time.Time `gorm:"column:created_at" json:"created_at"`             // 创建时间
	UpdatedAt      time.Time `gorm:"column:updated_at" json:"updated_at"`             // 更新时间
}

func (*RoomSceneHistory) TableName() string {
	return TableNameRoomSceneHistory
}

func (a *RoomSceneHistory) IsEmpty() bool {
	if a == nil {
		return true
	}
	return a.Id == 0
}

// DefaultScope 过滤已删除 status=3
func (a *RoomSceneHistory) DefaultScope(db *gorm.DB) *gorm.DB {
	return db.Where(FieldRoomSceneHistoryStatus.Neq(constant.StatusDel))
}

// RoomSceneHistoryWhereOption 条件项
type RoomSceneHistoryWhereOption func(*gorm.DB) *gorm.DB

// WhereRoomSceneHistoryId 查询直播场次id
func WhereRoomSceneHistoryId(value int) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryId.Eq(value))
	}
}

// WhereRoomSceneHistoryRoomId 查询直播间id
func WhereRoomSceneHistoryRoomId(value int) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryRoomId.Eq(value))
	}
}

// WhereRoomSceneHistoryUserId 查询主播id
func WhereRoomSceneHistoryUserId(value int) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryUserId.Eq(value))
	}
}

// WhereRoomSceneHistoryGameId 查询直播游戏id
func WhereRoomSceneHistoryGameId(value int) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryGameId.Eq(value))
	}
}

// WhereRoomSceneHistoryPayRules 查询付费规则 1- 免费 2-分钟付费 3-场次付费
func WhereRoomSceneHistoryPayRules(value int) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryPayRules.Eq(value))
	}
}

// WhereRoomSceneHistoryTrialDuration 查询试看时长/秒
func WhereRoomSceneHistoryTrialDuration(value int) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryTrialDuration.Eq(value))
	}
}

// WhereRoomSceneHistoryUnitPrice 查询单价 每分钟单价/每场单价，根据付费规则
func WhereRoomSceneHistoryUnitPrice(value int) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryUnitPrice.Eq(value))
	}
}

// WhereRoomSceneHistoryViewsNum 查询本次观看人数
func WhereRoomSceneHistoryViewsNum(value int) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryViewsNum.Eq(value))
	}
}

// WhereRoomSceneHistoryGiftAmount 查询本次收到礼物金额/钻石
func WhereRoomSceneHistoryGiftAmount(value int) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryGiftAmount.Eq(value))
	}
}

// WhereRoomSceneHistoryGiftNum 查询本场送礼物人数
func WhereRoomSceneHistoryGiftNum(value int) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryGiftNum.Eq(value))
	}
}

// WhereRoomSceneHistoryFollowersNum 查询本次关注数量
func WhereRoomSceneHistoryFollowersNum(value int) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryFollowersNum.Eq(value))
	}
}

// WhereRoomSceneHistoryLikeNum 查询本次点赞数量
func WhereRoomSceneHistoryLikeNum(value int) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryLikeNum.Eq(value))
	}
}

// WhereRoomSceneHistoryOpenAt 查询开播时间
func WhereRoomSceneHistoryOpenAt(value time.Time) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryOpenAt.Eq(value))
	}
}

// WhereRoomSceneHistoryEndAt 查询结束时间
func WhereRoomSceneHistoryEndAt(value time.Time) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryEndAt.Eq(value))
	}
}

// WhereRoomSceneHistoryStatus 查询状态 1- 直播中 2-结束
func WhereRoomSceneHistoryStatus(value int) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryStatus.Eq(value))
	}
}

// WhereRoomSceneHistoryCreatedAt 查询创建时间
func WhereRoomSceneHistoryCreatedAt(value time.Time) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryCreatedAt.Eq(value))
	}
}

// WhereRoomSceneHistoryUpdatedAt 查询更新时间
func WhereRoomSceneHistoryUpdatedAt(value time.Time) RoomSceneHistoryWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSceneHistoryUpdatedAt.Eq(value))
	}
}

// RoomSceneHistoryUpdateOption 修改项
type RoomSceneHistoryUpdateOption func(map[string]interface{})

// SetRoomSceneHistoryId 设置直播场次id
func SetRoomSceneHistoryId(value int) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryId.ColumnName().String()] = value
	}
}

// SetRoomSceneHistoryRoomId 设置直播间id
func SetRoomSceneHistoryRoomId(value int) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryRoomId.ColumnName().String()] = value
	}
}

// SetRoomSceneHistoryUserId 设置主播id
func SetRoomSceneHistoryUserId(value int) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryUserId.ColumnName().String()] = value
	}
}

// SetRoomSceneHistoryGameId 设置直播游戏id
func SetRoomSceneHistoryGameId(value int) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryGameId.ColumnName().String()] = value
	}
}

// SetRoomSceneHistoryPayRules 设置付费规则 1- 免费 2-分钟付费 3-场次付费
func SetRoomSceneHistoryPayRules(value int) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryPayRules.ColumnName().String()] = value
	}
}

// SetRoomSceneHistoryTrialDuration 设置试看时长/秒
func SetRoomSceneHistoryTrialDuration(value int) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryTrialDuration.ColumnName().String()] = value
	}
}

// SetRoomSceneHistoryUnitPrice 设置单价 每分钟单价/每场单价，根据付费规则
func SetRoomSceneHistoryUnitPrice(value float64) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryUnitPrice.ColumnName().String()] = value
	}
}

// SetRoomSceneHistoryViewsNum 设置本次观看人数
func SetRoomSceneHistoryViewsNum(value int) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryViewsNum.ColumnName().String()] = value
	}
}

// SetRoomSceneHistoryGiftAmount 设置本次收到礼物金额/钻石
func SetRoomSceneHistoryGiftAmount(value int) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryGiftAmount.ColumnName().String()] = value
	}
}

// SetRoomSceneHistoryGiftNum 设置本场送礼物人数
func SetRoomSceneHistoryGiftNum(value int) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryGiftNum.ColumnName().String()] = value
	}
}

// SetRoomSceneHistoryFollowersNum 设置本次关注数量
func SetRoomSceneHistoryFollowersNum(value int) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryFollowersNum.ColumnName().String()] = value
	}
}

// SetRoomSceneHistoryLikeNum 设置本次点赞数量
func SetRoomSceneHistoryLikeNum(value int) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryLikeNum.ColumnName().String()] = value
	}
}

// SetRoomSceneHistoryOpenAt 设置开播时间
func SetRoomSceneHistoryOpenAt(value time.Time) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryOpenAt.ColumnName().String()] = value
	}
}

// SetRoomSceneHistoryEndAt 设置结束时间
func SetRoomSceneHistoryEndAt(value time.Time) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryEndAt.ColumnName().String()] = value
	}
}

// SetRoomSceneHistoryStatus 设置状态 1- 直播中 2-结束
func SetRoomSceneHistoryStatus(value int) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryStatus.ColumnName().String()] = value
	}
}

// SetRoomSceneHistoryCreatedAt 设置创建时间
func SetRoomSceneHistoryCreatedAt(value time.Time) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryCreatedAt.ColumnName().String()] = value
	}
}

// SetRoomSceneHistoryUpdatedAt 设置更新时间
func SetRoomSceneHistoryUpdatedAt(value time.Time) RoomSceneHistoryUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSceneHistoryUpdatedAt.ColumnName().String()] = value
	}
}
