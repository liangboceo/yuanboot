package endpoints

import (
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"github.com/liangboceo/yuanboot/web"
	"github.com/liangboceo/yuanboot/web/router"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func UsePrometheus(router router.IRouterBuilder) {
	xlog.GetXLogger("Endpoint").Debugf("loaded prometheus endpoint.")

	router.GET("/actuator/metrics", web.WarpHttpHandlerFunc(promhttp.Handler()))
}
