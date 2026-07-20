package main

import (
	"embed"
	"fmt"

	"github.com/liangboceo/dependencyinjection"
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"github.com/liangboceo/yuanboot/pkg/swagger"
	"github.com/liangboceo/yuanboot/web"
	"github.com/liangboceo/yuanboot/web/actionresult/extension"
	"github.com/liangboceo/yuanboot/web/endpoints"
	"github.com/liangboceo/yuanboot/web/middlewares"
	"github.com/liangboceo/yuanboot/web/mvc"
	"github.com/liangboceo/yuanboot/web/router"
	"sendex-server/internal/controller"
	"sendex-server/internal/middleware"
	"sendex-server/internal/service"
	"sendex-server/version"
)

//go:embed conf/*.yml
var fs embed.FS

func main() {
	xlog.Fs = fs
	host := CreateHostBuilder().Build()
	host.SetAppMode(version.Env())
	host.Run()
}

func CreateHostBuilder() *abstractions.HostBuilder {
	config := abstractions.NewConfigurationBuilder().AddEmbedFs(fs).AddYamlFile("conf/config").Build()
	return web.NewWebHostBuilder().SetEnvironment(version.Env()).
		UseConfiguration(config).
		Configure(func(app *web.ApplicationBuilder) {
			app.SetJsonSerializer(extension.CamelJson())
			app.UseMiddleware(middlewares.NewCORS())
			app.UseStaticAssets()
			app.UseEndpoints(registerEndpoints)
			app.UseMvc(func(builder *mvc.ControllerBuilder) {
				builder.AddViewsByConfig()
				builder.EnableRouteAttributes()
				builder.AddController(controller.NewSysController)
			})
			app.UseMiddlewareFront(middleware.NewAuthMiddleware())
		}).
		ConfigureServices(func(sc *dependencyinjection.ServiceCollection) {
			sc.AddSingleton(service.NewDbService)
			sc.AddSingleton(service.NewCacheService)
			sc.AddSingleton(service.NewSysService)
		})
}

func registerEndpoints(rb router.IRouterBuilder) {
	endpoints.UseHealth(rb)
	endpoints.UseViz(rb)
	endpoints.UsePrometheus(rb)
	endpoints.UsePprof(rb)
	endpoints.UseReadiness(rb)
	endpoints.UseLiveness(rb)
	endpoints.UseRouteInfo(rb)
	endpoints.UseSwaggerDoc(rb, swagger.Info{
		Title:          "sendex-server 服务接口文档",
		Version:        version.Version(),
		Description:    fmt.Sprintf("Yuanboot 管理后台接口 %s", version.Version()),
		TermsOfService: "https://yuanboot.star2cloud.com",
		Contact:        swagger.Contact{Name: "yuanboot"},
		License:        swagger.License{Name: "MIT", Url: "https://opensource.org/licenses/MIT"},
	}, func(openapi *swagger.OpenApi) { openapi.AddSecurityBearerAuth() })
}
