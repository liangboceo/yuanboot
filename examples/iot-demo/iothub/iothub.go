package iothub

import (
	"context"
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/web/middlewares"
	"github.com/segmentio/kafka-go"
	"iot-demo/pkg/mq"
	"iot-demo/pkg/service"
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
