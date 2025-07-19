package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"serverApi/pkg/tools/cast"

	"github.com/redis/go-redis/v9"

	constsR "serverApi/pkg/constant/redis"
	"serverApi/pkg/zlogger"
)

type RoomCache struct {
	rdb     redis.UniversalClient // Redis客户端
	expires time.Duration         // 缓存过期时间(默认0不过期)
}

// NewRoomCache 创建一个 RoomCache 实例
func NewRoomCache(rdb redis.UniversalClient) *RoomCache {
	return &RoomCache{
		rdb:     rdb,
		expires: 0 * time.Second,
	}
}

type RoomCacheInfo struct {
	Id                int    `json:"id"`
	CountryCode       string `json:"countryCode"`
	UserId            int    `json:"userId"`
	Title             string `json:"title"`
	Tags              string `json:"tags"` // 直播间标签
	Cover             string `json:"cover"`
	VideoClarity      int    `json:"videoClarity"`
	PayRules          int    `json:"payRules"`
	TrialDuration     int    `json:"trialDuration"`
	UnitPrice         int    `json:"unitPrice"`
	GiftRatio         int    `json:"giftRatio"`
	PlatformRatio     int    `json:"platformRatio"`
	FamilyRatio       int    `json:"familyRatio"`
	Status            int    `json:"status"`
	ChatRoomId        string `json:"chatRoomId"` // 是播间临时聊天室id
	LiveStatus        int    `json:"liveStatus"`
	Summary           string `json:"summary"`        // 简介
	SceneHistoryId    int    `json:"sceneHistoryId"` // 场次id
	GameId            int    `json:"gameId"`
	StartLiveTime     int    `json:"startLiveTime"`
	PaidPurviewStatus int    `json:"paidPurviewStatus"` // 付费权限
	IsOpenVibrator    int    `json:"isOpenVibrator"`    // 是否开启跳蛋
	Sort              int    `json:"sort"`              // 置顶排序
	BottomSort        int    `json:"bottomSort"`        // 置底排序
	LastStartLiveTime int    `json:"lastStartLiveTime"` // 最后开播时间
	LastEndLiveTime   int    `json:"lastEndLiveTime"`   // 最后下播时间
	LevelId           int    `json:"levelId"`           // 自然等级
	SetLevelId        int    `json:"setLevelId"`        // 后台设置等级
}

func NewRoomCacheInfo(
	Id int,
	countryCode string,
	userId int,
	title string,
	tags string,
	cover string,
	videoClarity int,
	giftRatio int,
	platformRatio int,
	familyRatio int,
	status int,
	chatRoomId string,
	liveStatus int,
	summary string,
	sceneHistoryId int,
	gameId int,
	startLiveTime int,
	paidPurviewStatus int,
	levelId int,
	setLevelId int,
	payRules int,
	trialDuration int,
	unitPrice int,
) *RoomCacheInfo {
	return &RoomCacheInfo{
		Id:                Id,
		CountryCode:       countryCode,
		UserId:            userId,
		Title:             title,
		Tags:              tags,
		Cover:             cover,
		VideoClarity:      videoClarity,
		GiftRatio:         giftRatio,
		PlatformRatio:     platformRatio,
		FamilyRatio:       familyRatio,
		Status:            status,
		ChatRoomId:        chatRoomId,
		LiveStatus:        liveStatus,
		Summary:           summary,
		SceneHistoryId:    sceneHistoryId,
		GameId:            gameId,
		StartLiveTime:     startLiveTime,
		PaidPurviewStatus: paidPurviewStatus,
		LevelId:           levelId,
		SetLevelId:        setLevelId,
		PayRules:          payRules,
		TrialDuration:     trialDuration,
		UnitPrice:         unitPrice,
	}
}

// Save 将房间信息保存到 Redis 中
func (rc *RoomCache) Save(ctx context.Context, roomId int, value *RoomCacheInfo) error {
	cacheKey := fmt.Sprintf(constsR.RoomCacheInfoKey, cast.ToString(roomId))

	marshal, err := json.Marshal(value)
	if err != nil {
		zlogger.Errorf("marshal user cache info error: %v", err)
		return fmt.Errorf("failed to save user data to redis: %v", err)
	}

	err = rc.rdb.Set(ctx, cacheKey, marshal, rc.expires).Err()
	if err != nil {
		zlogger.Errorf("failed to save room data to redis failed, err: %v", err)
		return fmt.Errorf("failed to save room data to redis: %v", err)
	}

	return nil
}

// Update 更新 Redis 中的房间信息
func (rc *RoomCache) Update(ctx context.Context, roomId int, value *RoomCacheInfo) error {
	return rc.Save(ctx, roomId, value)
}

// Delete 删除 Redis 中的房间信息
func (rc *RoomCache) Delete(ctx context.Context, roomId int) error {
	err := rc.rdb.Del(ctx, fmt.Sprintf(constsR.RoomCacheInfoKey, cast.ToString(roomId))).Err()
	if err != nil {
		zlogger.Errorf("deletion of room data failed, err: %v", err)
		return fmt.Errorf("deletion of room data failed: %v", err)
	}
	return nil
}

// Get 从 Redis 中获取房间信息
func (rc *RoomCache) Get(ctx context.Context, roomId int) (*RoomCacheInfo, error) {
	var roomCacheInfo *RoomCacheInfo

	data, err := rc.rdb.Get(ctx, fmt.Sprintf(constsR.RoomCacheInfoKey, cast.ToString(roomId))).Result()
	if err := CheckErr(err); err != nil {
		zlogger.Errorf("roomCacheInfo Get | err: %v", err)
		return nil, fmt.Errorf("failed to get room data db redis: %v", err)
	}

	if data == "" {
		zlogger.Errorf("roomCacheInfo | err: failed to get room cache")
		return nil, errors.New("failed to get room cache")
	}

	if err := json.Unmarshal([]byte(data), &roomCacheInfo); err != nil {
		zlogger.Errorf("roomCacheInfo json.Unmarshal |data:%v| err: %v", data, err)
		return nil, fmt.Errorf("unable to deserialize room data: %v", err)
	}

	return roomCacheInfo, nil
}

// Renew 更新缓存有效期
func (rc *RoomCache) Renew(ctx context.Context, roomId int) error {
	err := rc.rdb.Expire(ctx, fmt.Sprintf(constsR.RoomCacheInfoKey, cast.ToString(roomId)), rc.expires).Err()
	if err != nil {
		zlogger.Errorf("failed to set room %v expiration time, err: %v", roomId, err)
		return fmt.Errorf("failed to set room expiration time, err: %v", err)
	}

	return nil
}

// RoomSceneIncome 直播间场收益
func (rc *RoomCache) RoomSceneIncome(ctx context.Context, roomId, income int) (int64, error) {
	return rc.incrInt(ctx, fmt.Sprintf(constsR.SceneIncome, roomId), cast.ToInt64(income))
}

// incrInt 整数自增
func (rc *RoomCache) incrInt(ctx context.Context, key string, value int64) (int64, error) {
	return rc.rdb.IncrBy(ctx, key, value).Result()
}

// decrInt 整数自减
func (rc *RoomCache) decrInt(ctx context.Context, key string, value int64) (int64, error) {
	return rc.rdb.DecrBy(ctx, key, value).Result()
}

// AddStartBroadcast 新增开播直播间
func (rc *RoomCache) AddStartBroadcast(ctx context.Context, roomId int) error {
	return rc.rdb.SAdd(ctx, constsR.RoomStartLiveSet, roomId).Err()
}

// RemStartBroadcast 移除开播
func (rc *RoomCache) RemStartBroadcast(ctx context.Context, roomId int) error {
	return rc.rdb.SRem(ctx, constsR.RoomStartLiveSet, roomId).Err()
}

type RoomUpdateOption func(*RoomCacheInfo)

// UpdateRoomOption 更新部分缓存
func (rc *RoomCache) UpdateRoomOption(ctx context.Context, roomId int, opts ...RoomUpdateOption) error {
	// 获取现有的 RoomCacheInfo
	roomInfo, err := rc.Get(ctx, roomId)
	if err != nil {
		return err
	}

	// 应用更新选项
	for _, opt := range opts {
		opt(roomInfo)
	}

	return rc.Save(ctx, roomId, roomInfo)
}

func UpdateRoomTitle(title string) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.Title = title
	}
}

func UpdateRoomCover(cover string) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.Cover = cover
	}
}

func UpdateRoomVideoClarity(videoClarity int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.VideoClarity = videoClarity
	}
}

func UpdateRoomPayRules(payRules int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.PayRules = payRules
	}
}

func UpdateRoomTrialDuration(trialDuration int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.TrialDuration = trialDuration
	}
}

func UpdateRoomUnitPrice(unitPrice int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.UnitPrice = unitPrice
	}
}

func UpdateRoomGiftRatio(giftRatio int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.GiftRatio = giftRatio
	}
}

func UpdateRoomFamilyRatio(familyRatio int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.FamilyRatio = familyRatio
	}
}

func UpdateRoomLiveStatus(status int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.LiveStatus = status
	}
}

func UpdateRoomSummary(summary string) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.Summary = summary
	}
}

func UpdateRoomSceneHistoryId(sceneHistoryId int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.SceneHistoryId = sceneHistoryId
	}
}

func UpdateRoomGameId(gameId int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.GameId = gameId
	}
}

func UpdateRoomStartLiveTime(startLiveTime int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.StartLiveTime = startLiveTime
	}
}

func UpdateRoomIsOpenVibrator(isOpenVibrator int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.IsOpenVibrator = isOpenVibrator
	}
}

func UpdateRoomSort(sort int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.Sort = sort
	}
}

func UpdateRoomBottomSort(bottomSort int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.BottomSort = bottomSort
	}
}

func UpdateRoomLastStartLiveTime(time int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.LastStartLiveTime = time
	}
}

func UpdateRoomLastEndLiveTime(time int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.LastEndLiveTime = time
	}
}

func UpdateRoomPaidPurviewStatus(paidPurviewStatus int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.PaidPurviewStatus = paidPurviewStatus
	}
}

func UpdateRoomChatRoomId(chatRoomId string) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.ChatRoomId = chatRoomId
	}
}
