package cache

import "context"

type CodeCache interface {
	CheckCode(ctx context.Context, key string, input string) error
	SetCode(ctx context.Context, key string, val string) error
}
