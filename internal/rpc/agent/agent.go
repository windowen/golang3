package agent

import (
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"serverApi/pkg/db/cache"
	"serverApi/pkg/db/dbconn"
	"serverApi/pkg/db/model/agent"
	pb "serverApi/pkg/protobuf/agent"
)

func RegisterFn(server *grpc.Server) {
	db, err := dbconn.NewGormDB()
	if err != nil {
		panic(err)
	}
	redisClient, err := cache.NewRedis()
	if err != nil {
		panic(err)
	}

	pb.RegisterAgentServerServer(server, &agentSvr{
		agentDB: agent.NewAgent(db),
		redis:   redisClient,
	})

	return
}

type agentSvr struct {
	pb.UnimplementedAgentServerServer
	agentDB *agent.Agent
	redis   redis.UniversalClient
}
