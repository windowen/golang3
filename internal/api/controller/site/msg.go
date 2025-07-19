package site

import (
	"github.com/gin-gonic/gin"

	"serverApi/pkg/protobuf/common"
	"serverApi/pkg/protobuf/site"
	siteRes "serverApi/pkg/response/site"
	"serverApi/pkg/tools/apiresp"
)

// RedPoint 消息红点以及未读条数
func (api *ApiSite) RedPoint(c *gin.Context) {
	var req common.Empty
	redPoint, err := api.siteRpcClient.RedPoint(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, siteRes.RedPointResp{
		IsShowRedPoint: redPoint.IsShowRedPoint,
		UnreadCount:    redPoint.UnreadCount,
	})
}

// ReadAll 一键已读
func (api *ApiSite) ReadAll(c *gin.Context) {
	var req common.Empty
	redPoint, err := api.siteRpcClient.ReadAll(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, redPoint)
}

// MsgList 消息列表
func (api *ApiSite) MsgList(c *gin.Context) {
	var req site.MsgListReq

	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	msgResp, err := api.siteRpcClient.MsgList(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	var list []siteRes.MsgListItem
	for _, item := range msgResp.List {
		msg := siteRes.MsgListItem{
			Id:        item.Id,
			Content:   item.Content,
			CreatedAt: item.CreatedAt,
		}
		if req.MsgType == 2 {
			msg.ReadStatus = item.ReadStatus
		}

		list = append(list, msg)
	}

	apiresp.GinSuccess(c, siteRes.MsgListResp{
		TotalRecord: msgResp.TotalRecord,
		TotalPage:   msgResp.TotalPage,
		PageSize:    msgResp.PageSize,
		PageNum:     msgResp.PageNum,
		List:        list,
	})
}

// MarkRead 多条交易消息标记已读
func (api *ApiSite) MarkRead(c *gin.Context) {
	var req site.MarkReadReq

	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	resp, err := api.siteRpcClient.MarkRead(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, resp)
}

// Rotating 跑马灯
func (api *ApiSite) Rotating(c *gin.Context) {
	var req site.RotatingReq

	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	resp, err := api.siteRpcClient.Rotating(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	var list []siteRes.RotatingItem
	for _, item := range resp.List {
		list = append(list, siteRes.RotatingItem{
			Id:   item.Id,
			Name: item.Name,
		})
	}

	apiresp.GinSuccess(c, siteRes.RotatingResp{
		List: list,
	})
}

// ClearTransactionMsgReadAll 清除交易已读消息
func (api *ApiSite) ClearTransactionMsgReadAll(c *gin.Context) {
	var req common.Empty
	redPoint, err := api.siteRpcClient.ClearTransactionMsgReadAll(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, redPoint)
}

// MsgSummary 消息汇总
func (api *ApiSite) MsgSummary(c *gin.Context) {
	// var req site.MsgSummaryReq
	var req common.Empty
	msgSummary, err := api.siteRpcClient.MsgSummary(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, msgSummary)
}
