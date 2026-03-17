package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type GoRedisClusterOps struct {
	GoRedisStandaloneOps
}

func NewClusterOps(options *Options) *GoRedisClusterOps {
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    options.Addrs,
		Password: options.Password,
	})
	return &GoRedisClusterOps{GoRedisStandaloneOps{client: client}}
}

// Publish posts a message to the channel
func (ops *GoRedisClusterOps) Publish(channel string, message interface{}) (int64, error) {
	return ops.client.(*redis.ClusterClient).Publish(context.Background(), channel, message).Result()
}

// Subscribe subscribes the client to the specified channels
func (ops *GoRedisClusterOps) Subscribe(channels ...string) (*Subscription, error) {
	ps := ops.client.(*redis.ClusterClient).Subscribe(context.Background(), channels...)
	return &Subscription{pubsub: ps}, nil
}

// PSubscribe subscribes the client to the given patterns
func (ops *GoRedisClusterOps) PSubscribe(patterns ...string) (*Subscription, error) {
	ps := ops.client.(*redis.ClusterClient).PSubscribe(context.Background(), patterns...)
	return &Subscription{pubsub: ps}, nil
}
