package config

// 配置信息
type Config struct {
	DB    DBConfig    `mapstructure:"db"`
	Redis RedisConfig `mapstructure:"redis"`
}
type DBConfig struct {
	DSN string `mapstructure:"dsn"`
}
type RedisConfig struct {
	Addr string `mapstructure:"addr"`
}
