package site

type (
	UserInfoResp struct {
		ID             int32  `json:"id"`          // ID
		CountryCode    string `json:"countryCode"` // 国家编码
		Avatar         string `json:"avatar"`      // 用户头像
		Nickname       string `json:"nickname"`    // 用户昵称
		Sex            int32  `json:"sex"`         // 性别 1-男 2-女
		Birthday       string `json:"birthday"`    // 生日
		Feeling        int32  `json:"feeling"`     // 感情
		Country        string `json:"country"`     // 国家
		Area           string `json:"area"`        // 地区
		Sign           string `json:"sign"`        // 签名
		LevelID        int32  `json:"levelId"`     // VIP等级ID
		Category       int32  `json:"category"`    // 用户类型
		InviteCode     string `json:"inviteCode"`  // 用户邀请码
		RoomId         int32  `json:"roomId"`      // 直播间id
		LiveStatus     int32  `json:"liveStatus"`  // 直播状态
		ChatUUID       string `json:"chatUUID"`
		Profession     int32  `json:"profession"`     // 职业
		IsFamilyMaster int32  `json:"isFamilyMaster"` // 是否家族长
		IsFollow       int32  `json:"isFollow"`       // 是否关注
	}
)
