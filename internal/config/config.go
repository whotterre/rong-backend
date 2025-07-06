package config

import "github.com/spf13/viper"

type Config struct {
	RedisAddr     string `mapstructure:"REDIS_ADDR"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       string `mapstructure:"REDIS_DB"`
	ServicePort   string `mapstructure:"SERVICE_PORT"`
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("../")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
