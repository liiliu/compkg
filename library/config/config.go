package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"time"
)

func Init(env string) {
	viper.AddConfigPath("./config/env")
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetDuration(key string) time.Duration {
	return viper.GetDuration(key)
}

func MustHas(key string) {
	if viper.GetString(key) == "" {
		log.Println(fmt.Sprintf("miss %s in config file", key))
		panic(-1)
	}
}
