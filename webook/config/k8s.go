//go:build k8s

// 在k8s环境下使用
package config

var Config = webookConfig{
	DB:    DB{DSN: "root:root@tcp(webook-mysql:3308)/webook"},
	Redis: RedisConfig{Addr: "webook-redis:6380"},
}
