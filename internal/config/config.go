package config

import (
    "os"
)

type Config struct {
    RedisAddr     string `mapstructure:"REDIS_ADDR"`
    RedisPassword string `mapstructure:"REDIS_PASSWORD"`
    RedisDB       string `mapstructure:"REDIS_DB"`
    ServicePort   string `mapstructure:"SERVICE_PORT"`
    ServiceName   string `mapstructure:"SERVICE_NAME"`
    DevEnv        string `mapstructure:"DEV_ENV"`
}

func LoadConfig() Config {
    return Config{
        RedisAddr:     getEnv("REDIS_ADDR", "redis-1:6379"),
        RedisPassword: getEnv("REDIS_PASSWORD", ""),
        RedisDB:       getEnv("REDIS_DB", "0"),
        ServicePort:   getEnv("SERVICE_PORT", ":3001"),
        ServiceName:   getEnv("SERVICE_NAME", "leaderboard"),
        DevEnv:        getEnv("DEV_ENV", "development"),
    }
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}