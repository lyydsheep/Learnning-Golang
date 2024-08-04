package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lyydsheep/Learnning-Golang/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type UserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) *UserCache {
	//设置过期时间
	return &UserCache{
		client:     client,
		expiration: 15 * time.Minute,
	}
}

func (u *UserCache) Get(ctx context.Context, id int) (domain.User, error) {
	key := u.Key(id)
	//获取缓存中的val
	//注意该val是序列化后的
	val, err := u.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	//反序列化
	var user domain.User
	err = json.Unmarshal(val, &user)
	return user, err
}

func (u *UserCache) Set(ctx context.Context, user domain.User) error {
	//设置key值
	key := u.Key(user.Id)
	//将val转成字节数据，因为redis只能存序列化后的数据，无法直接存domain.User
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	//存入缓存中
	return u.client.Set(ctx, key, val, u.expiration).Err()
}

func (u *UserCache) Key(id int) string {
	key := fmt.Sprintf("user:info:%d", id)
	return key
}
