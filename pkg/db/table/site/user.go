package site

import (
	"time"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"serverApi/pkg/constant"
)

const TableNameUser = "site_user"

var (
	FieldUserId                = field.NewInt(TableNameUser, "id")
	FieldUserChatUuid          = field.NewString(TableNameUser, "chat_uuid")
	FieldUserCountryCode       = field.NewString(TableNameUser, "country_code")
	FieldUserAvatar            = field.NewString(TableNameUser, "avatar")
	FieldUserNickname          = field.NewString(TableNameUser, "nickname")
	FieldUserSex               = field.NewInt(TableNameUser, "sex")
	FieldUserSign              = field.NewString(TableNameUser, "sign")
	FieldUserBirthday          = field.NewString(TableNameUser, "birthday")
	FieldUserFeeling           = field.NewInt(TableNameUser, "feeling")
	FieldUserCountry           = field.NewString(TableNameUser, "country")
	FieldUserArea              = field.NewString(TableNameUser, "area")
	FieldUserProfession        = field.NewInt(TableNameUser, "profession")
	FieldUserPayPassword       = field.NewString(TableNameUser, "pay_password")
	FieldUserCategory          = field.NewInt(TableNameUser, "category")
	FieldUserSiteId            = field.NewInt(TableNameUser, "site_id")
	FieldUserInviteCode        = field.NewString(TableNameUser, "invite_code")
	FieldUserParentId          = field.NewInt(TableNameUser, "parent_id")
	FieldUserLevelId           = field.NewInt(TableNameUser, "level_id")
	FieldUserSetLevelId        = field.NewInt(TableNameUser, "set_level_id")
	FieldUserLoginErrorTimes   = field.NewInt(TableNameUser, "login_error_times")
	FieldUserRemark            = field.NewString(TableNameUser, "remark")
	FieldUserGmStatus          = field.NewInt(TableNameUser, "gm_status")
	FieldUserStatus            = field.NewInt(TableNameUser, "status")
	FieldUserAgentRebateStatus = field.NewInt(TableNameUser, "agent_rebate_status")
	FieldUserBetStatus         = field.NewInt(TableNameUser, "bet_status")
	FieldUserDrawStatus        = field.NewInt(TableNameUser, "draw_status")
	FieldUserCreatedAt         = field.NewTime(TableNameUser, "created_at")
	FieldUserUpdatedAt         = field.NewTime(TableNameUser, "updated_at")
)

type User struct {
	Id                int       `gorm:"column:id" json:"id"`                                   // ID
	ChatUuid          string    `gorm:"column:chat_uuid" json:"chat_uuid"`                     // 聊天室用户id
	CountryCode       string    `gorm:"column:country_code" json:"country_code"`               // 注册国家code
	Avatar            string    `gorm:"column:avatar" json:"avatar"`                           // 用户头像
	Nickname          string    `gorm:"column:nickname" json:"nickname"`                       // 昵称
	Sex               int       `gorm:"column:sex" json:"sex"`                                 // 性别 1-男 2-女
	Sign              string    `gorm:"column:sign" json:"sign"`                               // 签名
	Birthday          string    `gorm:"column:birthday" json:"birthday"`                       // 生日
	Feeling           int       `gorm:"column:feeling" json:"feeling"`                         // 感情
	Country           string    `gorm:"column:country" json:"country"`                         // 国家
	Area              string    `gorm:"column:area" json:"area"`                               // 地区
	Profession        int       `gorm:"column:profession" json:"profession"`                   // 职业
	PayPassword       string    `gorm:"column:pay_password" json:"pay_password"`               // 支付密码
	Category          int       `gorm:"column:category" json:"category"`                       // 类型 1-用户 2-主播 3-机器人
	SiteId            int       `gorm:"column:site_id" json:"site_id"`                         // 站点ID
	InviteCode        string    `gorm:"column:invite_code" json:"invite_code"`                 // 邀请码
	ParentId          int       `gorm:"column:parent_id" json:"parent_id"`                     // 上级ID
	LevelId           int       `gorm:"column:level_id" json:"level_id"`                       // 等级ID
	SetLevelId        int       `gorm:"column:set_level_id" json:"set_level_id"`               // 后台设置等级id 优先展示
	LoginErrorTimes   int       `gorm:"column:login_error_times" json:"login_error_times"`     // 登陆错误次数
	Remark            string    `gorm:"column:remark" json:"remark"`                           // 备注
	GmStatus          int       `gorm:"column:gm_status" json:"gm_status"`                     // 超级管理员状态 1- 开启 2-关闭
	Status            int       `gorm:"column:status" json:"status"`                           // 账号状态 1-正常 2-禁用 3-删除
	AgentRebateStatus int       `gorm:"column:agent_rebate_status" json:"agent_rebate_status"` // 代理返点状态 1-正常 2-禁用
	BetStatus         int       `gorm:"column:bet_status" json:"bet_status"`                   // 投注状态 1-正常 2-禁用
	DrawStatus        int       `gorm:"column:draw_status" json:"draw_status"`                 // 出款状态 1-正常 2-禁用
	CreatedAt         time.Time `gorm:"column:created_at" json:"created_at"`                   // 创建时间
	UpdatedAt         time.Time `gorm:"column:updated_at" json:"updated_at"`                   // 更新时间
}

func (*User) TableName() string {
	return TableNameUser
}

func (a *User) IsEmpty() bool {
	if a == nil {
		return true
	}
	return a.Id == 0
}

// DefaultScope 过滤已删除 status=3
func (a *User) DefaultScope(db *gorm.DB) *gorm.DB {
	return db.Where(FieldUserStatus.Eq(constant.UserStatusNormal))
}

// UserWhereOption 条件项
type UserWhereOption func(*gorm.DB) *gorm.DB

// WhereUserId 查询ID
func WhereUserId(value int) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserId.Eq(value))
	}
}

// WhereUserChatUuid 查询聊天室用户id
func WhereUserChatUuid(value string) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserChatUuid.Eq(value))
	}
}

// WhereUserCountryCode 查询注册国家code
func WhereUserCountryCode(value string) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserCountryCode.Eq(value))
	}
}

// WhereUserAvatar 查询用户头像
func WhereUserAvatar(value string) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserAvatar.Eq(value))
	}
}

// WhereUserNickname 查询昵称
func WhereUserNickname(value string) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserNickname.Eq(value))
	}
}

// WhereUserSex 查询性别 1-男 2-女
func WhereUserSex(value int) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserSex.Eq(value))
	}
}

// WhereUserSign 查询签名
func WhereUserSign(value string) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserSign.Eq(value))
	}
}

// WhereUserBirthday 查询生日
func WhereUserBirthday(value string) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserBirthday.Eq(value))
	}
}

// WhereUserFeeling 查询感情
func WhereUserFeeling(value int) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserFeeling.Eq(value))
	}
}

// WhereUserCountry 查询国家
func WhereUserCountry(value string) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserCountry.Eq(value))
	}
}

// WhereUserArea 查询地区
func WhereUserArea(value string) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserArea.Eq(value))
	}
}

// WhereUserProfession 查询职业
func WhereUserProfession(value int) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserProfession.Eq(value))
	}
}

// WhereUserPayPassword 查询支付密码
func WhereUserPayPassword(value string) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserPayPassword.Eq(value))
	}
}

// WhereUserCategory 查询类型 1-用户 2-主播 3-机器人
func WhereUserCategory(value int) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserCategory.Eq(value))
	}
}

// WhereUserSiteId 查询站点ID
func WhereUserSiteId(value int) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserSiteId.Eq(value))
	}
}

// WhereUserInviteCode 查询邀请码
func WhereUserInviteCode(value string) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserInviteCode.Eq(value))
	}
}

// WhereUserParentId 查询上级ID
func WhereUserParentId(value int) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserParentId.Eq(value))
	}
}

// WhereUserLevelId 查询等级ID
func WhereUserLevelId(value int) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserLevelId.Eq(value))
	}
}

// WhereUserSetLevelId 查询后台设置等级id 优先展示
func WhereUserSetLevelId(value int) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserSetLevelId.Eq(value))
	}
}

// WhereUserLoginErrorTimes 查询登陆错误次数
func WhereUserLoginErrorTimes(value int) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserLoginErrorTimes.Eq(value))
	}
}

// WhereUserRemark 查询备注
func WhereUserRemark(value string) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserRemark.Eq(value))
	}
}

// WhereUserGmStatus 查询超级管理员状态 1- 开启 2-关闭
func WhereUserGmStatus(value int) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserGmStatus.Eq(value))
	}
}

// WhereUserStatus 查询账号状态 1-正常 2-禁用 3-删除
func WhereUserStatus(value int) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserStatus.Eq(value))
	}
}

// WhereUserAgentRebateStatus 查询代理返点状态 1-正常 2-禁用
func WhereUserAgentRebateStatus(value int) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserAgentRebateStatus.Eq(value))
	}
}

// WhereUserBetStatus 查询投注状态 1-正常 2-禁用
func WhereUserBetStatus(value int) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserBetStatus.Eq(value))
	}
}

// WhereUserDrawStatus 查询出款状态 1-正常 2-禁用
func WhereUserDrawStatus(value int) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserDrawStatus.Eq(value))
	}
}

// WhereUserCreatedAt 查询创建时间
func WhereUserCreatedAt(value time.Time) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserCreatedAt.Eq(value))
	}
}

// WhereUserUpdatedAt 查询更新时间
func WhereUserUpdatedAt(value time.Time) UserWhereOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(FieldUserUpdatedAt.Eq(value))
	}
}

// UserUpdateOption 修改项
type UserUpdateOption func(map[string]interface{})

// SetUserId 设置ID
func SetUserId(value int) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserId.ColumnName().String()] = value
	}
}

// SetUserChatUuid 设置聊天室用户id
func SetUserChatUuid(value string) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserChatUuid.ColumnName().String()] = value
	}
}

// SetUserCountryCode 设置注册国家code
func SetUserCountryCode(value string) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserCountryCode.ColumnName().String()] = value
	}
}

// SetUserAvatar 设置用户头像
func SetUserAvatar(value string) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserAvatar.ColumnName().String()] = value
	}
}

// SetUserNickname 设置昵称
func SetUserNickname(value string) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserNickname.ColumnName().String()] = value
	}
}

// SetUserSex 设置性别 1-男 2-女
func SetUserSex(value int) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserSex.ColumnName().String()] = value
	}
}

// SetUserSign 设置签名
func SetUserSign(value string) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserSign.ColumnName().String()] = value
	}
}

// SetUserBirthday 设置生日
func SetUserBirthday(value string) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserBirthday.ColumnName().String()] = value
	}
}

// SetUserFeeling 设置感情
func SetUserFeeling(value int) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserFeeling.ColumnName().String()] = value
	}
}

// SetUserCountry 设置国家
func SetUserCountry(value string) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserCountry.ColumnName().String()] = value
	}
}

// SetUserArea 设置地区
func SetUserArea(value string) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserArea.ColumnName().String()] = value
	}
}

// SetUserProfession 设置职业
func SetUserProfession(value int) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserProfession.ColumnName().String()] = value
	}
}

// SetUserPayPassword 设置支付密码
func SetUserPayPassword(value string) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserPayPassword.ColumnName().String()] = value
	}
}

// SetUserCategory 设置类型 1-用户 2-主播 3-机器人
func SetUserCategory(value int) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserCategory.ColumnName().String()] = value
	}
}

// SetUserSiteId 设置站点ID
func SetUserSiteId(value int) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserSiteId.ColumnName().String()] = value
	}
}

// SetUserInviteCode 设置邀请码
func SetUserInviteCode(value string) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserInviteCode.ColumnName().String()] = value
	}
}

// SetUserParentId 设置上级ID
func SetUserParentId(value int) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserParentId.ColumnName().String()] = value
	}
}

// SetUserLevelId 设置等级ID
func SetUserLevelId(value int) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserLevelId.ColumnName().String()] = value
	}
}

// SetUserSetLevelId 设置后台设置等级id 优先展示
func SetUserSetLevelId(value int) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserSetLevelId.ColumnName().String()] = value
	}
}

// SetUserLoginErrorTimes 设置登陆错误次数
func SetUserLoginErrorTimes(value int) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserLoginErrorTimes.ColumnName().String()] = value
	}
}

// SetUserRemark 设置备注
func SetUserRemark(value string) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserRemark.ColumnName().String()] = value
	}
}

// SetUserGmStatus 设置超级管理员状态 1- 开启 2-关闭
func SetUserGmStatus(value int) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserGmStatus.ColumnName().String()] = value
	}
}

// SetUserStatus 设置账号状态 1-正常 2-禁用 3-删除
func SetUserStatus(value int) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserStatus.ColumnName().String()] = value
	}
}

// SetUserAgentRebateStatus 设置代理返点状态 1-正常 2-禁用
func SetUserAgentRebateStatus(value int) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserAgentRebateStatus.ColumnName().String()] = value
	}
}

// SetUserBetStatus 设置投注状态 1-正常 2-禁用
func SetUserBetStatus(value int) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserBetStatus.ColumnName().String()] = value
	}
}

// SetUserDrawStatus 设置出款状态 1-正常 2-禁用
func SetUserDrawStatus(value int) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserDrawStatus.ColumnName().String()] = value
	}
}

// SetUserCreatedAt 设置创建时间
func SetUserCreatedAt(value time.Time) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserCreatedAt.ColumnName().String()] = value
	}
}

// SetUserUpdatedAt 设置更新时间
func SetUserUpdatedAt(value time.Time) UserUpdateOption {
	return func(updates map[string]interface{}) {
		updates[FieldUserUpdatedAt.ColumnName().String()] = value
	}
}
