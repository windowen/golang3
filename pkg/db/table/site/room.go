package site

import (
	"time"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"serverApi/pkg/constant"
)

const TableNameRoom = "live_room"

var (
	FieldRoomId                = field.NewInt(TableNameRoom, "id")
	FieldRoomSort              = field.NewInt(TableNameRoom, "sort")
	FieldRoomUserId            = field.NewInt(TableNameRoom, "user_id")
	FieldRoomChatRoomId        = field.NewString(TableNameRoom, "chat_room_id")
	FieldRoomCountryCode       = field.NewString(TableNameRoom, "country_code")
	FieldRoomTitle             = field.NewString(TableNameRoom, "title")
	FieldRoomCountryTagsJson   = field.NewString(TableNameRoom, "tags_json")
	FieldRoomCover             = field.NewString(TableNameRoom, "cover")
	FieldRoomSummary           = field.NewString(TableNameRoom, "summary")
	FieldRoomVideoClarity      = field.NewInt(TableNameRoom, "video_clarity")
	FieldRoomGiftRatio         = field.NewInt(TableNameRoom, "gift_ratio")
	FieldRoomPlatformRatio     = field.NewInt(TableNameRoom, "platform_ratio")
	FieldRoomFamilyRatio       = field.NewInt(TableNameRoom, "family_ratio")
	FieldRoomRemark            = field.NewString(TableNameRoom, "remark")
	FieldRoomStatus            = field.NewInt(TableNameRoom, "status")
	FieldRoomLiveStatus        = field.NewInt(TableNameRoom, "live_status")
	FieldRoomPaidPurviewStatus = field.NewInt(TableNameRoom, "paid_purview_status")
	FieldRoomBottomSort        = field.NewInt(TableNameRoom, "bottom_sort")
	FieldRoomLevelId           = field.NewInt(TableNameRoom, "level_id")
	FieldRoomSetLevelId        = field.NewInt(TableNameRoom, "set_level_id")
	FieldRoomCreatedAt         = field.NewTime(TableNameRoom, "created_at")
	FieldRoomUpdatedAt         = field.NewTime(TableNameRoom, "updated_at")
)

type Room struct {
	Id                int       `gorm:"column:id" json:"id"`
	Sort              int       `gorm:"column:sort" json:"sort"`                                                               // 排序 不能存在相同的排序，如果已存在先设置为0
	UserId            int       `gorm:"column:user_id" json:"user_id"`                                                         // 主播id
	ChatRoomId        string    `gorm:"column:chat_room_id" json:"chat_room_id"`                                               // chat房间id(创建主播时调用rpc返回)
	CountryCode       string    `gorm:"column:country_code" json:"country_code"`                                               // 国家code
	Title             string    `gorm:"column:title" json:"title"`                                                             // 直播间标题
	TagsJson          string    `gorm:"column:tags_json" json:"tags_json"`                                                     // 标签 ["1", "2", "3"] 二级标签id
	Cover             string    `gorm:"column:cover" json:"cover"`                                                             // 封面图片
	Summary           string    `gorm:"column:summary" json:"summary"`                                                         // 简介
	VideoClarity      int       `gorm:"column:video_clarity" json:"video_clarity"`                                             // 视屏清晰度 1- 480p 2-540p 3-720p
	GiftRatio         int       `gorm:"column:gift_ratio" json:"gift_ratio"`                                                   // 主播抽成比例  实际*100
	PlatformRatio     int       `gorm:"column:platform_ratio" json:"platform_ratio"`                                           // 平台抽成比例  实际*100
	FamilyRatio       int       `gorm:"column:family_ratio" json:"family_ratio"`                                               // 家族抽成比例  实际*100
	Remark            string    `gorm:"column:remark" json:"remark"`                                                           // 备注
	Status            int       `gorm:"column:status" json:"status"`                                                           // 房间状态 1-正常  2-禁用 3-删除
	LiveStatus        int       `gorm:"column:live_status" json:"live_status"`                                                 // 直播状态 1-直播 2-下播
	PaidPurviewStatus int       `gorm:"column:paid_purview_status" json:"paid_purview_status"`                                 // 是否有可以开启付费直播
	BottomSort        int       `gorm:"column:bottomSort" json:"bottomSort"`                                                   // 置底
	LevelId           int       `gorm:"column:level_id" json:"level_id"`                                                       // 等级ID
	SetLevelId        int       `gorm:"column:set_level_id" json:"set_level_id"`                                               // 后台设置等级id 优先展示
	LastStartLiveTime time.Time `gorm:"column:last_start_live_time;default:'1971-01-01 00:00:00'" json:"last_start_live_time"` // 最后开播时间
	LastEndLiveTime   time.Time `gorm:"column:last_end_live_time;default:'1971-01-01 00:00:00'" json:"last_end_live_time"`     // 最后下播时间
	CreatedAt         time.Time `gorm:"column:created_at" json:"created_at"`                                                   // 创建时间
	UpdatedAt         time.Time `gorm:"column:updated_at" json:"updated_at"`                                                   // 更新时间
}

func (*Room) TableName() string {
	return TableNameRoom
}

func (a *Room) IsEmpty() bool {
	if a == nil {
		return true
	}
	return a.Id == 0
}

// DefaultScope 过滤已删除 status=3
func (a *Room) DefaultScope(db *gorm.DB) *gorm.DB {
	return db.Where(FieldRoomStatus.Neq(constant.StatusDel))
}

// RoomWhereOption 条件项
type RoomWhereOption func(*gorm.DB) *gorm.DB

// WhereRoomId 查询
func WhereRoomId(value int) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomId.Eq(value))
	}
}

// WhereRoomSort 查询排序 不能存在相同的排序，如果已存在先设置为0
func WhereRoomSort(value int) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSort.Eq(value))
	}
}

// WhereRoomUserId 查询主播id
func WhereRoomUserId(value int) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomUserId.Eq(value))
	}
}

// WhereRoomChatRoomId 查询chat房间id(创建主播时调用rpc返回)
func WhereRoomChatRoomId(value string) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomChatRoomId.Eq(value))
	}
}

// WhereRoomCountryCode 查询国家code
func WhereRoomCountryCode(value string) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomCountryCode.Eq(value))
	}
}

// WhereRoomTitle 查询直播间标题
func WhereRoomTitle(value string) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomTitle.Eq(value))
	}
}

// WhereRoomCover 查询封面图片
func WhereRoomCover(value string) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomCover.Eq(value))
	}
}

// WhereRoomSummary 查询简介
func WhereRoomSummary(value string) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomSummary.Eq(value))
	}
}

// WhereRoomVideoClarity 查询视屏清晰度 1- 480p 2-540p 3-720p
func WhereRoomVideoClarity(value int) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomVideoClarity.Eq(value))
	}
}

// WhereRoomGiftRatio 查询主播抽成比例  实际*100
func WhereRoomGiftRatio(value int) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomGiftRatio.Eq(value))
	}
}

// WhereRoomPlatformRatio 查询平台抽成比例  实际*100
func WhereRoomPlatformRatio(value int) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomPlatformRatio.Eq(value))
	}
}

// WhereRoomFamilyRatio 查询家族抽成比例  实际*100
func WhereRoomFamilyRatio(value int) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomFamilyRatio.Eq(value))
	}
}

// WhereRoomRemark 查询备注
func WhereRoomRemark(value string) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomRemark.Eq(value))
	}
}

// WhereRoomStatus 查询房间状态 1-正常  2-禁用 3-删除
func WhereRoomStatus(value int) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomStatus.Eq(value))
	}
}

// WhereRoomLiveStatus 查询直播状态 1-直播 2-下播
func WhereRoomLiveStatus(value int) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomLiveStatus.Eq(value))
	}
}

// WhereRoomIsBottom 查询是否设置底部
func WhereRoomIsBottom(value int) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomBottomSort.Eq(value))
	}
}

// WhereRoomCreatedAt 查询创建时间
func WhereRoomCreatedAt(value time.Time) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomCreatedAt.Eq(value))
	}
}

// WhereRoomUpdatedAt 查询更新时间
func WhereRoomUpdatedAt(value time.Time) RoomWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldRoomUpdatedAt.Eq(value))
	}
}

// RoomUpdateOption 修改项
type RoomUpdateOption func(map[string]interface{})

// SetRoomId 设置
func SetRoomId(value int) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomId.ColumnName().String()] = value
	}
}

// SetRoomSort 设置排序 不能存在相同的排序，如果已存在先设置为0
func SetRoomSort(value int) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSort.ColumnName().String()] = value
	}
}

// SetRoomUserId 设置主播id
func SetRoomUserId(value int) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomUserId.ColumnName().String()] = value
	}
}

// SetRoomChatRoomId 设置chat房间id(创建主播时调用rpc返回)
func SetRoomChatRoomId(value string) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomChatRoomId.ColumnName().String()] = value
	}
}

// SetRoomCountryCode 设置国家code
func SetRoomCountryCode(value string) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomCountryCode.ColumnName().String()] = value
	}
}

// SetRoomTitle 设置直播间标题
func SetRoomTitle(value string) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomTitle.ColumnName().String()] = value
	}
}

// SetRoomCover 设置封面图片
func SetRoomCover(value string) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomCover.ColumnName().String()] = value
	}
}

// SetRoomSummary 设置简介
func SetRoomSummary(value string) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomSummary.ColumnName().String()] = value
	}
}

// SetRoomVideoClarity 设置视屏清晰度 1- 480p 2-540p 3-720p
func SetRoomVideoClarity(value int) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomVideoClarity.ColumnName().String()] = value
	}
}

// SetRoomGiftRatio 设置主播抽成比例  实际*100
func SetRoomGiftRatio(value int) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomGiftRatio.ColumnName().String()] = value
	}
}

// SetRoomPlatformRatio 设置平台抽成比例  实际*100
func SetRoomPlatformRatio(value int) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomPlatformRatio.ColumnName().String()] = value
	}
}

// SetRoomFamilyRatio 设置家族抽成比例  实际*100
func SetRoomFamilyRatio(value int) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomFamilyRatio.ColumnName().String()] = value
	}
}

// SetRoomRemark 设置备注
func SetRoomRemark(value string) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomRemark.ColumnName().String()] = value
	}
}

// SetRoomStatus 设置房间状态 1-正常  2-禁用 3-删除
func SetRoomStatus(value int) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomStatus.ColumnName().String()] = value
	}
}

// SetRoomLiveStatus 设置直播状态 1-直播 2-下播
func SetRoomLiveStatus(value int) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomLiveStatus.ColumnName().String()] = value
	}
}

// SetRoomIsBottom 设置是否设置底部
func SetRoomIsBottom(value int) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomBottomSort.ColumnName().String()] = value
	}
}

// SetRoomCreatedAt 设置创建时间
func SetRoomCreatedAt(value time.Time) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomCreatedAt.ColumnName().String()] = value
	}
}

// SetRoomUpdatedAt 设置更新时间
func SetRoomUpdatedAt(value time.Time) RoomUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldRoomUpdatedAt.ColumnName().String()] = value
	}
}
