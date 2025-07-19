package site

import (
	"context"
	"fmt"

	"serverApi/pkg/queue"
	"serverApi/pkg/tools/cast"

	"serverApi/pkg/captcha"
	"serverApi/pkg/constant"
	constsR "serverApi/pkg/constant/redis"
	table "serverApi/pkg/db/table/site"
	commonPb "serverApi/pkg/protobuf/common"
	pb "serverApi/pkg/protobuf/site"
	"serverApi/pkg/rocketmq"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/zlogger"
)

func (s *siteSrv) Register(ctx context.Context, req *pb.RegisterReq) (*commonPb.Empty, error) {
	var (
		parentId int
		sendReq  = captcha.NewSendReq(cast.ToInt(req.GetAccountType()), constant.SceneRegister, req.GetAreaCode(), req.GetMobile(), req.GetEmail())
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

	// 校验验证码
	if err := s.checkSmsCode(ctx, sendReq, req.GetCode()); err != nil {
		return nil, errs.ErrInternalServer.Wrap("captcha_error")
	}

	// 校验用户是否已经注册
	if err := s.checkUserAlreadyRegistered(ctx, req); err != nil {
		return nil, err
	}

	// 查询邀请码
	parentId, err := s.checkInviteCode(ctx, req.GetInviteCode())
	if err != nil {
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 获取国家编码
	code, err := s.siteDB.GetGlobalAreaByMobileCode(req.GetAreaCode())
	if err != nil || code.IsEmpty() {
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 注册用户
	userId, err := s.siteDB.Register(ctx, req, parentId, code.CountryCode, s.SysDefaultAvatar(ctx))
	if err != nil {
		zlogger.Errorf("Register Register |req:%v| err: %v", cast.ToString(req), err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	// 发送创建聊天用户消息
	rocketmq.PublishJson(rocketmq.SiteUserRegister, &table.User{
		Id: userId,
	})

	// 用户注册事件
	rocketmq.PublishJson(rocketmq.StatsEvent, &queue.EventStats{EventType: queue.EventUserRegistrations})

	// 删除短信验证码
	captcha.DelCode(ctx, s.redis, sendReq)

	return &commonPb.Empty{}, nil
}

// 校验短信验证码
func (s *siteSrv) checkSmsCode(ctx context.Context, sendReq *captcha.SendReq, code string) error {
	if err := captcha.CheckCode(ctx, s.redis, sendReq, code); err != nil {
		zlogger.Errorf("CheckSmsCode |sendReq:%v| err: %v", sendReq, err)
		return err
	}

	return nil
}

// 校验用户是否已经注册
func (s *siteSrv) checkUserAlreadyRegistered(ctx context.Context, req *pb.RegisterReq) error {
	auth, err := s.siteDB.GetUserAuth(ctx,
		table.WhereUserAuthAreaCode(req.GetAreaCode()),
		table.WhereUserAuthMobile(req.GetMobile()),
		table.WhereUserAuthEmail(req.GetEmail()),
	)
	if err != nil {
		zlogger.Errorf("Register GetUserAuth |req:%v| err: %v", req, err)
		return errs.ErrInternalServer.Wrap("operation_failed")
	}

	if !auth.IsEmpty() {
		zlogger.Errorf("Register GetUserAuth |req:%v| err: mobile phone number has been registered", req)
		return errs.ErrInternalServer.Wrap("user_exist")
	}

	return nil
}

// 校验邀请码
func (s *siteSrv) checkInviteCode(ctx context.Context, inviteCode string) (int, error) {
	if inviteCode == "" {
		return 0, nil
	}

	inviteUser, err := s.siteDB.FindUserInfo(ctx, table.WhereUserInviteCode(inviteCode))
	if err != nil {
		zlogger.Errorf("Register CheckInviteCode |inviteCode:%v| err: %v", inviteCode, err)
		return 0, err
	}

	if inviteUser.IsEmpty() {
		zlogger.Errorf("Register CheckInviteCode |inviteCode:%v| err: invitation code does not exist", inviteCode)
		return 0, nil
	}

	return inviteUser.Id, nil
}
