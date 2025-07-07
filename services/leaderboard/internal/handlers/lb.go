package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/whotterre/tiermaster/internal/models"
	"github.com/whotterre/tiermaster/internal/services"
	"go.uber.org/zap" // Import zap for logging
)

type LeaderboardHandler struct {
	lbService services.LeaderboardService
	logger    *zap.Logger // Add logger to the handler
}

func NewLeaderboardHandler(lbService services.LeaderboardService, logger *zap.Logger) *LeaderboardHandler {
	return &LeaderboardHandler{
		lbService: lbService,
		logger:    logger.With(zap.String("handler", "leaderboard")), // Initialize logger
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
		l.logger.Error("Failed to parse request body for AddScore", zap.Error(err)) // Log error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	newEntry := models.LeaderboardEntry{
		UserID:        req.UserID,
		Score:         req.Score,
		Username:      req.Username,
		LastUpdatedAt: time.Now(),
	}

	err := l.lbService.SubmitScore(newEntry)
	if err != nil {
		l.logger.Error("Failed to submit score to service", zap.Error(err), zap.String("userID", req.UserID)) // Log error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Something went wrong while submitting the scores",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Successfully submitted score",
	})
}

func (l *LeaderboardHandler) GetTopNPlayers(c *fiber.Ctx) error {
	limitStr := c.Query("limit")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		l.logger.Error("Invalid 'limit' query parameter", zap.String("limitStr", limitStr), zap.Error(err)) // Log error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid 'limit' query parameter. Must be an integer.", // Corrected error message
		})
	}

	leaderboard, err := l.lbService.GetTopPlayers(limit)
	if err != nil {
		l.logger.Error("Failed to get top players from service", zap.Int("limit", limit), zap.Error(err)) // Log error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get top %d scores", limit),
		})
	}

	// --- NEW DEBUGGING LINE: Log the actual data being sent ---
	l.logger.Debug("Sending leaderboard response",
		zap.Int("limit", limit),
		zap.Any("leaderboardData", leaderboard), // Log the actual data structure
	)
	// --- END NEW DEBUGGING LINE ---

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     fmt.Sprintf("Top %d users", limit),
		"leaderboard": leaderboard,
	})
}

func (l *LeaderboardHandler) GetHighestScore(c *fiber.Ctx) error {
	highestScore, err := l.lbService.GetHighestScore()
	if err != nil {
		l.logger.Error("Failed to get highest score from service", zap.Error(err)) // Log error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get highest score",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"highestScore": highestScore,
	})
}
