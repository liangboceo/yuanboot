package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"github.com/liangboceo/yuanboot/utils"
	"github.com/liangboceo/yuanboot/utils/jwt"
	"github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/middlewares"
)

type AuthMiddleware struct {
	*middlewares.BaseMiddleware
	Log       xlog.ILogger
	SecretKey string
	appId     string
	SkipPath  []interface{}
}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{BaseMiddleware: &middlewares.BaseMiddleware{},
		Log: xlog.GetXLogger("AuthMiddleware")}
}

type JwtCustomClaims struct {
	jwt.StandardClaims

	// addition
	UserName  string `json:"username"`
	UserId    uint   `json:"userid"`
	IsAdmin   int    `json:"isAdmin"`
	TokenType string `json:"tokenType"`
}

func (middleware *AuthMiddleware) SetConfiguration(config abstractions.IConfiguration) {
	var hasSecretKey bool
	if config != nil {
		middleware.SecretKey, hasSecretKey = config.Get("yuanboot.application.server.uas.auth.jwt-secret").(string)
		middleware.SkipPath, _ = config.Get("yuanboot.application.server.uas.auth.anon-urls").([]interface{})
	}

	if !hasSecretKey {
		middleware.SecretKey = "5Zk2Qx8LpW7rT3eY9uB1vF4sH6dG2jK8mN3bV7cX1zA9sD4fG7hJ2kL5pR8tY3"
	}
}

func (middleware *AuthMiddleware) Inovke(ctx *context.HttpContext, next func(ctx *context.HttpContext)) {
	defer func() {
		if err := recover(); err != nil {
			middleware.Log.Errorf("panic: %v", err)
			middleware.sendUnauthorizedResponse(ctx, "认证失败")
		}
	}()
	middleware.Log.Debug("AuthMiddleware Invoke")
	// 1、原有逻辑：如果JWT未启用或路径在跳过列表中，则跳过验证
	if utils.LikeContains(ctx.Input.Path(), middleware.SkipPath) {
		next(ctx)
		return
	}
	// 2. 跨域预检请求 OPTIONS 直接放行
	if ctx.Input.Request.Method == http.MethodOptions {
		next(ctx)
		return
	}
	// 3. 获取 Authorization Header（WebSocket 连接通过 query param 传递 token）
	authHeader := ctx.Input.Request.Header.Get("Authorization")
	if authHeader == "" {
		// Fallback: check token query parameter for WebSocket connections
		tokenFromQuery := ctx.Input.Request.URL.Query().Get("token")
		if tokenFromQuery != "" {
			authHeader = "Bearer " + tokenFromQuery
		}
	}
	if authHeader == "" {
		middleware.Log.Debug("无token，请重新登录")
		middleware.sendUnauthorizedResponse(ctx, "无token，请重新登录")
		return
	}

	// 4. 验证 Bearer 前缀
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		middleware.Log.Debug("未提供有效的token令牌")
		middleware.sendUnauthorizedResponse(ctx, "未提供有效的token令牌")
		return
	}

	// 5. 提取 Token 并验证
	token := parts[1]
	// 解析 Token
	info, err := jwt.ParseToken(token, []byte(middleware.SecretKey))
	if err != nil {
		middleware.sendUnauthorizedResponse(ctx, "认证失败")
		return
	}
	mapClaims := info.(jwt.MapClaims)
	userInfo := make(map[string]interface{})
	for k, v := range mapClaims {
		userInfo[k] = v
	}
	ctx.SetItem("userinfo", userInfo)

	next(ctx)
}

// sendUnauthorizedResponse 统一返回 401 未授权 JSON 响应
func (middleware *AuthMiddleware) sendUnauthorizedResponse(ctx *context.HttpContext, message string) {
	// 记录调试日志
	middleware.Log.Debug(message)

	// 设置 401 状态码
	ctx.Output.SetStatusCode(http.StatusUnauthorized)

	// 设置响应头为 JSON 格式
	ctx.Output.Response.Header().Set("Content-Type", "application/json; charset=utf-8")

	// 构造 JSON 响应体
	resp := fmt.Sprintf(`{"code":401,"msg":"%s"}`, message)

	// 写入响应体
	_, err := ctx.Output.Response.Write([]byte(resp))
	if err != nil {
		middleware.Log.Error("写入响应失败", err)
	}
}
