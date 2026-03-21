package redis

import (
	"crypto/md5"
	"fmt"
	"math"
	"time"
)

type Client struct {
	//------- data struct -----------
	kv         KV
	list       List
	hash       Hash
	set        Set
	zset       ZSet
	geo        Geo
	lock       *Lock
	pubsub     PubSub
	pipeline   Pipeline
	txPipeline Pipeline
	//-----------------------
	ops             Ops
	valueSerializer ISerializer
}

// NewClient new redis client by Options
func NewClient(options *Options) IClient {
	var ops Ops
	if options.Addrs == nil {
		ops = NewStandaloneOps(options)
	} else {
		ops = NewClusterOps(options)
	}

	kv := KV{ops: ops}
	list := List{ops: ops}
	hash := Hash{ops: ops}
	set := Set{ops: ops}
	zset := ZSet{ops: ops}
	geo := Geo{ops: ops}
	lock := NewLock(ops)
	pubsub := PubSub{ops: ops}
	pipeline := ops.Pipeline()
	txPipeline := ops.TxPipeline()

	return &Client{ops: ops, valueSerializer: DefaultSerializer, kv: kv, list: list, hash: hash, set: set, zset: zset, geo: geo, lock: lock, pubsub: pubsub, pipeline: pipeline, txPipeline: txPipeline}
}

/*
	getBucketId

https://www.cnblogs.com/colorfulkoala/p/5783556.html
过程变化简单描述为：get(key1) -> hget(md5(key1), key1) 从而得到value1。
如果我们通过预先计算，让很多key可以在BucketId空间里碰撞，那么可以认为一个BucketId下面挂了多个key。比如平均每个BucketId下面挂10个key，那么理论上我们将会减少超过90%的redis key的个数。
具体实现起来有一些麻烦，而且用这个方法之前你要想好容量规模。我们通常使用的md5是32位的hexString（16进制字符），它的空间是128bit，这个量级太大了，我们需要存储的是百亿级，大约是33bit，所以我们需要有一种机制计算出合适位数的散列，而且为了节约内存，我们需要利用全部字符类型（ASCII码在0~127之间）来填充，而不用HexString，这样Key的长度可以缩短到一半。
*/
func GetBucketId(key []byte, bit int) []byte {
	mdBytes := md5.Sum(key)
	md5str1 := fmt.Sprintf("%x", mdBytes) //将[]byte转成16进制
	md := []byte(md5str1)
	r := make([]byte, (bit-1)/7+1)
	a := byte(math.Pow(2, float64(bit%7)) - 2)
	md[len(r)-1] = byte(md[len(r)-1] & a)
	copy(r, md)
	for i := 0; i < len(r); i++ {
		r[i] &= 127
	}
	return r
}

// SetSerializer set value serializer
func (c *Client) SetSerializer(serializer ISerializer) {
	c.valueSerializer = serializer
}

// WrapClient 包装从连接池获取的 Redis 客户端
// 这个方法用于将数据源连接池中获取的客户端包装为统一的 IClient 接口
func WrapClient(client IClient) IClient {
	return client
}

// NewClientFromOps 使用已有的 Ops 创建客户端
// 这个方法允许使用外部创建的 Ops（如从连接池获取的）
func NewClientFromOps(ops Ops) IClient {
	kv := KV{ops: ops}
	list := List{ops: ops}
	hash := Hash{ops: ops}
	set := Set{ops: ops}
	zset := ZSet{ops: ops}
	geo := Geo{ops: ops}
	lock := NewLock(ops)
	pubsub := PubSub{ops: ops}

	return &Client{ops: ops, valueSerializer: DefaultSerializer, kv: kv, list: list, hash: hash, set: set, zset: zset, geo: geo, lock: lock, pubsub: pubsub}
}

// GetKVOps Returns the operations performed on simple values (or Strings in Redis terminology).
func (c *Client) GetKVOps() KV {
	c.kv.serializer = c.valueSerializer
	return c.kv
}

// GetListOps Returns the operations performed on list values.
func (c *Client) GetListOps() List {
	c.list.serializer = c.valueSerializer
	return c.list
}

// GetHashOps Returns the operations performed on hash values.
func (c *Client) GetHashOps() Hash {
	c.hash.serializer = c.valueSerializer
	return c.hash
}

// GetSetOps Returns the operations performed on set values.
func (c *Client) GetSetOps() Set {
	return c.set
}

// GetZSetOps Returns the operations performed on zset values (also known as sorted sets).
func (c *Client) GetZSetOps() ZSet {
	return c.zset
}

// GetGeoOps Geo Returns the operations performed on geo values (also known GIS system).
func (c *Client) GetGeoOps() Geo {
	return c.geo
}

// GetLockOps Returns the operations performed on locker values.
func (c *Client) GetLockOps() *Lock {
	return c.lock
}

// GetPubSubOps Returns the operations performed on publish/subscribe values.
func (c *Client) GetPubSubOps() PubSub {
	return c.pubsub
}

// GetPipeLineOps creates a pipeline for batch commands execution
func (c *Client) GetPipeLineOps() Pipeline {
	return c.pipeline
}

// GetTxPipeLineOps creates a transaction  pipeline for batch commands execution
func (c *Client) GetTxPipeLineOps() Pipeline {
	return c.txPipeline
}

// Ping return PONG and error
func (c *Client) Ping() (string, error) {
	return c.ops.Ping()
}

// SetExpire cmd by expire
func (c *Client) SetExpire(key string, expiration time.Duration) (bool, error) {
	return c.ops.SetExpire(key, expiration)
}

// SetExpire cmd by TTL
func (c *Client) GetExpire(key string) (time.Duration, error) {
	return c.ops.TTL(key)
}

// Delete delete the key, cmd by del
func (c *Client) Delete(key string) bool {
	n, _ := c.ops.DeleteKey(key)
	return n > 0
}

// HasKey exists the key, cmd by exists
func (c *Client) HasKey(key string) bool {
	h, _ := c.ops.Exists(key)
	return h
}

// RandomKey return random Key for db
func (c *Client) RandomKey() (string, error) {
	return c.ops.RandomKey()
}

// RandomKey return random Key for db
func (c *Client) Close() error {
	return c.ops.Close()
}

func (c *Client) ListKeys(page uint64, pattern string, pageSize int64) ([]string, int, error) {
	return c.ops.ListKeys(page, pattern, pageSize)
}
func (c *Client) Info() (string, error) {
	return c.ops.Info()
}

// GetPipeline creates a  pipeline
func (c *Client) GetPipeline() Pipeline {
	return c.ops.Pipeline()
}

// GetTxPipeline creates a transaction pipeline
func (c *Client) GetTxPipeline() Pipeline {
	return c.ops.TxPipeline()
}
