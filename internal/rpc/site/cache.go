package site

import (
	"context"
	"errors"
	"fmt"
	"time"

	table "serverApi/pkg/db/table/site"
	"serverApi/pkg/tools/cast"
	"serverApi/pkg/zlogger"

	"serverApi/pkg/constant"
	"serverApi/pkg/db/cache"
)

// 获取用户信息
func (s *siteSrv) findUserCacheInfo(ctx context.Context, userId int) (*cache.UserCacheInfo, error) {
	userCacheInfo, err := s.userCache.Get(ctx, userId)
	if err != nil || userCacheInfo == nil {
		var (
			roomId         = 0
			isFamilyMaster = 0
			familyId       = 0
			familyMasterId = 0
			mountsId       = 0
		)

		if err != nil {
			zlogger.Errorf("findUserCacheInfo |userId:%v| err: %v", userId, err)
		}

		// 获取用户数据
		userInfo, err := s.siteDB.FindUserInfo(ctx, table.WhereUserId(userId))
		if err != nil {
			zlogger.Errorf("findUserCacheInfo FindUserInfo |userId:%v| err: %v", userId, err)
			return nil, err
		}

		if userInfo.IsEmpty() {
			zlogger.Errorf("findUserCacheInfo |userId:%v| err: user not exist", userId)
			return nil, errors.New("user not exist")
		}

		// 获取用户认证信息
		userAuth, err := s.siteDB.GetUserAuth(ctx, table.WhereUserAuthUserId(userId))
		if err != nil {
			zlogger.Errorf("findUserCacheInfo GetUserAuth |userId:%v| err: %v", userId, err)
			return nil, err
		}

		if userAuth.IsEmpty() {
			zlogger.Errorf("findUserCacheInfo |userId:%v| err: user auth not exist", userId)
			return nil, errors.New("user auth not exist")
		}

		// 获取直播间信息
		if userInfo.Category == constant.UserCategoryAnchor {
			roomInfo, err := s.siteDB.FindRoomInfo(ctx, table.WhereRoomUserId(userId))
			if err != nil {
				zlogger.Errorf("findUserCacheInfo FindRoomInfo |userId:%v| err: %v", userId, err)
				return nil, err
			}

			roomId = roomInfo.Id
		}

		// 家族信息
		familyInfo, err := s.siteDB.FindFamily(ctx, table.WhereFamilyUserId(userId))
		if err != nil {
			zlogger.Errorf("findUserCacheInfo FindFamily |userId:%v| err: %v", userId, err)
			return nil, err
		}

		// 家族长
		if !familyInfo.IsEmpty() {
			isFamilyMaster = constant.Yes
			familyId = familyInfo.Id
			familyMasterId = userId
		}

		// 是否加入家族
		familyAnchor, err := s.siteDB.FindFamilyAnchors(ctx, table.WhereFamilyAnchorsUserId(userId))
		if err != nil {
			zlogger.Errorf("findUserCacheInfo FindFamilyAnchors |userId:%v| err: %v", userId, err)
			return nil, err
		}
		if !familyAnchor.IsEmpty() {
			familyId = familyAnchor.FamilyId
			familyMasterId = familyAnchor.PatriarchId
		}

		// 查询坐骑信息
		mountInfo, err := s.siteDB.FindUserUseMount(ctx,
			table.WhereUserMountsUserId(userId),
			table.WhereUserMountsIsSelected(constant.Yes),
			table.WhereUserMountsExpiredTime(time.Now()),
		)
		if err != nil {
			zlogger.Errorf("findUserCacheInfo FindUserUseMount |userId:%v| err: %v", userId, err)
			return nil, err
		}

		if !mountInfo.IsEmpty() {
			mountsId = mountInfo.MountsId
		}

		var userCache = cache.NewUserCacheInfo(
			userInfo.Id,
			userInfo.CountryCode,
			userAuth.AreaCode,
			userAuth.Mobile,
			userAuth.Email,
			userInfo.Nickname,
			userInfo.Avatar,
			userInfo.Sign,
			userInfo.Sex,
			userInfo.Birthday,
			userInfo.Feeling,
			userInfo.Country,
			userInfo.Area,
			userInfo.Profession,
			userInfo.Category,
			userInfo.InviteCode,
			userInfo.ParentId,
			userInfo.LevelId,
			userInfo.SetLevelId,
			userInfo.Remark,
			userInfo.Status,
			userAuth.Password,
			userInfo.PayPassword,
			roomId,
			userInfo.ChatUuid,
			userInfo.GmStatus,
			isFamilyMaster,
			familyId,
			familyMasterId,
			mountsId,
		)

		err = s.userCache.Save(ctx, userId, userCache)
		if err != nil {
			zlogger.Errorf("findUserCacheInfo userCache.Save |userId:%v| err: %v", userId, err)
		}

		return userCache, nil
	}

	return userCacheInfo, nil
}

// 获取直播间信息
func (s *siteSrv) findRoomCacheInfo(ctx context.Context, roomId int) (*cache.RoomCacheInfo, error) {
	roomCacheInfo, err := s.roomCache.Get(ctx, roomId)
	if err != nil || roomCacheInfo == nil {
		if err != nil {
			zlogger.Errorf("findRoomCacheInfo |roomId:%v| err: %v", roomId, err)
		}

		// 获取直播间信息
		roomInfo, err := s.siteDB.FindRoomInfo(ctx, table.WhereRoomId(roomId))
		if err != nil {
			zlogger.Errorf("findRoomCacheInfo FindRoomInfo |roomId:%v| err: %v", roomId, err)
			return nil, err
		}

		if roomInfo.IsEmpty() {
			return nil, fmt.Errorf("findRoomCacheInfo |roomId:%v| roomInfo not exist", roomId)
		}

		// 开播信息
		history, err := s.siteDB.FindRoomSceneHistory(ctx, table.WhereRoomSceneHistoryRoomId(roomId), table.WhereRoomSceneHistoryStatus(constant.RoomSceneStatusDo))
		if err != nil {
			zlogger.Errorf("findRoomCacheInfo FindRoomSceneHistory |roomId:%v| err: %v", roomId, err)
			return nil, err
		}

		sceneHistoryId := constant.Zero
		startLiveTime := constant.Zero
		payRules := constant.RoomChargingRulesFree
		trialDuration := constant.Zero
		unitPrice := constant.Zero
		if !history.IsEmpty() {
			sceneHistoryId = history.Id
			payRules = history.PayRules
			trialDuration = history.TrialDuration
			unitPrice = history.UnitPrice
			startLiveTime = cast.ToInt(history.OpenAt.Unix())
		}

		roomCache := cache.NewRoomCacheInfo(
			roomInfo.Id,
			roomInfo.CountryCode,
			roomInfo.UserId,
			roomInfo.Title,
			roomInfo.TagsJson,
			roomInfo.Cover,
			roomInfo.VideoClarity,
			roomInfo.GiftRatio,
			roomInfo.PlatformRatio,
			roomInfo.FamilyRatio,
			roomInfo.Status,
			roomInfo.ChatRoomId,
			roomInfo.LiveStatus,
			roomInfo.Summary,
			sceneHistoryId,
			constant.Zero,
			startLiveTime,
			roomInfo.PaidPurviewStatus,
			roomInfo.LevelId,
			roomInfo.SetLevelId,
			payRules,
			trialDuration,
			unitPrice,
		)

		err = s.roomCache.Save(ctx, roomId, roomCache)
		if err != nil {
			zlogger.Errorf("findRoomCacheInfo roomCache.Save |roomId:%v| err: %v", roomId, err)
		}

		return roomCache, nil
	}

	return roomCacheInfo, nil
}
