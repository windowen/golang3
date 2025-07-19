package site

import (
	"context"
	"fmt"

	"serverApi/pkg/tools/cast"
	"serverApi/pkg/tools/utils"

	"serverApi/pkg/captcha"
	"serverApi/pkg/common/mctx"
	"serverApi/pkg/constant"
	constsR "serverApi/pkg/constant/redis"
	"serverApi/pkg/db/cache"
	table "serverApi/pkg/db/table/site"
	commonPb "serverApi/pkg/protobuf/common"
	pb "serverApi/pkg/protobuf/site"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/tools/passwordhelper"
	"serverApi/pkg/zlogger"
)

// ModifyProfile 修改用户信息
func (s *siteSrv) ModifyProfile(ctx context.Context, req *pb.ModifyProfileReq) (*commonPb.Empty, error) {
	uid, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("ModifyProfile CheckUser | err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 获取锁
	lockKey := fmt.Sprintf(constsR.LockUserOperate, uid)
	isLock := s.redisLock.TryLock(ctx, lockKey, lockKey)
	if !isLock {
		zlogger.Errorf("ModifyProfile TryLock |lockKey:%v| err: failed to acquire lock", lockKey)
		return nil, nil
	}

	defer func() {
		ret := s.redisLock.ReleaseLock(ctx, lockKey, lockKey)
		if !ret {
			zlogger.Errorf("ModifyProfile ReleaseLock |lockKey:%v| err: failed to release lock", lockKey)
		}
	}()

	err = s.siteDB.UpdateUserInfo(ctx, uid,
		table.SetUserAvatar(req.GetAvatar()),
		table.SetUserNickname(req.GetNickname()),
		table.SetUserSign(req.GetSign()),
		table.SetUserSex(cast.ToInt(req.GetSex())),
		table.SetUserBirthday(req.GetBirthday()),
		table.SetUserFeeling(cast.ToInt(req.GetFeeling())),
		table.SetUserCountry(req.GetCountry()),
		table.SetUserArea(req.GetArea()),
		table.SetUserProfession(cast.ToInt(req.GetProfession())),
	)
	if err != nil {
		zlogger.Errorf("ModifyProfile UpdateUserInfo |userId:%v| err: %v", uid, err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	err = s.userCache.UpdateUserOption(ctx, uid,
		cache.UpdateUserAvatar(req.GetAvatar()),
		cache.UpdateUserNickname(req.GetNickname()),
		cache.UpdateUserSign(req.GetSign()),
		cache.UpdateUserSex(cast.ToInt(req.GetSex())),
		cache.UpdateUserBirthday(req.GetBirthday()),
		cache.UpdateUserFeeling(cast.ToInt(req.GetFeeling())),
		cache.UpdateUserCountry(req.GetCountry()),
		cache.UpdateUserArea(req.GetArea()),
		cache.UpdateUserProfession(cast.ToInt(req.GetProfession())),
	)
	if err != nil {
		zlogger.Errorf("ModifyProfile UserCache.UpdateUserOption |userId:%v| err: %v", uid, err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	return nil, nil
}

// ModifyPaymentPassword 修改用户支付密码
func (s *siteSrv) ModifyPaymentPassword(ctx context.Context, req *pb.ModifyPaymentReq) (*commonPb.Empty, error) {
	uid, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("ModifyPaymentPassword CheckUser | err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 获取锁
	lockKey := fmt.Sprintf(constsR.LockUserOperate, uid)
	isLock := s.redisLock.TryLock(ctx, lockKey, lockKey)
	if !isLock {
		zlogger.Errorf("ModifyPaymentPassword TryLock |lockKey:%v| err: failed to acquire lock", lockKey)
		return nil, nil
	}

	defer func() {
		ret := s.redisLock.ReleaseLock(ctx, lockKey, lockKey)
		if !ret {
			zlogger.Errorf("ModifyPaymentPassword ReleaseLock |lockKey:%v| err: failed to release lock", lockKey)
		}
	}()

	userCacheInfo, err := s.findUserCacheInfo(ctx, uid)
	if err != nil || userCacheInfo == nil {
		zlogger.Errorf("ModifyPaymentPassword UserCacheInfo |req:%v,err: %v", cast.ToString(req), err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	sendType := constant.UserAccountTypeMobile
	if userCacheInfo.Email != "" {
		sendType = constant.UserAccountTypeEmail
	}

	sendReq := captcha.NewSendReq(sendType, constant.SceneModifyPaymentPassword, userCacheInfo.AreaCode, userCacheInfo.Mobile, userCacheInfo.Email)

	if err := captcha.CheckCode(ctx, s.redis, sendReq, req.GetCode()); err != nil {
		zlogger.Errorf("ModifyPaymentPassword CheckSmsCode |req:%v| err: %v", cast.ToString(req), err)
		return nil, errs.ErrInternalServer.Wrap("invalid_captcha_data")
	}

	err = s.siteDB.UpdateUserInfo(
		ctx,
		uid,
		table.SetUserPayPassword(passwordhelper.GenPassword(req.GetPassword())),
	)
	if err != nil {
		zlogger.Errorf("ModifyPaymentPassword UpdateUserInfo |req:%v| err: %v", cast.ToString(req), err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	err = s.userCache.UpdateUserOption(ctx, uid,
		cache.UpdateUserPayPassword(passwordhelper.GenPassword(req.GetPassword())),
	)
	if err != nil {
		zlogger.Errorf("ModifyProfile UserCache.UpdateUserOption |userId:%v| err: %v", uid, err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	captcha.DelCode(ctx, s.redis, sendReq)

	return nil, err
}

// ModifyLoginPassword 修改用户登陆密码
func (s *siteSrv) ModifyLoginPassword(ctx context.Context, req *pb.ModifyLoginPasswordReq) (*commonPb.Empty, error) {
	uid, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("ModifyLoginPassword CheckUser | err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 获取锁
	lockKey := fmt.Sprintf(constsR.LockUserOperate, uid)
	isLock := s.redisLock.TryLock(ctx, lockKey, lockKey)
	if !isLock {
		zlogger.Errorf("ModifyLoginPassword TryLock |lockKey:%v| err: failed to acquire lock", lockKey)
		return nil, nil
	}

	defer func() {
		ret := s.redisLock.ReleaseLock(ctx, lockKey, lockKey)
		if !ret {
			zlogger.Errorf("ModifyLoginPassword ReleaseLock |lockKey:%v| err: failed to release lock", lockKey)
		}
	}()

	userCacheInfo, err := s.findUserCacheInfo(ctx, uid)
	if err != nil || userCacheInfo == nil {
		zlogger.Errorf("ModifyLoginPassword UserCacheInfo |req:%v| err: %v", cast.ToString(req), err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	sendType := constant.UserAccountTypeMobile
	if userCacheInfo.Email != "" {
		sendType = constant.UserAccountTypeEmail
	}
	sendReq := captcha.NewSendReq(sendType, constant.SceneModifyPaymentPassword, userCacheInfo.AreaCode, userCacheInfo.Mobile, userCacheInfo.Email)

	if err := captcha.CheckCode(ctx, s.redis, sendReq, req.GetCode()); err != nil {
		zlogger.Errorf("ModifyLoginPassword CheckSmsCode |req:%v| err: %v", cast.ToString(req), err)
		return nil, errs.ErrInternalServer.Wrap("invalid_captcha_data")
	}

	err = s.siteDB.UpdateUserAuth(ctx, uid, table.SetUserAuthPassword(passwordhelper.GenPassword(req.GetPassword())))
	if err != nil {
		zlogger.Errorf("ModifyLoginPassword UpdateUserAuth |req:%v| err: %v", cast.ToString(req), err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	s.jwt.DelToken(ctx, uid)
	captcha.DelCode(ctx, s.redis, sendReq)

	return nil, err
}

// Logout 登出
func (s *siteSrv) Logout(ctx context.Context, _ *commonPb.Empty) (*commonPb.Empty, error) {
	uid, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("Logout CheckUser | err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 获取锁
	lockKey := fmt.Sprintf(constsR.LockUserOperate, uid)
	isLock := s.redisLock.TryLock(ctx, lockKey, lockKey)
	if !isLock {
		zlogger.Errorf("Logout TryLock |lockKey:%v| err: failed to acquire lock", lockKey)
		return nil, nil
	}

	defer func() {
		ret := s.redisLock.ReleaseLock(ctx, lockKey, lockKey)
		if !ret {
			zlogger.Errorf("Logout ReleaseLock |lockKey:%v| err: failed to release lock", lockKey)
		}
	}()

	s.jwt.DelToken(ctx, uid)

	return nil, err
}

// ForgetPassword 忘记密码
func (s *siteSrv) ForgetPassword(ctx context.Context, req *pb.ForgetPasswordReq) (*commonPb.Empty, error) {
	var (
		err     error
		sendReq = captcha.NewSendReq(cast.ToInt(req.GetAccountType()), constant.SceneForgetPassword, req.GetAreaCode(), req.GetMobile(), req.GetEmail())
	)

	// 获取锁
	lockKey := fmt.Sprintf(constsR.LockUserOperate, req.GetAreaCode()+req.GetMobile())
	isLock := s.redisLock.TryLock(ctx, lockKey, lockKey)
	if !isLock {
		zlogger.Errorf("ForgetPassword TryLock |lockKey:%v| err: failed to acquire lock", lockKey)
		return nil, nil
	}

	defer func() {
		ret := s.redisLock.ReleaseLock(ctx, lockKey, lockKey)
		if !ret {
			zlogger.Errorf("ForgetPassword ReleaseLock |lockKey:%v| err: failed to release lock", lockKey)
		}
	}()

	if err = captcha.CheckCode(ctx, s.redis, sendReq, req.GetCode()); err != nil {
		zlogger.Errorf("ForgetPassword CheckSmsCode |req:%v| err: %v", cast.ToString(req), err)
		return nil, errs.ErrInternalServer.Wrap("invalid_captcha_data")
	}

	auth, err := s.siteDB.GetUserAuth(
		ctx,
		table.WhereUserAuthAreaCode(req.GetAreaCode()),
		table.WhereUserAuthMobile(req.GetMobile()),
		table.WhereUserAuthEmail(req.GetEmail()),
	)
	if err != nil {
		zlogger.Errorf("ForgetPassword GetUserAuth |req:%v| err: %v", cast.ToString(req), err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	if auth.IsEmpty() {
		zlogger.Errorf("ForgetPassword |req:%v| err: account information does not exist", cast.ToString(req))
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	err = s.siteDB.UpdateUserAuth(ctx, auth.UserId, table.SetUserAuthPassword(passwordhelper.GenPassword(req.GetPassword())))
	if err != nil {
		zlogger.Errorf("ForgetPassword UpdateUserAuth |req:%v| err: %v", cast.ToString(req), err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	captcha.DelCode(ctx, s.redis, sendReq)

	return nil, nil
}

// GetProfile 获取用户信息
func (s *siteSrv) GetProfile(ctx context.Context, req *pb.ProfileReq) (*pb.ProfileResp, error) {
	var (
		userId     = cast.ToInt(req.GetUserId())
		roomId     = constant.ON
		liveStatus = constant.RoomSceneStatusEnd
		isFollow   = constant.ON
	)

	uid, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("GetProfile CheckUser | err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 限制游客
	if uid == constant.Zero {
		zlogger.Errorf("GetProfile |userId:%v| err: visitors are prohibited from operating", uid)
		return nil, errs.ErrNoPermission.Wrap("token_not_permissions")
	}

	// 查看他们信息
	if userId == constant.Zero {
		userId = uid
	}

	// 获取用户缓存
	userCacheInfo, err := s.findUserCacheInfo(ctx, userId)
	if err != nil || userCacheInfo == nil {
		zlogger.Errorf("GetProfile UserCacheInfo |userId:%v| err: %v", userId, err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 是否主播
	if userCacheInfo.Category == constant.UserCategoryAnchor {
		// 获取直播间信息
		roomCacheInfo, err := s.findRoomCacheInfo(ctx, userCacheInfo.RoomId)
		if err != nil {
			zlogger.Errorf("GetProfile RoomCacheInfo |roomId:%v| err: %v", userCacheInfo.RoomId, err)
			return nil, errs.ErrInternalServer.Wrap("query_failed")
		}

		roomId = roomCacheInfo.Id
		liveStatus = roomCacheInfo.LiveStatus
	}

	if userId != uid {
		isEx, err := s.userCache.IsFollow(ctx, constsR.UserFollowCacheKey, uid, userId)
		if err != nil {
			zlogger.Errorf("GetProfile IsFollow |userId:%v| err: %v", userId, err)
			return nil, errs.ErrInternalServer.Wrap("query_failed")
		}
		if isEx {
			isFollow = constant.Yes
		}
	}

	return &pb.ProfileResp{
		UserId:         cast.ToInt32(userCacheInfo.Id),
		CountryCode:    userCacheInfo.CountryCode,
		Avatar:         userCacheInfo.Avatar,
		Nickname:       userCacheInfo.Nickname,
		Sex:            cast.ToInt32(userCacheInfo.Sex),
		Sign:           userCacheInfo.Sign,
		Birthday:       userCacheInfo.Birthday,
		Feeling:        cast.ToInt32(userCacheInfo.Feeling),
		Country:        userCacheInfo.Country,
		Area:           userCacheInfo.Area,
		Category:       cast.ToInt32(userCacheInfo.Category),
		Status:         cast.ToInt32(userCacheInfo.Status),
		ChatUUID:       userCacheInfo.ChatUuid,
		RoomId:         cast.ToInt32(roomId),
		LiveStatus:     cast.ToInt32(liveStatus),
		LevelId:        cast.ToInt32(utils.CompareMax(userCacheInfo.SetLevelId, userCacheInfo.LevelId)),
		InviteCode:     userCacheInfo.InviteCode,
		Profession:     cast.ToInt32(userCacheInfo.Profession),
		IsFamilyMaster: cast.ToInt32(userCacheInfo.IsFamilyMaster),
		IsFollow:       cast.ToInt32(isFollow),
	}, nil
}

// FindUserInfo 查询用户信息
func (s *siteSrv) FindUserInfo(ctx context.Context, req *pb.FindUserInfoReq) (*pb.FindUserInfoResp, error) {
	// 获取缓存用户信息
	userCacheInfo, err := s.findUserCacheInfo(ctx, cast.ToInt(req.GetUserId()))
	if err != nil {
		zlogger.Errorf("FindUserInfo UserCacheInfo|req:%v| err: %v", cast.ToString(req), err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	return &pb.FindUserInfoResp{
		UserId:      cast.ToInt32(userCacheInfo.Id),
		CountryCode: userCacheInfo.CountryCode,
		Avatar:      userCacheInfo.Avatar,
		Nickname:    userCacheInfo.Nickname,
		Category:    cast.ToInt32(userCacheInfo.Category),
		Status:      cast.ToInt32(userCacheInfo.Status),
		ChatUUID:    userCacheInfo.ChatUuid,
	}, nil
}
