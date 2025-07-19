package sms

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"serverApi/pkg/constant"
	"serverApi/pkg/tools/httpclient"
	"serverApi/pkg/tools/strhelper"
	"serverApi/pkg/zlogger"
)

type paasoo struct {
	config *PaasooConfig
}

func newPaasoo(cfg *PaasooConfig) SMS {
	return &paasoo{
		config: cfg,
	}
}

type PaasooConfig struct {
	Endpoint string `json:"endpoint"`
	Key      string `json:"key"`
	Secret   string `json:"secret"`
	From     string `json:"from"`
}

type PaaSooMessageResp struct {
	Status     string `json:"status"`
	MessageId  string `json:"messageid"`
	StatusCode string `json:"status_code"`
}

func (m *PaaSooMessageResp) successful() bool {
	return m.Status == "0"
}

func NewPaasooConfig(endpoint string, key string, secret string, from string) *PaasooConfig {
	return &PaasooConfig{Endpoint: endpoint, Key: key, Secret: secret, From: from}
}

func (a *paasoo) Name() string {
	return constant.VerifyPlatformPaasoo
}

func (a *paasoo) SendSms(ctx context.Context, req *SendSmsRequest) error {
	params := url.Values{}
	params.Add("key", a.config.Key)
	params.Add("secret", a.config.Secret)
	params.Add("from", a.config.From)
	params.Add("to", fmt.Sprintf("%s%s", req.AreaCode, req.Mobile))
	params.Add("text", req.Message)

	uri := fmt.Sprintf("%s?%s", a.config.Endpoint, params.Encode())
	resp, err := httpclient.ProxyGet(uri, nil)
	if err != nil {
		zlogger.Errorf("ProxyGet | err: %v", err)
		return err
	}

	var psResp PaaSooMessageResp
	if err = strhelper.Json2Struct(string(resp), &psResp); err != nil {
		zlogger.Errorf("Json2Struct  |body:%v| err: %v", string(resp), err)
		return err
	}

	if !psResp.successful() {
		zlogger.Errorf("sms sending failed, err: %s,%s", psResp.Status, psResp.StatusCode)
		return errors.New("msg_send_fail")
	}

	return nil
}
