package site

import (
	"github.com/gin-gonic/gin"
	"serverApi/pkg/protobuf/site"
	siteRes "serverApi/pkg/response/site"
	"serverApi/pkg/tools/apiresp"
)

func (api *ApiSite) FollowFocusOn(c *gin.Context) {
	var req site.FollowReq

	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	if err := apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	_, err := api.siteRpcClient.FollowUser(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, nil)
}

func (api *ApiSite) FollowUnlock(c *gin.Context) {
	var req site.FollowReq

	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	if err := apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	_, err := api.siteRpcClient.FollowUnlockUser(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, nil)
}

func (api *ApiSite) FollowFans(c *gin.Context) {
	var req site.FollowsReq

	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	if err := apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	fans, err := api.siteRpcClient.FollowFans(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	var res []siteRes.FollowsItem
	for _, val := range fans.List {
		res = append(res, siteRes.FollowsItem{
			UserId:   val.UserId,
			Nickname: val.Nickname,
			Avatar:   val.Avatar,
			Sex:      val.Sex,
			LevelId:  val.LevelId,
			Sign:     val.Sign,
		})
	}

	apiresp.GinSuccess(c, siteRes.FollowsResp{List: res})
}

func (api *ApiSite) FollowFollowing(c *gin.Context) {
	var req site.FollowsReq

	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	if err := apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	fans, err := api.siteRpcClient.FollowFollowing(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	var res []siteRes.FollowsItem
	for _, val := range fans.List {
		res = append(res, siteRes.FollowsItem{
			UserId:   val.UserId,
			Nickname: val.Nickname,
			Avatar:   val.Avatar,
			Sex:      val.Sex,
			LevelId:  val.LevelId,
			Sign:     val.Sign,
		})
	}

	apiresp.GinSuccess(c, siteRes.FollowsResp{List: res})
}

func (api *ApiSite) FollowStatistical(c *gin.Context) {
	var req site.FollowCountReq

	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	if err := apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	statistical, err := api.siteRpcClient.FollowCount(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, siteRes.FollowCountResp{
		FollowingCount: statistical.FollowingCount,
		FansCount:      statistical.FansCount,
	})
}
