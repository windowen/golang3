package router

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"serverApi/internal/api/controller/site"
	"serverApi/pkg/tools/mw"
	"serverApi/pkg/tools/upload"
)

func InitRouter(router gin.IRouter, rdb redis.UniversalClient) {

	// 站点服务
	siteApi := site.NewSiteApi()

	// 公共接口分组
	siteRouter := router.Group("/api/site")
	siteRouter.POST("/login", siteApi.Login)                       // 登陆
	siteRouter.POST("/register", siteApi.Register)                 // 注册
	siteRouter.POST("/sendVerifyCode", siteApi.SendValidationCode) // 发送验证码(未登录)
	siteRouter.POST("/forgetPassword", siteApi.ForgetPassword)     // 忘记密码
	siteRouter.POST("/globalAreas", siteApi.GlobalAreas)           // 地区码
	siteRouter.POST("/carousel", siteApi.SiteCarousel)             // 站点轮播图
	siteRouter.POST("/configs", siteApi.SiteConfigs)               // 站点系统配置
	siteRouter.POST("/rotating", siteApi.Rotating)                 // 跑马灯
	siteRouter.POST("/startupImage", siteApi.StartupImage)         // 启动图
	siteRouter.POST("/stay", siteApi.Stay)                         // 页面停留时长
	siteRouter.POST("/bannerClick", siteApi.BannerClick)           // 页面点击

	// 鉴权接口分组

	authCommonRouter := siteRouter.Group("/", mw.UserAuthCheck(rdb))
	authCommonRouter.POST("/upload", upload.GinHandler)                                // 上传
	authCommonRouter.POST("/auth/sendVerifyCode", siteApi.SendCertifiedValidationCode) // 发送验证码

	// 用户接口
	authUserRouter := siteRouter.Group("/user", mw.UserAuthCheck(rdb))
	authUserRouter.POST("/profile", siteApi.GetProfile)                          // 获取个人信息
	authUserRouter.POST("/modifyProfile", siteApi.ModifyProfile)                 // 修改个人信息
	authUserRouter.POST("/modifyPaymentPassword", siteApi.ModifyPaymentPassword) // 修改支付密码
	authUserRouter.POST("/modifyLoginPassword", siteApi.ModifyLoginPassword)     // 修改密码
	authUserRouter.POST("/logout", siteApi.Logout)                               // 登出

	// 账号安全
	authUserRouter.POST("/security/authInfo", siteApi.SecurityAuthInfo) // 认证信息
	authUserRouter.POST("/security/binding", siteApi.SecurityBinding)   // 绑定

	authUserRouter.POST("/follow/focusOn", siteApi.FollowFocusOn)         // 关注
	authUserRouter.POST("/follow/unlock", siteApi.FollowUnlock)           // 取关
	authUserRouter.POST("/follow/fans", siteApi.FollowFans)               // 获取粉丝列表
	authUserRouter.POST("/follow/following", siteApi.FollowFollowing)     // 获取关注列表
	authUserRouter.POST("/follow/statistical", siteApi.FollowStatistical) // 获取关注/粉丝数量

	// 平台消息、活动消息、交易消息
	authUserRouter.POST("/msg/redPoint", siteApi.RedPoint)                                     // 消息红点以及未读条数
	authUserRouter.POST("/msg/readAll", siteApi.ReadAll)                                       // 一键已读
	authUserRouter.POST("/msg/msgList", siteApi.MsgList)                                       // 交易消息列表
	authUserRouter.POST("/msg/markRead", siteApi.MarkRead)                                     // 单条交易消息标记已读
	authUserRouter.POST("/msg/clearTransactionMsgReadAll", siteApi.ClearTransactionMsgReadAll) // 清除交易已读消息
	authUserRouter.POST("/msg/msgSummary", siteApi.MsgSummary)                                 // 消息汇总

	// 代理服务
	// agentApi := agent.NewAgent(discov)
	// routerAgent := router.Group("/api/agent")
	// routerAgent.POST("/info", agentApi.AgentInfo)

}
