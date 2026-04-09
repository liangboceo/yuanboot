package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Pipeline 封装 Redis Pipeline 操作，支持批量执行命令
// v9 的 Pipeline 性能相比 v8 提升了约 38%
type Pipeline interface {
	// Set 批量设置值
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	// Get 批量获取值
	Get(key string) *redis.StringCmd

	// SetNX 批量设置值（仅当键不存在）
	SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd

	// Del 批量删除键
	Del(keys ...string) *redis.IntCmd

	// Incr 批量递增
	Incr(key string) *redis.IntCmd

	// IncrBy 批量递增指定步长
	IncrBy(key string, value int64) *redis.IntCmd

	// HSet 批量设置哈希字段
	HSet(key string, field string, value interface{}) *redis.IntCmd

	// HGet 批量获取哈希字段
	HGet(key string, field string) *redis.StringCmd

	// ZAdd 批量添加有序集合成员
	ZAdd(key string, members ...redis.Z) *redis.IntCmd

	// SAdd 批量添加无序集合成员
	SAdd(key string, members ...interface{}) *redis.IntCmd

	// LPush 批量列表左推入
	LPush(key string, values ...interface{}) *redis.IntCmd

	// RPush 批量列表右推入
	RPush(key string, values ...interface{}) *redis.IntCmd

	// MGet 批量获取多个键
	MGet(keys ...string) *redis.SliceCmd

	// MSet 批量设置多个键
	MSet(pairs ...interface{}) *redis.StatusCmd

	// Exec 执行所有命令
	Exec(ctx context.Context) ([]redis.Cmder, error)

	// Discard 丢弃 Pipeline
	Discard()
}

// GoRedisPipeline 实现 Pipeline 接口
type GoRedisPipeline struct {
	client     redis.Pipeliner
	serializer ISerializer
}

// NewPipeline 创建新的 Pipeline
func NewPipeline(client redis.Pipeliner, serializer ISerializer) Pipeline {
	return &GoRedisPipeline{client: client, serializer: serializer}
}

func (p *GoRedisPipeline) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	ss, _ := p.serializer.Serialization(value)
	return p.client.Set(context.Background(), key, ss, expiration)
}

func (p *GoRedisPipeline) Get(key string) *redis.StringCmd {
	return p.client.Get(context.Background(), key)
}

func (p *GoRedisPipeline) SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	return p.client.SetNX(context.Background(), key, value, expiration)
}

func (p *GoRedisPipeline) Del(keys ...string) *redis.IntCmd {
	return p.client.Del(context.Background(), keys...)
}

func (p *GoRedisPipeline) Incr(key string) *redis.IntCmd {
	return p.client.Incr(context.Background(), key)
}

func (p *GoRedisPipeline) IncrBy(key string, value int64) *redis.IntCmd {
	return p.client.IncrBy(context.Background(), key, value)
}

func (p *GoRedisPipeline) HSet(key string, field string, value interface{}) *redis.IntCmd {
	ss, _ := p.serializer.Serialization(value)
	return p.client.HSet(context.Background(), key, field, ss)
}

func (p *GoRedisPipeline) HGet(key string, field string) *redis.StringCmd {
	return p.client.HGet(context.Background(), key, field)
}

func (p *GoRedisPipeline) ZAdd(key string, members ...redis.Z) *redis.IntCmd {
	return p.client.ZAdd(context.Background(), key, members...)
}

func (p *GoRedisPipeline) SAdd(key string, members ...interface{}) *redis.IntCmd {
	return p.client.SAdd(context.Background(), key, members...)
}

func (p *GoRedisPipeline) LPush(key string, values ...interface{}) *redis.IntCmd {
	return p.client.LPush(context.Background(), key, values...)
}

func (p *GoRedisPipeline) RPush(key string, values ...interface{}) *redis.IntCmd {
	return p.client.RPush(context.Background(), key, values...)
}

func (p *GoRedisPipeline) MGet(keys ...string) *redis.SliceCmd {
	return p.client.MGet(context.Background(), keys...)
}

func (p *GoRedisPipeline) MSet(pairs ...interface{}) *redis.StatusCmd {
	return p.client.MSet(context.Background(), pairs...)
}

func (p *GoRedisPipeline) Exec(ctx context.Context) ([]redis.Cmder, error) {
	return p.client.Exec(ctx)
}

func (p *GoRedisPipeline) Discard() {
	p.client.Discard()
}
