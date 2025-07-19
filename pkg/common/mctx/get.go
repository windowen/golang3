package mctx

import (
	"context"
	"strconv"
	"time"

	"serverApi/pkg/tools/cast"
	"serverApi/pkg/zlogger"

	"serverApi/pkg/constant"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/tools/utils"
)

func HaveOpUser(ctx context.Context) bool {
	return ctx.Value(constant.RpcOpUserId) != nil
}

func Check(ctx context.Context) (int, int32, error) {
	opUserIDVal := ctx.Value(constant.RpcOpUserId)
	opUserID := cast.ToInt(opUserIDVal)
	if opUserID == 0 {
		return 0, 0, errs.ErrNoPermission.Wrap("opuser_id_empty")
	}
	opUserTypeArr, ok := ctx.Value(constant.RpcOpUserType).([]string)
	if !ok {
		return 0, 0, errs.ErrNoPermission.Wrap("missing_user_type")
	}
	if len(opUserTypeArr) == 0 {
		return 0, 0, errs.ErrNoPermission.Wrap("user type empty")
	}
	userType, err := strconv.Atoi(opUserTypeArr[0])
	if err != nil {
		return 0, 0, errs.ErrNoPermission.Wrap("user type invalid " + err.Error())
	}
	if !(userType == constant.AdminUser || userType == constant.NormalUser) {
		return 0, 0, errs.ErrNoPermission.Wrap("user type invalid")
	}
	return opUserID, int32(userType), nil
}

func CheckUser(ctx context.Context) (int, error) {
	userID, userType, err := Check(ctx)
	if err != nil {
		return 0, err
	}
	if userType != constant.NormalUser {
		return 0, errs.ErrNoPermission.Wrap("not user")
	}
	return userID, nil
}

func GetOpUserId(ctx context.Context) int {
	userID, _ := ctx.Value(constant.OpUserId).(int)
	return userID
}

func GetCountryCode(ctx context.Context) string {
	countryCode, _ := ctx.Value(constant.CountryCode).(string)
	return countryCode
}

func GetLanguage(ctx context.Context) string {
	language, _ := ctx.Value(constant.Language).(string)
	return language
}

func GetUserType(ctx context.Context) (int, error) {
	userTypeArr, _ := ctx.Value(constant.RpcOpUserType).([]string)
	userType, err := strconv.Atoi(userTypeArr[0])
	if err != nil {
		return 0, errs.ErrNoPermission.Wrap("user type invalid " + err.Error())
	}
	return userType, nil
}

// GetLocation 获取时区
func GetLocation(ctx context.Context) *time.Location {
	locationName, _ := ctx.Value(constant.Location).(string)

	location, err := time.LoadLocation(locationName)
	if err != nil {
		zlogger.Errorf("GetLocation |location:%v| err: %v", locationName, err)
		return time.Local
	}

	return location
}

func WithOpUserId(ctx context.Context, opUserID string, userType int) context.Context {
	headers, _ := ctx.Value(constant.RpcCustomHeader).([]string)
	ctx = context.WithValue(ctx, constant.RpcOpUserId, opUserID)
	ctx = context.WithValue(ctx, constant.RpcOpUserType, []string{strconv.Itoa(userType)})
	if utils.IndexOf(constant.RpcOpUserType, headers...) < 0 {
		ctx = context.WithValue(ctx, constant.RpcCustomHeader, append(headers, constant.RpcOpUserType))
	}
	return ctx
}

func WithApiToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, constant.CtxApiToken, token)
}
