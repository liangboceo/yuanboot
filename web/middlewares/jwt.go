package middlewares

import (
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/utils"
	"github.com/liangboceo/yuanboot/utils/jwt"
	"github.com/liangboceo/yuanboot/web/context"
	"net/http"
)

type JwtMiddleware struct {
	*BaseMiddleware

	Enable    bool
	SecretKey string
	Prefix    string
	Header    string
	SkipPath  []interface{}
}

type JwtRequest struct {
	Id   string `json:"id"`
	Name string `extension:"name"`
}

func NewJwt() *JwtMiddleware {
	return &JwtMiddleware{BaseMiddleware: &BaseMiddleware{}}
}

func (jwtmdw *JwtMiddleware) SetConfiguration(config abstractions.IConfiguration) {
	var hasEnable, hasSecret, hasPrefix, hasHeader bool
	if config != nil {
		jwtmdw.Enable, hasEnable = config.Get("yuanboot.application.server.jwt.enable").(bool)
		jwtmdw.SecretKey, hasSecret = config.Get("yuanboot.application.server.jwt.secret").(string)
		jwtmdw.Prefix, hasPrefix = config.Get("yuanboot.application.server.jwt.prefix").(string)
		jwtmdw.Header, hasHeader = config.Get("yuanboot.application.server.jwt.header").(string)
		jwtmdw.SkipPath, _ = config.Get("yuanboot.application.server.jwt.skip_path").([]interface{})
	}

	if !hasEnable {
		jwtmdw.Enable = false
	}

	if !hasSecret {
		jwtmdw.SecretKey = "12391JdeOW^%$#@"
	}
	if !hasPrefix {
		jwtmdw.Prefix = "Bearer"
	}
	if !hasHeader {
		jwtmdw.Header = "Authorization"
	}

	jwtmdw.SkipPath = append(jwtmdw.SkipPath, "/")
	jwtmdw.SkipPath = append(jwtmdw.SkipPath, "/auth/token")
	jwtmdw.SkipPath = append(jwtmdw.SkipPath, "/actuator/health")

}

func (jwtmdw *JwtMiddleware) Inovke(ctx *context.HttpContext, next func(ctx *context.HttpContext)) {
	// 原有逻辑：如果JWT未启用或路径在跳过列表中，则跳过验证
	if !jwtmdw.Enable || utils.LikeContains(ctx.Input.Path(), jwtmdw.SkipPath) {
		next(ctx)
		return
	}

	// 原有JWT验证逻辑
	auth := ctx.Input.Header(jwtmdw.Header)
	if auth == "" {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		return
	}
	token := auth[len(jwtmdw.Prefix)+1:]
	info, err := jwt.ParseToken(token, []byte(jwtmdw.SecretKey))

	if err != nil {
		ctx.Output.SetStatus(http.StatusUnauthorized)
		ctx.Output.Error(http.StatusUnauthorized, "Unauthorized")
		return
	} else {
		mapClaims := info.(jwt.MapClaims)
		userInfo := make(map[string]interface{})
		for k, v := range mapClaims {
			userInfo[k] = v
		}
		ctx.SetItem("userinfo", userInfo)
		next(ctx)
	}
}
