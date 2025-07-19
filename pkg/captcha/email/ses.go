package email

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"

	"serverApi/pkg/constant"
	"serverApi/pkg/zlogger"
)

type (
	awsSes struct {
		config *SesCfg
		svc    *ses.SES
	}

	SesCfg struct {
		Sender          string `json:"sender"`          // 发件人
		CharSet         string `json:"charSet"`         // 发件人
		AccessKeyID     string `json:"accessKeyID"`     // Key
		SecretAccessKey string `json:"secretAccessKey"` // 密钥
	}
)

func newAwsSes(config *SesCfg) (Email, error) {
	// 创建 AWS 会话
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(config.AccessKeyID, config.SecretAccessKey, ""),
	})
	if err != nil {
		zlogger.Errorf("Email newAwsSes | err: %v", err)
		return nil, err
	}

	// 创建 SES 服务客户端
	svc := ses.New(sess)

	return &awsSes{config: config, svc: svc}, nil
}

func NewSesCfg(sender string, charSet string, accessKeyID string, secretAccessKey string) *SesCfg {
	return &SesCfg{Sender: sender, CharSet: charSet, AccessKeyID: accessKeyID, SecretAccessKey: secretAccessKey}
}

func (s *awsSes) Name() string {
	return constant.VerifyPlatformSES
}

func (s *awsSes) SendEmail(ctx context.Context, req *SendEmailReq) error {
	var (
		text = ""
		html = ""
	)

	switch req.ContentType {
	case constant.EmailContentText:
		text = req.Content
	case constant.EmailContentHtml:
		html = req.Content
	}

	// 构建邮件内容
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(req.To),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(s.config.CharSet),
					Data:    aws.String(html),
				},
				Text: &ses.Content{
					Charset: aws.String(s.config.CharSet),
					Data:    aws.String(text),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(s.config.CharSet),
				Data:    aws.String(req.Subject),
			},
		},
		Source: aws.String(s.config.Sender),
	}

	// 发送邮件
	_, err := s.svc.SendEmail(input)
	if err != nil {
		var sesErr awserr.Error
		if errors.As(err, &sesErr) {
			switch sesErr.Code() {
			case ses.ErrCodeMessageRejected:
				zlogger.Errorf("Email ses.SendEmail |%v| err: %v", ses.ErrCodeMessageRejected, err)
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				zlogger.Errorf("Email ses.SendEmail |%v| err: %v", ses.ErrCodeMailFromDomainNotVerifiedException, err)
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				zlogger.Errorf("Email ses.SendEmail |%v| err: %v", ses.ErrCodeConfigurationSetDoesNotExistException, err)
			default:
				zlogger.Errorf("Email ses.SendEmail | err: %v", err)
			}
		}

		return err
	}

	return nil
}
