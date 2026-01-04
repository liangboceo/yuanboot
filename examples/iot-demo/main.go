package main

import (
	"embed"
	"github.com/liangboceo/dependencyinjection"
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/abstractions/hosting"
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"github.com/liangboceo/yuanboot/web"
	"github.com/liangboceo/yuanboot/web/actionresult/extension"
	"github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/endpoints"
	"github.com/liangboceo/yuanboot/web/middlewares"
	"github.com/liangboceo/yuanboot/web/mvc"
	"github.com/liangboceo/yuanboot/web/router"
	"iot-demo/controller"
	"iot-demo/iothub"
	"iot-demo/pkg/mq"
	"iot-demo/pkg/service"
	"iot-demo/version"
	"os"
)

//go:embed conf/*.yml
var fs embed.FS

func main() {
	xlog.Fs = fs
	_ = os.Setenv("TZ", "Asia/Shanghai")
	_ = os.Setenv("YUANBOOT_PROFILE", version.Env())
	host := CreateMVCBuilder().Build()
	host.SetAppMode(version.Env())
	host.Run()
}

// * Create the builder of Web host
func CreateMVCBuilder() *abstractions.HostBuilder {
	configuration := abstractions.NewConfigurationBuilder().
		AddEnvironment().
		AddEmbedFs(fs).AddYamlFile("conf/bootstrap").Build()
	return web.NewWebHostBuilder().SetEnvironment(version.Env()).
		UseConfiguration(configuration).
		Configure(func(app *web.ApplicationBuilder) {
			app.SetJsonSerializer(extension.CamelJson())
			app.UseMiddleware(middlewares.NewCORS())
			app.UseStaticAssets()
			app.UseEndpoints(registerEndpointRouterConfig)
			app.UseMvc(func(builder *mvc.ControllerBuilder) {
				builder.AddViewsByConfig() //视图
				builder.EnableRouteAttributes()
				builder.AddController(controller.NewDemoController) // 注册默认
			})
		}).
		ConfigureServices(func(serviceCollection *dependencyinjection.ServiceCollection) {
			// ioc
			serviceCollection.AddSingleton(service.NewCacheService)
			serviceCollection.AddSingleton(mq.NewKafkaClient)
			hosting.AddHostService(serviceCollection, iothub.InitIotHub)

		})
}

func registerEndpointRouterConfig(rb router.IRouterBuilder) {
	//运维监控等配置
	endpoints.UseHealth(rb)
	endpoints.UsePprof(rb)
	endpoints.UseReadiness(rb)
	endpoints.UseLiveness(rb)
	rb.GET("/", func(ctx *context.HttpContext) {
		panic("home")
	})
	rb.GET("/error", func(ctx *context.HttpContext) {
		panic("http get error")
	})

}
