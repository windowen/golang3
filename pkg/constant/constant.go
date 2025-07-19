package constant

const (
	ConfigPath               = "/config/config.yaml"
	CtxApiToken              = "api-token"
	LocalHost                = "127.0.0.1"
	ServerAPICommonConfigKey = "ServerAPICommonConfig"
)

const (
	LangChinese = "zh-CN" // 简体中文
	LangTaiwan  = "zh-TW" // 繁体中文
	LangEnglish = "en"    // 英语
	LangId      = "id"    // 印尼语
	LangVi      = "vi"    // 越南语
	LangPT      = "pt"    // 巴西语
)

var (
	LangPack = []string{LangChinese, LangEnglish, LangId, LangVi, LangPT}
)

// 基础响应
const (
	SuccessResponseCode = 10000                      // 成功响应code
	ErrorResponseCode   = SuccessResponseCode + iota // 失败响应code
	LoginResponseCode                                // 需要登录时响应code
	KicOutResponseCode                               // 被踢出
)

const (
	OperationId          = "operationId"
	OpUserId             = "opUserId"
	Token                = "token"
	RpcCustomHeader      = "customHeader" // rpc中间件自定义ctx参数
	ConnID               = "connId"
	CountryCode          = "countryCode"
	OpUserPlatform       = "platform"
	Language             = "language"
	Authorization        = "Authorization"
	RefreshAuthorization = "Refresh-Authorization"
	OpTourist            = "isTourist"    // 是否游客
	Location             = "timeZone"     // 时区
	RefreshToken         = "refreshToken" // 刷新令牌
)

const (
	RpcOperationId = OperationId
	RpcOpUserId    = OpUserId
	RpcOpUserType  = "opUserType"
)

const (
	NormalUser = 1
	AdminUser  = 2
)

const PasswordIteratorCount = 3

const (
	Zero          = 0
	ON            = 0
	Yes           = 1
	StatusNormal  = 1 // 正常
	StatusDisable = 2 // 禁用
	StatusDel     = 3 // 删除
)
