package endpoints

import (
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/router"
)

func UseReadiness(router router.IRouterBuilder) {
	xlog.GetXLogger("Endpoint").Debugf("loaded health-readiness endpoint.")

	router.GET("/actuator/health/readiness", func(ctx *context.HttpContext) {
		var appLife *abstractions.ApplicationLife
		_ = ctx.RequiredServices.GetService(&appLife)
		statusCode := 500
		if appLife.State == "up" {
			statusCode = 200
		}

		ctx.JSON(statusCode, context.H{"status": appLife.State})
	})
}
