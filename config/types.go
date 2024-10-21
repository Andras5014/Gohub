package config

// 配置信息
type Config struct {
	DB    DBConfig    `mapstructure:"db"`
	Redis RedisConfig `mapstructure:"redis"`

	Kafka KafkaConfig `mapstructure:"kafka"`
}
type DBConfig struct {
	DSN string `mapstructure:"dsn"`
}
type RedisConfig struct {
	Addr string `mapstructure:"addr"`
}
type KafkaConfig struct {
	Addrs []string `mapstructure:"addr"`
}
