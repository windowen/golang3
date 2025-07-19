package site

type RedPointResp struct {
	IsShowRedPoint bool  `json:"isShowRedPoint"`
	UnreadCount    int32 `json:"unreadCount"`
}

// MsgListItem 消息列表项
type MsgListItem struct {
	Id         int32  `json:"id"`
	Content    string `json:"content"`
	CreatedAt  string `json:"createdAt"`
	ReadStatus int32  `json:"readStatus"`
}

// MsgListResp 消息列表响应
type MsgListResp struct {
	TotalRecord int32         `json:"totalRecord"`
	TotalPage   int32         `json:"totalPage"`
	PageSize    int32         `json:"pageSize"`
	PageNum     int32         `json:"pageNum"`
	List        []MsgListItem `json:"list"`
}

// RotatingReq 跑马灯请求
type RotatingReq struct {
	RotatingType int32 `json:"rotatingType"`
}

// RotatingItem 跑马灯项
type RotatingItem struct {
	Id   int32  `json:"id"`
	Name string `json:"name"`
}

// RotatingResp 跑马灯响应
type RotatingResp struct {
	List []RotatingItem `json:"list"`
}
