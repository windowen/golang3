package constant

var (
	UserAccountTypes   = []int{UserAccountTypeMobile, UserAccountTypeEmail}
	UserLoginAuthTypes = []int{UserAuthTypeCode, UserAuthTypePassword}

	UserPasswordRegexp = `^[a-zA-Z0-9]{6,15}$`
	UserMobileRegexp   = `^0?\d{5,17}$`
	UserEmailRegexp    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

const (
	UserAccountTypeMobile = 1 // 手机
	UserAccountTypeEmail  = 2 // 邮箱
	UserAuthTypeCode      = 1 // 验证码
	UserAuthTypePassword  = 2 // 密码

	UserInviteCodeLen   = 6 // 邀请码长度
	UserCategoryGeneral = 1 // 普通用户
	UserCategoryAnchor  = 2 // 主播
	UserCategoryRobot   = 3 // 机器人

	UserStatusNormal  = 1 // 正常
	UserStatusDisable = 2 // 禁用
	UserStatusDel     = 3 // 删除

	UserFollowFocusOn = 1 // 关注
	UserFollowUnlock  = 2 // 取关
)
