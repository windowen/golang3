package redis

const (
	RoomCacheInfoKey          = "room_cache_info_%v"           // 直播间信息缓存key
	RoomManageCacheKey        = "room_manage_cache_%v"         // 直播间房管缓存key
	RoomCacheOnlineKey        = "room_cache_online_%v"         // 直播间在线用户缓存key
	RoomCacheLiveKey          = "room_cache_like_%v"           // 点赞
	RoomCacheStatKey          = "room_cache_stat_%v"           // 房间数据统计
	RoomCacheChatHisKey       = "room_cache_chat_history_%v"   // 聊天历史
	RoomCacheChatChatRooms    = "room_cache_chat_rooms"        // 房间缓存
	RoomCacheChatBarrageKey   = "room_cache_barrage_%v"        // 弹幕消息
	RoomCacheKickOutKey       = "room_cache_kick_out_%v"       // T出
	RoomCacheDefaultChatRooms = "room_cache_default_chat_room" // 房间缓存
	RoomStartLiveSet          = "room_start_live"              // 开播直播间缓存
)

// 场数据
const (
	SceneFollow      = "room_scene_follow_%v"        // 直播间场关注数据缓存key  v1-房间id
	SceneTrial       = "room_scene_trial_%v"         // 直播间试看信息缓存key   v1-房间id
	SceneIncome      = "room_scene_income_%v"        // 直播间场收益数据缓存key  v1-房间id
	SceneMute        = "room_scene_mute_%v"          // 直播间场禁言数据缓存key  v1-房间id
	SceneKickOut     = "room_scene_kick_out_%v"      // 直播间场踢出数据缓存key  v1-房间id
	SceneBlock       = "room_scene_block_%v"         // 直播间拉黑数据缓存key  v1-房间id
	ScenePayUsers    = "room_scene_pay_users_%v"     // 直播间付费用户数据缓存key  v1-房间id
	ScenePayUsersSet = "room_scene_pay_users_set_%v" // 直播间付费用户数据缓存key  v1-房间id
)
