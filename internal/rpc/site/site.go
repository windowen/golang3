package site

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/grpc/credentials/insecure"

	"serverApi/pkg/gozero/discov"
	"serverApi/pkg/gozero/zrpc"
	"serverApi/pkg/tools/cast"
	"serverApi/pkg/tools/mw"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"serverApi/pkg/captcha"
	"serverApi/pkg/common/config"
	"serverApi/pkg/db/cache"
	"serverApi/pkg/db/dbconn"
	"serverApi/pkg/db/model/site"
	"serverApi/pkg/jwt"
	commonPb "serverApi/pkg/protobuf/common"
	"serverApi/pkg/protobuf/live"
	pb "serverApi/pkg/protobuf/site"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/tools/redis_lock"
	"serverApi/pkg/zlogger"
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

	siteDB := site.NewSite(db)

	rpcKey := strings.ToLower(fmt.Sprintf("%s:///%s", config.Config.Etcd.Schema, config.Config.RpcName.LiveRPCName))
	client := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: config.Config.Etcd.Addr,
			Key:   rpcKey,
		},
	}, zrpc.WithDialOption(mw.AddUserType()), zrpc.WithDialOption(mw.GrpcClient()), zrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))

	pb.RegisterSiteServiceServer(server, &siteSrv{
		liveApiRpcClient: live.NewLiveServerClient(client.Conn()),
		siteDB:           siteDB,
		redis:            redisClient,
		jwt:              jwt.NewJWT(redisClient),
		userCache:        cache.NewUserCache(redisClient),
		roomCache:        cache.NewRoomCache(redisClient),
		redisLock:        redis_lock.NewRedisLock(redisClient),
	})

	return
}

type siteSrv struct {
	pb.UnimplementedSiteServiceServer
	liveApiRpcClient live.LiveServerClient
	siteDB           *site.Site
	redis            redis.UniversalClient
	jwt              *jwt.JWT
	userCache        *cache.UserCache
	roomCache        *cache.RoomCache
	redisLock        *redis_lock.RedisLock
}

// SendValidationCode 发送短信消息(未登录)
func (s *siteSrv) SendValidationCode(ctx context.Context, req *pb.SendValidationReq) (*commonPb.Empty, error) {
	var err error

	err = captcha.SendCode(ctx, s.redis, captcha.NewSendReq(cast.ToInt(req.GetAccountType()), req.GetScene(), req.GetAreaCode(), req.GetMobile(), req.GetEmail()))
	if err != nil {
		zlogger.Errorf("SendValidationCode SendSmsCode |req:%v| err: %v", cast.ToString(req), err)
		return nil, err
	}

	return nil, nil
}

// GlobalAreas 全球地区码
func (s *siteSrv) GlobalAreas(ctx context.Context, _ *commonPb.Empty) (*pb.GlobalAreaResp, error) {
	var globalResp pb.GlobalAreaResp

	// 获取缓存地区语言
	areaCache, err := s.userCache.CommonGlobalAreaGet(ctx)
	if err == nil && areaCache != "" {
		if err := json.Unmarshal([]byte(areaCache), &globalResp); err == nil {
			return &globalResp, nil
		}
	}

	areas, err := s.siteDB.GlobalAreas(ctx)
	if err != nil {
		zlogger.Errorf("GlobalAreas | err:%v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	var globalAreas []*pb.GlobalAreaResp_GlobalArea
	for _, area := range areas {
		globalAreas = append(globalAreas,
			&pb.GlobalAreaResp_GlobalArea{
				Name:        area.Name,
				MobileCode:  area.MobileCode,
				CountryCode: area.CountryCode,
			})
	}

	langList, err := s.siteDB.GlobalLang(ctx)
	if err != nil {
		zlogger.Errorf("GlobalAreas | err:%v", err)
		return nil, errs.ErrInternalServer.Wrap("query_failed")
	}

	var globalLang []*pb.GlobalAreaResp_GlobalLang
	for _, lang := range langList {
		globalLang = append(globalLang,
			&pb.GlobalAreaResp_GlobalLang{
				LanguageCode: lang.LanguageCode,
				LanguageName: lang.LanguageName,
			})
	}

	globalResp.Area = globalAreas
	globalResp.Lang = globalLang

	err = s.userCache.CommonGlobalAreaSet(ctx, &globalResp)
	if err != nil {
		zlogger.Errorf("GlobalAreas CommonGlobalAreaSet |data:%v| err:%v", cast.ToString(&globalResp), err)
		return nil, nil
	}

	return &globalResp, nil
}
