package cache

import (
	"context"
	_ "embed"
	"errors"
	"github.com/redis/go-redis/v9"
)

// 引入lua脚本
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/check_code.lua
var luaCheckCode string

type RedisCodeCache struct {
	client redis.Cmdable
}

func NewCodeRedis(c redis.Cmdable) CodeCache {
	return &RedisCodeCache{client: c}
}

func (cc *RedisCodeCache) CheckCode(ctx context.Context, key string, input string) error {
	res, err := cc.client.Eval(ctx, luaCheckCode, []string{key}, input).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return ErrNotMatch
	case -3:
		return ErrExceed
	default:
		return errors.New("系统错误")
	}
}

func (cc *RedisCodeCache) SetCode(ctx context.Context, key string, val string) error {
	//获取返回值
	res, err := cc.client.Eval(ctx, luaSetCode, []string{key}, val).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return ErrTooFrequent
	default:
		return errors.New("系统错误")

	}
}
