
<img src="https://yuanboot.star2cloud.com/images/logo.png" width = "380px" height = "120px" alt="" align=center />[中文](https://github.com/liangboceo/yuanboot/blob/master/README.md)  / [English](https://github.com/liangboceo/yuanboot/blob/master/README_En.md)

yuanboot 简单、轻量、快速、基于依赖注入的微服务框架

* 文档： https://yuanboot.star2cloud.com

![Release](https://img.shields.io/github/v/tag/liangboceo/yuanboot.svg?color=24B898&label=release&logo=github&sort=semver)
![Go](https://github.com/liangboceo/yuanboot/workflows/Go/badge.svg)
![GoVersion](https://img.shields.io/github/go-mod/go-version/liangboceo/yuanboot)
[![Report](https://goreportcard.com/badge/github.com/liangboceo/yuanboot)](https://goreportcard.com/report/github.com/liangboceo/yuanboot)
[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?color=24B898&logo=go&logoColor=ffffff)](https://godoc.org/github.com/liangboceo/yuanboot)
![Contributors](https://img.shields.io/github/contributors/liangboceo/yuanboot.svg)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

# yuanboot 特色
- 漂亮又快速的路由器 & MVC 模式 .
- 丰富的中间件支持 (handler func & custom middleware) .
- 微服务框架抽象了分层，在一个框架体系兼容各种server实现，如 rest,grpc等 .
- 充分运用依赖注入DI，管理运行时生命周期，为框架提供了强大的扩展性 .
- 功能强大的微服务集成能力 (Nacos, Eureka, Consul, ETCD) .
- 受到许多出色的 Go Web 框架的启发，并实现了多种 server : **fasthttp** 和 **net.http** 和 **grpc** .

![framework desgin](https://yuanboot.star2cloud.com/images/framework-desgin.jpg)



# 框架安装
```bash
go get github.com/liangboceo/yuanboot
```
# 安装依赖 (由于某些原因国内下载不了依赖)
##  go version < 1.13
```bash
window 下在 cmd 中执行：
set GO111MODULE=on
set GOPROXY=https://goproxy.cn,direct
linux  下执行：
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct
```
##  go version >= 1.13
```
go env -w GOPROXY=https://goproxy.cn,direct
```
### vendor
```
go mod vendor       // 将依赖包拷贝到项目目录中去
```
# 简单的例子
```go
package main
import ...

func main() {
	WebApplication.CreateDefaultBuilder(func(rb router.IRouterBuilder) {
        rb.GET("/info",func (ctx *context.HttpContext) {    // 支持Group方式
            ctx.JSON(200, context.H{"info": "ok"})
        })
    }).Build().Run()       //默认端口号 :8080
}
```
![](https://yuanboot.star2cloud.com/images/starup.png)

# 实现进度
## 标准功能
* [☑️] 打印Logo和日志（yuanboot）
* [️☑️] 统一程序输入参数和环境变量 (yuanboot)
* [☑️] 简单路由器绑定句柄功能
* [☑️] HttpContext 上下文封装(请求，响应)
* [☑️] 静态文件端点（静态文件服务器）
* [☑️] JSON 序列化结构（Context.H）
* [☑️] 获取请求文件并保存
* [☑️] 获取请求数据（form-data，x-www-form-urlencoded，Json ，XML，Protobuf 等）
* [☑️] Http 请求的绑定模型（Url, From，JSON，XML，Protobuf）
### 响应渲染功能
* [☑️] Render Interface
* [☑️] JSON Render
* [☑️] JSONP Render
* [☑️] Indented Json Render
* [☑️] Secure Json Render
* [☑️] Ascii Json Render
* [☑️] Pure Json Render
* [☑️] Binary Data Render
* [☑️] TEXT
* [☑️] Protobuf
* [☑️] MessagePack
* [☑️] XML
* [☑️] YAML
* [☑️] File
* [☑️] Image
* [☑️] Template
* [☑️] Auto formater Render

## 中间件
* [☑️] Logger
* [☑️] StaticFile
* [☑️] Router Middleware
* [☑️] CORS
* [☑️] Binding
* [☑️] JWT
* [☑️] RequestId And Tracker for SkyWorking

## 路由
* [☑️] GET，POST，HEAD，PUT，DELETE 方法支持
* [☑️] 路由解析树与表达式支持
* [☑️] RouteData路由数据 (/api/:version/) 与 Binding的集成
* [☑️] 路由组功能
* [☑️] MVC默认模板功能
* [☑️] 路由过滤器 Filter

## MVC
* [☑️] 路由请求触发Controller&Action
* [☑️] Action方法参数绑定
* [☑️] 内部对象的DI化
* [☑️] 关键对象的参数传递

## Dependency injection
* [☑️] 抽象集成第三方DI框架
* [☑️] MVC模式集成
* [☑️] 框架级的DI支持功能

## 扩展
* [☑️] 配置
* [☑️] WebSocket
* [☑️] JWT
* [☑️] swagger
* [☑️] GRpc
* [☑️] Prometheus
# 进阶范例

```golang
package main

import
import "github.com/liangboceo/dependencyinjection"
...

func main() {
	webHost := CreateCustomWebHostBuilder().Build()
	webHost.Run()
}

// 自定义HostBuilder并支持 MVC 和 自动参数绑定功能，简单情况也可以直接使用CreateDefaultBuilder 。
func CreateCustomBuilder() *abstractions.HostBuilder {

	configuration := abstractions.NewConfigurationBuilder().
		AddEnvironment().
		AddYamlFile("config").Build()

	return WebApplication.NewWebHostBuilder().
		UseConfiguration(configuration).
		Configure(func(app *WebApplication.WebApplicationBuilder) {
			app.UseMiddleware(middlewares.NewCORS())
			//WebApplication.UseMiddleware(middlewares.NewRequestTracker())
			app.UseStaticAssets()
			app.UseEndpoints(registerEndpointRouterConfig)
			app.UseMvc(func(builder *mvc.ControllerBuilder) {
				//builder.AddViews(&view.Option{Path: "./static/templates"})
				builder.AddViewsByConfig()
				builder.AddController(contollers.NewUserController)
				builder.AddFilter("/v1/user/info", &contollers.TestActionFilter{})
			})
		}).
		ConfigureServices(func(serviceCollection *dependencyinjection.ServiceCollection) {
			serviceCollection.AddTransientByImplements(models.NewUserAction, new(models.IUserAction))
			//eureka.UseServiceDiscovery(serviceCollection)
			//consul.UseServiceDiscovery(serviceCollection)
			nacos.UseServiceDiscovery(serviceCollection)
		}).
		OnApplicationLifeEvent(getApplicationLifeEvent)
}

//region endpoint 路由绑定函数
func registerEndpoints(rb router.IRouterBuilder) {
	Endpoints.UseHealth(rb)
	Endpoints.UseViz(rb)
	Endpoints.UsePrometheus(rb)
	Endpoints.UsePprof(rb)
	Endpoints.UseJwt(rb)

	//swagger api document
	endpoints.UseSwaggerDoc(rb,
		swagger.Info{
			Title:          "yuanboot 框架文档演示",
			Version:        "v1.0.0",
			Description:    "框架文档演示swagger文档 v1.0 [ #yuanboot](https://github.com/liangboceo/yuanboot).",
			TermsOfService: "https://yuanboot.star2cloud.com/",
			Contact: swagger.Contact{
				Email: "173120209@qq.com",
				Name:  "yuanboot",
			},
			License: swagger.License{
				Name: "MIT",
				Url:  "https://opensource.org/licenses/MIT",
			},
		},
		func(openapi *swagger.OpenApi) {
			openapi.AddSecurityBearerAuth()
		})

	rb.GET("/error", func(ctx *context.HttpContext) {
		panic("http get error")
	})

	//POST 请求: /info/:id ?q1=abc&username=123
	rb.POST("/info/:id", func(ctx *context.HttpContext) {
		qs_q1 := ctx.Query("q1")
		pd_name := ctx.Param("username")

		userInfo := &UserInfo{}

		_ = ctx.Bind(userInfo) // 手动绑定请求对象

		strResult := fmt.Sprintf("Name:%s , Q1:%s , bind: %s", pd_name, qs_q1, userInfo)

		ctx.JSON(200, context.H{"info": "hello world", "result": strResult})
	})

	// 路由组功能实现绑定 GET 请求:  /v1/api/info
	rb.Group("/v1/api", func(router *router.RouterGroup) {
		router.GET("/info", func(ctx *context.HttpContext) {
			ctx.JSON(200, context.H{"info": "ok"})
		})
	})

	// GET 请求: HttpContext.RequiredServices获取IOC对象
	rb.GET("/ioc", func(ctx *context.HttpContext) {
		var userAction models.IUserAction
		_ = ctx.RequiredServices.GetService(&userAction)
		ctx.JSON(200, context.H{"info": "ok " + userAction.Login("zhang")})
	})
}

//endregion

//region 请求对象
type UserInfo struct {
	UserName string `param:"username"`
	Number   string `param:"q1"`
	Id       string `param:"id"`
}

// ----------------------------------------- MVC 定义 ------------------------------------------------------

// 定义Controller
type UserController struct {
	*mvc.ApiController
	userAction models.IUserAction // IOC 对象参数
}

// 构造器依赖注入
func NewUserController(userAction models.IUserAction) *UserController {
	return &UserController{userAction: userAction}
}

// 请求对象的参数化绑定 , 使用 doc属性标注 支持swagger文档
type RegisterRequest struct {
	mvc.RequestBody `route:"/api/users/register" doc:"用户注册"`
	UserName        string `uri:"userName" doc:"用户名"`
	Password        string `uri:"password" doc:"密码"`
	TestNumber      uint64 `uri:"num" doc:"数字"`
}

// Register函数自动绑定参数
func (this *UserController) Register(ctx *context.HttpContext, request *RegiserRequest) actionresult.IActionResult {
	result := mvc.ApiResult{Success: true, Message: "ok", Data: request}
	return actionresult.Json{Data: result}
}

// use userAction interface by ioc  
func (this *UserController) GetInfo() mvc.ApiResult {
	return this.OK(this.userAction.Login("zhang"))
}

// DocumentResponse custom document response , use doc tag for swagger
type DocumentResponse struct {
	Message string        `json:"message" doc:"消息"`
	List    []DocumentDto `json:"list" doc:"文档列表"`
	Success bool          `json:"success" doc:"是否成功"`
}

// Swagger API 文档支持
func (controller UserController) GetDocumentList(request *struct {
	mvc.RequestGET `route:"/v1/user/doc/list" doc:"获取全部文档列表"`
}) DocumentResponse {

	return DocumentResponse{Message: "GetDocumentList", List: []DocumentDto{
		{Id: 1, Name: "test1", Time: time.Now()}, {Id: 2, Name: "test2", Time: time.Now()},
		{Id: 3, Name: "test3", Time: time.Now()}, {Id: 4, Name: "test4", Time: time.Now()},
		{Id: 5, Name: "test5", Time: time.Now()}, {Id: 6, Name: "test6", Time: time.Now()},
	}, Success: true}
}

// Web程序的开始与停止事件
func fireApplicationLifeEvent(life *abstractions.ApplicationLife) {
	printDataEvent := func(event abstractions.ApplicationEvent) {
		fmt.Printf("[yuanboot] Topic: %s; Event: %v\n", event.Topic, event.Data)
	}
	for {
		select {
		case ev := <-life.ApplicationStarted:
			go printDataEvent(ev)
		case ev := <-life.ApplicationStopped:
			go printDataEvent(ev)
			break
		}
	}
}

```
