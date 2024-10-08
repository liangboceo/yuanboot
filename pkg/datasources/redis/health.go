package redis

import (
	"github.com/liangboceo/yuanboot/abstractions/health"
	"github.com/liangboceo/yuanboot/pkg/cache/redis"
)

type RedisHealthIndicator struct {
	redis redis.IClient
}

func NewRedisHealthIndicator(client redis.IClient) *RedisHealthIndicator {
	return &RedisHealthIndicator{redis: client}
}

func (h *RedisHealthIndicator) Health() health.ComponentStatus {
	status := health.Up("redisHealth")
	val, err := h.redis.Ping()
	if val != "PONG" {
		status.SetStatus("down")
	}
	if err != nil {
		val = err.Error()
	}
	return status.WithDetail("message", val)
}
