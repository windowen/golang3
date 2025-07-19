package mw

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"serverApi/pkg/constant"
	"serverApi/pkg/db/cache"
	"serverApi/pkg/jwt"
	"serverApi/pkg/tools/apiresp"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/zlogger"
)

func UserAuthCheck(rdb redis.UniversalClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			authorization = strings.TrimSpace(c.GetHeader(constant.Authorization))
			userId        = constant.Zero
		)

		if authorization == "" {
			apiresp.GinError(c, errs.ErrNoPermission.WithDetail("token_not_permissions"))
			c.Abort()
			return
		}

		jwtTools := jwt.NewJWT(rdb)

		claims, err := jwtTools.ParseToken(c, authorization)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				// 生成新的Token
				result, err := jwtTools.RefreshToken(c, strings.TrimSpace(c.GetHeader(constant.RefreshToken)))
				if err != nil {
					zlogger.Errorf("UserAuthCheck jwtTools.GenToken |token:%v| err: %v", authorization, err)
					apiresp.GinError(c, err)
					c.Abort()
					return
				}

				userId = result.UserId
				// 返回新的令牌
				c.Header(constant.RefreshAuthorization, result.Token)
			} else {
				zlogger.Errorf("UserAuthCheck UserCache.Get |token:%v| ,err: %v", authorization, err)
				apiresp.GinError(c, err)
				c.Abort()
				return
			}
		}

		if claims != nil {
			userId = claims.UserId
		}

		userCache := cache.NewUserCache(rdb)
		// 检查用户状态是否正常
		userCacheInfo, err := userCache.Get(c, userId)
		if err != nil {
			zlogger.Errorf("UserAuthCheck UserCache.Get |userId:%v| ,err: %v", userId, err)
			apiresp.GinError(c, errs.ErrNoPermission.WithDetail("token_expired"))
			c.Abort()
			return
		}

		if userCacheInfo.Status != constant.UserStatusNormal {
			apiresp.GinError(c, errs.ErrNoPermission.WithDetail("user_disabled"))
			c.Abort()
			return
		}

		setRequireParamsWithOpts(c,
			WithUserId(userId),
		)

		c.Next()
	}
}
