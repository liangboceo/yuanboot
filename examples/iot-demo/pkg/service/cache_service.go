package service

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
