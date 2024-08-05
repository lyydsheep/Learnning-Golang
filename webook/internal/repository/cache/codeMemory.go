package cache

import (
	"context"
	"fmt"
	"github.com/coocood/freecache"
	"strconv"
)

type CodeMemory struct {
	client *freecache.Cache
	expire int
}

func NewCodeMemory(c *freecache.Cache) *CodeMemory {
	return &CodeMemory{
		client: c,
		expire: 60,
	}
}

func (c *CodeMemory) CheckCode(ctx context.Context, key string, input string) error {
	fmt.Println("this is freeCache checkCode")
	// 有 or 没有
	keyByte, cntByte := c.toByte(key), c.toByte(key+"cnt")
	valByte, err := c.client.Get(keyByte)
	if err != nil {
		return err
	}
	cnt, err := c.client.Get(cntByte)
	if err != nil {
		return err
	}
	remain, err := strconv.Atoi(c.toString(cnt))
	if err != nil {
		return err
	}
	if remain <= 0 {
		return ErrExceed
	}
	remain--
	val := c.toString(valByte)
	//对 or 不对
	if val == input {
		return nil
	}
	err = c.client.Set(cntByte, c.toByte(strconv.Itoa(remain)), c.expire)
	if err != nil {
		return err
	}
	return ErrNotMatch
}

func (c *CodeMemory) SetCode(ctx context.Context, key string, val string) error {
	//要上锁
	fmt.Println("this is freeCache setCode")
	keyByte, valByte, cntByte := c.toByte(key), c.toByte(val), c.toByte(key+"cnt")
	//存在 or 不存在
	_, err := c.client.Get(keyByte)
	if err == nil {
		//存在
		return nil
	}
	//不存在
	err = c.client.Set(keyByte, valByte, c.expire)
	if err != nil {
		return err
	}
	err = c.client.Set(cntByte, []byte("3"), c.expire)
	return err
}

func (c *CodeMemory) toByte(s string) []byte {
	return []byte(s)
}

func (c *CodeMemory) toString(x []byte) string {
	return string(x)
}
