package cache

import "errors"

var (
	ErrTooFrequent = errors.New("发送过于频繁")
	ErrNotMatch    = errors.New("不对哦")
	ErrExceed      = errors.New("太多次了")
)
