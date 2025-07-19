package sms

import (
	"context"
	"fmt"
	"strings"

	"serverApi/pkg/common/config"
	"serverApi/pkg/constant"
	"serverApi/pkg/zlogger"
)

type SMS interface {
	Name() string
	SendSms(ctx context.Context, req *SendSmsRequest) error
}

type SendSmsRequest struct {
	AreaCode string `json:"area_code"` // 地区编码
	Mobile   string `json:"mobile"`    // 手机号
	Message  string `json:"message"`   // 消息
}

func NewSms(isFinancialScene bool) (SMS, error) {
	var (
		smsCfg          = config.Config.Captcha.Sms
		smsPlatformName = smsCfg.DefaultPlfName
	)

	// 财务场景使用特定通道
	if isFinancialScene {
		smsPlatformName = smsCfg.FinancePlfName
	}

	switch strings.ToLower(smsPlatformName) {
	case strings.ToLower(constant.VerifyPlatformPaasoo):
		return newPaasoo(NewPaasooConfig(
			smsCfg.Platform.Paasoo.Endpoint,
			smsCfg.Platform.Paasoo.Key,
			smsCfg.Platform.Paasoo.Secret,
			smsCfg.Platform.Paasoo.From,
		)), nil
	default:
		zlogger.Errorf("Sms NewSms | err: platform initialization failed")
		return nil, fmt.Errorf("not support sms: %s", smsPlatformName)
	}
}
