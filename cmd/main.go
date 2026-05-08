package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	reqLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/whotterre/tiermaster/internal/config"
	"github.com/whotterre/tiermaster/internal/conn"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)

	}
	defer logger.Sync()
	// Load config
	config := config.LoadConfig()
	logger.Info("Starting leaderboard microservice",
		zap.String("port", config.ServicePort),
		zap.String("service", config.ServiceName),
	)
	// Initialize Redis client
	client, err := conn.GetRedisClient(config)
	if err != nil {
		logger.Error("Failed to initialize Redis client")
	}

	app := fiber.New()

	app.Use(reqLogger.New())
	SetupRoutes(app, client, logger)

	// Health check 
	app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "healthy",
            "service": config.ServiceName,
        })
    })
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
