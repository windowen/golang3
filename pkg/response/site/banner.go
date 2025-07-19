package site

type (
	BannerResp struct {
		List []BannerItem `json:"list"`
	}

	BannerItem struct {
		Uri      string `json:"uri"`
		ShowType int32  `json:"showType"`
		ExtInfo  string `json:"extInfo"`
	}
)
