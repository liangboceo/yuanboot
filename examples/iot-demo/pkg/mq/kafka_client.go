package mq

import (
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/abstractions/xlog"
	kafka "github.com/segmentio/kafka-go"
	"strings"
	"sync"
	"time"
)

type Kafka struct {
	Brokers string `mapstructure:"brokers" config:"brokers"`
	Topic   string `mapstructure:"topic" config:"topic"`
	GroupID string `mapstructure:"groupId" config:"groupId"`
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
