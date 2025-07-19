package i18nh

import (
	"serverApi/pkg/constant"
)

func LangToString(langId int) string {
	var langMap map[int]string
	langMap = make(map[int]string)
	langMap[1] = constant.LangChinese
	langMap[2] = constant.LangEnglish
	langMap[3] = constant.LangId
	langMap[4] = constant.LangVi
	langMap[5] = constant.LangPT
	langMap[5] = constant.LangPT
	langMap[6] = constant.LangTaiwan
	if l, ok := langMap[langId]; ok {
		return l
	}
	return ""
}
