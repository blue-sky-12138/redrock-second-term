package DBFirst

import (
	"github.com/go-redis/redis"
)

type DialOption interface {
	apply(options *redis.Options)
}

func newFuncDialOption(f func(*redis.Options)) *funcDialOption {
	return &funcDialOption{
		f: f,
	}
}

type funcDialOption struct {
	f func(*redis.Options)
}

func (f *funcDialOption) apply(options *redis.Options) {
	f.f(options)
}

func WithAddr(addr string) DialOption {
	return newFuncDialOption(func(options *redis.Options) {
		options.Addr = addr
	})
}

func WithPoolSize(pz int) DialOption {
	return newFuncDialOption(func(options *redis.Options) {
		options.PoolSize = pz
	})
}

//由于go-redis包里会对值为零值的变量初始化，故这里只需要处理个性化数据。
//其实能改哪些，结构体里的变量代表什么，我的相关知识<5
func NewClient(options []DialOption) *redis.Client {
	o := redis.Options{}
	for _, opt := range options {
		opt.apply(&o)
	}
	return redis.NewClient(&o)
}
