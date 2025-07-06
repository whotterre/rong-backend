package main

import (
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/whotterre/tiermaster/internal/handlers"
	"github.com/whotterre/tiermaster/internal/repositories"
	"github.com/whotterre/tiermaster/internal/services"
	"go.uber.org/zap"
)

func SetupRoutes(app *fiber.App, client *redis.Client, logger *zap.Logger){

	lbGroup := app.Group("/api/v1/leaderboard")
	lbRepo := repositories.NewLeaderBoardRepo(client, logger)
	lbService := services.NewLeaderboardService(lbRepo, logger)
	lbHandler := handlers.NewLeaderboardHandler(lbService)

	// Leaderboard logic 
	lbGroup.Post("/score/:userID", lbHandler.AddScore)
	lbGroup.Get("/scores", lbHandler.GetTopNPlayers)
	lbGroup.Get("/highest", lbHandler.GetHighestScore)

}