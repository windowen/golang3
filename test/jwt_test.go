package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"

	"serverApi/pkg/jwt"
)

func TestJWT_GenToken(t *testing.T) {
	type fields struct {
		key []byte
		rdb redis.UniversalClient
	}
	type args struct {
		ctx    context.Context
		userId int
	}

	InitConfig()

	// 模拟 Redis 客户端
	mockRdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis 在本地运行
	})

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Success - Valid token generation",
			fields: fields{
				key: []byte("QQYnRFerJTSEcrfB89fw8prOaObmrch8"), // 使用实际的密钥
				rdb: mockRdb,                                    // 此处可以模拟 Redis 客户端或传递实际客户端
			},
			args: args{
				ctx:    context.TODO(),
				userId: 1,
			},
			want:    "expected-token", // 这里替换为预期生成的 token，或者将 `want` 设为空，忽略比较
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := jwt.NewJWT(tt.fields.rdb)

			got, err := j.GenToken(tt.args.ctx, tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != "" && got != tt.want {
				fmt.Println(got)
			}
		})
	}
}
