package ioc

import (
	"github.com/Andras5014/gohub/config"
	"github.com/spf13/viper"
)

func InitConfig() *config.Config {
	cfg := &config.Config{}
	err := viper.Unmarshal(cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
