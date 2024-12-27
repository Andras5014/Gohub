package startup

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func init() {
	initViper()
}
func initViper() {
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	// 实时监听配置变更
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println(in.Name, in.Op)
		fmt.Println("config file changed")
	})

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
