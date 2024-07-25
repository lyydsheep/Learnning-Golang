//go:build !k8s

// 在非k8s环境下使用
package config

var Config = webookConfig{
	DB:    DB{DSN: "root:root@tcp(localhost:11316)/webook"},
	Redis: RedisConfig{Addr: "localhost:6379"},
}
