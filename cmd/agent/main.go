package main

import (
	"github.com/google/uuid"
	"go.uber.org/zap"

	"serverApi/pkg/rocketmq"
	"serverApi/pkg/service"
	"serverApi/pkg/tools/utils"
	"serverApi/pkg/zlogger"

	"serverApi/internal/rpc/agent"
	"serverApi/pkg/common/config"
	rpc "serverApi/pkg/rpcstart"
	"serverApi/pkg/tools/component"
	"serverApi/pkg/tools/i18nh"
)

func main() {
	configFile, logFile, err := config.FlagParse("agent")
	if err != nil {
		panic(err)
	}

	if err = config.InitConfig(configFile); err != nil {
		panic(err)
	}

	// 初始化日志
	svcID := uuid.NewString()
	zlogger.SetGlobalFields(map[string]string{
		"ssid": svcID,
		"snm":  "agent",
	})
	zlogger.InitLogConfig(logFile)

	// 设置全局时区为北京时间
	utils.SetGlobalTimeZone(utils.GetBjTimeLoc())

	// 初始化国际化
	i18nh.New()

	if err = component.ComponentCheck(configFile, config.Config.App.Discovery, true); err != nil {
		panic(err)
	}

	// 初始化rocketmq
	rocketmq.Init()

	if err = rpc.Register(config.Config.RpcPort.AgentRPCPort, config.Config.RpcName.AgentRPCName, 0, agent.RegisterFn); err != nil {
		zlogger.Errorw("agentRpc启动失败！", zap.Error(err))
		panic(err)
	}

	service.Start("agent")
	defer service.Stop("agent")
}
