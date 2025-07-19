package apiresp

import (
	"encoding/json"
	"errors"
	"reflect"

	"serverApi/pkg/common/config"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/tools/i18nh"
	"serverApi/pkg/tools/utils"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
	Data    any    `json:"data"`
}

func isAllFieldsPrivate(v any) bool {
	typeOf := reflect.TypeOf(v)
	if typeOf == nil {
		return false
	}
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	if typeOf.Kind() != reflect.Struct {
		return false
	}
	num := typeOf.NumField()
	for i := 0; i < num; i++ {
		c := typeOf.Field(i).Name[0]
		if c >= 'A' && c <= 'Z' {
			return false
		}
	}
	return true
}

func ApiSuccess(data any) *ApiResponse {
	if format, ok := data.(ApiFormat); ok {
		format.ApiFormat()
	}
	if isAllFieldsPrivate(data) {
		return &ApiResponse{}
	}

	if config.Config.App.OpenEncrypt {
		resData, err := json.Marshal(data)
		if err != nil {
			return &ApiResponse{}
		}

		resultData, _ := utils.AesCBCPk7EncryptHex(resData, []byte(config.Config.App.EncryptKey), []byte(config.Config.App.EncryptIV))
		return &ApiResponse{
			Data: resultData,
		}
	}

	return &ApiResponse{
		Data: data,
	}

}

func ParseError(c *gin.Context, err error) *ApiResponse {
	if err == nil {
		return ApiSuccess(nil)
	}

	unwrap := errs.Unwrap(err)
	var codeErr errs.CodeError
	if errors.As(unwrap, &codeErr) {
		resp := ApiResponse{ErrCode: codeErr.Code(), ErrMsg: codeErr.Msg(), ErrDlt: i18nh.T(c, codeErr.Detail())}
		if resp.ErrDlt == "" {
			resp.ErrDlt = err.Error()
		}
		return &resp
	}

	return &ApiResponse{ErrCode: errs.InternalSystemError, ErrMsg: i18nh.T(c, err.Error())}
}
