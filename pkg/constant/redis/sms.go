package redis

const (
	VerifyCodeCacheBasisKey    = "verify_code_%v_%v"        // 验证码发送场景key v1: 验证场景 v2:手机号/邮箱
	VerifyCodeCache24HLimitKey = "verify_code_24h_limit_%v" // 24H发送限制key v1:手机号/邮箱
	VerifyCodeCoolDownCacheKey = "verify_code_cool_down_%v" // 验证码发送60秒冷却key v1:手机号/邮箱
)
