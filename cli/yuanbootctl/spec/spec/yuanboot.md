# Yuanboot Framework Specification

## 框架概述

Yuanboot 是一个简单、轻量、快速、基于依赖注入的 Go 微服务框架。

- **版本**: v1.9.2747
- **文档**: https://yuanboot.star2cloud.com
- **仓库**: https://github.com/liangboceo/yuanboot

### 核心特色

- 漂亮又快速的路由器 & MVC 模式
- 丰富的中间件支持 (handler func & custom middleware)
- 微服务框架抽象分层，兼容多种 server 实现（REST、gRPC）
- 充分运用依赖注入（DI），管理运行时生命周期
- 功能强大的微服务集成能力（Nacos, Eureka, Consul, ETCD）
- 支持多种 HTTP 服务器实现：**fasthttp**、**net.http**、**grpc**

---

## 目录结构

```
yuanboot/
├── abstractions/          # 核心抽象层
│   ├── configuration/     # 配置系统
│   ├── health/            # 健康检查
│   ├── hosting/           # 主机托管
│   ├── pool/              # 连接池
│   ├── servicediscovery/ # 服务发现
│   ├── xlog/              # 日志系统
│   └── *.go               # 核心接口
├── pkg/                   # 功能包
│   ├── cache/             # 缓存（Redis）
│   ├── captcha/           # 验证码
│   ├── configuration/     # 配置管理
│   ├── datasources/       # 数据源（MySQL、Redis）
│   ├── email/             # 邮件服务
│   ├── httpclient/        # HTTP客户端
│   ├── mq/                # 消息队列
│   ├── scheduler/         # 定时任务
│   ├── servicediscovery/  # 服务发现实现
│   ├── swagger/           # Swagger文档
│   └── task/              # 任务管理
├── web/                   # Web框架
│   ├── actionresult/       # 响应结果
│   ├── binding/           # 数据绑定
│   ├── captcha/           # Web验证码
│   ├── context/           # HTTP上下文
│   ├── endpoints/         # 端点（健康检查、监控）
│   ├── middlewares/       # 中间件
│   ├── mvc/               # MVC框架
│   ├── router/            # 路由系统
│   ├── session/           # 会话管理
│   └── view/              # 视图模板
├── grpc/                  # gRPC框架
│   ├── conn/              # 连接管理
│   └── interceptors/      # 拦截器
├── utils/                 # 工具函数
│   ├── jwt/               # JWT认证
│   ├── cast/              # 类型转换
│   └── typeconverter/     # 类型转换器
├── cli/                   # CLI命令行工具
├── examples/              # 示例代码
└── tests/                 # 测试用例
```

---

## 快速开始

### 安装框架

```bash
go get github.com/liangboceo/yuanboot
```

### 简单示例

```go
package main

import (
    "github.com/liangboceo/yuanboot/web"
    "github.com/liangboceo/yuanboot/web/context"
    "github.com/liangboceo/yuanboot/web/router"
)

func main() {
    web.WebApplication.CreateDefaultBuilder(func(rb router.IRouterBuilder) {
        rb.GET("/info", func(ctx *context.HttpContext) {
            ctx.JSON(200, context.H{"info": "ok"})
        })
    }).Build().Run()  // 默认端口 :8080
}
```

---

## abstractions 核心抽象层

### 配置系统 (configuration/)

配置系统提供统一的配置管理，支持 YAML、环境变量、嵌入式文件系统。

```go
// 创建配置构建器
config := abstractions.NewConfigurationBuilder().
    AddEnvironment().                    // 添加环境变量
    AddYamlFile("config").              // 添加 YAML 文件
    Build()

// 获取配置值
config.GetString("app.name")
config.GetInt("app.port")
config.GetBool("app.debug")
config.GetSection("app")               // 获取配置节
```

**配置接口方法**:

| 方法 | 说明 |
|------|------|
| `Get(name)` | 获取配置值 |
| `GetString/GetBool/GetInt` | 按类型获取值 |
| `GetSection(name)` | 获取配置节 |
| `Unmarshal(v)` | 反序列化到结构体 |
| `GetProfile()` | 获取环境配置 (Dev/Prod) |
| `GetConfDir()` | 获取配置目录 |
| `RefreshAll()` | 刷新所有配置 |

### 日志系统 (xlog/)

```go
import "github.com/liangboceo/yuanboot/abstractions/xlog"

// 获取日志器
log := xlog.GetXLogger("MyService")

// 记录日志
log.Debug("debug message: %s", "info")
log.Info("info message: %s", "info")
log.Warn("warn message: %s", "info")
log.Error("error message: %s", "info")
```

### 连接池 (pool/)

```go
import "github.com/liangboceo/yuanboot/abstractions/pool"

// 创建连接池
p, err := pool.NewChannelPool(&pool.Config{
    InitialCap: 5,           // 初始连接数
    MaxCap: 30,               // 最大连接数
    Factory: func() (interface{}, error) {
        return createConnection()
    },
    Close: func(v interface{}) error {
        return v.(*sql.DB).Close()
    },
    Ping: func(v interface{}) error {
        return v.(*sql.DB).Ping()
    },
    IdleTimeout: 5 * time.Minute,
})
```

**Pool 接口**:

| 方法 | 说明 |
|------|------|
| `Get()` | 获取连接 |
| `Put(interface{})` | 归还连接 |
| `Close()` | 关闭连接池 |
| `Ping(interface{})` | 检测连接 |

### 服务发现 (servicediscovery/)

```go
// 服务发现接口
type IServiceDiscovery interface {
    GetName() string
    Register() error
    Update() error
    Unregister() error
    GetHealthyInstances(name string) []ServiceInstance
    GetAllInstances(name string) []ServiceInstance
    Watch(opts ...WatchOption) (Watcher, error)
    GetAllServices() ([]*Service, error)
}
```

**支持的服务注册中心**:

- Nacos
- Eureka
- Consul
- ETCD

### 健康检查 (health/)

```go
// 健康指示器接口
type IHealthIndicator interface {
    Health() (HealthStatus, error)
}
```

---

## web Web框架

### 创建应用

```go
// 简单构建器
app := web.WebApplication.CreateDefaultBuilder(func(rb router.IRouterBuilder) {
    // 路由配置
}).Build()

// 自定义构建器
app := web.WebApplication.NewWebHostBuilder().
    UseConfiguration(config).
    Configure(func(builder *web.WebApplicationBuilder) {
        builder.UseMiddleware(middlewares.NewCORS())
        builder.UseEndpoints(registerEndpoints)
        builder.UseMvc(func(cb *mvc.ControllerBuilder) {
            cb.AddController(user.NewUserController)
        })
    }).
    ConfigureServices(func(sc *dependencyinjection.ServiceCollection) {
        sc.AddTransientByImplements(services.NewUserService, new(services.IUserService))
    }).
    Build()
```

### HTTP上下文 (context/)

```go
func handler(ctx *context.HttpContext) {
    // 获取请求参数
    id := ctx.Param("id")
    name := ctx.Query("name")
    token := ctx.Header("Authorization")
    
    // 绑定请求数据
    var req UserRequest
    ctx.Bind(&req)
    
    // 响应数据
    ctx.JSON(200, context.H{
        "id":   id,
        "user": req,
    })
}
```

**请求绑定类型**:

| 类型 | Tag | 说明 |
|------|-----|------|
| URI参数 | `param:"name"` | 路由参数 |
| Query参数 | `query:"name"` | URL查询参数 |
| Header参数 | `header:"name"` | HTTP头 |
| JSON Body | `json:"name"` | JSON请求体 |
| Form表单 | `form:"name"` | 表单数据 |
| Path路径 | `path:"name"` | URL路径 |

**响应方法**:

| 方法 | 说明 |
|------|------|
| `JSON(code int, data)` | JSON响应 |
| `HTML(code int, html)` | HTML响应 |
| `Text(code int, text)` | 文本响应 |
| `XML(code int, data)` | XML响应 |
| `File(filePath)` | 文件下载 |
| `Redirect(url)` | 重定向 |
| `Render(view, data)` | 模板渲染 |

### 路由系统 (router/)

```go
// 基础路由
rb.GET("/users", handler)
rb.POST("/users", handler)
rb.PUT("/users/:id", handler)
rb.DELETE("/users/:id", handler)

// 路由组
rb.Group("/api/v1", func(g *router.RouterGroup) {
    g.GET("/users", listHandler)
    g.POST("/users", createHandler)
    g.GET("/users/:id", getHandler)
})

// 路由过滤器
rb.Filter("/admin/*", adminFilter)
```

### 中间件 (middlewares/)

```go
// 使用内置中间件
builder.UseMiddleware(middlewares.NewRecovery())
builder.UseMiddleware(middlewares.NewLogger())
builder.UseMiddleware(middlewares.NewCORS())
builder.UseMiddleware(middlewares.NewRequestTracker())
builder.UseMiddleware(middlewares.NewPrometheus())
builder.UseMiddleware(middlewares.NewPprof())

// 自定义中间件
builder.UseMiddleware(func(next func(ctx *context.HttpContext)) func(ctx *context.HttpContext) {
    return func(ctx *context.HttpContext) {
        // 前置处理
        start := time.Now()
        
        next(ctx)  // 调用下一个
        
        // 后置处理
        log.Printf("request took %v", time.Since(start))
    }
})
```

**内置中间件**:

| 中间件 | 说明 |
|--------|------|
| `Recovery` | 异常恢复 |
| `Logger` | 请求日志 |
| `CORS` | 跨域支持 |
| `RequestTracker` | 请求追踪 |
| `Prometheus` | 监控指标 |
| `Pprof` | 性能分析 |
| `JWT` | JWT认证 |
| `StaticFile` | 静态文件 |

### MVC框架 (mvc/)

```go
// Controller定义
type UserController struct {
    *mvc.ApiController
    userService services.IUserService  // 依赖注入
}

func NewUserController(userService services.IUserService) *UserController {
    return &UserController{
        ApiController: mvc.NewApiController(),
        userService: userService,
    }
}

// Action方法 - 自动参数绑定
func (c *UserController) GetUser(ctx *context.HttpContext, req *GetUserRequest) mvc.ApiResult {
    user := c.userService.GetById(req.Id)
    return c.OK(user)
}

func (c *UserController) CreateUser(req *CreateUserRequest) mvc.ApiResult {
    user := c.userService.Create(req)
    return c.Created(user)
}

// 注册请求对象
type GetUserRequest struct {
    mvc.RequestGET `route:"/api/users/:id"`
    Id            int64 `param:"id"`
}

type CreateUserRequest struct {
    mvc.RequestBody `route:"/api/users" doc:"创建用户"`
    UserName       string `json:"username" doc:"用户名"`
    Password       string `json:"password" doc:"密码"`
    Email          string `json:"email" doc:"邮箱"`
}
```

**请求类型**:

| 类型 | 说明 |
|------|------|
| `mvc.RequestGET` | GET请求 |
| `mvc.RequestPOST` | POST请求 |
| `mvc.RequestPUT` | PUT请求 |
| `mvc.RequestDELETE` | DELETE请求 |
| `mvc.RequestPATCH` | PATCH请求 |
| `mvc.RequestBody` | 带请求体的请求 |

### 端点 (endpoints/)

```go
func registerEndpoints(rb router.IRouterBuilder) {
    // 健康检查
    endpoints.UseHealth(rb)
    
    // 可视化
    endpoints.UseViz(rb)
    
    // Prometheus监控
    endpoints.UsePrometheus(rb)
    
    // 性能分析
    endpoints.UsePprof(rb)
    
    // JWT端点
    endpoints.UseJwt(rb)
    
    // Swagger文档
    endpoints.UseSwaggerDoc(rb, swagger.Info{
        Title:       "API文档",
        Version:     "v1.0.0",
        Description: "API描述",
    }, func(openapi *swagger.OpenApi) {
        openapi.AddSecurityBearerAuth()
    })
}
```

### 会话管理 (session/)

```go
// 获取会话
session := ctx.GetSession()
session.Set("user_id", userId)
userId := session.Get("user_id")
session.Delete("user_id")

// 配置会话
builder.UseSession(session.Options{
    Name:   "yuanboot_session",
    MaxAge: 86400,
})
```

### WebSocket

```go
// WebSocket处理
ws := websocket.NewWebSocketUpgrader()
ws.Handle(func(conn *websocket.Conn) {
    for {
        message, err := conn.ReadMessage()
        if err != nil {
            break
        }
        conn.WriteMessage(message)
    }
})
```

---

## grpc gRPC框架

### 创建gRPC服务

```go
func main() {
    host := grpc.NewHostBuilder().
        UseConfiguration(config).
        Configure(func(builder *grpc.ServerBuilder) {
            builder.AddServices(registerServices)
            builder.AddUnaryInterceptor(grpcMiddleware.Logging)
            builder.AddInterceptor(grpcAuth.ServerInterceptor(authFunc))
        }).
        Build()
    
    host.Run()
}

func registerServices(builder *grpc.ServiceBuilder) {
    builder.AddService(func registrar grpc.ServiceRegistrar) {
        pb.RegisterUserServiceServer(registrar, &userServer{})
    })
}
```

### 拦截器

```go
// 一元拦截器
builder.AddUnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    // 前置处理
    err := authenticate(ctx)
    if err != nil {
        return nil, err
    }
    
    // 调用处理器
    resp, err := handler(ctx, req)
    
    // 后置处理
    logResponse(resp)
    
    return resp, err
})

// 流拦截器
builder.AddStreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
    return handler(srv, ss)
})
```

---

## pkg 功能包

### 缓存 (cache/redis/)

```go
import "github.com/liangboceo/yuanboot/pkg/cache/redis"

// 创建Redis客户端
client := redis.NewClient(redis.Options{
    Addr: "localhost:6379",
    Password: "",
    DB: 0,
})

// 基本操作
client.Set("key", "value", time.Hour)
val, _ := client.Get("key")
client.Del("key")

// 分布式锁
lock := redis.NewLock(client)
err, ok := lock.TryLock("mykey", 30*time.Second)
if ok {
    defer lock.Unlock("mykey")
    // 业务逻辑
}

// 等待锁
err, ok = lock.TryLockWithWait("mykey", 10*time.Second, 30*time.Second)

// 带上下文的锁
ctx, _ := lock.LockWithContext(context.Background(), "mykey", 30*time.Second)
defer ctx.Value("distributed_lock").(interface{ Unlock() }).Unlock()
```

### 数据源 (datasources/)

#### MySQL

```go
import "github.com/liangboceo/yuanboot/pkg/datasources/mysql"

// 配置YAML
/*
yuanboot:
  datasource:
    db:
      name: default
      url: tcp(localhost:3306)/mydb
      username: root
      password: password
      debug: true
      pool:
        initCap: 5
        maxCap: 30
        idletimeout: 300
    pool:
      max_open_conns: 100
      max_idle_conns: 10
      conn_max_lifetime: 3600
*/

// 使用数据源
ds := mysql.NewMysqlDataSource(config)
conn, put, err := ds.Open()
if err == nil {
    defer put()
    db := conn.(*sql.DB)
    // 执行SQL
}
```

#### Redis

```go
import "github.com/liangboceo/yuanboot/pkg/datasources/redis"

// 配置YAML
/*
yuanboot:
  datasource:
    redis:
      default:
        addr: localhost:6379
        password: ""
        db: 0
        pool_size: 10
*/

// 使用Redis
rds := redis.NewRedisDataSource(config)
client := rds.GetClient("default")
```

### 配置管理 (configuration/)

```go
import "github.com/liangboceo/yuanboot/pkg/configuration"

// 监听配置变化
watcher, _ := config.Watch("app.name", func(event *WatchEvent) {
    fmt.Printf("config changed: %s -> %s\n", event.OldValue, event.NewValue)
})
defer watcher.Stop()

// 刷新配置
config.RefreshAll()
```

### HTTP客户端 (httpclient/)

```go
import "github.com/liangboceo/yuanboot/pkg/httpclient"

// 创建客户端
client := httpclient.NewClient(httpclient.Options{
    Timeout: 30 * time.Second,
})

// GET请求
resp, _ := client.Get("https://api.example.com/users", nil)

// POST请求
resp, _ := client.Post("https://api.example.com/users", user)

// 设置请求头
req := &httpclient.Request{
    URL: "https://api.example.com/users",
    Method: "POST",
    Headers: map[string]string{
        "Authorization": "Bearer token",
    },
    Body: userData,
}
resp, _ := client.Do(req)
```

### 邮件服务 (email/)

```go
import "github.com/liangboceo/yuanboot/pkg/email"

// 发送邮件
emailClient := email.NewClient(email.Options{
    Host: "smtp.example.com",
    Port: 587,
    Username: "user@example.com",
    Password: "password",
})

err := emailClient.Send(&email.Message{
    From:    "user@example.com",
    To:      []string{"recipient@example.com"},
    Subject: "Test Email",
    Body:    "Hello World",
})
```

### 消息队列 (mq/)

```go
import "github.com/liangboceo/yuanboot/pkg/mq"

// Kafka生产者
producer, _ := mq.NewKafkaProducer(config)
producer.Send("topic", "message")

// Kafka消费者
consumer, _ := mq.NewKafkaConsumer(config, "topic", "group")
consumer.Consume(func(msg *mq.Message) {
    fmt.Println(string(msg.Body))
})
```

### 定时任务 (scheduler/)

```go
import "github.com/liangboceo/yuanboot/pkg/scheduler"

// 创建调度器
sched := scheduler.NewScheduler()

// 添加定时任务
sched.AddJob("* * * * *", func() {
    fmt.Println("每分钟执行")
})

sched.AddJob("0 2 * * *", func() {
    fmt.Println("每天凌晨2点执行")
})

sched.Start()
```

### Swagger文档 (swagger/)

```go
import "github.com/liangboceo/yuanboot/pkg/swagger"

// 定义Swagger信息
swaggerInfo := swagger.Info{
    Title:       "API Documentation",
    Version:     "1.0.0",
    Description: "API Description",
    Contact: swagger.Contact{
        Email: "support@example.com",
        Name:  "API Support",
    },
    License: swagger.License{
        Name: "MIT",
        URL:  "https://opensource.org/licenses/MIT",
    },
}

// 注册文档
swagger.UseSwaggerDoc(rb, swaggerInfo)

// 添加安全认证
func(openapi *swagger.OpenApi) {
    openapi.AddSecurityBearerAuth()
    openapi.AddSecurityApiKeyAuth("Authorization", "header")
}
```

---

## utils 工具函数

### JWT认证

```go
import "github.com/liangboceo/yuanboot/utils/jwt"

// 创建Token
claims := jwt.MapClaims{
    "user_id": 123,
    "username": "admin",
    "exp": time.Now().Add(24 * time.Hour).Unix(),
}

token, _ := jwt.NewToken(claims, "secret")
tokenString, _ := token.SignedString()

// 验证Token
parsedToken, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    return []byte("secret"), nil
})

claims := parsedToken.Claims.(jwt.MapClaims)
userId := claims["user_id"].(float64)
```

### 类型转换

```go
import "github.com/liangboceo/yuanboot/utils/cast"

// 转字符串
str := cast.ToString(123)
str = cast.ToStringE(123)  // 带错误处理

// 转Int
num := cast.ToInt("123")
num = cast.ToIntE("123")

// 转Float
f := cast.ToFloat64("3.14")

// 转Bool
b := cast.ToBool("true")
```

---

## 配置文件示例

### config_dev.yml

```yaml
yuanboot:
  app:
    name: myapp
    port: 8080
    profile: dev
    debug: true
  
  datasource:
    db:
      name: default
      url: tcp(localhost:3306)/mydb
      username: root
      password: password
      debug: true
      pool:
        initCap: 5
        maxCap: 30
        idletimeout: 300
    pool:
      max_open_conns: 100
      max_idle_conns: 10
      conn_max_lifetime: 3600
    redis:
      default:
        addr: localhost:6379
        password: ""
        db: 0
        pool_size: 10
  
  servicediscovery:
    nacos:
      enabled: true
      server_addr: localhost:8848
      namespace: public
      group: DEFAULT_GROUP
  
  mq:
    kafka:
      brokers:
        - localhost:9092
      topic: mytopic
  
  scheduler:
    enabled: true
  
  email:
    host: smtp.example.com
    port: 587
    username: user@example.com
    password: password
```

### config_prod.yml

```yaml
yuanboot:
  app:
    name: myapp
    port: 8080
    profile: prod
    debug: false
  
  datasource:
    db:
      name: default
      url: tcp(prod-mysql:3306)/mydb
      username: ${DB_USERNAME}
      password: ${DB_PASSWORD}
      debug: false
      pool:
        initCap: 10
        maxCap: 100
        idletimeout: 600
    pool:
      max_open_conns: 200
      max_idle_conns: 20
      conn_max_lifetime: 7200
    redis:
      default:
        addr: prod-redis:6379
        password: ${REDIS_PASSWORD}
        db: 0
        pool_size: 20
  
  servicediscovery:
    nacos:
      enabled: true
      server_addr: prod-nacos:8848
      namespace: public
      group: DEFAULT_GROUP
```

---

## CLI工具

使用yuanbootctl命令行工具创建项目：

```bash
# 创建标准Web项目
yuanbootctl new web mywebapp

# 创建gRPC项目
yuanbootctl new grpc mygrpcapp

# 创建MVC项目
yuanbootctl new mvc mymvcapp

# 创建微服务项目
yuanbootctl new microservice mymicroservice

# 查看帮助
yuanbootctl help
```

---

## 最佳实践

### 项目结构

```
myapp/
├── cmd/
│   └── server/
│       └── main.go          # 入口文件
├── config/
│   ├── config_dev.yml       # 开发配置
│   └── config_prod.yml      # 生产配置
├── internal/
│   ├── controller/          # 控制器
│   ├── service/             # 服务层
│   ├── repository/          # 数据访问层
│   ├── model/               # 数据模型
│   └── middleware/          # 中间件
├── pkg/                     # 内部包
├── api/                     # API定义
│   └── proto/               # Proto文件
├── go.mod
├── go.sum
└── Makefile
```

### 依赖注入

```go
// 使用构造器注入
func NewUserService(repo repository.IUserRepository) *UserService {
    return &UserService{repo: repo}
}

// 接口编程
type IUserService interface {
    GetById(id int64) (*User, error)
    Create(user *User) (*User, error)
}

// 注册服务
ConfigureServices(func(sc *dependencyinjection.ServiceCollection) {
    sc.AddTransientByImplements(
        services.NewUserService, 
        new(services.IUserService),
    )
})
```

### 错误处理

```go
// 使用统一错误响应
func (c *UserController) GetUser(req *GetUserRequest) mvc.ApiResult {
    user, err := c.userService.GetById(req.Id)
    if err != nil {
        return c.Error(err)  // 自动转换为API错误响应
    }
    return c.OK(user)
}
```

---

## 版本历史

- {{.Version}}- 当前版本
- 完整版本历史请查看 [CHANGELOG](https://github.com/liangboceo/yuanboot/releases)
