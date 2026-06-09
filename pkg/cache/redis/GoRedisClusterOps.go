package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type GoRedisClusterOps struct {
	GoRedisStandaloneOps
	clusterClient   *redis.ClusterClient
	valueSerializer ISerializer
}

func NewClusterOps(options *Options, serializer ISerializer) *GoRedisClusterOps {
	if options == nil {
		options = &Options{}
	}

	// 设置默认值
	if options.PoolSize == 0 {
		options.PoolSize = 10
	}
	if options.MinIdleConns == 0 {
		options.MinIdleConns = 2
	}
	if options.MaxRetries == 0 {
		options.MaxRetries = 3
	}
	if options.MinRetryBackoff == 0 {
		options.MinRetryBackoff = 8 * time.Millisecond
	}
	if options.MaxRetryBackoff == 0 {
		options.MaxRetryBackoff = 512 * time.Millisecond
	}
	if options.DialTimeout == 0 {
		options.DialTimeout = 5 * time.Second
	}
	if options.ReadTimeout == 0 {
		options.ReadTimeout = 3 * time.Second
	}
	if options.WriteTimeout == 0 {
		options.WriteTimeout = 3 * time.Second
	}
	if options.PoolTimeout == 0 {
		options.PoolTimeout = 4 * time.Second
	}
	if options.ConnMaxIdleTime == 0 {
		options.ConnMaxIdleTime = 1 * time.Minute
	}

	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:           options.Addrs,
		Password:        options.Password,
		PoolSize:        options.PoolSize,
		MinIdleConns:    options.MinIdleConns,
		MaxRetries:      options.MaxRetries,
		MinRetryBackoff: options.MinRetryBackoff,
		MaxRetryBackoff: options.MaxRetryBackoff,
		DialTimeout:     options.DialTimeout,
		ReadTimeout:     options.ReadTimeout,
		WriteTimeout:    options.WriteTimeout,
		PoolTimeout:     options.PoolTimeout,
		ConnMaxIdleTime: options.ConnMaxIdleTime,
	})
	return &GoRedisClusterOps{
		GoRedisStandaloneOps: GoRedisStandaloneOps{
			client:          client,
			valueSerializer: serializer,
		},
		clusterClient: client,
	}
}

// Close closes the cluster client
func (ops *GoRedisClusterOps) Close() error {
	return ops.clusterClient.Close()
}

// Publish posts a message to the channel
func (ops *GoRedisClusterOps) Publish(channel string, message interface{}) (int64, error) {
	return ops.clusterClient.Publish(context.Background(), channel, message).Result()
}

// Subscribe subscribes the client to the specified channels
func (ops *GoRedisClusterOps) Subscribe(channels ...string) (*Subscription, error) {
	ps := ops.clusterClient.Subscribe(context.Background(), channels...)
	return &Subscription{pubsub: ps}, nil
}

// PSubscribe subscribes the client to the given patterns
func (ops *GoRedisClusterOps) PSubscribe(patterns ...string) (*Subscription, error) {
	ps := ops.clusterClient.PSubscribe(context.Background(), patterns...)
	return &Subscription{pubsub: ps}, nil
}

// Pipeline creates a pipeline for batch commands execution
func (ops *GoRedisClusterOps) Pipeline() Pipeline {
	return NewPipeline(ops.clusterClient.Pipeline(), ops.valueSerializer)
}

// TxPipeline creates a transaction pipeline
func (ops *GoRedisClusterOps) TxPipeline() Pipeline {
	return NewPipeline(ops.clusterClient.TxPipeline(), ops.valueSerializer)
}
