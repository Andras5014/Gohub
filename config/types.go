package config

// 配置信息
type config struct {
	DB    DBConfig
	Redis RedisConfig
}
type DBConfig struct {
	DSN string
}
type RedisConfig struct {
	Addr string
}
