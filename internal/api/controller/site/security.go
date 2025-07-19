package site

import (
	"github.com/gin-gonic/gin"

	commonPb "serverApi/pkg/protobuf/common"
	"serverApi/pkg/protobuf/site"
	siteRes "serverApi/pkg/response/site"
	"serverApi/pkg/tools/apiresp"
)

// SecurityAuthInfo 认证信息
func (api *ApiSite) SecurityAuthInfo(c *gin.Context) {
	info, err := api.siteRpcClient.SecurityAuthInfo(c, &commonPb.Empty{})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, siteRes.SecurityAuthInfoResp{
		AreaCode: info.GetAreaCode(),
		Mobile:   info.GetMobile(),
		Email:    info.GetEmail(),
	})
}

// SecurityBinding 绑定
func (api *ApiSite) SecurityBinding(c *gin.Context) {
	var req site.SecurityBindingReq

	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	if err := apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	_, err := api.siteRpcClient.SecurityBinding(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, nil)
}
