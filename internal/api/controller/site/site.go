package site

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"serverApi/pkg/common/config"
	"serverApi/pkg/gozero/discov"
	"serverApi/pkg/gozero/zrpc"
	"serverApi/pkg/protobuf/common"
	"serverApi/pkg/protobuf/site"
	siteRes "serverApi/pkg/response/site"
	"serverApi/pkg/tools/apiresp"
	"serverApi/pkg/tools/mw"
)

func NewSiteApi() *ApiSite {
	rpcKey := strings.ToLower(fmt.Sprintf("%s:///%s", config.Config.Etcd.Schema, config.Config.RpcName.SiteRpcName))
	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: config.Config.Etcd.Addr,
			Key:   rpcKey,
		},
	}, zrpc.WithDialOption(mw.GrpcClient()), zrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))

	return &ApiSite{
		siteRpcClient: site.NewSiteServiceClient(client.Conn()),
	}
}

type ApiSite struct {
	siteRpcClient site.SiteServiceClient
}

// Login 登录
func (api *ApiSite) Login(c *gin.Context) {
	var req site.LoginReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	if err := apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := api.siteRpcClient.Login(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// Register 注册
func (api *ApiSite) Register(c *gin.Context) {
	var (
		err error
		req site.RegisterReq
	)
	if err = c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	if err = apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	_, err = api.siteRpcClient.Register(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, nil)
}

// GetProfile 获取用户信息
func (api *ApiSite) GetProfile(c *gin.Context) {
	var (
		err error
		req site.ProfileReq
	)

	if err = c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	if err = apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	resp, err := api.siteRpcClient.GetProfile(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, siteRes.UserInfoResp{
		ID:             resp.UserId,
		CountryCode:    resp.CountryCode,
		Avatar:         resp.Avatar,
		Nickname:       resp.Nickname,
		Sex:            resp.Sex,
		Birthday:       resp.Birthday,
		Feeling:        resp.Feeling,
		Country:        resp.Country,
		Area:           resp.Area,
		Sign:           resp.Sign,
		LevelID:        resp.LevelId,
		Category:       resp.Category,
		InviteCode:     resp.InviteCode,
		RoomId:         resp.RoomId,
		LiveStatus:     resp.LiveStatus,
		ChatUUID:       resp.ChatUUID,
		Profession:     resp.Profession,
		IsFamilyMaster: resp.IsFamilyMaster,
		IsFollow:       resp.IsFollow,
	})
}

// ModifyPaymentPassword 修改支付密码
func (api *ApiSite) ModifyPaymentPassword(c *gin.Context) {
	var req site.ModifyPaymentReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	if err := apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := api.siteRpcClient.ModifyPaymentPassword(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// ModifyLoginPassword 修改登陆密码
func (api *ApiSite) ModifyLoginPassword(c *gin.Context) {
	var req site.ModifyLoginPasswordReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	if err := apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := api.siteRpcClient.ModifyLoginPassword(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// Logout 登出
func (api *ApiSite) Logout(c *gin.Context) {
	resp, err := api.siteRpcClient.Logout(c, &common.Empty{})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// SendValidationCode 发送验证码消息
func (api *ApiSite) SendValidationCode(c *gin.Context) {
	var req site.SendValidationReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	if err := apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := api.siteRpcClient.SendValidationCode(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// SendCertifiedValidationCode 发送验证码消息-已登陆
func (api *ApiSite) SendCertifiedValidationCode(c *gin.Context) {
	var req site.SendValidationCertifiedReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	if err := apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	resp, err := api.siteRpcClient.SendValidationCode(c, &site.SendValidationReq{
		AccountType: req.AccountType,
		Mobile:      req.Mobile,
		AreaCode:    req.AreaCode,
		Email:       req.Email,
		Scene:       req.Scene,
	})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// ForgetPassword 忘记密码
func (api *ApiSite) ForgetPassword(c *gin.Context) {
	var req site.ForgetPasswordReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	if err := apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	resp, err := api.siteRpcClient.ForgetPassword(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}

// GlobalAreas 获取国家相关信息
func (api *ApiSite) GlobalAreas(c *gin.Context) {
	resp, err := api.siteRpcClient.GlobalAreas(c, &common.Empty{})
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, resp)
}

func (api *ApiSite) ModifyProfile(c *gin.Context) {
	var req site.ModifyProfileReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	if err := apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	_, err := api.siteRpcClient.ModifyProfile(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
	}

	apiresp.GinSuccess(c, nil)
}

// SiteConfigs 系统站点配置
func (api *ApiSite) SiteConfigs(c *gin.Context) {
	resp, err := api.siteRpcClient.SiteConfigs(c, &common.Empty{})
	if err != nil {
		apiresp.GinError(c, err)
	}

	apiresp.GinSuccess(c, resp)
}

// SiteCarousel 站点轮播图
func (api *ApiSite) SiteCarousel(c *gin.Context) {
	var req site.SiteBannerReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	if err := apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	resp, err := api.siteRpcClient.SiteBanner(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
	}

	var res []siteRes.BannerItem
	for _, val := range resp.List {
		res = append(res, siteRes.BannerItem{
			Uri:      val.Uri,
			ShowType: val.ShowType,
			ExtInfo:  val.ExtInfo,
		})
	}

	apiresp.GinSuccess(c, siteRes.BannerResp{List: res})
}

// StartupImage 启动图
func (api *ApiSite) StartupImage(c *gin.Context) {
	resp, err := api.siteRpcClient.StartupImage(c, &common.Empty{})
	if err != nil {
		apiresp.GinError(c, err)
	}

	apiresp.GinSuccess(c, resp)
}

// Stay 页面停留时长
func (api *ApiSite) Stay(c *gin.Context) {
	var req site.PageStayReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	resp, err := api.siteRpcClient.Stay(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
	}

	apiresp.GinSuccess(c, resp)
}

// BannerClick 页面点击
func (api *ApiSite) BannerClick(c *gin.Context) {
	var req site.BannerClickReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}

	resp, err := api.siteRpcClient.BannerClick(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
	}

	apiresp.GinSuccess(c, resp)
}
