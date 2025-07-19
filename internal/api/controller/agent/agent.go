package agent

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"serverApi/pkg/common/config"
	"serverApi/pkg/gozero/discov"
	"serverApi/pkg/gozero/zrpc"
	"serverApi/pkg/protobuf/agent"
	"serverApi/pkg/tools/apiresp"
	"serverApi/pkg/tools/mw"
)

func NewAgent() *Api {
	rpcKey := strings.ToLower(fmt.Sprintf("%s:///%s", config.Config.Etcd.Schema, config.Config.RpcName.SiteRpcName))
	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: config.Config.Etcd.Addr,
			Key:   rpcKey,
		},
	}, zrpc.WithDialOption(mw.GrpcClient()), zrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))

	return &Api{
		Client: agent.NewAgentServerClient(client.Conn()),
	}
}

type Api struct {
	Client agent.AgentServerClient
}

func (o *Api) AgentInfo(c *gin.Context) {
	var req agent.AgentInfoReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, err)
		return
	}
	if err := apiresp.Validate(&req); err != nil {
		apiresp.GinError(c, err) // 参数校验失败
		return
	}
	resp, err := o.Client.AgentInfo(c, &req)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}
