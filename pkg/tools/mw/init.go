package mw

var (
	empty         = struct{}{}
	ignoreEncrypt = make(map[string]struct{})
)

func init() {

	// 忽略http加密解密
	ignoreEncrypt["/api/site/upload"] = empty

}
