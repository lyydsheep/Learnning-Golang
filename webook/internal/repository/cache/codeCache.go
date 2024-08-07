package cache

import "context"

type codeRedisCache interface {
	CheckCode(ctx context.Context, key string, input string) error
	SetCode(ctx context.Context, key string, val string) error
}
