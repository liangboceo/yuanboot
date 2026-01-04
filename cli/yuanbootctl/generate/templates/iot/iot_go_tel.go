package iot

const DemoController_Tel = `
package {{.CurrentModelName}}

import (
	"github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/mvc"
)

type DemoController struct {
	mvc.ApiController // 必须继承
}

func NewDemoController() *DemoController {
	return &DemoController{}
}

//-------------------------------------------------------------------------------
type RegisterRequest struct {
	mvc.RequestBody
	UserName string` + " `param:\"UserName\"`\n" +
	"    Password string" + "`param:\"Password\"`\n" +
	`}

//GET URL  http://localhost:8080/app/v1/demo/register?UserName=max&Password=123
func (controller DemoController) Register(ctx *context.HttpContext, request *RegisterRequest) mvc.ApiResult {
	return mvc.ApiResult{Success: true, Message: "ok", Data: request}
}

//GET URL http://localhost:8080/app/v1/demo/getinfo
func (controller DemoController) GetInfo() mvc.ApiResult {
	return controller.OK("ok")
}

`

const Main_Tel = `
package {{.CurrentModelName}}

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
	"{{.ModelName}}/controller"
	"{{.ModelName}}/iothub"
	"{{.ModelName}}/pkg/mq"
	"{{.ModelName}}/pkg/service"
	"{{.ModelName}}/version"
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

//* Create the builder of Web host
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
`
const Mod_Tel = `

module {{.ModelName}}


go 1.16

require (
	github.com/liangboceo/dependencyinjection v1.0.0
	github.com/liangboceo/yuanboot {{.Version}}
	github.com/segmentio/kafka-go v0.4.49
)
`

const Config_Tel = `
yuanboot:
  application:
    name: {{.ModelName}}
    metadata: "develop"
    server:
      type: "fasthttp"
      address: ":8080"
      path: "app"
      max_request_size: 2096157
      session:
        name: "yuanboot_SESSIONID"
        timeout: 3600
      tls:
        cert: ""
        key: ""
      mvc:
        template: "v1/{controller}/{action}"
        views:
          path: "./static/templates"
          includes: [ "","" ]
      static:
        patten: "/"
        webroot: "./static"
      jwt:
        header: "Authorization"
        secret: "12391JdeOW^%$#@"
        prefix: "Bearer"
        expires: 3
        enable: false
        skip_path: [
            "/info",
            "/v1/user/GetInfo",
            "/v1/user/GetSD"
        ]
      cors:
        allow_origins: ["*"]
        allow_methods: ["POST","GET","PUT", "PATCH"]
        allow_credentials: true
  datasource:
    redis:
      name: reids1
      url: 10.5.215.33:6379
      password: rhzl@2014
      db: 11
    pool:
      init_cap: 2
      max_cap: 5
      idle_timeout: 5
  kafka:
    brokers: 10.5.215.34:9092,10.5.215.34:9093,10.5.215.34:9094
    topic: default_topic
    groupId: default_group
`

const Log_Config_Tel = `
yuanboot:
    log:
      log_level: info
      app_name: {{.ModelName}}
      log_path: /mnt/data/log/platform/
`

const Version_Tel = `
package {{.CurrentModelName}}

const version = "1.0.0"

var (
	env = "dev"
)

// Version return the version string
func Version() string {
	return version
}

// Env return the env string
func Env() string {
	return env
}


`

const Kafka_Tel = `
package {{.CurrentModelName}}

import (
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	kafka "github.com/segmentio/kafka-go"
	"strings"
	"sync"
	"time"
)

type Kafka struct {
	Brokers string` + " `mapstructure:\"brokers\" config:\"brokers\" `\n" +
	"    Topic string" + "`mapstructure:\"topic\" config:\"topic\" `\n" +
	"    GroupID string" + "`mapstructure:\"groupId\" config:\"groupId\" `\n" +
	`
}

type KafkaClient struct {
	config        *Kafka
	log           xlog.ILogger
	producerPool  sync.Map
	once          sync.Once
	cleanupTicker *time.Ticker
}

func NewKafkaClient(configuration abstractions.IConfiguration) *KafkaClient {
	kafkaConfig := &Kafka{}
	configuration.GetConfigObject("yuanboot.kafka", kafkaConfig)
	client := &KafkaClient{
		config: kafkaConfig,
		log:    xlog.GetXLogger("KafkaClient"),
	}
	return client
}

func NewTestKafkaClient(broker string) *KafkaClient {
	kafkaConfig := &Kafka{Brokers: broker}
	client := &KafkaClient{
		config: kafkaConfig,
		log:    xlog.GetXLogger("KafkaClient"),
	}
	return client
}

func (k *KafkaClient) GetProducerWithBroker(topic string, brokers string) *kafka.Writer {
	if topic == "" {
		topic = k.config.Topic
	}
	if brokers == "" {
		brokers = k.config.Brokers
	}
	k.once.Do(func() {
		k.cleanupTicker = time.NewTicker(5 * time.Minute)
		go k.cleanupProducers()
	})
	key := topic + brokers
	if writer, ok := k.producerPool.Load(key); ok {
		return writer.(*kafka.Writer)
	}
	newWriter := &kafka.Writer{
		Addr:                   kafka.TCP(strings.Split(brokers, ",")...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		Async:                  true,
		AllowAutoTopicCreation: true,
	}
	k.producerPool.Store(key, newWriter)
	return newWriter
}

func (k *KafkaClient) GetProducer(topic string) *kafka.Writer {
	if topic == "" {
		topic = k.config.Topic
	}

	k.once.Do(func() {
		k.cleanupTicker = time.NewTicker(5 * time.Minute)
		go k.cleanupProducers()
	})

	if writer, ok := k.producerPool.Load(topic); ok {
		return writer.(*kafka.Writer)
	}

	newWriter := &kafka.Writer{
		Addr:                   kafka.TCP(strings.Split(k.config.Brokers, ",")...),
		Topic:                  topic,
		Async:                  true,
		AllowAutoTopicCreation: true,
	}
	k.producerPool.Store(topic, newWriter)
	return newWriter
}

func (k *KafkaClient) cleanupProducers() {
	for range k.cleanupTicker.C {
		k.producerPool.Range(func(key, value interface{}) bool {
			writer := value.(*kafka.Writer)
			if writer.Stats().Writes == 0 {
				_ = writer.Close()
				k.producerPool.Delete(key)
			}
			return true
		})
	}
}

// DestroyProducers 主动销毁生产者
func (k *KafkaClient) DestroyProducers(topic string, brokers string) {
	if topic == "" {
		topic = k.config.Topic
	}
	if brokers == "" {
		brokers = k.config.Brokers
	}
	key := topic + brokers
	if writer, ok := k.producerPool.Load(key); ok {
		_ = writer.(*kafka.Writer).Close()
		k.producerPool.Delete(key)
	}
}

func (k *KafkaClient) GetAckProducer(topic string) *kafka.Writer {
	if topic == "" {
		topic = k.config.Topic
	}

	k.once.Do(func() {
		k.cleanupTicker = time.NewTicker(5 * time.Minute)
		go k.cleanupProducers()
	})

	if writer, ok := k.producerPool.Load(topic); ok {
		return writer.(*kafka.Writer)
	}
	newWriter := &kafka.Writer{
		Addr:                   kafka.TCP(strings.Split(k.config.Brokers, ",")...),
		Topic:                  topic,
		RequiredAcks:           kafka.RequireAll, // ack模式
		Async:                  true,             // 异步
		AllowAutoTopicCreation: true,             //自动创建topi
	}
	k.producerPool.Store(topic, newWriter)
	return newWriter
}

func (k *KafkaClient) GetConsumerWithBroker(brokers string, topics string, groupId string) *kafka.Reader {
	if topics == "" {
		topics = k.config.Topic
	}
	if groupId == "" {
		groupId = k.config.GroupID
	}
	if brokers == "" {
		brokers = k.config.Brokers
	}
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        strings.Split(brokers, ","),
		GroupID:        groupId,
		GroupTopics:    strings.Split(topics, ","),
		CommitInterval: time.Second, // 每秒刷新一次提交给 Kafka
		MinBytes:       1e3,         // 1KB（降低最小读取阈值）
		MaxBytes:       1e5,         // 100KB（限制最大批量）
		QueueCapacity:  10000,       // 增加内部队列容量
		StartOffset:    kafka.LastOffset,
		MaxWait:        100 * time.Millisecond, // 缩短轮询等待时间
	})
}

func (k *KafkaClient) GetConsumer(topic string, groupId string) *kafka.Reader {
	if topic == "" {
		topic = k.config.Topic
	}
	if groupId == "" {
		groupId = k.config.GroupID
	}
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        strings.Split(k.config.Brokers, ","),
		GroupID:        groupId,
		Topic:          topic,
		CommitInterval: time.Second, // 每秒刷新一次提交给 Kafka
		MinBytes:       1e3,         // 1KB（降低最小读取阈值）
		MaxBytes:       1e5,         // 100KB（限制最大批量）
		QueueCapacity:  10000,       // 增加内部队列容量
		StartOffset:    kafka.LastOffset,
		MaxWait:        100 * time.Millisecond, // 缩短轮询等待时间
	})
}
`
const Cache_Tel = `
package {{.CurrentModelName}}

import (
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"github.com/liangboceo/yuanboot/pkg/cache/redis"
	redisdb "github.com/liangboceo/yuanboot/pkg/datasources/redis"
)

type CacheService struct {
	RedisClient redis.IClient
	Log         xlog.ILogger
}

func NewCacheService(redisDataSource *redisdb.RedisDataSource) *CacheService {
	conn, _, _ := redisDataSource.Open()
	client := conn.(redis.IClient)
	log := xlog.GetXLogger("CacheService")
	return &CacheService{RedisClient: client, Log: log}
}

`

const IotHub_Tel = `
package {{.CurrentModelName}}

import (
	"context"
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/web/middlewares"
	"github.com/segmentio/kafka-go"
	"{{.ModelName}}/pkg/mq"
	"{{.ModelName}}/pkg/service"
)

type Hub struct {
	abstractions.IHostService
	CacheService *service.CacheService
	KafkaClient  *mq.KafkaClient
	Config       abstractions.IConfiguration
	Log          *middlewares.Logger
	KafkaReader  *kafka.Reader
}

func InitIotHub(cacheService *service.CacheService, kafkaClient *mq.KafkaClient, configuration abstractions.IConfiguration) *Hub {
	return &Hub{CacheService: cacheService, KafkaClient: kafkaClient, Config: configuration, Log: middlewares.NewLogger(), KafkaReader: kafkaClient.GetConsumer("test_topic", "test_topic_group")}
}

func (hub *Hub) Run() error {
	//监听kafka用于更新程序
	go hub.ListenKafka()
	return nil
}

func (hub *Hub) ListenKafka() {
	ctx := context.Background()
	for {
		// 从 Kafka 中获取下一条消息
		m, err := hub.KafkaReader.ReadMessage(ctx)
		go func() {
			defer hub.ErrorHandle()
			if err != nil {
				hub.Log.ALogger.Error("获取消息时出错: %v", err)
				return
			}
			hub.Log.ALogger.Info("<UNK>: %s", string(m.Value))
		}()
	}
}

func (hub *Hub) Stop() error {
	return nil
}

func (hub *Hub) ErrorHandle() {
	if e := recover(); e != nil {
		hub.Log.ALogger.Error("", e)
		return
	}
}

`
