package test

import (
	"context"
	"testing"

	"github.com/redis/go-redis/v9"

	"serverApi/pkg/captcha"
	"serverApi/pkg/constant"
)

func TestSendCaptcha(t *testing.T) {
	InitLog()

	InitConfig()

	// 模拟 Redis 客户端
	mockRdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis 在本地运行
	})

	type args struct {
		ctx context.Context
		rdb redis.UniversalClient
		req *captcha.SendReq
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Sms send",
			args: args{
				ctx: context.Background(),
				rdb: mockRdb,
				req: &captcha.SendReq{
					Lang:     "zh-CN",
					SendType: constant.UserAccountTypeMobile,
					Scene:    constant.SceneLogin,
					AreaCode: "855",
					Mobile:   "081349704",
				},
			},
		},
		{
			name: "Email send",
			args: args{
				ctx: context.Background(),
				rdb: mockRdb,
				req: &captcha.SendReq{
					Lang:     "zh-CN",
					SendType: constant.UserAccountTypeEmail,
					Scene:    constant.SceneLogin,
					Email:    "yolowork6@gmail.com",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := captcha.SendCode(tt.args.ctx, tt.args.rdb, tt.args.req)
			if err != nil {
				t.Errorf("Sendcaptcha error = %v", err)
			}
		})
	}
}
