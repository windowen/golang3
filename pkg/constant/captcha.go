package constant

var (
	CaptchaNotAuthScene = []string{SceneLogin, SceneRegister, SceneForgetPassword}
	CaptchaAuthScene    = []string{SceneModifyLoginPassword, SceneModifyPaymentPassword, SceneDep, SceneDraw, SceneBind}
)

const (
	SceneLogin                 = "Login"                 // 登录
	SceneRegister              = "Register"              // 注册
	SceneModifyLoginPassword   = "ModifyLoginPassword"   // 修改密码
	SceneModifyPaymentPassword = "ModifyPaymentPassword" // 修改支付密码
	SceneForgetPassword        = "ForgetPassword"        // 忘记密码
	SceneDep                   = "Dep"                   // 充值
	SceneDraw                  = "Draw"                  // 提现
	SceneBind                  = "Bind"                  // 绑定

	VerifyPlatformPaasoo = "paasoo"
	VerifyPlatformSES    = "ses"

	EmailContentText = 1 // 文本
	EmailContentHtml = 2 // Html

	TmplTypeSms   = 8 // 短信
	TmplTypeEmail = 9 // 邮件
)
