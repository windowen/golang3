package site

import (
	"context"
	"fmt"
	"time"

	"serverApi/pkg/tools/cast"
	"serverApi/pkg/tools/utils"

	"serverApi/pkg/common/mctx"
	"serverApi/pkg/constant"
	constsR "serverApi/pkg/constant/redis"
	table "serverApi/pkg/db/table/site"
	commonPb "serverApi/pkg/protobuf/common"
	"serverApi/pkg/protobuf/live"
	pb "serverApi/pkg/protobuf/site"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/zlogger"
)

// FollowUser 关注
func (s *siteSrv) FollowUser(ctx context.Context, req *pb.FollowReq) (*commonPb.Empty, error) {
	uid, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("FollowUser Check user err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	if uid == constant.Zero {
		zlogger.Errorf("FollowUser |userId:%v,followId:%v| err: visitors are prohibited from operating", uid, req.GetFollowId())
		return nil, errs.ErrNoPermission.Wrap("token_not_permissions")
	}

	// 获取锁
	lockKey := fmt.Sprintf(constsR.LockUserOperate, uid)
	isLock := s.redisLock.TryLock(ctx, lockKey, lockKey)
	if !isLock {
		zlogger.Errorf("FollowUser TryLock |lockKey:%v| err: failed to acquire lock", lockKey)
		return nil, nil
	}

	defer func() {
		ret := s.redisLock.ReleaseLock(ctx, lockKey, lockKey)
		if !ret {
			zlogger.Errorf("FollowUser ReleaseLock |lockKey:%v| err: failed to release lock", lockKey)
		}
	}()

	if uid == cast.ToInt(req.GetFollowId()) {
		zlogger.Errorf("FollowUser |userId:%v,followId:%v| err: CAN To FOCUS ON MYSELF", uid, req.GetFollowId())
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	// 被关注者用户信息
	followCacheInfo, err := s.findUserCacheInfo(ctx, cast.ToInt(req.GetFollowId()))
	if err != nil {
		zlogger.Errorf("FollowUser findUserCacheInfo |followId:%v| err: %v", req.GetFollowId(), err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 是否是主播
	if followCacheInfo.Category != constant.UserCategoryAnchor {
		zlogger.Errorf("FollowUser findUserCacheInfo |followId:%v| err: the person being followed is not the anchor", req.GetFollowId())
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	userCacheInfo, err := s.findUserCacheInfo(ctx, uid)
	if err != nil {
		zlogger.Errorf("FollowUser findUserCacheInfo |userId:%v| err: %v", uid, err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 检查是否关注
	isFollow, err := s.userCache.IsFollow(ctx, constsR.UserFollowCacheKey, uid, followCacheInfo.Id)
	if err != nil {
		zlogger.Errorf("FollowUser UserCache.IsFollow | err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	if isFollow {
		return nil, errs.ErrFollow.Wrap("user_follow_exist")
	}

	err = s.siteDB.CreateFollow(ctx, &table.Follow{
		FollowId:  cast.ToInt(req.FollowId),
		UserId:    uid,
		Status:    constant.UserFollowFocusOn,
		CreatedAt: time.Now(),
	})

	err = s.userCache.FollowUser(ctx, uid, followCacheInfo.Id)
	if err != nil {
		zlogger.Errorf("FollowUser UserCache.FollowUser | err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	if followCacheInfo.Category == constant.UserCategoryAnchor {
		// 获取主播直播间
		RoomCacheInfo, err := s.findRoomCacheInfo(ctx, followCacheInfo.RoomId)
		if err != nil {
			zlogger.Errorf("FollowUser findRoomCacheInfo |roomId:%v|", followCacheInfo.RoomId)
			return nil, errs.ErrInternalServer.Wrap("query_failed")
		}

		if RoomCacheInfo.LiveStatus == constant.RoomSceneStatusDo {
			_, err = s.liveApiRpcClient.FollowUser(ctx, &live.FollowUserReq{
				UserId:     cast.ToInt32(uid),
				UserName:   userCacheInfo.Nickname,
				ChatRoomId: RoomCacheInfo.ChatRoomId,
			})
			if err != nil {
				zlogger.Errorf("FollowUser FollowUser err: %v", err)
				return nil, errs.ErrInternalServer.Wrap("operation_failed")
			}

			err := s.roomCache.AddRoomScene(ctx, constsR.SceneFollow, followCacheInfo.RoomId, cast.ToString(uid))
			if err != nil {
				zlogger.Errorf("FollowUser AddRoomSceneFollow | roomId:%v, sceneI的:%v,userId:%v | err: %v", followCacheInfo.RoomId, RoomCacheInfo.SceneHistoryId, uid, err)
				return nil, errs.ErrInternalServer.Wrap("operation_failed")
			}
		}
	}

	return nil, nil
}

// FollowUnlockUser 取消关注
func (s *siteSrv) FollowUnlockUser(ctx context.Context, req *pb.FollowReq) (*commonPb.Empty, error) {
	uid, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("FollowUnlockUser Check user err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	if uid == constant.Zero {
		zlogger.Errorf("FollowUser |userId:%v,followId:%v| err: visitors are prohibited from operating", uid, req.GetFollowId())
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	// 获取锁
	lockKey := fmt.Sprintf(constsR.LockUserOperate, uid)
	isLock := s.redisLock.TryLock(ctx, lockKey, lockKey)
	if !isLock {
		zlogger.Errorf("FollowUnlockUser TryLock |lockKey:%v| err: failed to acquire lock", lockKey)
		return nil, nil
	}

	defer func() {
		ret := s.redisLock.ReleaseLock(ctx, lockKey, lockKey)
		if !ret {
			zlogger.Errorf("FollowUnlockUser ReleaseLock |lockKey:%v| err: failed to release lock", lockKey)
		}
	}()

	if uid == cast.ToInt(req.GetFollowId()) {
		zlogger.Errorf("FollowUser |userId:%v,followId:%v| err: CAN To FOCUS ON MYSELF", uid, req.GetFollowId())
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	isFollow, err := s.userCache.IsFollow(ctx, constsR.UserFollowCacheKey, uid, cast.ToInt(req.GetFollowId()))
	if err != nil {
		zlogger.Errorf("FollowUser UserCache.IsFollow | err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	if !isFollow {
		return nil, errs.ErrFollow.Wrap("user_follow_not")
	}

	follow, err := s.siteDB.FindFollow(ctx, table.WhereFollowFollowId(cast.ToInt(req.GetFollowId())), table.WhereFollowUserId(uid))
	if err != nil {
		zlogger.Errorf("FollowUnlockUser FindFollow | err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	if follow.IsEmpty() {
		return nil, errs.ErrFollow.Wrap("user_follow_not")
	}

	err = s.siteDB.UpdateFollow(ctx, follow.Id, table.SetFollowStatus(constant.UserFollowUnlock))
	if err != nil {
		zlogger.Errorf("FollowUnlockUser UpdateFollow | err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	err = s.userCache.UnfollowUser(ctx, uid, cast.ToInt(req.GetFollowId()))
	if err != nil {
		zlogger.Errorf("FollowUnlockUser Unfollow | err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("operation_failed")
	}

	// 被关注者用户信息
	followCacheInfo, err := s.findUserCacheInfo(ctx, cast.ToInt(req.GetFollowId()))
	if err != nil {
		zlogger.Errorf("FollowUnlockUser UserCache.Get |followId:%v| err: %v", req.GetFollowId(), err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	if followCacheInfo.Category == constant.UserCategoryAnchor {
		// 获取主播直播间
		RoomCacheInfo, err := s.findRoomCacheInfo(ctx, followCacheInfo.RoomId)
		if err != nil {
			zlogger.Errorf("FollowUnlockUser findRoomCacheInfo |roomId:%v|", followCacheInfo.RoomId)
			return nil, errs.ErrInternalServer.Wrap("query_failed")
		}

		if RoomCacheInfo.LiveStatus == constant.RoomSceneStatusDo {
			err := s.roomCache.RemRoomScene(ctx, constsR.SceneFollow, followCacheInfo.RoomId, cast.ToString(uid))
			if err != nil {
				zlogger.Errorf("FollowUser RemRoomSceneFollow | roomId:%v, sceneI的:%v,userId:%v | err: %v", followCacheInfo.RoomId, RoomCacheInfo.SceneHistoryId, uid, err)
				return nil, errs.ErrInternalServer.Wrap("operation_failed")
			}
		}
	}

	return nil, nil
}

// FollowFans 粉丝列表
func (s *siteSrv) FollowFans(ctx context.Context, req *pb.FollowsReq) (*pb.FollowsResp, error) {
	uid, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("FollowFans CheckUser err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	if uid == 0 {
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	list, err := s.userCache.FollowList(ctx, constsR.UserFansSortedCacheKey, uid, cast.ToInt(req.GetLastId()), cast.ToInt(req.GetPageSize()))
	if err != nil {
		zlogger.Errorf("FollowFans UserCache.FollowList | err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	var resp []*pb.FollowsItem
	for _, val := range list {
		// 获取用户缓存信息
		userCacheInfo, err := s.findUserCacheInfo(ctx, cast.ToInt(val))
		if err != nil {
			zlogger.Errorf("FollowFans UserCache.Get |userId:%v| err: %v", val, err)
			continue
		}

		resp = append(resp, &pb.FollowsItem{
			UserId:   cast.ToInt32(userCacheInfo.Id),
			Nickname: userCacheInfo.Nickname,
			Avatar:   userCacheInfo.Avatar,
			Sex:      cast.ToInt32(userCacheInfo.Sex),
			LevelId:  cast.ToInt32(utils.CompareMax(userCacheInfo.SetLevelId, userCacheInfo.LevelId)),
			Sign:     userCacheInfo.Sign,
		})
	}

	return &pb.FollowsResp{List: resp}, nil
}

// FollowFollowing 关注列表
func (s *siteSrv) FollowFollowing(ctx context.Context, req *pb.FollowsReq) (*pb.FollowsResp, error) {
	uid, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("FollowFollowing CheckUser err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	if uid == 0 {
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	list, err := s.userCache.FollowList(ctx, constsR.UserFollowSortedCacheKey, uid, cast.ToInt(req.GetLastId()), cast.ToInt(req.GetPageSize()))
	if err != nil {
		zlogger.Errorf("FollowFollowing UserCache.FollowList | err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	var resp []*pb.FollowsItem
	for _, val := range list {
		// 获取用户缓存信息
		userCacheInfo, err := s.findUserCacheInfo(ctx, cast.ToInt(val))
		if err != nil {
			zlogger.Errorf("FollowFollowing UserCache.Get |userId:%v| err: %v", val, err)
			continue
		}

		resp = append(resp, &pb.FollowsItem{
			UserId:   cast.ToInt32(userCacheInfo.Id),
			Nickname: userCacheInfo.Nickname,
			Avatar:   userCacheInfo.Avatar,
			Sex:      cast.ToInt32(userCacheInfo.Sex),
			LevelId:  cast.ToInt32(utils.CompareMax(userCacheInfo.SetLevelId, userCacheInfo.LevelId)),
			Sign:     userCacheInfo.Sign,
		})
	}

	return &pb.FollowsResp{List: resp}, nil
}

// FollowCount 粉丝/关注数量
func (s *siteSrv) FollowCount(ctx context.Context, req *pb.FollowCountReq) (*pb.FollowCountResp, error) {
	uid, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("FollowCount CheckUser err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	if req.GetUserId() > 0 {
		uid = cast.ToInt(req.GetUserId())
	}

	statistics, err := s.userCache.FollowStatistics(ctx, uid)
	if err != nil {
		zlogger.Errorf("FollowCount FollowStatistics err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	return &pb.FollowCountResp{FollowingCount: cast.ToInt32(statistics.FollowCount), FansCount: cast.ToInt32(statistics.FansCount)}, nil
}

// FollowCheck 是否关注
func (s *siteSrv) FollowCheck(ctx context.Context, req *pb.FollowReq) (*pb.FollowCheckResp, error) {
	uid, err := mctx.CheckUser(ctx)
	if err != nil {
		zlogger.Errorf("FollowCheck CheckUser err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	// 检查是否关注
	isFollow, err := s.userCache.IsFollow(ctx, constsR.UserFollowCacheKey, uid, cast.ToInt(req.GetFollowId()))
	if err != nil {
		zlogger.Errorf("FollowCheck IsFollow err: %v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	var whether int32 = 1
	if !isFollow {
		whether = 0
	}

	return &pb.FollowCheckResp{Whether: whether}, nil
}
