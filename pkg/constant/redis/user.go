package redis

const (
	UserCacheInfoKey         = "user_cache_info_%v"          // 用户信息缓存key v1-用户id
	UserLoginToken           = "user_login_token_%v"         // 用户登陆令牌Key
	UserLoginTokenVersion    = "user_login_token_version_%v" // 用户登陆令牌版本Key
	UserFollowSortedCacheKey = "user_follow_sorted_%v"       // 用户关注排序集合
	UserFollowCacheKey       = "user_follow_%v"              // 用户关注列表
	UserFansSortedCacheKey   = "user_fans_sorted_%v"         // 用户粉丝排序集合
	UserFansCacheKey         = "user_fans_%v"                // 用户粉丝列表
)
