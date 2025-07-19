package jwt

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"

	"serverApi/pkg/tools/cast"
	"serverApi/pkg/tools/errs"

	"serverApi/pkg/common/config"
	constsR "serverApi/pkg/constant/redis"
	"serverApi/pkg/db/cache"
	"serverApi/pkg/zlogger"
)

var ErrTokenExpired = errors.New("token is expired")

type JWT struct {
	key []byte
	rdb redis.UniversalClient
}

type MyCustomClaims struct {
	UserId  int `json:"userId"`
	Version int `json:"version"`
	jwt.RegisteredClaims
}

type GenTokenResult struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshTokenResult struct {
	Token  string `json:"token"`
	UserId int    `json:"userId"`
}

func NewJWT(redisClient redis.UniversalClient) *JWT {
	return &JWT{
		key: []byte(config.Config.Jwt.Key),
		rdb: redisClient,
	}
}

// generateToken 生成Token
func (j *JWT) generateToken(userId, version int, validityPeriod time.Duration) (string, error) {
	var (
		timeNow   = time.Now()
		expiresAt = timeNow.Add(validityPeriod)
	)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		UserId:  userId,
		Version: version,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "",
			ExpiresAt: jwt.NewNumericDate(expiresAt), // 7天有效期
			IssuedAt:  jwt.NewNumericDate(timeNow),   // 签发时间
			NotBefore: jwt.NewNumericDate(timeNow),   // 生效时间
		},
	})

	tokenString, err := token.SignedString(j.key)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Bearer %s", tokenString), nil
}

// ParseToken 解析token
func (j *JWT) ParseToken(ctx context.Context, tokenString string) (*MyCustomClaims, error) {
	// 移除 Bearer 前缀
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			zlogger.Errorf("ParseToken Method | err: unexpected signing method")
			return nil, errs.ErrNoPermission.WithDetail("token_match_fail")
		}
		return j.key, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		zlogger.Errorf("ParseToken ParseWithClaims | err: %v", err)
		return nil, errs.ErrNoPermission.WithDetail("token_match_fail")
	}

	// 检查 token 是否有效
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		// 版本检查
		storedVersion, err := j.rdb.Get(ctx, fmt.Sprintf(constsR.UserLoginTokenVersion, claims.UserId)).Result()
		if err := cache.CheckErr(err); err != nil {
			zlogger.Errorf("ParseToken UserLoginTokenVersion | err: %v", err)
			return nil, errs.ErrRemotePlaceError.WithDetail("token_remote_login")
		}

		if cast.ToInt(storedVersion) != claims.Version {
			zlogger.Errorf("ParseToken UserLoginTokenVersion |userId:%v,cacheVersion:%v,tokenVersion:%v| err: user remote login", claims.UserId, storedVersion, claims.Version)
			return nil, errs.ErrRemotePlaceError.WithDetail("token_remote_login")
		}

		return claims, nil
	}

	return nil, errs.ErrNoPermission.WithDetail("token_expired")
}

// GenToken 获取携带版本的Token
func (j *JWT) GenToken(ctx context.Context, userId int) (*GenTokenResult, error) {
	// 获取用户token当前版本
	storedVersion, err := j.rdb.Get(ctx, fmt.Sprintf(constsR.UserLoginTokenVersion, userId)).Result()
	if err := cache.CheckErr(err); err != nil {
		return nil, fmt.Errorf("failed to get token version: %v", err)
	}

	newVersion := 1
	if storedVersion != "" {
		var version int
		version = cast.ToInt(storedVersion)

		newVersion = version + 1
	}

	token, err := j.generateToken(userId, newVersion, constsR.TokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	err = j.rdb.Set(ctx, fmt.Sprintf(constsR.UserLoginTokenVersion, userId), newVersion, constsR.RefreshTokenTTL).Err()
	if err := cache.CheckErr(err); err != nil {
		return nil, fmt.Errorf("failed to save token: %v", err)
	}

	err = j.rdb.Set(ctx, fmt.Sprintf(constsR.UserLoginToken, userId), token, constsR.RefreshTokenTTL).Err()
	if err := cache.CheckErr(err); err != nil {
		return nil, fmt.Errorf("failed to save token version: %v", err)
	}

	refreshToken, err := j.generateToken(userId, newVersion, constsR.RefreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %v", err)
	}

	return &GenTokenResult{
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (j *JWT) DelToken(ctx context.Context, userId int) {
	err := j.rdb.Del(ctx, fmt.Sprintf(constsR.UserLoginToken, userId)).Err()
	if err := cache.CheckErr(err); err != nil {
		zlogger.Errorf("failed to del user token cache, err: %v", err)
		return
	}

	err = j.rdb.Del(ctx, fmt.Sprintf(constsR.UserLoginTokenVersion, userId)).Err()
	if err := cache.CheckErr(err); err != nil {
		zlogger.Errorf("failed to del user token version cache, err: %v", err)
		return
	}
}

// RefreshToken 刷新并生成新的 Token
func (j *JWT) RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResult, error) {
	// 解析并验证刷新令牌
	claims, err := j.ParseToken(ctx, refreshToken)
	if err != nil {
		return nil, errs.ErrNoPermission.WithDetail("token_expired")
	}

	// 获取当前用户的 token 版本
	storedVersion, err := j.rdb.Get(ctx, fmt.Sprintf(constsR.UserLoginTokenVersion, claims.UserId)).Result()
	if err := cache.CheckErr(err); err != nil {
		return nil, errs.ErrNoPermission.WithDetail("token_expired")
	}

	if cast.ToInt(storedVersion) != claims.Version {
		return nil, errs.ErrRemotePlaceError.WithDetail("token_remote_login")
	}

	// 生成新的访问令牌
	newToken, err := j.generateToken(claims.UserId, claims.Version, constsR.TokenTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new token: %v", err)
	}

	err = j.rdb.Set(ctx, fmt.Sprintf(constsR.UserLoginTokenVersion, claims.UserId), claims.Version, constsR.RefreshTokenTTL).Err()
	if err := cache.CheckErr(err); err != nil {
		return nil, errs.ErrNoPermission.WithDetail("token_expired")
	}

	err = j.rdb.Set(ctx, fmt.Sprintf(constsR.UserLoginToken, claims.UserId), newToken, constsR.RefreshTokenTTL).Err()
	if err := cache.CheckErr(err); err != nil {
		return nil, errs.ErrNoPermission.WithDetail("token_expired")
	}

	return &RefreshTokenResult{
		Token:  newToken,
		UserId: claims.UserId,
	}, nil
}
