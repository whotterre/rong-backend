package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/whotterre/tiermaster/internal/models"
	"github.com/whotterre/tiermaster/internal/services"
)

type LeaderboardHandler struct {
	lbService services.LeaderboardService
}

func NewLeaderboardHandler(lbService services.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{
		lbService: lbService,
	}
}

type CreateLeaderboardEntry struct {
	UserID        string    `json:"userId"`
	Username      string    `json:"username"`
	Score         float64   `json:"score"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
}

func (l *LeaderboardHandler) AddScore(c *fiber.Ctx) error {
	var req CreateLeaderboardEntry
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Failed to parse request body",
		})
	}
	
	newEntry := models.LeaderboardEntry{
		UserID: req.UserID,
		Score: req.Score,
		Username: req.Username,
		LastUpdatedAt: time.Now(),
	}


	err := l.lbService.SubmitScore(newEntry)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong while submitting the scores",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully submitted score",
	})
}

func (l *LeaderboardHandler) GetTopNPlayers (c *fiber.Ctx) error {
	limitStr := c.Query("limit")

	// Convert limit to int 
	limit, err := strconv.Atoi(limitStr) 
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Failed to parse request body",
		})
	}


	leaderboard, err := l.lbService.GetTopPlayers(limit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : fmt.Sprintf("Failed to get top %d scores", limit),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Top %d users", limit),
		"leaderboard": leaderboard,
	})

}

func (l *LeaderboardHandler) GetHighestScore(c *fiber.Ctx) error {
	highestScore, err := l.lbService.GetHighestScore()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get highest score",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"highestScore": highestScore,
	})
}

