package main

import (
	"github.com/liangboceo/dependencyinjection"
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"github.com/liangboceo/yuanboot/pkg/configuration"
	"github.com/liangboceo/yuanboot/pkg/servicediscovery/nacos"
	"github.com/liangboceo/yuanboot/web"
	"github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/endpoints"
	"github.com/liangboceo/yuanboot/web/router"
)

func main() {
	host := CreateGatewayBuilder().Build()
	host.Run()
}

// CreateGatewayBuilder 创建网关构建器
func CreateGatewayBuilder() *abstractions.HostBuilder {
	config := configuration.LocalConfig("config")

	return web.NewWebHostBuilder().
		UseConfiguration(config).
		Configure(func(app *web.ApplicationBuilder) {
			app.UseEndpoints(registerGatewayRoutes)
		}).
		ConfigureServices(func(serviceCollection *dependencyinjection.ServiceCollection) {
			// 注册 nacos 服务发现
			nacos.UseServiceDiscovery(serviceCollection)

			// 注册网关服务
			RegisterGatewayServices(serviceCollection)
		}).
		OnApplicationLifeEvent(getApplicationLifeEvent)
}

// registerGatewayRoutes 注册网关路由
func registerGatewayRoutes(rb router.IRouterBuilder) {
	// 注册健康检查端点
	endpoints.UseHealth(rb)
	endpoints.UseReadiness(rb)
	endpoints.UseLiveness(rb)
	endpoints.UsePrometheus(rb)

	// 网关管理接口
	rb.GET("/gateway/services", GetServices)
	rb.GET("/gateway/routes", GetRoutes)
	rb.POST("/gateway/refresh", RefreshCache)
	rb.GET("/gateway/service/:name", GetServiceInstances)

	// 动态路由 - 捕获所有请求并转发到后端服务
	rb.Any("/api/{path:*}", ProxyRequest)
	rb.Any("/user/{path:*}", ProxyRequest)
	rb.Any("/order/{path:*}", ProxyRequest)
	rb.Any("/product/{path:*}", ProxyRequest)
}

// GetServices 获取已注册的服务列表
func GetServices(ctx *context.HttpContext) {
	var gatewayHandler *GatewayHandler
	ctx.RequiredServices.GetService(&gatewayHandler)

	services := gatewayHandler.GetRegisteredServices()
	ctx.JSON(200, map[string]interface{}{
		"services": services,
		"count":    len(services),
	})
}

// GetRoutes 获取路由配置
func GetRoutes(ctx *context.HttpContext) {
	var gatewayHandler *GatewayHandler
	ctx.RequiredServices.GetService(&gatewayHandler)

	ctx.JSON(200, map[string]interface{}{
		"routes": gatewayHandler.GetRoutes(),
		"count":  len(gatewayHandler.GetRoutes()),
	})
}

// RefreshCache 刷新缓存
func RefreshCache(ctx *context.HttpContext) {
	var gatewayHandler *GatewayHandler
	ctx.RequiredServices.GetService(&gatewayHandler)

	gatewayHandler.ForceRefresh()
	ctx.JSON(200, map[string]interface{}{
		"message": "Cache refreshed successfully",
	})
}

// GetServiceInstances 获取服务实例
func GetServiceInstances(ctx *context.HttpContext) {
	var gatewayHandler *GatewayHandler
	ctx.RequiredServices.GetService(&gatewayHandler)

	serviceName := ctx.Input.Param("name")
	instances := gatewayHandler.GetServiceInstances(serviceName)
	ctx.JSON(200, map[string]interface{}{
		"serviceName": serviceName,
		"instances":   instances,
		"count":       len(instances),
	})
}

// ProxyRequest 代理请求处理器
func ProxyRequest(ctx *context.HttpContext) {
	var gatewayHandler *GatewayHandler
	ctx.RequiredServices.GetService(&gatewayHandler)

	// 获取完整的转发路径
	path := "/" + ctx.Input.Param("path")
	gatewayHandler.HandleProxyRequest(ctx, path)
}

// getApplicationLifeEvent 应用生命周期事件
func getApplicationLifeEvent(life *abstractions.ApplicationLife) {
	logger := xlog.GetXLogger("Gateway Life Event:")

	for {
		select {
		case ev := <-life.ApplicationStarted:
			logger.Infof("Gateway started: %v", ev.Data)
		case ev := <-life.ApplicationStopped:
			logger.Infof("Gateway stopped: %v", ev.Data)
			return
		}
	}
}
