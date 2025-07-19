package mw

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"

	"serverApi/pkg/common/config"
	"serverApi/pkg/constant"
	"serverApi/pkg/tools/apiresp"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/tools/utils"
	"serverApi/pkg/zlogger"

	"github.com/gin-gonic/gin"
)

type RequestHeader struct {
	Language             string
	Platform             string
	CountryCode          string
	OperationId          string
	TimeZone             string
	DeviceId             string
	UserId               int
	TokenExpireAt        string
	RefreshTokenExpireAt string
}

type responseWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.buf.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinLog() gin.HandlerFunc {
	return func(c *gin.Context) {

		language := strings.TrimSpace(c.GetHeader(constant.Language))
		platform := strings.TrimSpace(c.GetHeader(constant.OpUserPlatform))
		operationId := strings.TrimSpace(c.GetHeader(constant.OperationId))
		countryCode := strings.TrimSpace(c.GetHeader(constant.CountryCode))
		location := strings.TrimSpace(c.GetHeader(constant.Location))

		if language == "" || platform == "" || operationId == "" || countryCode == "" || location == "" {
			c.Abort()
			apiresp.GinError(c, errs.ErrArgs.Wrap("para_err"))
			return
		}

		setRequireParamsWithOpts(c,
			WithPlatform(platform),
			WithLanguage(language),
			WithOperationId(operationId),
			WithCountryCode(countryCode),
			WithLocation(location),
		)

		_, encrypt := ignoreEncrypt[c.Request.RequestURI]

		if config.Config.App.OpenEncrypt && !encrypt {
			bodyBytes, err := c.GetRawData()
			if err != nil {
				c.Abort()
				apiresp.GinError(c, errs.ErrArgs.Wrap("para_err"))
				return
			}

			var bodyReq BodyReq
			err = sonic.Unmarshal(bodyBytes, &bodyReq)
			if err != nil {
				c.Abort()
				apiresp.GinError(c, errs.ErrArgs.Wrap("para_err"))
				return
			}

			decryptHex, err := utils.AesCBCPk7DecryptHex(bodyReq.Data, []byte(config.Config.App.EncryptKey), []byte(config.Config.App.EncryptIV))
			if err != nil {
				c.Abort()
				apiresp.GinError(c, errs.ErrArgs.Wrap("para_err"))
				return
			}

			c.Request.Body = io.NopCloser(bytes.NewBuffer([]byte(decryptHex)))

			if config.Config.App.Env == "dev" {
				zlogger.Infow("ServerApi Request", zap.String("operationId", operationId), zap.String("method", c.Request.Method), zap.String("uri", c.Request.RequestURI), zap.String("req", decryptHex))
			}
		} else {
			req, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.String(http.StatusBadRequest, err.Error())
				c.Abort()
				return
			}

			if config.Config.App.Env == "dev" {
				requestHeader := &RequestHeader{
					Language:    c.GetHeader("language"),
					Platform:    c.GetHeader("platform"),
					CountryCode: c.GetHeader("countryCode"),
					OperationId: c.GetHeader("operationId"),
					TimeZone:    c.GetHeader("timeZone"),
					DeviceId:    c.GetHeader("deviceId"),
				}

				token, _ := ParseToken(c.GetHeader("Authorization"))
				if token != nil {
					requestHeader.UserId = token.UserId
					requestHeader.TokenExpireAt = token.ExpiresAt.String()
				}

				refreshToken, _ := ParseToken(c.GetHeader("refreshToken"))
				if refreshToken != nil {
					requestHeader.RefreshTokenExpireAt = refreshToken.ExpiresAt.String()
				}

				zlogger.Infow("ServerApi Request", zap.String("method", c.Request.Method), zap.String("uri", c.Request.RequestURI), zap.String("req", string(req)), zap.Any("headers", requestHeader))
			}

			c.Request.Body = io.NopCloser(bytes.NewReader(req))
		}

		writer := &responseWriter{
			ResponseWriter: c.Writer,
			buf:            bytes.NewBuffer(nil),
		}
		c.Writer = writer
		c.Next()

		resp := writer.buf.Bytes()

		if config.Config.App.Env == "dev" {
			zlogger.Infow("ServerApi Response", zap.String("operationId", operationId), zap.String("resp", string(resp)))
		}
	}
}

type MyCustomClaims struct {
	UserId  int `json:"userId"`
	Version int `json:"version"`
	jwt.RegisteredClaims
}

// ParseToken 解析token
func ParseToken(tokenString string) (*MyCustomClaims, error) {
	// 移除 Bearer 前缀
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, _ := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			zlogger.Errorf("ParseToken Method | err: unexpected signing method")
			return nil, errs.ErrNoPermission.WithDetail("token_match_fail")
		}
		return []byte(config.Config.Jwt.Key), nil
	})

	if token != nil {
		claims := token.Claims.(*MyCustomClaims)
		return claims, nil
	}

	return nil, errs.ErrNoPermission.WithDetail("token_match_fail")
}
