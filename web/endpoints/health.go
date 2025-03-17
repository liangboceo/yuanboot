package endpoints

import (
	"github.com/liangboceo/yuanboot/abstractions/health"
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/router"
)

func UseHealth(router router.IRouterBuilder) {
	xlog.GetXLogger("Endpoint").Debugf("loaded health endpoint.")

	router.GET("/actuator/health/detail", func(ctx *context.HttpContext) {
		var indicatorList []health.Indicator
		_ = ctx.RequiredServices.GetService(&indicatorList)
		builder := health.NewHealthIndicator(indicatorList)
		root := builder.Build()
		statusCode := 200
		if root["status"] != "up" {
			statusCode = 500
		}

		ctx.JSON(statusCode, builder.Build())
	})
}
