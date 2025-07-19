package main

import (
	"github.com/google/uuid"
	"go.uber.org/zap"

	"serverApi/internal/rpc/site"
	"serverApi/pkg/common/config"
	"serverApi/pkg/rocketmq"
	siteRpc "serverApi/pkg/rpcstart"
	"serverApi/pkg/service"
	"serverApi/pkg/tools/component"
	"serverApi/pkg/tools/i18nh"
	"serverApi/pkg/tools/utils"
	"serverApi/pkg/zlogger"
)

func main() {
	configFile, logFile, err := config.FlagParse("site")
	if err != nil {
		panic(err)
	}

	if err = config.InitConfig(configFile); err != nil {
		panic(err)
	}

	// 设置全局时区为北京时间
	utils.SetGlobalTimeZone(utils.GetBjTimeLoc())

	// 初始化国际化
	i18nh.New()

	// 初始化日志
	svcID := uuid.NewString()
	zlogger.SetGlobalFields(map[string]string{
		"ssid": svcID,
		"snm":  "site",
	})
	zlogger.InitLogConfig(logFile)

	// etcd检测
	if err = component.ComponentCheck(configFile, config.Config.App.Discovery, true); err != nil {
		panic(err)
	}

	// 初始化rocketmq
	rocketmq.Init()

	// 站点服务启动
	if err = siteRpc.Register(config.Config.RpcPort.SiteRPCPort, config.Config.RpcName.SiteRpcName, 0, site.RegisterFn); err != nil {
		zlogger.Errorw("siteRpc启动失败！", zap.Error(err))
		panic(err)
	}

	service.Start("site")
	defer service.Stop("site")
}
