package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"serverApi/pkg/tools/network"
	"serverApi/pkg/tools/utils"
	"serverApi/pkg/zlogger"

	"github.com/gin-gonic/gin"

	"serverApi/internal/api/router"
	"serverApi/pkg/common/config"
	"serverApi/pkg/db/cache"
	"serverApi/pkg/tools/component"
	"serverApi/pkg/tools/i18nh"
	"serverApi/pkg/tools/mw"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			tmpStr := fmt.Sprintf("err=%v panic ==> %s\n", err, string(buf[:n]))
			fmt.Println(tmpStr)
		}
	}()

	configFile, logFile, err := config.FlagParse("serverApi-gateway")
	if err != nil {
		panic(err)
	}

	if err = config.InitConfig(configFile); err != nil {
		panic(err)
	}

	if err = component.ComponentCheck(configFile, config.Config.App.Discovery, true); err != nil {
		panic(err)
	}

	// 初始化日志
	svcID := uuid.NewString()
	zlogger.SetGlobalFields(map[string]string{
		"ssid": svcID,
		"snm":  "serverApi-gateway",
	})
	zlogger.InitLogConfig(logFile)

	// 设置全局时区为北京时间
	utils.SetGlobalTimeZone(utils.GetBjTimeLoc())

	// 初始化国际化
	i18nh.New()

	// 默认RPC中间件
	engine := gin.Default()

	rdb, err := cache.NewRedis()
	if err != nil {
		zlogger.Errorw("Failed to initialize Redis", zap.Error(err))
		panic(err)
	}
	engine.Use(mw.CorsHandler(), mw.GinLog())

	// 初始化route
	router.InitRouter(engine, rdb)

	ip, err := network.GetLocalIP()
	if err != nil {
		panic(err)
	}

	host := fmt.Sprintf("%s:%d", ip, config.Config.RpcPort.ServerApiGatePort)
	srv := &http.Server{
		Addr:    host,
		Handler: engine,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zlogger.Errorw("listen", zap.Error(err))
		}
	}()

	zlogger.Infow("server api gateway started", zap.String("host", host))

	// 优雅停机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		zlogger.Errorw("listen", zap.Error(err))
	}
}
