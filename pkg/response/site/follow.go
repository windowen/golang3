package site

type (
	FollowCountResp struct {
		FollowingCount int32 `json:"followingCount"`
		FansCount      int32 `json:"fansCount"`
	}

	FollowsResp struct {
		List []FollowsItem `json:"list"`
	}

	FollowsItem struct {
		UserId   int32  `json:"userId"`
		Nickname string `json:"nickname"`
		Avatar   string `json:"avatar"`
		Sex      int32  `json:"sex"`
		LevelId  int32  `json:"levelId"`
		Sign     string `json:"sign"`
	}
)
