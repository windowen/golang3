package site

import (
	"regexp"
	"strings"

	"serverApi/pkg/tools/cast"

	"serverApi/pkg/common/config"
	"serverApi/pkg/constant"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/tools/utils"
)

// 去除所有 HTML 标签、换行符、回车
func stripHTMLTagsAndNewlines(input string) string {
	// 使用正则表达式匹配并替换 HTML 标签
	re := regexp.MustCompile(`<[^>]*>`)
	// 去除 HTML 标签
	plainText := re.ReplaceAllString(input, "")
	// 将换行符替换为空格
	plainText = strings.ReplaceAll(plainText, "\n", " ")
	// 替换回车符
	plainText = strings.ReplaceAll(plainText, "\r", " ")
	// 去除首尾空白
	plainText = strings.TrimSpace(plainText)

	return plainText
}

// Check 注册参数验证
func (req *RegisterReq) Check() error {
	if !utils.SliceHas(constant.UserAccountTypes, cast.ToInt(req.GetAccountType())) {
		return errs.ErrArgs.WithDetail("parameter_error")
	}

	if cast.ToInt(req.GetAccountType()) == constant.UserAccountTypeMobile {
		if len(strings.TrimSpace(req.GetMobile())) == 0 {
			return errs.ErrArgs.WithDetail("user_mobile_empty")
		}

		if len(strings.TrimSpace(req.GetAreaCode())) == 0 {
			return errs.ErrArgs.WithDetail("user_mobile_empty")
		}

		re := regexp.MustCompile(constant.UserMobileRegexp)
		if !re.MatchString(req.GetMobile()) {
			return errs.ErrArgs.WithDetail("user_mobile_error")
		}
	}

	if cast.ToInt(req.GetAccountType()) == constant.UserAccountTypeEmail {
		if len(strings.TrimSpace(req.GetAreaCode())) == 0 {
			return errs.ErrArgs.WithDetail("user_area_empty")
		}

		if len(strings.TrimSpace(req.GetEmail())) == 0 {
			return errs.ErrArgs.WithDetail("user_email_empty")
		}

		re := regexp.MustCompile(constant.UserEmailRegexp)
		if !re.MatchString(req.GetEmail()) {
			return errs.ErrArgs.WithDetail("user_email_error")
		}
	}

	if len(strings.TrimSpace(req.GetPassword())) == 0 {
		return errs.ErrArgs.WithDetail("user_password_empty")
	}

	re := regexp.MustCompile(constant.UserPasswordRegexp)
	if !re.MatchString(req.GetPassword()) {
		return errs.ErrArgs.WithDetail("user_password_format")
	}

	if len(strings.TrimSpace(req.GetCode())) != config.Config.Captcha.CodeLen {
		return errs.ErrArgs.WithDetail("captcha_error")
	}

	return nil
}

// Check 修改支付密码参数验证
func (req *ModifyPaymentReq) Check() error {
	if len(strings.TrimSpace(req.GetPassword())) == 0 {
		return errs.ErrArgs.WithDetail("user_password_empty")
	}

	re := regexp.MustCompile(constant.UserPasswordRegexp)
	if !re.MatchString(req.GetPassword()) {
		return errs.ErrArgs.WithDetail("user_password_format")
	}

	if len(strings.TrimSpace(req.GetCode())) != config.Config.Captcha.CodeLen {
		return errs.ErrArgs.WithDetail("captcha_error")
	}

	return nil
}

// Check 修改密码参数验证
func (req *ModifyLoginPasswordReq) Check() error {
	if len(strings.TrimSpace(req.GetPassword())) == 0 {
		return errs.ErrArgs.WithDetail("user_password_empty")
	}

	re := regexp.MustCompile(constant.UserPasswordRegexp)
	if !re.MatchString(req.GetPassword()) {
		return errs.ErrArgs.WithDetail("user_password_format")
	}

	if len(strings.TrimSpace(req.GetCode())) != config.Config.Captcha.CodeLen {
		return errs.ErrArgs.WithDetail("captcha_error")
	}
	return nil
}

// Check 登录参数验证
func (req *LoginReq) Check() error {
	if !utils.SliceHas(constant.UserLoginAuthTypes, cast.ToInt(req.GetAuthType())) {
		return errs.ErrArgs.WithDetail("parameter_error")
	}

	if !utils.SliceHas(constant.UserAccountTypes, cast.ToInt(req.GetAccountType())) {
		return errs.ErrArgs.WithDetail("parameter_error")
	}

	if req.GetAccountType() == constant.UserAccountTypeMobile {
		if len(strings.TrimSpace(req.GetMobile())) == 0 {
			return errs.ErrArgs.WithDetail("user_mobile_empty")
		}

		if len(strings.TrimSpace(req.GetAreaCode())) == 0 {
			return errs.ErrArgs.WithDetail("user_mobile_empty")
		}

		re := regexp.MustCompile(constant.UserMobileRegexp)
		if !re.MatchString(req.GetMobile()) {
			return errs.ErrArgs.WithDetail("user_mobile_error")
		}
	}

	if req.GetAccountType() == constant.UserAccountTypeEmail {
		if cast.ToInt(req.GetAuthType()) != constant.UserAuthTypePassword {
			return errs.ErrArgs.WithDetail("parameter_error")
		}

		if len(strings.TrimSpace(req.GetEmail())) == 0 {
			return errs.ErrArgs.WithDetail("user_email_empty")
		}

		re := regexp.MustCompile(constant.UserEmailRegexp)
		if !re.MatchString(req.GetEmail()) {
			return errs.ErrArgs.WithDetail("user_email_error")
		}
	}

	// 验证码登陆
	if req.GetAuthType() == constant.UserAuthTypeCode {
		if len(strings.TrimSpace(req.GetCode())) != config.Config.Captcha.CodeLen {
			return errs.ErrArgs.WithDetail("captcha_error")
		}
	}

	// 密码登录
	if req.GetAuthType() == constant.UserAuthTypePassword {
		if len(strings.TrimSpace(req.GetPassword())) == 0 {
			return errs.ErrArgs.WithDetail("user_password_empty")
		}

		re := regexp.MustCompile(constant.UserPasswordRegexp)
		if !re.MatchString(req.GetPassword()) {
			return errs.ErrArgs.WithDetail("user_password_format")
		}
	}

	return nil
}

// Check 验证码发送参数验证
func (req *SendValidationReq) Check() error {
	if !utils.SliceHas(constant.CaptchaNotAuthScene, req.GetScene()) {
		return errs.ErrArgs.WithDetail("parameter_error")
	}

	if cast.ToInt(req.GetAccountType()) == constant.UserAccountTypeMobile {
		if len(strings.TrimSpace(req.GetMobile())) == 0 {
			return errs.ErrArgs.WithDetail("user_mobile_empty")
		}

		if len(strings.TrimSpace(req.GetAreaCode())) == 0 {
			return errs.ErrArgs.WithDetail("user_mobile_empty")
		}

		re := regexp.MustCompile(constant.UserMobileRegexp)
		if !re.MatchString(req.GetMobile()) {
			return errs.ErrArgs.WithDetail("user_mobile_error")
		}
	}

	if cast.ToInt(req.GetAccountType()) == constant.UserAccountTypeEmail {
		if len(strings.TrimSpace(req.GetEmail())) == 0 {
			return errs.ErrArgs.WithDetail("user_email_empty")
		}

		re := regexp.MustCompile(constant.UserEmailRegexp)
		if !re.MatchString(req.GetEmail()) {
			return errs.ErrArgs.WithDetail("user_email_error")
		}
	}

	return nil
}

func (req *SendValidationCertifiedReq) Check() error {
	if !utils.SliceHas(constant.CaptchaAuthScene, req.GetScene()) {
		return errs.ErrArgs.WithDetail("parameter_error")
	}

	if cast.ToInt(req.GetAccountType()) == constant.UserAccountTypeMobile {
		if len(strings.TrimSpace(req.GetMobile())) == 0 {
			return errs.ErrArgs.WithDetail("user_mobile_empty")
		}

		if len(strings.TrimSpace(req.GetAreaCode())) == 0 {
			return errs.ErrArgs.WithDetail("user_mobile_empty")
		}

		re := regexp.MustCompile(constant.UserMobileRegexp)
		if !re.MatchString(req.GetMobile()) {
			return errs.ErrArgs.WithDetail("user_mobile_error")
		}
	}

	if cast.ToInt(req.GetAccountType()) == constant.UserAccountTypeEmail {
		if len(strings.TrimSpace(req.GetEmail())) == 0 {
			return errs.ErrArgs.WithDetail("user_email_empty")
		}

		re := regexp.MustCompile(constant.UserEmailRegexp)
		if !re.MatchString(req.GetEmail()) {
			return errs.ErrArgs.WithDetail("user_email_error")
		}
	}

	return nil
}

func (req *ForgetPasswordReq) Check() error {
	if cast.ToInt(req.GetAccountType()) == constant.UserAccountTypeMobile {
		if len(strings.TrimSpace(req.GetMobile())) == 0 {
			return errs.ErrArgs.WithDetail("user_mobile_empty")
		}

		if len(strings.TrimSpace(req.GetAreaCode())) == 0 {
			return errs.ErrArgs.WithDetail("user_mobile_empty")
		}

		re := regexp.MustCompile(constant.UserMobileRegexp)
		if !re.MatchString(req.GetMobile()) {
			return errs.ErrArgs.WithDetail("user_mobile_error")
		}
	}

	if cast.ToInt(req.GetAccountType()) == constant.UserAccountTypeEmail {
		if len(strings.TrimSpace(req.GetEmail())) == 0 {
			return errs.ErrArgs.WithDetail("user_email_empty")
		}

		re := regexp.MustCompile(constant.UserEmailRegexp)
		if !re.MatchString(req.GetEmail()) {
			return errs.ErrArgs.WithDetail("user_email_error")
		}
	}

	if len(strings.TrimSpace(req.GetPassword())) == 0 {
		return errs.ErrArgs.WithDetail("user_password_empty")
	}

	re := regexp.MustCompile(constant.UserPasswordRegexp)
	if !re.MatchString(req.GetPassword()) {
		return errs.ErrArgs.WithDetail("user_password_format")
	}

	if len(strings.TrimSpace(req.GetCode())) != config.Config.Captcha.CodeLen {
		return errs.ErrArgs.WithDetail("captcha_error")
	}

	return nil
}

func (req *FollowReq) Check() error {
	if req.GetFollowId() <= 0 {
		return errs.ErrArgs.WithDetail("parameter_error")
	}

	return nil
}

func (req *FollowsReq) Check() error {
	if req.GetLastId() < 0 {
		return errs.ErrArgs.WithDetail("parameter_error")
	}
	if req.GetPageSize() <= 0 || req.GetPageSize() > 999 {
		return errs.ErrArgs.WithDetail("parameter_error")
	}

	return nil
}

func (req *ProfileReq) Check() error {
	if req.GetUserId() < 0 {
		return errs.ErrArgs.WithDetail("parameter_error")
	}

	return nil
}

func (req *FollowCountReq) Check() error {
	if req.GetUserId() < 0 {
		return errs.ErrArgs.WithDetail("parameter_error")
	}

	return nil
}

func (req *ModifyProfileReq) Check() error {
	req.Nickname = stripHTMLTagsAndNewlines(req.GetNickname())
	if len(req.Nickname) == 0 || len(req.Nickname) > 30 {
		return errs.ErrArgs.WithDetail("parameter_error")
	}

	if req.GetSex() != 0 && !utils.SliceHas([]int32{1, 2}, req.GetSex()) {
		return errs.ErrArgs.WithDetail("parameter_error")
	}

	req.Sign = stripHTMLTagsAndNewlines(req.GetSign())
	if len(req.Sign) != 0 && len(req.Sign) > 255 {
		return errs.ErrArgs.WithDetail("parameter_error")
	}

	return nil
}

func (req *SecurityBindingReq) Check() error {
	if !utils.SliceHas(constant.UserAccountTypes, cast.ToInt(req.GetAccountType())) {
		return errs.ErrArgs.WithDetail("parameter_error")
	}

	if req.GetAccountType() == constant.UserAccountTypeMobile {
		if len(strings.TrimSpace(req.GetMobile())) == 0 {
			return errs.ErrArgs.WithDetail("user_mobile_empty")
		}

		if len(strings.TrimSpace(req.GetAreaCode())) == 0 {
			return errs.ErrArgs.WithDetail("user_mobile_empty")
		}

		re := regexp.MustCompile(constant.UserMobileRegexp)
		if !re.MatchString(req.GetMobile()) {
			return errs.ErrArgs.WithDetail("user_mobile_error")
		}
	}

	if req.GetAccountType() == constant.UserAccountTypeEmail {
		if len(strings.TrimSpace(req.GetEmail())) == 0 {
			return errs.ErrArgs.WithDetail("user_email_empty")
		}

		re := regexp.MustCompile(constant.UserEmailRegexp)
		if !re.MatchString(req.GetEmail()) {
			return errs.ErrArgs.WithDetail("user_email_error")
		}
	}

	if len(strings.TrimSpace(req.GetCode())) != config.Config.Captcha.CodeLen {
		return errs.ErrArgs.WithDetail("captcha_error")
	}

	return nil
}

func (req *SiteBannerReq) Check() error {
	if req.GetCategory() <= 0 {
		return errs.ErrArgs.WithDetail("parameter_error")
	}

	return nil
}
