package DBSecond

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

const (
	Addr        = "localhost:6379"
	IdLeTimeout = 5
	MaxIdle     = 20
	MaxActive   = 8
)

type Client struct {
	Option OptionPool
	pool   *redis.Pool
}

var DefaultOption = OptionPool{
	addr:        Addr,
	idLeTimeout: IdLeTimeout,
	maxIdle:     MaxIdle,
	maxActive:   MaxActive,
}

type OptionPool struct {
	addr        string
	idLeTimeout int
	maxIdle     int
	maxActive   int
}

type PoolExt interface {
	apply(*OptionPool)
}

type tempFunc func(pool *OptionPool)

type funcPoolExt struct {
	f tempFunc
}

func (f *funcPoolExt) apply(p *OptionPool) {
	f.f(p)
}

func newFuncPoolExt(f tempFunc) *funcPoolExt {
	return &funcPoolExt{
		f: f,
	}
}

func WithAddr(addr string) PoolExt {
	return newFuncPoolExt(func(pool *OptionPool) {
		pool.addr = addr
	})
}

func WithIdLeTimeout(d int) PoolExt {
	return newFuncPoolExt(func(pool *OptionPool) {
		pool.idLeTimeout = d
	})
}

func WithMaxIdle(maxId int) PoolExt {
	return newFuncPoolExt(func(pool *OptionPool) {
		pool.maxIdle = maxId
	})
}

func WithMaxActive(maxAct int) PoolExt {
	return newFuncPoolExt(func(pool *OptionPool) {
		pool.maxActive = maxAct
	})
}

func NewClient(op ...PoolExt) *Client {
	c := &Client{Option: DefaultOption}
	for _, p := range op {
		p.apply(&c.Option)
	}
	c.setRedisPool()
	return c
}

func (pc *Client) setRedisPool() {
	pc.pool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", pc.Option.addr)
			if conn == nil || err != nil {
				return nil, err
			}
			return conn, nil
		},
		MaxIdle:     pc.Option.maxIdle,                                  // 最大空闲连接数
		MaxActive:   pc.Option.maxActive,                                // 最大活跃连接数
		IdleTimeout: time.Second * time.Duration(pc.Option.idLeTimeout), // 连接等待时间
	}
}
func (pc *Client) newRedisOperate(commandName string, args ...interface{}) (interface{}, error) {
	c := pc.pool.Get()
	defer c.Close()
	values, err := c.Do(commandName, args...)
	if err != nil {
		return nil, err
	}
	return values, nil
}

func (pc *Client) Set(args ...interface{}) error {
	_, err := pc.newRedisOperate("SET", args...)
	return err
}

func (pc *Client) Get(key interface{}) (interface{}, error) {
	return pc.newRedisOperate("GET", key)
}

func (pc *Client) TTL(key interface{}) (interface{}, error) {
	return pc.newRedisOperate("TTL", key)
}

func (pc *Client) MSet(args ...interface{}) error {
	_, err := pc.newRedisOperate("MSET", args...)
	return err
}

func (pc *Client) MGet(args ...interface{}) ([]interface{}, error) {
	values, err := pc.newRedisOperate("MGET", args...)
	if err != nil {
		return nil, err
	}
	v, _ := redis.Values(values, nil)
	return v, nil
}

func (pc *Client) Exists(key interface{}) int {
	value, err := pc.newRedisOperate("EXISTS", key)
	if err != nil {
		return 0
	}
	v, _ := redis.Int(value, nil)
	return v
}

func (pc *Client) SetEX(time int, args ...interface{}) error {
	args = append(args, "EX", time)
	return pc.Set(args...)
}

func (pc *Client) SetNX(args ...interface{}) error {
	args = append(args, "NX")
	return pc.Set(args...)
}

func (pc *Client) GetSet(args ...interface{}) (interface{}, error) {
	return pc.newRedisOperate("GETSET", args...)
}

func (pc *Client) Strlen(key interface{}) (int, error) {
	value, err := pc.newRedisOperate("STRLEN", key)
	if err != nil {
		return 0, err
	}
	v, _ := redis.Int(value, nil)
	return v, nil
}

func (pc *Client) Append(key interface{}, value interface{}) error {
	_, err := pc.newRedisOperate("APPEND", key, value)
	return err
}

func (pc *Client) SetRange(key interface{}, offset int, value interface{}) error {
	_, err := pc.newRedisOperate("SETRANGE", key, offset, value)
	return err
}

func (pc *Client) GetRange(key interface{}, start int, end int) (interface{}, error) {
	value, err := pc.newRedisOperate("GETRANGE", key, start, end)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (pc *Client) Del(args ...interface{}) error {
	_, err := pc.newRedisOperate("Del", args...)
	return err
}

func (pc *Client) INCR(key interface{}) (int, error) {
	value, err := pc.newRedisOperate("INCR", key)
	if err != nil {
		return 0, err
	}
	v, _ := redis.Int(value, nil)
	return v, nil
}

func (pc *Client) HSet(hash interface{}, field interface{}, value interface{}) error {
	_, err := pc.newRedisOperate("HSET", hash, field, value)
	return err
}

func (pc *Client) HGet(hash interface{}, field interface{}) (interface{}, error) {
	return pc.newRedisOperate("HGET", hash, field)
}

func (pc *Client) HStrlen(hash interface{}, field interface{}) (int, error) {
	value, err := pc.newRedisOperate("HSTRLEN", hash, field)
	if err != nil {
		return 0, err
	}
	v, _ := redis.Int(value, nil)
	return v, nil
}

func (pc *Client) HDel(hash interface{}, field interface{}) error {
	_, err := pc.newRedisOperate("HDEL", hash, field)
	return err
}

func (pc *Client) LPush(args ...interface{}) error {
	_, err := pc.newRedisOperate("LPUSH", args...)
	return err
}

func (pc *Client) LPop(args ...interface{}) (interface{}, error) {
	return pc.newRedisOperate("LPUSH", args...)
}

func (pc *Client) LSet(key interface{}, index int, value interface{}) (interface{}, error) {
	return pc.newRedisOperate("LSET", key, index, value)
}

func (pc *Client) LRangeAll(key interface{}) ([]interface{}, error) {
	return pc.LRange(key, 0, -1)
}

func (pc *Client) LRange(key interface{}, start int, end int) ([]interface{}, error) {
	values, err := pc.newRedisOperate("LRANGE", key, start, end)
	if err != nil {
		return nil, err
	}
	v, _ := redis.Values(values, err)
	return v, nil
}

func (pc *Client) SAdd(args ...interface{}) error {
	_, err := pc.newRedisOperate("SADD", args...)
	return err
}

func (pc *Client) SPop(key interface{}, count int) ([]interface{}, error) {
	if count < 1 {
		count = 1
	}
	values, err := pc.newRedisOperate("SPOP", key, count)
	if err != nil {
		return nil, err
	}
	v, _ := redis.Values(values, nil)
	return v, nil
}

func (pc *Client) SMove(args ...interface{}) error {
	_, err := pc.newRedisOperate("SMOVE", args...)
	return err
}

func (pc *Client) SMembers(key interface{}) ([]interface{}, error) {
	values, err := pc.newRedisOperate("SMEMBERS", key)
	if err != nil {
		return nil, err
	}
	v, _ := redis.Values(values, nil)
	return v, nil
}

func (pc *Client) ZAdd(args ...interface{}) error {
	_, err := pc.newRedisOperate("ZADD", args...)
	return err
}

func (pc *Client) ZRangeAll(key interface{}) ([]interface{}, error) {
	return pc.ZRange(key, 0, -1)
}

func (pc *Client) ZRange(key interface{}, start int, end int) ([]interface{}, error) {
	values, err := pc.newRedisOperate("ZRANGE", key, start, end)
	if err != nil {
		return nil, err
	}
	v, _ := redis.Values(values, nil)
	return v, nil
}

func (pc *Client) ZRangeWithScores(key interface{}, start int, end int) ([]interface{}, error) {
	values, err := pc.newRedisOperate("ZRANGE", key, start, end, "WITHSCORES")
	if err != nil {
		return nil, err
	}
	v, _ := redis.Values(values, nil)
	return v, nil
}

func (pc *Client) ZRem(args ...interface{}) error {
	_, err := pc.newRedisOperate("ZREM", args...)
	return err
}
