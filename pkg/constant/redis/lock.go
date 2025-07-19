package redis

const (
	LockSendGift      = "lock_send_gift_%v"        // 送礼物锁 v1-用户id
	LockAnchorOperate = "lock_anchor_operate_%v"   // 主播操作锁 v1-用户id
	LockUserOperate   = "lock_user_operate_%v"     // 用户操作锁 v1-用户id
	LockLiveCallBack  = "lock_live_callback_%v_%v" // 回调操作锁 v1-事件类型  v2-直播间id
)
