package repository

import (
	"context"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/repository/cache"
)

var (
	ErrTooFrequent = cache.ErrTooFrequent
	ErrNotMatch    = cache.ErrNotMatch
	ErrExceed      = cache.ErrExceed
)

type CodeRepository interface {
	Store(ctx context.Context, key, val string) error
	Check(ctx context.Context, key, input string) error
}

type CachedCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return &CachedCodeRepository{cache: c}
}

func (cr *CachedCodeRepository) Store(ctx context.Context, key, val string) error {
	return cr.cache.SetCode(ctx, key, val)
}

func (cr *CachedCodeRepository) Check(ctx context.Context, key, input string) error {
	return cr.cache.CheckCode(ctx, key, input)
}
