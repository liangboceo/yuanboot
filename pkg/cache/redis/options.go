package redis

import "time"

type Options struct {
	Addr     string
	Addrs    []string
	Password string
	DB       int
	// 连接池配置 (v9 优化)
	PoolSize        int           // 连接池大小，默认 10
	MinIdleConns    int           // 最小空闲连接数，默认 2
	MaxRetries      int           // 最大重试次数，默认 3
	MinRetryBackoff time.Duration // 最小重试间隔，默认 8ms
	MaxRetryBackoff time.Duration // 最大重试间隔，默认 512ms
	DialTimeout     time.Duration // 连接超时，默认 5s
	ReadTimeout     time.Duration // 读取超时，默认 3s
	WriteTimeout    time.Duration // 写入超时，默认 3s
	PoolTimeout     time.Duration // 连接池超时，默认 4s
	ConnMaxIdleTime time.Duration // 连接间隔时间，默认 1m
}
