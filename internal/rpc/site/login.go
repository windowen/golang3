package site

import (
	"context"
	"errors"
	"fmt"
	"time"

	"serverApi/pkg/jwt"
	"serverApi/pkg/tools/cast"

	"serverApi/pkg/captcha"
	"serverApi/pkg/constant"
	constsR "serverApi/pkg/constant/redis"
	"serverApi/pkg/db/cache"
	table "serverApi/pkg/db/table/site"
	pb "serverApi/pkg/protobuf/site"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/tools/passwordhelper"
	"serverApi/pkg/zlogger"
)

func (s *siteSrv) Login(ctx context.Context, req *pb.LoginReq) (*pb.LoginResp, error) {
	var (
		err         error
		auth        *table.UserAuth
		tokenResult *jwt.GenTokenResult
		sendReq     = captcha.NewSendReq(cast.ToInt(req.GetAccountType()), constant.SceneLogin, req.GetAreaCode(), req.GetMobile(), req.GetEmail())
	)

	// 获取锁
	lockKey := fmt.Sprintf(constsR.LockUserOperate, req.GetAreaCode()+req.GetMobile())
	isLock := s.redisLock.TryLock(ctx, lockKey, lockKey)
	if !isLock {
		zlogger.Errorf("Login TryLock |lockKey:%v| err: failed to acquire lock", lockKey)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	defer func() {
		ret := s.redisLock.ReleaseLock(ctx, lockKey, lockKey)
		if !ret {
			zlogger.Errorf("Login ReleaseLock |lockKey:%v| err: failed to release lock", lockKey)
		}
	}()

	// 手机号登录校验
	if req.GetAuthType() == constant.UserAuthTypeCode {
		if err := s.checkSmsCode(ctx, sendReq, req.GetCode()); err != nil {
			return nil, errs.ErrInternalServer.Wrap("captcha_error")
		}
	}

	// 获取用户认证信息
	auth, err = s.getUserAuth(ctx, table.WhereUserAuthAreaCode(req.GetAreaCode()), table.WhereUserAuthMobile(req.GetMobile()), table.WhereUserAuthEmail(req.GetEmail()))
	if err != nil {
		return nil, errs.ErrInternalServer.Wrap("user_password_error")
	}

	// 密码登录校验
	if req.GetAuthType() == constant.UserAuthTypePassword && !passwordhelper.VerifyPassword(auth.Password, req.GetPassword()) {
		return nil, errs.ErrInternalServer.Wrap("user_password_error")
	}

	// 获取用户信息
	userInfo, err := s.getUserInfo(ctx, auth.UserId)
	if err != nil {
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 生成Token
	if tokenResult, err = s.generateToken(ctx, userInfo.Id); err != nil {
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	// 删除验证码（如果是手机登录）
	if req.GetAuthType() == constant.UserAuthTypeCode {
		captcha.DelCode(ctx, s.redis, sendReq)
	}

	// 缓存用户数据
	if err = s.cacheUserInfo(ctx, auth, userInfo); err != nil {
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	return &pb.LoginResp{Token: tokenResult.Token, RefreshToken: tokenResult.RefreshToken}, nil
}

// 获取用户认证信息
func (s *siteSrv) getUserAuth(ctx context.Context, opsWhere ...table.UserAuthWhereOption) (*table.UserAuth, error) {
	auth, err := s.siteDB.GetUserAuth(ctx, opsWhere...)
	if err != nil {
		zlogger.Errorf("Login GetUserAuth |opsWhere:%v| err: %v", opsWhere, err)
		return nil, err
	}

	if auth.IsEmpty() {
		zlogger.Errorf("Login GetUserAuth |opsWhere:%v| err: user does not exist", opsWhere)
		return nil, errors.New("user does not exist")
	}

	return auth, nil
}

// 获取用户信息
func (s *siteSrv) getUserInfo(ctx context.Context, userId int) (*table.User, error) {
	userInfo, err := s.siteDB.FindUserInfo(ctx, table.WhereUserId(userId))
	if err != nil {
		zlogger.Errorf("Login GetUserInfo |userId:%v| err: %v", userId, err)
		return nil, err
	}

	if userInfo.IsEmpty() {
		zlogger.Errorf("Login GetUserInfo |userId:%v| err: user does not exist", userId)
		return nil, errors.New("user does not exist")
	}

	return userInfo, nil
}

// 生成Token
func (s *siteSrv) generateToken(ctx context.Context, userId int) (*jwt.GenTokenResult, error) {
	result, err := s.jwt.GenToken(ctx, userId)
	if err != nil {
		zlogger.Errorf("Login GenToken |userId:%v| err: %v", userId, err)
		return nil, err
	}

	return result, nil
}

// 缓存用户数据
func (s *siteSrv) cacheUserInfo(ctx context.Context, auth *table.UserAuth, userInfo *table.User) error {
	var (
		isFamilyMaster, familyId, familyMasterId int
		err                                      error
	)

	// 获取并缓存家族信息
	if familyId, familyMasterId, isFamilyMaster, err = s.getFamilyInfo(ctx, userInfo); err != nil {
		zlogger.Errorf("Login getFamilyInfo |userId:%v| err: %v", userInfo.Id, err)
		return err
	}

	// 获取用户坐骑信息
	mountInfo, err := s.siteDB.FindUserUseMount(ctx, table.WhereUserMountsUserId(userInfo.Id), table.WhereUserMountsIsSelected(1), table.WhereUserMountsExpiredTime(time.Now()))
	if err != nil {
		zlogger.Errorf("Login FindUserUseMount |userId:%v| err: %v", userInfo.Id, err)
		return err
	}
	userMountId := 0
	if !mountInfo.IsEmpty() {
		userMountId = mountInfo.MountsId
	}

	// 普通用户和主播的缓存处理
	if userInfo.Category != constant.UserCategoryAnchor {
		err = s.cacheNormalUser(ctx, auth, userInfo, 0, isFamilyMaster, familyId, familyMasterId, userMountId)
	} else {
		err = s.cacheAnchorUser(ctx, auth, userInfo, isFamilyMaster, familyId, familyMasterId, userMountId)
	}

	return err
}

// 获取家族信息
func (s *siteSrv) getFamilyInfo(ctx context.Context, userInfo *table.User) (int, int, int, error) {
	familyInfo, err := s.siteDB.FindFamilyAnchors(ctx, table.WhereFamilyAnchorsUserId(userInfo.Id))
	if err != nil {
		return 0, 0, 0, err
	}

	if !familyInfo.IsEmpty() {
		return familyInfo.FamilyId, familyInfo.PatriarchId, 1, err
	}

	family, err := s.siteDB.FindFamily(ctx, table.WhereFamilyUserId(userInfo.Id))
	if err != nil {
		return 0, 0, 0, err
	}

	if !family.IsEmpty() {
		return family.Id, family.UserId, 1, nil
	}

	return 0, 0, 0, nil
}

// 缓存普通用户
func (s *siteSrv) cacheNormalUser(ctx context.Context, auth *table.UserAuth, userInfo *table.User, roomId, isFamilyMaster, familyId, familyMasterId, userMountId int) error {
	err := s.userCache.Save(ctx, userInfo.Id, cache.NewUserCacheInfo(
		userInfo.Id, userInfo.CountryCode, auth.AreaCode, auth.Mobile, auth.Email, userInfo.Nickname,
		userInfo.Avatar, userInfo.Sign, userInfo.Sex, userInfo.Birthday, userInfo.Feeling,
		userInfo.Country, userInfo.Area, userInfo.Profession, userInfo.Category, userInfo.InviteCode,
		userInfo.ParentId, userInfo.LevelId, userInfo.SetLevelId, userInfo.Remark, userInfo.Status, auth.Password,
		userInfo.PayPassword, roomId, userInfo.ChatUuid, userInfo.GmStatus, isFamilyMaster, familyId, familyMasterId, userMountId,
	))
	if err != nil {
		zlogger.Errorf("Login cacheNormalUser userCache.Save |userId:%v| err: %v", userInfo.Id, err)
		return err
	}

	return nil
}

// 缓存主播用户
func (s *siteSrv) cacheAnchorUser(ctx context.Context, auth *table.UserAuth, userInfo *table.User, isFamilyMaster, familyId, familyMasterId, userMountId int) error {
	roomInfo, err := s.siteDB.FindRoomInfo(ctx, table.WhereRoomUserId(userInfo.Id))
	if err != nil {
		zlogger.Errorf("Login cacheAnchorUser FindRoomInfo |userId:%v| err: %v", userInfo.Id, err)
		return err
	}

	if roomInfo.IsEmpty() {
		zlogger.Errorf("Login cacheAnchorUser FindRoomInfo |userId:%v| err: the live broadcast room does not exist", userInfo.Id)
		return errors.New("room does not exist")
	}

	err = s.cacheNormalUser(ctx, auth, userInfo, roomInfo.Id, isFamilyMaster, familyId, familyMasterId, userMountId)
	if err != nil {
		return err
	}

	err = s.roomCache.Save(ctx, roomInfo.Id, cache.NewRoomCacheInfo(
		roomInfo.Id, roomInfo.CountryCode, roomInfo.UserId, roomInfo.Title, roomInfo.TagsJson, roomInfo.Cover,
		roomInfo.VideoClarity, roomInfo.GiftRatio, roomInfo.PlatformRatio, roomInfo.FamilyRatio,
		roomInfo.Status, roomInfo.ChatRoomId, constant.RoomSceneStatusEnd, roomInfo.Summary, 0, 0, 0,
		roomInfo.PaidPurviewStatus, roomInfo.LevelId, roomInfo.SetLevelId, constant.RoomChargingRulesFree, constant.Zero, constant.Zero,
	))
	if err != nil {
		zlogger.Errorf("Login cacheAnchorUser roomCache.Save |userId:%v| err: %v", userInfo.Id, err)
		return err
	}

	return nil
}
