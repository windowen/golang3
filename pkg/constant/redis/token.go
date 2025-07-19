package redis

import "time"

const (
	TokenTTL        = 24 * time.Hour     // token过期时间
	RefreshTokenTTL = 7 * 24 * time.Hour // 刷新token过期时间
)
