package service

import (
	"fmt"
	"strconv"

	"github.com/liangboceo/yuanboot/abstractions/xlog"
	"github.com/liangboceo/yuanboot/pkg/cache/redis"
	redisdb "github.com/liangboceo/yuanboot/pkg/datasources/redis"
	"time"
)

type CacheService struct {
	redisClient redis.IClient
	Log         xlog.ILogger
}

func NewCacheService(redisDataSource *redisdb.RedisDataSource) *CacheService {
	conn, _, _ := redisDataSource.Open()
	client := conn.(redis.IClient)
	log := xlog.GetXLogger("CacheService")
	return &CacheService{redisClient: client, Log: log}
}

const (
	DATALOCK         = "sendex:data:lock:%s:%s:%s"
	SYSUSERLOGININFO = "sendex:sys:user:login:%s"
	SYSCAPTCHA       = "sendex:sys:captcha:%s"
	SYSTEM_CONFIG    = "sendex:sys:config:%s"

	// PubSub channels
	WorkOrderEventChannel        = "sendex:workorder:event"
	DiagnosisReportEventChannel  = "sendex:evaluation:diagnosis:report:event"
	SessionStreamChannel         = "sendex:session:stream:%d"
	AgentStatusChannel           = "sendex:agent:status:event"
	DynamicApiUpdateChannel      = "sendex:dynamicapi:update:event"
	RedisDataSourceUpdateChannel = "sendex:redis-datasource:update:event"
	SkillCompileChannel          = "sendex:skill:compile:event"
)

func (cache *CacheService) GetClient() redis.IClient {
	return cache.redisClient
}

// ==================== Pub/Sub 通知 ====================

func (cache *CacheService) PublishDynamicApiUpdate(apiID uint) (int64, error) {
	return cache.redisClient.GetPubSubOps().Publish(DynamicApiUpdateChannel, strconv.FormatUint(uint64(apiID), 10))
}

func (cache *CacheService) SubscribeDynamicApiUpdate(msgChan chan<- string) (*redis.Subscription, error) {
	proxyChan := make(chan *redis.Message, 100)
	sub, err := cache.redisClient.GetPubSubOps().ReceiveMessages(DynamicApiUpdateChannel, proxyChan)
	if err != nil {
		return nil, err
	}
	go func() {
		for msg := range proxyChan {
			msgChan <- msg.Payload
		}
	}()
	return sub, nil
}

func (cache *CacheService) PublishRedisDataSourceUpdate(resourceKey string) (int64, error) {
	return cache.redisClient.GetPubSubOps().Publish(RedisDataSourceUpdateChannel, resourceKey)
}

func (cache *CacheService) SubscribeRedisDataSourceUpdate(msgChan chan<- string) (*redis.Subscription, error) {
	proxyChan := make(chan *redis.Message, 100)
	sub, err := cache.redisClient.GetPubSubOps().ReceiveMessages(RedisDataSourceUpdateChannel, proxyChan)
	if err != nil {
		return nil, err
	}
	go func() {
		for msg := range proxyChan {
			msgChan <- msg.Payload
		}
	}()
	return sub, nil
}

// PublishWorkOrderEvent 发布工单处理事件到 Redis Pub/Sub
func (cache *CacheService) PublishWorkOrderEvent(workOrderID uint) (int64, error) {
	return cache.redisClient.GetPubSubOps().Publish(WorkOrderEventChannel, strconv.FormatUint(uint64(workOrderID), 10))
}

// SubscribeWorkOrderEvent 订阅工单处理事件，消息会发送到 msgChan
// 返回的 Subscription 需要调用 Close() 来取消订阅
func (cache *CacheService) SubscribeWorkOrderEvent(msgChan chan<- string) (*redis.Subscription, error) {
	proxyChan := make(chan *redis.Message, 100)
	sub, err := cache.redisClient.GetPubSubOps().ReceiveMessages(WorkOrderEventChannel, proxyChan)
	if err != nil {
		return nil, err
	}
	go func() {
		for msg := range proxyChan {
			msgChan <- msg.Payload
		}
	}()
	return sub, nil
}

// PublishDiagnosisReportEvent 发布诊断报告生成事件到 Redis Pub/Sub
func (cache *CacheService) PublishDiagnosisReportEvent(submissionID uint) (int64, error) {
	return cache.redisClient.GetPubSubOps().Publish(DiagnosisReportEventChannel, strconv.FormatUint(uint64(submissionID), 10))
}

// SubscribeDiagnosisReportEvent 订阅诊断报告生成事件，消息会发送到 msgChan
// 返回的 Subscription 需要调用 Close() 来取消订阅
func (cache *CacheService) SubscribeDiagnosisReportEvent(msgChan chan<- string) (*redis.Subscription, error) {
	proxyChan := make(chan *redis.Message, 100)
	sub, err := cache.redisClient.GetPubSubOps().ReceiveMessages(DiagnosisReportEventChannel, proxyChan)
	if err != nil {
		return nil, err
	}
	go func() {
		for msg := range proxyChan {
			msgChan <- msg.Payload
		}
	}()
	return sub, nil
}

// ==================== 会话流式事件 Pub/Sub ====================

// PublishSessionEvent 发布会话流式事件到 Redis Pub/Sub（跨节点广播）
func (cache *CacheService) PublishSessionEvent(sessionID uint, payload string) (int64, error) {
	return cache.redisClient.GetPubSubOps().Publish(fmt.Sprintf(SessionStreamChannel, sessionID), payload)
}

// SubscribeSessionEvent 订阅会话流式事件，消息会发送到 msgChan
// 返回的 Subscription 需要调用 Close() 来取消订阅
func (cache *CacheService) PublishAgentStatusEvent(payload string) (int64, error) {
	return cache.redisClient.GetPubSubOps().Publish(AgentStatusChannel, payload)
}

func (cache *CacheService) SubscribeAgentStatusEvent(msgChan chan<- string) (*redis.Subscription, error) {
	proxyChan := make(chan *redis.Message, 256)
	sub, err := cache.redisClient.GetPubSubOps().ReceiveMessages(AgentStatusChannel, proxyChan)
	if err != nil {
		return nil, err
	}
	go func() {
		for msg := range proxyChan {
			msgChan <- msg.Payload
		}
	}()
	return sub, nil
}

func (cache *CacheService) SubscribeSessionEvent(sessionID uint, msgChan chan<- string) (*redis.Subscription, error) {
	proxyChan := make(chan *redis.Message, 256)
	sub, err := cache.redisClient.GetPubSubOps().ReceiveMessages(fmt.Sprintf(SessionStreamChannel, sessionID), proxyChan)
	if err != nil {
		return nil, err
	}
	go func() {
		for msg := range proxyChan {
			msgChan <- msg.Payload
		}
	}()
	return sub, nil
}

// ==================== 锁操作 ====================

func (cache *CacheService) Lock(name string, timeOut int) bool {
	ops := cache.redisClient.GetLockOps()
	err, flag := ops.GetDLock(name, timeOut)
	if err != nil {
		return false
	}
	return flag
}
func (cache *CacheService) Unlock(name string) bool {
	ops := cache.redisClient.GetLockOps()
	err, flag := ops.DisposeLock(name)
	if err != nil {
		return false
	}
	return flag
}

func (cache *CacheService) GetSysUserLoginInfo(username string, data interface{}) error {
	return cache.redisClient.GetKVOps().Get(fmt.Sprintf(SYSUSERLOGININFO, username), data)
}

func (cache *CacheService) PutSysUserLoginInfo(username string, data interface{}, duration time.Duration) error {
	return cache.redisClient.GetKVOps().Set(fmt.Sprintf(SYSUSERLOGININFO, username), data, duration)
}

func (cache *CacheService) DelSysUserLoginInfo(username string) bool {
	return cache.redisClient.Delete(fmt.Sprintf(SYSUSERLOGININFO, username))
}

func (cache *CacheService) GetCaptcha(captchaId string) (string, error) {
	return cache.redisClient.GetKVOps().GetString(fmt.Sprintf(SYSCAPTCHA, captchaId))
}

func (cache *CacheService) PutCaptcha(captchaId string, data interface{}, duration time.Duration) error {
	return cache.redisClient.GetKVOps().Set(fmt.Sprintf(SYSCAPTCHA, captchaId), data, duration)
}

func (cache *CacheService) DelCaptcha(captchaId string) bool {
	return cache.redisClient.Delete(fmt.Sprintf(SYSCAPTCHA, captchaId))
}

// ──────────── Generic cache helpers ────────────

// CacheGetJSON retrieves a JSON-serialized object from cache into dest.
// Returns true if found, false if miss/error.
func (cache *CacheService) CacheGetJSON(key string, dest interface{}) bool {
	if err := cache.redisClient.GetKVOps().Get(key, dest); err != nil {
		return false
	}
	return true
}

// CacheSetJSON stores an object as JSON in cache with the given TTL.
func (cache *CacheService) CacheSetJSON(key string, data interface{}, ttl time.Duration) error {
	return cache.redisClient.GetKVOps().Set(key, data, ttl)
}

// CacheDelete removes a key from cache.
func (cache *CacheService) CacheDelete(key string) bool {
	return cache.redisClient.Delete(key)
}

// CacheIncr atomically increments a key and sets its TTL.
// If the key does not exist, it is created with value 0 before incrementing (standard Redis INCR behavior).
func (cache *CacheService) CacheIncr(key string, ttl time.Duration) (int64, error) {
	val, err := cache.redisClient.GetKVOps().Increment(key, 1)
	if err != nil {
		return 0, err
	}
	// Set expiry (best-effort, non-critical if it fails)
	if ttl > 0 {
		_, _ = cache.redisClient.SetExpire(key, ttl)
	}
	return val, nil
}

// ==================== Skill Compile Pub/Sub ====================

// PublishSkillCompileEvent publishes a skill compile event to Redis Pub/Sub.
// All nodes (including this one) will pick up the event and recompile the skill.
func (cache *CacheService) PublishSkillCompileEvent(skillCode string) (int64, error) {
	return cache.redisClient.GetPubSubOps().Publish(SkillCompileChannel, skillCode)
}

// SubscribeSkillCompileEvent subscribes to skill compile events, messages are sent to msgChan.
// Returns a Subscription that should be closed via Close() to unsubscribe.
func (cache *CacheService) SubscribeSkillCompileEvent(msgChan chan<- string) (*redis.Subscription, error) {
	proxyChan := make(chan *redis.Message, 100)
	sub, err := cache.redisClient.GetPubSubOps().ReceiveMessages(SkillCompileChannel, proxyChan)
	if err != nil {
		return nil, err
	}
	go func() {
		for msg := range proxyChan {
			msgChan <- msg.Payload
		}
	}()
	return sub, nil
}
