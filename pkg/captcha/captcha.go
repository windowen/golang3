package captcha

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"time"

	"serverApi/pkg/db/dbconn"
	"serverApi/pkg/db/model/site"
	"serverApi/pkg/tools/cast"
	"serverApi/pkg/tools/utils"

	"github.com/redis/go-redis/v9"

	"serverApi/pkg/captcha/email"
	"serverApi/pkg/captcha/sms"
	"serverApi/pkg/common/config"
	"serverApi/pkg/common/mctx"
	"serverApi/pkg/constant"
	constsR "serverApi/pkg/constant/redis"
	"serverApi/pkg/db/cache"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/tools/strhelper"
	"serverApi/pkg/zlogger"
)

type (
	SendReq struct {
		Lang     string `json:"lang"`      // 语言(用于信息模版获取)
		Scene    string `json:"scene"`     // 发送场景 如: Login
		SendType int    `json:"sendType"`  // 发送类型 1-手机 2-邮箱
		AreaCode string `json:"area_code"` // 地区编码
		Mobile   string `json:"mobile"`    // 手机号
		Email    string `json:"email"`     // 邮箱
	}

	TemplateCodeData struct {
		Code string
	}

	message struct {
		Subject string `json:"subject"`
		Content string `json:"content"`
	}
)

func NewSendReq(sendType int, scene string, areaCode string, mobile string, email string) *SendReq {
	return &SendReq{SendType: sendType, Scene: scene, AreaCode: areaCode, Mobile: mobile, Email: email}
}

func SendCode(ctx context.Context, rdb redis.UniversalClient, req *SendReq) error {
	var (
		captchaCfg      = config.Config.Captcha
		accountKey      = genAccountCacheKey(req)
		basisKey        = fmt.Sprintf(constsR.VerifyCodeCacheBasisKey, req.Scene, accountKey)
		coolDownKey     = fmt.Sprintf(constsR.VerifyCodeCoolDownCacheKey, accountKey)
		code24HLimitKey = fmt.Sprintf(constsR.VerifyCodeCache24HLimitKey, accountKey)
		code            = captchaCfg.DefaultCode
		code24HCount    = constant.Zero
	)

	req.Lang = mctx.GetLanguage(ctx)

	// 是否冷却期
	cmd := rdb.Get(ctx, coolDownKey)
	if err := cache.CheckErr(cmd.Err()); err != nil {
		return errs.ErrInternalServer.Wrap("operation_failed")
	}

	if cmd.Val() != "" {
		return errs.ErrInternalServer.Wrap("captcha_limit")
	}

	// 是否24H限制
	cmd = rdb.Get(ctx, code24HLimitKey)
	if err := cache.CheckErr(cmd.Err()); err != nil {
		return errs.ErrInternalServer.Wrap("operation_failed")
	}

	code24HCount = cast.ToInt(cmd.Val())
	if code24HCount >= captchaCfg.MaxSendTimesOf24Hours && captchaCfg.MaxSendTimesOf24Hours != 0 {
		return errs.ErrNoPermission.Wrap("captcha_limit")
	}
	code24HCount++

	// 是否是开发环境
	if config.Config.App.Env != "dev" && captchaCfg.IsOpen {
		code = strhelper.GenVerifyCode()

		// 获取模版
		tmp, err := findTmpl(ctx, rdb, req, code)
		if err != nil {
			zlogger.Errorf("SendCode findTmpl |req:%v| err: %v", req, err)
			return errs.ErrInternalServer.Wrap("operation_failed")
		}

		switch req.SendType {
		// 短信
		case constant.UserAccountTypeMobile:
			var isFinancialScene bool
			// 是否财务场景(充、提)
			if utils.SliceHas([]string{constant.SceneDep, constant.SceneDraw}, req.Scene) {
				isFinancialScene = true
			}

			// 获取短信渠道
			smsSer, err := sms.NewSms(isFinancialScene)
			if err != nil {
				zlogger.Errorf("SendCode NewSms | err: %v", err)
				return errs.ErrInternalServer.Wrap("operation_failed")
			}

			err = smsSer.SendSms(ctx, &sms.SendSmsRequest{
				AreaCode: req.AreaCode,
				Mobile:   req.Mobile,
				Message:  tmp.Content,
			})
			if err != nil {
				zlogger.Errorf("SendCode SendSms | err: %v", err)
				return errs.ErrInternalServer.Wrap("operation_failed")
			}

		// 邮件
		case constant.UserAccountTypeEmail:
			// 获取邮件渠道
			emailSer, err := email.NewEmail()
			if err != nil {
				zlogger.Errorf("SendCode NewEmail | err: %v", err)
				return errs.ErrInternalServer.Wrap("operation_failed")
			}

			err = emailSer.SendEmail(ctx, &email.SendEmailReq{
				ContentType: constant.EmailContentHtml,
				To:          req.Email,
				Subject:     tmp.Subject,
				Content:     tmp.Content,
			})
			if err != nil {
				zlogger.Errorf("SendCode SendEmail | err: %v", err)
				return errs.ErrInternalServer.Wrap("operation_failed")
			}
		}
	}

	// 缓存短信 却冷 24H限制
	if err := rdb.Set(ctx, basisKey, code, time.Duration(captchaCfg.ValidTime)*time.Second).Err(); err != nil {
		zlogger.Errorf("verification code caching failed, err: %v", err)
		return errs.ErrInternalServer.Wrap("operation_failed")
	}
	if err := rdb.Set(ctx, code24HLimitKey, code24HCount, 24*time.Hour).Err(); err != nil {
		zlogger.Errorf("verification code 24H limit caching failed, err: %v", err)
		return errs.ErrInternalServer.Wrap("operation_failed")
	}

	if err := rdb.Set(ctx, coolDownKey, time.Now().String(), time.Duration(captchaCfg.SendCooldown)*time.Second).Err(); err != nil {
		zlogger.Errorf("verification code cool downKey caching failed, err: %v", err)
		return errs.ErrInternalServer.Wrap("operation_failed")
	}

	return nil
}

func CheckCode(ctx context.Context, rdb redis.UniversalClient, req *SendReq, code string) error {
	var (
		basisKey = fmt.Sprintf(constsR.VerifyCodeCacheBasisKey, genAccountCacheKey(req), req.Mobile)
	)

	if config.Config.App.Env == "dev" {
		return nil
	}

	req.Lang = mctx.GetLanguage(ctx)

	cmd := rdb.Get(ctx, basisKey)
	if err := cache.CheckErr(cmd.Err()); err != nil {
		zlogger.Errorf("failed to obtain %s scenario verification code, err: %v", req.Scene, err)
		return errs.ErrInternalServer.Wrap("operation_failed")
	}

	if !strings.EqualFold(cmd.Val(), code) {
		return errs.ErrInternalServer.Wrap("captcha_error")
	}

	return nil
}

func DelCode(ctx context.Context, rdb redis.UniversalClient, req *SendReq) {
	var (
		basisKey = fmt.Sprintf(constsR.VerifyCodeCacheBasisKey, req.Scene, genAccountCacheKey(req))
	)

	req.Lang = mctx.GetLanguage(ctx)

	err := rdb.Del(ctx, basisKey).Err()
	if err != nil {
		zlogger.Errorf("failed to del user sms code cache, err: %v", err)
		return
	}
}

// 生成缓存key
func genAccountCacheKey(req *SendReq) string {
	switch req.SendType {
	case constant.UserAccountTypeMobile:
		return fmt.Sprintf("%s-%s", req.AreaCode, req.Mobile)
	case constant.UserAccountTypeEmail:
		return req.Email
	}

	return ""
}

// 获取模板
func findTmpl(ctx context.Context, rdb redis.UniversalClient, req *SendReq, code string) (*message, error) {
	tmpCache, err := findTmpCache(ctx, rdb, req)
	if err != nil {
		zlogger.Errorf("findTmpl findTmpCache | err: %v", err)
		return nil, err
	}

	if tmpCache == nil {
		zlogger.Errorf("findTmpl findTmpCache | err: failed to obtain template")
		return nil, errors.New("failed to obtain template")
	}

	// 创建模板并解析模板字符串
	tmpl, err := template.New("notification").Parse(tmpCache.Content)
	if err != nil {
		zlogger.Errorf("findTmpl template.New | err: %v", err)
		return nil, err
	}

	// 生成随机验证码
	data := TemplateCodeData{
		Code: code,
	}

	// 创建一个缓冲区来存储模板执行的结果
	var buf bytes.Buffer

	// 执行模板，将结果写入缓冲区
	err = tmpl.Execute(&buf, data)
	if err != nil {
		zlogger.Errorf("findTmpl Execute | err: %v", err)
		return nil, err
	}

	return &message{
		Subject: tmpCache.Subject,
		Content: buf.String(),
	}, nil
}

// 获取缓存模版
func findTmpCache(ctx context.Context, rdb redis.UniversalClient, req *SendReq) (*message, error) {
	var (
		field = ""
		msg   *message
	)

	switch req.SendType {
	case constant.UserAccountTypeMobile:
		field = fmt.Sprintf("Sms_%v_%v", req.Lang, req.Scene)
	case constant.UserAccountTypeEmail:
		field = fmt.Sprintf("Email_%v_%v", req.Lang, req.Scene)
	default:
		return nil, errors.New("template does not exist")
	}

	result, err := rdb.HGet(ctx, constsR.CommonCaptchaTmplHash, field).Result()
	if err := cache.CheckErr(err); err != nil || result == "" {
		if err != nil {
			zlogger.Errorf("findTmpl HGet | err: %v", err)
		}

		// db获取模版
		tmpCache, err := fetchTmpCache(ctx, field)
		if err != nil {
			zlogger.Errorf("findTmpl fetchTmpCache |field:%v| err: %v", field, err)
			return nil, err
		}

		return tmpCache, err
	}

	err = json.Unmarshal([]byte(result), &msg)
	if err != nil {
		zlogger.Errorf("findTmpl json.Unmarshal |value:%v| err: %v", result, err)
		return nil, err
	}

	return msg, nil
}

// 获取模版
func fetchTmpCache(ctx context.Context, cacheField string) (*message, error) {
	var (
		msg *message
	)

	redisClient, err := cache.NewRedis()
	if err != nil {
		return nil, err
	}

	db, err := dbconn.NewGormDB()
	if err != nil {
		return nil, err
	}

	siteDB := site.NewSite(db)

	// 获取模版
	tmpl, err := siteDB.SysTmpl(ctx)
	if err != nil {
		zlogger.Errorf("FetchTmpCache SysTmpl | err: %v", err)
		return nil, err
	}

	redisClient.Del(ctx, constsR.CommonCaptchaTmplHash)

	for _, val := range tmpl {
		if val.LanguageCode == "" {
			continue
		}

		message := &message{
			Subject: val.Subject,
			Content: val.TemplateContent,
		}
		jsonData, err := json.Marshal(message)
		if err != nil {
			zlogger.Errorf("FetchTmpCache json.Marshal | err: %v", err)
			continue
		}

		category := "Sms"
		if val.Category == constant.TmplTypeEmail {
			category = "Email"
		}

		field := fmt.Sprintf("%v_%v_%v", category, val.LanguageCode, val.TemplateCode)
		if field == cacheField {
			msg = message
		}

		err = redisClient.HSet(ctx, constsR.CommonCaptchaTmplHash, field, jsonData).Err()
		if err != nil {
			zlogger.Errorf("FetchTmpCache HSet |field:%v,value:%v| err: %v", field, cast.ToString(jsonData), err)
			continue
		}
	}

	return msg, nil
}
