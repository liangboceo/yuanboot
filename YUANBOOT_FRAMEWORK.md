# YuanBoot 框架文档

## 概述

YuanBoot 是一个简单、轻量、快速、基于依赖注入的微服务框架，采用 Go 语言开发。框架设计灵感来源于 .NET Core 的架构模式，提供了强大的依赖注入、MVC 模式、微服务集成等功能。

### 核心特性

- **依赖注入 (DI)**：基于第三方 DI 框架的深度集成，管理运行时生命周期
- **MVC 模式**：支持传统的 MVC 架构，自动参数绑定和路由解析
- **多服务器支持**：支持 fasthttp、net.http 和 grpc 三种服务器实现
- **微服务集成**：内置 Nacos、Eureka、Consul、ETCD 等服务发现和注册中心
- **丰富的中间件**：提供日志、CORS、JWT、请求追踪等多种中间件
- **配置管理**：支持 YAML 配置文件、环境变量、远程配置中心
- **API 文档**：集成 Swagger 自动生成 API 文档
- **监控支持**：内置 Prometheus、健康检查、性能分析等运维特性

## 架构设计

### 分层架构

```
┌─────────────────────────────────────────────────────────────┐
│                    Application Layer                          │
│  (Console / Web / Grpc / XXL-Job / IoT)                     │
├─────────────────────────────────────────────────────────────┤
│                    Framework Layer                           │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │  Web     │  │  Grpc    │  │ Console  │  │  Scheduler│  │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  │
├─────────────────────────────────────────────────────────────┤
│                    Abstraction Layer                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │   DI     │  │  Config  │  │   Host   │  │  Logging  │  │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  │
├─────────────────────────────────────────────────────────────┤
│                    Infrastructure Layer                       │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │ Service  │  │  Cache   │  │ Database │  │  Message  │  │
│  │ Discovery│  │          │  │          │  │   Queue   │  │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 核心模块

#### 1. abstractions (抽象层)

抽象层是框架的核心，提供了统一的接口和基础功能：

**主要组件：**

- **HostBuilder**：应用构建器，负责配置和构建应用
  - `UseConfiguration()` - 配置管理
  - `ConfigureServices()` - 依赖注入配置
  - `Configure()` - 应用配置
  - `OnApplicationLifeEvent()` - 生命周期事件

- **Configuration**：配置管理
  - 支持 YAML 配置文件
  - 环境变量支持
  - 远程配置中心（Nacos、Apollo）
  - 配置热更新

- **Dependency Injection**：依赖注入
  - 集成第三方 DI 框架
  - 支持生命周期管理（Singleton、Transient、Scoped）
  - 接口和实现分离

- **Service Discovery**：服务发现抽象
  - 统一的服务发现接口
  - 负载均衡策略
  - 服务缓存机制

#### 2. web (Web 模块)

Web 模块提供了完整的 HTTP 服务能力：

**核心组件：**

- **WebHostBuilder**：Web 应用构建器
- **ApplicationBuilder**：应用配置器
- **Router**：路由系统
  - 支持多种 HTTP 方法（GET、POST、PUT、DELETE 等）
  - 路由参数绑定 (`/api/:id`)
  - 路由组功能
  - 路由模板解析

- **HttpContext**：HTTP 上下文
  - 请求处理（Request）
  - 响应处理（Response）
  - 参数绑定
  - 服务获取

- **Middleware**：中间件系统
  - 中间件链式调用
  - 自定义中间件支持
  - 内置中间件：
    - Logger - 日志记录
    - CORS - 跨域处理
    - JWT - 身份认证
    - RequestTracker - 请求追踪
    - StaticFile - 静态文件服务

- **MVC**：MVC 框架
  - Controller - 控制器
  - Action - 动作方法
  - 参数自动绑定
  - 过滤器（Filter）
  - 视图引擎

- **ActionResult**：响应结果
  - JSON、XML、YAML
  - Protobuf、MessagePack
  - 文件下载
  - 模板渲染

#### 3. grpc (gRPC 模块)

gRPC 模块提供了完整的 gRPC 服务支持：

**核心组件：**

- **GrpcHostBuilder**：gRPC 应用构建器
- **ApplicationBuilder**：gRPC 应用配置器
- **Server**：gRPC 服务器
  - 支持 TLS
  - 服务注册
  - 拦截器（Interceptor）

- **ServiceContext**：服务上下文
  - 服务注册
  - 依赖注入集成

- **Interceptor**：拦截器
  - UnaryServerInterceptor - 一元调用拦截
  - StreamServerInterceptor - 流式调用拦截

#### 4. console (控制台模块)

控制台模块用于构建后台服务应用：

**核心组件：**

- **ConsoleHostBuilder**：控制台应用构建器
- **Startup**：启动配置
- **HostService**：后台服务
  - Run() - 启动服务
  - Stop() - 停止服务

#### 5. pkg (扩展包)

扩展包提供了丰富的功能组件：

**服务发现：**
- **Nacos**：阿里 Nacos 服务发现
- **Eureka**：Netflix Eureka 服务发现
- **Consul**：HashiCorp Consul 服务发现
- **ETCD**：CoreOS ETCD 服务发现

**数据源：**
- **MySQL**：MySQL 数据库支持（基于 GORM）
- **Redis**：Redis 缓存支持

**缓存：**
- **Redis**：Redis 缓存实现
- **Memory**：内存缓存

**调度器：**
- **XXL-Job**：分布式任务调度

**配置中心：**
- **Nacos**：Nacos 配置中心
- **Apollo**：Apollo 配置中心

**监控：**
- **Prometheus**：指标监控
- **SkyWalking**：APM 监控

**其他：**
- **Swagger**：API 文档生成
- **WebSocket**：WebSocket 支持
- **Session**：会话管理
- **HTTP Client**：HTTP 客户端（支持服务发现）

## 快速开始

### 安装

```bash
go get github.com/liangboceo/yuanboot
```

### 配置代理（国内用户）

```bash
go env -w GOPROXY=https://goproxy.cn,direct
```

### 创建 Web API 应用

```go
package main

import (
    "github.com/liangboceo/yuanboot/web"
    "github.com/liangboceo/yuanboot/web/context"
    "github.com/liangboceo/yuanboot/web/router"
)

func main() {
    web.CreateHttpBuilder(func(rb router.IRouterBuilder) {
        rb.GET("/info", func(ctx *context.HttpContext) {
            ctx.JSON(200, context.H{"info": "ok"})
        })
    }).Build().Run()
}
```

### 创建 MVC 应用

```go
package main

import (
    "github.com/liangboceo/dependencyinjection"
    "github.com/liangboceo/yuanboot/abstractions"
    "github.com/liangboceo/yuanboot/web"
    "github.com/liangboceo/yuanboot/web/mvc"
    "yourapp/controllers"
)

func main() {
    configuration := abstractions.NewConfigurationBuilder().
        AddEnvironment().
        AddYamlFile("config").Build()

    web.NewWebHostBuilder().
        UseConfiguration(configuration).
        Configure(func(app *web.ApplicationBuilder) {
            app.UseMvc(func(builder *mvc.ControllerBuilder) {
                builder.AddViewsByConfig()
                builder.AddController(controllers.NewUserController)
            })
        }).
        ConfigureServices(func(serviceCollection *dependencyinjection.ServiceCollection) {
            // 配置依赖注入
        }).
        Build().Run()
}
```

### 创建 gRPC 应用

```go
package main

import (
    "github.com/liangboceo/dependencyinjection"
    "github.com/liangboceo/yuanboot/abstractions"
    yrpc "github.com/liangboceo/yuanboot/grpc"
    "google.golang.org/grpc"
    pb "yourapp/proto/helloworld"
    "yourapp/services"
)

func main() {
    configuration := abstractions.NewConfigurationBuilder().
        AddEnvironment().
        AddYamlFile("config").Build()

    yrpc.NewHostBuilder().
        UseConfiguration(configuration).
        Configure(func(app *yrpc.ApplicationBuilder) {
            app.AddGrpcService(func(server *grpc.Server, ctx *yrpc.ServiceContext) {
                ctx.Register(pb.RegisterGreeterServer)
            })
        }).
        ConfigureServices(func(collection *dependencyinjection.ServiceCollection) {
            collection.AddSingleton(services.NewGreeterServer)
        }).
        Build().
        Run()
}
```

### 创建控制台应用

```go
package main

import (
    "github.com/liangboceo/yuanboot/abstractions"
    "github.com/liangboceo/yuanboot/console"
)

func main() {
    config := abstractions.NewConfigurationBuilder().
        AddEnvironment().
        AddYamlFile("config").Build()

    console.NewHostBuilder().
        UseConfiguration(config).
        UseStartup(Startup).
        Build().
        Run()
}
```

## 核心概念

### 依赖注入

YuanBoot 深度集成了依赖注入框架，支持三种生命周期：

**Singleton（单例）：**
```go
serviceCollection.AddSingleton(NewUserService)
```

**Transient（瞬时）：**
```go
serviceCollection.AddTransient(NewUserService)
```

**Scoped（作用域）：**
```go
serviceCollection.AddScoped(NewUserService)
```

**接口注册：**
```go
serviceCollection.AddSingletonByImplements(NewUserService, new(IUserService))
```

### 路由系统

**基础路由：**
```go
rb.GET("/api/users", handler)
rb.POST("/api/users", handler)
rb.PUT("/api/users/:id", handler)
rb.DELETE("/api/users/:id", handler)
```

**路由组：**
```go
rb.Group("/api/v1", func(rg *router.RouterGroup) {
    rg.GET("/users", handler)
    rg.POST("/users", handler)
})
```

**路由参数：**
```go
rb.GET("/api/users/:id", func(ctx *context.HttpContext) {
    id := ctx.Param("id")
    ctx.JSON(200, context.H{"id": id})
})
```

### MVC 模式

**Controller 定义：**
```go
type UserController struct {
    *mvc.ApiController
    userService IUserService
}

func NewUserController(userService IUserService) *UserController {
    return &UserController{userService: userService}
}
```

**Action 方法：**
```go
func (c *UserController) GetUser(id string) mvc.ApiResult {
    user := c.userService.GetUser(id)
    return c.OK(user)
}
```

**参数绑定：**
```go
type CreateUserRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (c *UserController) CreateUser(request *CreateUserRequest) mvc.ApiResult {
    user := c.userService.CreateUser(request)
    return c.OK(user)
}
```

**过滤器：**
```go
type AuthFilter struct{}

func (f *AuthFilter) OnActionExecuting(ctx *mvc.ActionFilterContext) bool {
    // 验证逻辑
    return true
}

func (f *AuthFilter) OnActionExecuted(ctx *mvc.ActionFilterContext) {
    // 执行后逻辑
}
```

### 中间件

**使用内置中间件：**
```go
app.UseMiddleware(middlewares.NewCORS())
app.UseMiddleware(middlewares.NewLogger())
app.UseMiddleware(middlewares.NewJWT())
```

**自定义中间件：**
```go
func CustomMiddleware() middlewares.MiddlewareHandler {
    return middlewares.MiddlewareHandlerFunc(func(ctx *context.HttpContext, next func(ctx *context.HttpContext)) {
        // 前置处理
        next(ctx)
        // 后置处理
    })
}
```

### 配置管理

**配置文件 (config.yml)：**
```yaml
yuanboot:
  application:
    name: myapp
    metadata: "production"
    server:
      type: "web"
      address: ":8080"
  datasource:
    db:
      type: mysql
      url: tcp(localhost:3306)/mydb
      username: root
      password: password
    redis:
      url: localhost:6379
      db: 0
```

**读取配置：**
```go
type AppConfig struct {
    Name string `mapstructure:"name"`
}

var config AppConfig
configuration.UnmarshalKey("yuanboot.application", &config)
```

### 服务发现

**配置服务发现：**
```yaml
yuanboot:
  cloud:
    discovery:
      type: "nacos"
      metadata:
        url: "localhost:8848"
        namespace: "public"
        group: ""
```

**使用服务发现：**
```go
nacos.UseServiceDiscovery(serviceCollection)
```

**调用其他服务：**
```go
type MyController struct {
    discoveryClient servicediscovery.IServiceDiscoveryClient
    httpFactory     httpclient.IDiscoveryClientFactory
}

func (c *MyController) CallOtherService() {
    instance, _ := c.discoveryClient.Select("service-name")
    client := c.httpFactory.CreateClient(instance)
    // 使用 client 调用服务
}
```

## 项目模板

YuanBoot 提供了 CLI 工具 `yuanbootctl` 来快速创建项目：

### 安装 CLI 工具

```bash
go install github.com/liangboceo/yuanboot/cli/yuanbootctl
```

### 可用模板

1. **console** - 控制台应用
2. **webapi** - Web API 应用
3. **mvc** - MVC 应用
4. **grpc** - gRPC 应用
5. **xxl-job** - XXL-Job 任务调度应用
6. **iot** - IoT 应用

### 创建项目

```bash
# 创建 Web API 项目
yuanbootctl new webapi -n myproject -p ./projects

# 创建 MVC 项目
yuanbootctl new mvc -n myproject -p ./projects

# 创建 gRPC 项目
yuanbootctl new grpc -n myproject -p ./projects
```

### 运行项目

```bash
cd myproject
yuanbootctl run
```

### 构建项目

```bash
yuanbootctl build
```

## 运维特性

### 健康检查

```go
endpoints.UseHealth(router)
```

访问：`/actuator/health`

### Prometheus 监控

```go
endpoints.UsePrometheus(router)
```

访问：`/actuator/prometheus`

### 性能分析

```go
endpoints.UsePprof(router)
```

访问：`/actuator/pprof`

### 路由信息

```go
endpoints.UseRouteInfo(router)
```

访问：`/actuator/routers`

### Swagger 文档

```go
endpoints.UseSwaggerDoc(router, swagger.Info{
    Title:          "My API",
    Version:        "v1.0.0",
    Description:    "API Documentation",
})
```

访问：`/swagger/index.html`

## 最佳实践

### 1. 项目结构

```
myapp/
├── main.go                 # 应用入口
├── config.yml              # 配置文件
├── go.mod                  # Go 模块文件
├── controllers/            # 控制器
│   ├── user_controller.go
│   └── product_controller.go
├── services/               # 服务层
│   ├── user_service.go
│   └── product_service.go
├── models/                 # 数据模型
│   ├── user.go
│   └── product.go
├── repositories/           # 数据访问层
│   ├── user_repository.go
│   └── product_repository.go
├── middleware/             # 中间件
│   └── auth_middleware.go
└── filters/               # 过滤器
    └── auth_filter.go
```

### 2. 依赖注入配置

```go
func ConfigureServices(serviceCollection *dependencyinjection.ServiceCollection) {
    // Repository 层
    serviceCollection.AddScoped(NewUserRepository)
    
    // Service 层
    serviceCollection.AddScoped(NewUserService)
    
    // Controller 层
    serviceCollection.AddTransient(NewUserController)
}
```

### 3. 错误处理

```go
func (c *UserController) GetUser(id string) mvc.ApiResult {
    user, err := c.userService.GetUser(id)
    if err != nil {
        return c.Error("用户不存在", err)
    }
    return c.OK(user)
}
```

### 4. 日志记录

```go
logger := xlog.GetXLogger("UserController")
logger.Infof("Getting user: %s", id)
```

### 5. 配置管理

```go
// 使用环境变量
configuration := abstractions.NewConfigurationBuilder().
    AddEnvironment().
    AddYamlFile("config").
    Build()

// 指定配置文件
configuration := abstractions.NewConfigurationBuilder().
    AddYamlFile("config_dev").
    Build()
```

## 高级特性

### 1. WebSocket 支持

YUANBOOT 框架提供了完整的 WebSocket 支持，基于标准的 `fasthttp/websocket` 库实现。WebSocket 允许服务器和客户端之间建立长连接，实现实时双向通信。

#### 核心组件

**1. Hub 中心**

Hub 是 WebSocket 连接的中央管理器，负责：
- 管理所有活跃的客户端连接
- 处理客户端注册和注销
- 广播消息给所有客户端
- 维护连接状态

**2. Client 客户端**

Client 代表单个 WebSocket 连接，负责：
- 读取客户端消息并发送到 Hub
- 从 Hub 接收消息并发送给客户端
- 处理连接的生命周期
- 管理消息的读写操作

**3. 连接升级**

使用 `web.Upgrade()` 函数将 HTTP 连接升级为 WebSocket 连接。

#### 完整示例

**1. 定义 Hub 结构体**

```go
package hubs

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				_ = client.Conn.Close()
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}
```

**2. 定义 Client 结构体**

```go
package hubs

import (
	"bytes"
	"github.com/fasthttp/websocket"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 65535
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client is a middleman between the socket connection and the hub.
type Client struct {
	Id  string
	Hub *Hub

	// The socket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	MaxMessageSize int64
}

// ReadPump pumps messages from the socket connection to the Hub.
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
	}()
	c.Conn.SetReadLimit(c.MaxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { _ = c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.Hub.broadcast <- message
	}
}

// WritePump pumps messages from the Hub to the socket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				if c.Conn != nil {
					_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				}
				return
			}
			if c.Conn != nil {
				w, err := c.Conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}
				_, _ = w.Write(message)

				// Add queued chat messages to the current socket message.
				n := len(c.Send)
				for i := 0; i < n; i++ {
					_, _ = w.Write(newline)
					_, _ = w.Write(<-c.Send)
				}

				if err := w.Close(); err != nil {
					return
				}
			}
		case <-ticker.C:
			if c.Conn != nil {
				_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}
}
```

**3. 注册 Hub 到依赖注入容器**

```go
func ConfigureServices(serviceCollection *dependencyinjection.ServiceCollection) {
	// 注册 WebSocket Hub
	serviceCollection.AddSingleton(hubs.NewHub)
	
	// 注册其他服务...
}
```

**4. 创建 WebSocket 控制器**

```go
package controllers

import (
	"github.com/fasthttp/websocket"
	"github.com/liangboceo/yuanboot/web"
	"github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/mvc"
	"yourapp/hubs"
)

// websocket hub
type HubController struct {
	mvc.ApiController

	hub *hubs.Hub
}

func NewHubController(hub *hubs.Hub) *HubController {
	go hub.Run() // 启动 Hub 处理循环
	return &HubController{hub: hub}
}

// WebSocket 连接端点
// URL: ws://localhost:8080/app/v1/hub/ws
func (controller HubController) GetWs(ctx *context.HttpContext) {
	web.Upgrade(ctx, func(conn *websocket.Conn) {
		client := &hubs.Client{
			Hub:            controller.hub,
			Conn:           conn,
			Send:           make(chan []byte, 256),
			MaxMessageSize: 65535,
		}
		client.Hub.Register <- client
		go client.WritePump()
		client.ReadPump()
	})
}
```

**5. 配置 WebSocket 路由**

```go
app.UseMvc(func(builder *mvc.ControllerBuilder) {
	builder.AddController(controllers.NewHubController)
})
```

#### 高级特性

**1. 客户端身份标识**

```go
func (controller HubController) GetWs(ctx *context.HttpContext) {
	web.Upgrade(ctx, func(conn *websocket.Conn) {
		// 从请求中获取用户身份
		userId := ctx.Query("userId")
		
		client := &hubs.Client{
			Id:             userId,
			Hub:            controller.hub,
			Conn:           conn,
			Send:           make(chan []byte, 256),
			MaxMessageSize: 65535,
		}
		client.Hub.Register <- client
		go client.WritePump()
		client.ReadPump()
	})
}
```

**2. 定向消息发送**

```go
// 向特定客户端发送消息
func (h *Hub) SendToClient(clientId string, message []byte) {
	for client := range h.clients {
		if client.Id == clientId {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
			break
		}
	}
}
```

**3. 消息处理和业务逻辑**

```go
// 自定义消息处理
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
	}()
	// ... 省略部分代码 ...
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			// 错误处理
			break
		}
		
		// 处理消息
		c.processMessage(message)
	}
}

func (c *Client) processMessage(message []byte) {
	// 解析消息
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		// 处理解析错误
		return
	}
	
	// 根据消息类型处理
	switch msg.Type {
	case "chat":
		// 处理聊天消息
		c.Hub.broadcast <- message
	case "join":
		// 处理加入消息
		c.handleJoin(msg)
	case "leave":
		// 处理离开消息
		c.handleLeave(msg)
	}
}
```

**4. 与 Redis 集成**

```go
// 使用 Redis 作为消息存储和分发
func (h *Hub) Run() {
	// 初始化 Redis 连接
	redisClient := redis.NewClient(/* 配置 */)
	
	// 订阅消息
	pubsub := redisClient.Subscribe("websocket:messages")
	ch := pubsub.Channel()
	
	go func() {
		for msg := range ch {
			h.broadcast <- []byte(msg.Payload)
		}
	}()
	
	// 处理客户端消息
	for {
		select {
		// ... 省略部分代码 ...
		case message := <-h.broadcast:
			// 广播给所有客户端
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
			
			// 发布到 Redis
			redisClient.Publish("websocket:messages", string(message))
		}
	}
}
```

#### 最佳实践

1. **连接管理**
   - 实现合理的心跳机制
   - 处理连接超时和异常关闭
   - 限制单个客户端的消息大小

2. **性能优化**
   - 使用缓冲通道避免阻塞
   - 合理设置消息队列大小
   - 考虑使用连接池管理WebSocket连接

3. **安全性**
   - 实现WebSocket认证
   - 防止消息注入攻击
   - 限制连接频率和消息速率

4. **可扩展性**
   - 使用Redis等中间件实现跨实例消息分发
   - 考虑使用消息队列处理高并发场景
   - 实现连接状态的持久化

#### 常见问题

**Q: WebSocket 连接失败怎么办？**
A: 检查网络连接、防火墙设置、服务端口是否正确，以及客户端和服务端的WebSocket实现是否兼容。

**Q: 如何处理大量WebSocket连接？**
A: 优化服务器配置，增加最大文件描述符限制，使用连接池，考虑水平扩展。

**Q: 如何实现消息的可靠传递？**
A: 实现消息确认机制，使用持久化存储，考虑使用MQ保证消息投递。

**Q: WebSocket 与 HTTP 性能比较？**
A: WebSocket 在长连接场景下性能更好，减少了HTTP握手开销，适合实时通信场景。

通过以上实现，YUANBOOT 框架提供了完整的 WebSocket 支持，可用于构建实时聊天、实时数据监控、游戏等需要双向通信的应用场景。

### 2. 会话管理

```go
session.UseSession(serviceCollection, func(options *session.Options) {
    options.AddSessionStoreFactory(store.NewRedis)
    options.AddSessionIdentity(identity.NewCookie())
})
```

### 3. JWT 认证

```go
app.UseMiddleware(middlewares.NewJWT(func(ctx *context.HttpContext) (interface{}, error) {
    token := ctx.GetHeader("Authorization")
    // 验证 token
    return claims, nil
}))
```

### 4. 分布式任务调度

```go
scheduler.NewXxlJobBuilder(config).
    ConfigureServices(func(collection *dependencyinjection.ServiceCollection) {
        scheduler.AddJobs(collection, NewDemoJob)
    }).
    Build().
    Run()
```

### 5. APM 监控

```yaml
yuanboot:
  cloud:
    apm:
      skywalking:
        address: localhost:11800
```

## 性能优化

### 1. 连接池配置

```yaml
yuanboot:
  datasource:
    pool:
      init_cap: 2
      max_cap: 5
      idle_timeout: 5
```

### 2. 缓存策略

```go
// 使用 Redis 缓存
cache := redis.NewRedisCache(config)
cache.Set("key", value, time.Minute)
```

### 3. 负载均衡

```go
// 随机策略
serviceCollection.AddSingletonByImplements(strategy.NewRandom, new(servicediscovery.Strategy))

// 轮询策略
serviceCollection.AddSingletonByImplements(strategy.NewRoundRobin, new(servicediscovery.Strategy))
```

## 故障排查

### 1. 日志级别配置

```yaml
log:
  level: debug
  format: json
```

### 2. 性能分析

```bash
# 启用 pprof
curl http://localhost:8080/actuator/pprof/heap > heap.prof
go tool pprof heap.prof
```

### 3. 健康检查

```bash
curl http://localhost:8080/actuator/health
```

## 总结

YuanBoot 是一个功能完善、设计优雅的 Go 微服务框架，具有以下优势：

1. **架构清晰**：分层设计，职责明确
2. **扩展性强**：基于依赖注入，易于扩展
3. **功能丰富**：内置微服务所需的各种功能
4. **易于使用**：提供 CLI 工具和项目模板
5. **生产就绪**：包含完整的监控、日志、健康检查等运维特性

适用于构建各种类型的 Go 应用，从简单的 Web API 到复杂的微服务架构。

## 相关资源

- **GitHub**: https://github.com/liangboceo/yuanboot
- **文档**: https://yuanboot.star2cloud.com
- **示例**: https://github.com/liangboceo/yuanboot/tree/master/examples

## 许可证

MIT License
