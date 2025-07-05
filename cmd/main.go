package main

import (
	"github.com/whotterre/tiermaster/internal/config"
	"github.com/whotterre/tiermaster/internal/conn"
	"go.uber.org/zap"
)

func main(){
	// Initialize logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	// Load config 
	config, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load config from .env file")
		panic(err)
	}
	// Initialize Redis client 
	client, err := conn.GetRedisClient(config)
	if err != nil {
		logger.Error("Failed to initialize Redis client")
	}
	// Pass client down to subsequent services
	println(client)

}