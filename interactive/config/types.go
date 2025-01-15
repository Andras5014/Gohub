package config

// 配置信息
type Config struct {
	DB    DBConfig    `mapstructure:"db"`
	Redis RedisConfig `mapstructure:"redis"`
	Kafka KafkaConfig `mapstructure:"kafka"`
	Grpc  GrpcConfig  `mapstructure:"grpc"`
}
type DBConfig struct {
	DSN string `mapstructure:"dsn"`
}
type RedisConfig struct {
	Addr string `mapstructure:"addr"`
}
type KafkaConfig struct {
	Addrs []string `mapstructure:"addrs"`
}

type GrpcConfig struct {
	Addr string `mapstructure:"addr"`
}
