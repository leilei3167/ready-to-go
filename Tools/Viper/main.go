package main

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

func main() {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	fmt.Println("当前目录为:", dir)

	viper.SetConfigType("toml")
	viper.SetConfigName("config")
	viper.AddConfigPath(dir)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("读取配置文件出错:", err)
	}
	fmt.Println(viper.Get("etcd.Host"))

	var c CommonConfig

	if err := viper.Unmarshal(&c); err != nil {
		log.Fatal("解码出错:", err)
	}
	fmt.Printf("%#v", c)

}

type CommonConfig struct {
	Etcd  `mapstructure:"etcd"`
	Redis `mapstructure:"redis"`
}

type Etcd struct {
	Host              string `mapstructure:"host"`
	BasePath          string `mapstructure:"basePath"`
	ServerPathLogic   string `mapstructure:"serverPathLogic"`
	ServerPathConnect string `mapstructure:"serverPathConnect"`
	UserName          string `mapstructure:"userName"`
	Password          string `mapstructure:"password"`
	ConnectionTimeout int    `mapstructure:"connectionTimeout"`
}

type Redis struct {
	RedisAddress  string `mapstructure:"redisAddress"`
	RedisPassword string `mapstructure:"redisPassword"`
	Db            int    `mapstructure:"db"`
}
