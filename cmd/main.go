package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/whotterre/tiermaster/internal/config"
	"github.com/whotterre/tiermaster/internal/conn"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app := fiber.New()
	app.Use(cors.New(cors.Config{
        AllowOrigins:     "http://localhost:8080", 
        AllowHeaders:     "Origin, Content-Type, Accept, Authorization, Content-Disposition",
        AllowCredentials: true,
        AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
    }))
	SetupRoutes(app,client, logger)
	// Dynamic port alloc
	var port string
	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = envPort
	} else {
		port = config.ServicePort
	}


	// Pass client down to subsequent services
	log.Fatal(app.Listen(port))

}
