package rpcstart

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"serverApi/pkg/gozero/zrpc"
	"serverApi/pkg/service"

	"serverApi/pkg/common/config"
	"serverApi/pkg/gozero/discov"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/tools/mw"
	"serverApi/pkg/tools/network"
	"serverApi/pkg/zlogger"
)

func Register(rpcPort int, rpcRegisterName string, prometheusPort int, rpcFn func(server *grpc.Server)) error {
	rpcKey := strings.ToLower(fmt.Sprintf("%s:///%s", config.Config.Etcd.Schema, rpcRegisterName))
	zlogger.Infow("start", zap.String("register name", rpcKey), zap.Int("server port", rpcPort), zap.Int("prometheusPort:", prometheusPort))

	localIp, err := network.GetLocalIP()
	if err != nil {
		return errs.Wrap(err)
	}

	server := zrpc.MustNewServer(zrpc.RpcServerConf{
		ListenOn: fmt.Sprintf("%s:%d", localIp, rpcPort),
		Etcd: discov.EtcdConf{
			Hosts: config.Config.Etcd.Addr,
			Key:   rpcKey,
		},
	}, rpcFn)

	server.AddOptions(mw.GrpcServer())

	service.RegisterService(server)
	return nil
}
