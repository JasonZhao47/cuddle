package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	//go:embed lua/set_code.lua
	luaSetCode string
	//go:embed lua/verify_code.lua
	luaVerifyCode         string
	ErrNextCodeTooSoon    = errors.New("验证码发送间隔过短")
	ErrInvalidCode        = errors.New("未知错误")
	ErrVerifyTooManyTimes = errors.New("验证码验证次数过多")
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, input string) (bool, error)
}

type RedisCodeCache struct {
	cmd redis.Cmdable
}

func NewRedisCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	status, err := c.cmd.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch status {
	case -1:
		return ErrNextCodeTooSoon
	case 0:
		return nil
	default:
		return ErrInvalidCode
	}
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, input string) (bool, error) {
	status, err := c.cmd.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, input).Int()
	if err != nil {
		return false, err
	}
	switch status {
	case 0:
		return true, nil
	case -1:
		return false, ErrVerifyTooManyTimes
	default:
		return false, nil
	}

}

func (c *RedisCodeCache) key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
