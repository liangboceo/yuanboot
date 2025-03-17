package endpoints

import (
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/router"
)

func UseLiveness(router router.IRouterBuilder) {
	xlog.GetXLogger("Endpoint").Debugf("loaded health endpoint.")

	router.GET("/actuator/health/liveness", func(ctx *context.HttpContext) {
		ctx.JSON(200, context.H{
			"status": "UP",
		})
	})
}
