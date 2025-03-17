package endpoints

import (
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"github.com/liangboceo/yuanboot/utils/jwt"
	"github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/middlewares"
	"github.com/liangboceo/yuanboot/web/router"
	"strconv"
	"time"
)

func UseJwt(router router.IRouterBuilder) {
	xlog.GetXLogger("Endpoint").Debugf("loaded jwt endpoint.")
	config := router.GetConfiguration()
	var secretKey string
	var expires int64
	var hasSecret, hasExpires bool
	if config != nil {
		secretKey, hasSecret = config.Get("yuanboot.application.server.jwt.secret").(string)
		expires, hasExpires = config.Get("yuanboot.application.server.jwt.expires").(int64)
	}
	if !hasSecret {
		secretKey = "12391JdeOW^%$#@"
	}
	if !hasExpires {
		expires = 3
	}
	if config != nil {
		router.POST("/auth/token", func(ctx *context.HttpContext) {
			name := ctx.Input.Param("name")
			id := ctx.Input.Param("id")
			if name == "" || id == "" {
				request := &middlewares.JwtRequest{}
				err := ctx.Bind(request)
				if err == nil {
					id = request.Id
					name = request.Name
				}
			}
			if name == "" || id == "" {
				xlog.GetXLogger("Jwt Endpoint").Debugf("Create Token: name: %s , id: %v , token: %s")
				ctx.JSON(200, context.H{
					"token":   "",
					"expires": 0,
					"success": false,
				})
				return
			}

			uid, _ := strconv.Atoi(id)
			token, expires := jwt.CreateToken([]byte(secretKey), name, uint(uid), int64(time.Now().Add(time.Hour*time.Duration(expires)).Unix()))
			xlog.GetXLogger("Jwt Endpoint").Debugf("Create Token: ( name: %s , id: %s , token: %s )", name, id, token)
			ctx.JSON(200, context.H{
				"token":   token,
				"expires": expires,
				"success": true,
			})
		})
	} else {
		xlog.GetXLogger("Jwt Endpoint").Errorf("config load error.")
	}

}
