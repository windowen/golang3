package site

import (
	"context"
	"fmt"
	"serverApi/pkg/captcha"
	constsR "serverApi/pkg/constant/redis"

	"serverApi/pkg/common/mctx"
	"serverApi/pkg/constant"
	"serverApi/pkg/db/cache"
	table "serverApi/pkg/db/table/site"
	commonPb "serverApi/pkg/protobuf/common"
	pb "serverApi/pkg/protobuf/site"
	"serverApi/pkg/tools/cast"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/zlogger"
)

// SecurityAuthInfo 认证信息
func (s *siteSrv) SecurityAuthInfo(ctx context.Context, empty *commonPb.Empty) (*pb.SecurityAuthInfoResp, error) {
	userId, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("SecurityAuthInfo CheckUser | err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 游客限制
	if userId == 0 {
		return nil, errs.ErrNoPermission.Wrap("token_not_permissions")
	}

	// 获取用户缓存信息
	userCacheInfo, err := s.findUserCacheInfo(ctx, userId)
	if err != nil {
		zlogger.Errorf("SecurityAuthInfo findUserCacheInfo |userId:%v| err: %v", userId, err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	return &pb.SecurityAuthInfoResp{
		AreaCode: userCacheInfo.AreaCode,
		Mobile:   userCacheInfo.Mobile,
		Email:    userCacheInfo.Email,
	}, nil
}

// SecurityBinding 绑定
func (s *siteSrv) SecurityBinding(ctx context.Context, req *pb.SecurityBindingReq) (*commonPb.Empty, error) {
	userId, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("SecurityBinding CheckUser | err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 游客限制
	if userId == constant.Zero {
		return nil, errs.ErrNoPermission.Wrap("token_not_permissions")
	}

	// 获取用户缓存信息
	userCacheInfo, err := s.findUserCacheInfo(ctx, userId)
	if err != nil {
		zlogger.Errorf("SecurityBinding findUserCacheInfo |userId:%v| err: %v", userId, err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 获取锁
	lockKey := fmt.Sprintf(constsR.LockUserOperate, userId)
	isLock := s.redisLock.TryLock(ctx, lockKey, lockKey)
	if !isLock {
		zlogger.Errorf("SecurityBinding TryLock |lockKey:%v| err: failed to acquire lock", lockKey)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	defer func() {
		ret := s.redisLock.ReleaseLock(ctx, lockKey, lockKey)
		if !ret {
			zlogger.Errorf("SecurityBinding ReleaseLock |lockKey:%v| err: failed to release lock", lockKey)
		}
	}()

	// 验证码验证
	if err := s.checkSmsCode(ctx, captcha.NewSendReq(cast.ToInt(req.GetAccountType()), constant.SceneBind, req.GetAreaCode(), req.GetMobile(), req.GetEmail()), req.GetCode()); err != nil {
		return nil, errs.ErrInternalServer.Wrap("captcha_error")
	}

	// 查询换绑配置
	securityChangeBind, err := cache.FindSysSiteConfig(ctx, s.redis, constant.SysSiteCfgSecurityChangeBind)
	if err != nil {
		zlogger.Errorf("SecurityBinding FindSysSiteConfig |key:%v| err: %v", constant.SysSiteCfgSecurityChangeBind, err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 是否开启换绑
	if cast.ToInt(securityChangeBind) == constant.ON {
		// 检查是否绑定
		if cast.ToInt(req.GetAccountType()) == constant.UserAccountTypeMobile {
			if userCacheInfo.AreaCode != "" && userCacheInfo.Mobile != "" {
				zlogger.Errorf("SecurityBinding |userId:%v,req:%v| err: mobile phone bound", userId, req)
				return nil, errs.ErrInternalServer.Wrap("user_band_mobile")
			}
		}

		if cast.ToInt(req.GetAccountType()) == constant.UserAccountTypeEmail {
			if userCacheInfo.Email != "" {
				zlogger.Errorf("SecurityBinding |userId:%v,req:%v| err: email address bound", userId, req)
				return nil, errs.ErrInternalServer.Wrap("user_band_email")
			}
		}
	}

	// 检查手机邮箱是否使用
	auth, err := s.siteDB.GetUserAuth(ctx,
		table.WhereUserAuthAreaCode(req.GetAreaCode()),
		table.WhereUserAuthMobile(req.GetMobile()),
		table.WhereUserAuthEmail(req.GetEmail()),
	)
	if err != nil {
		zlogger.Errorf("SecurityBinding GetUserAuth |req:%v| err: %v", req, err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	if !auth.IsEmpty() {
		switch cast.ToInt(req.AccountType) {
		case constant.UserAccountTypeMobile:
			return nil, errs.ErrInternalServer.Wrap("user_band_mobile_exist")
		case constant.UserAccountTypeEmail:
			return nil, errs.ErrInternalServer.Wrap("user_band_email_exist")
		}
	}

	err = s.siteDB.UpdateUserAuth(ctx, userId,
		table.SetUserAuthAreaCode(req.GetAreaCode()),
		table.SetUserAuthMobile(req.GetMobile()),
		table.SetUserAuthEmail(req.GetEmail()),
	)
	if err != nil {
		zlogger.Errorf("SecurityBinding UpdateUserAuth |userId:%v,req:%v| err: %v", userId, req, err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	err = s.userCache.UpdateUserOption(ctx, userId,
		cache.UpdateUserAreaCode(req.GetAreaCode(), false),
		cache.UpdateUserMobile(req.GetMobile(), false),
		cache.UpdateUserEmail(req.GetEmail(), false),
	)
	if err != nil {
		zlogger.Errorf("SecurityBinding userCache.UpdateUserOption |userId:%v,req:%v| err: %v", userId, req, err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	return nil, nil
}
