package email

import (
	"context"
	"fmt"
	"strings"

	"serverApi/pkg/common/config"
	"serverApi/pkg/constant"
	"serverApi/pkg/zlogger"
)

type Email interface {
	Name() string
	SendEmail(ctx context.Context, req *SendEmailReq) error
}

type SendEmailReq struct {
	ContentType int    `json:"contentType"` // 发送类型 1-文本 2-Html
	To          string `json:"to"`          // 发给谁
	Subject     string `json:"subject"`     // 邮件主题
	Content     string `json:"content"`     // 邮件内容
}

func NewEmail() (Email, error) {
	var emailCfg = config.Config.Captcha.Email

	switch strings.ToLower(emailCfg.PlatformName) {
	case strings.ToLower(constant.VerifyPlatformSES):
		return newAwsSes(NewSesCfg(
			emailCfg.Platform.Ses.Sender,
			emailCfg.Platform.Ses.CharSet,
			emailCfg.Platform.Ses.AccessKeyID,
			emailCfg.Platform.Ses.SecretAccessKey,
		))
	default:
		zlogger.Errorf("Email NewEmail | err: platform initialization failed")
		return nil, fmt.Errorf("not support email: %s", emailCfg.PlatformName)
	}
}
