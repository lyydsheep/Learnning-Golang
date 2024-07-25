package config

type webookConfig struct {
	DB    DB
	Redis RedisConfig
}

type DB struct {
	DSN string
}

type RedisConfig struct {
	Addr string
}
