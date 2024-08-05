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

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{cache: c}
}

func (cr *CodeRepository) Store(ctx context.Context, key, val string) error {
	return cr.cache.SetCode(ctx, key, val)
}

func (cr *CodeRepository) Check(ctx context.Context, key, input string) error {
	return cr.cache.CheckCode(ctx, key, input)
}
