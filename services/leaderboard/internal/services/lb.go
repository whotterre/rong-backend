package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/whotterre/tiermaster/internal/models"
	"github.com/whotterre/tiermaster/internal/repositories"
	"go.uber.org/zap"
)

// Service interface defines the business capabilities
type LeaderboardService interface {
	SubmitScore(entry models.LeaderboardEntry) error
	GetTopPlayers(limit int) ([]repositories.EntryWithRank, error)
	GetHighestScore() (int64, error)
}

type leaderboardService struct {
	repo   repositories.LeaderboardRepo
	logger *zap.Logger
}

func NewLeaderboardService(repo repositories.LeaderboardRepo, logger *zap.Logger) LeaderboardService {
	return &leaderboardService{
		repo:   repo,
		logger: logger.With(zap.String("service", "leaderboard")),
	}
}

// SubmitScore handles business logic before storage
func (s *leaderboardService) SubmitScore(entry models.LeaderboardEntry) error {
	// Validation
	if entry.Score < 0 {
		return errors.New("score cannot be negative")
	}
	if entry.UserID == "" {
		return errors.New("user ID is required")
	}

	// Set timestamp
	entry.LastUpdatedAt = time.Now().UTC()

	// Business logic example: Minimum score threshold
	const minScore = 2
	if entry.Score < minScore {
		s.logger.Warn("Score below threshold",
			zap.String("userID", entry.UserID),
			zap.Float64("score", entry.Score),
		)
		return fmt.Errorf("score must be at least %d", minScore)
	}

	// Call repository
	if err := s.repo.AddScore(entry); err != nil {
		s.logger.Error("Failed to submit score",
			zap.String("userID", entry.UserID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to submit score: %w", err)
	}

	s.logger.Info("Score submitted",
		zap.String("userID", entry.UserID),
		zap.Float64("score", entry.Score),
	)
	return nil
}

// GetTopPlayers retrieves and enriches leaderboard data
func (s *leaderboardService) GetTopPlayers(limit int) ([]repositories.EntryWithRank, error) {
	if limit <= 0 || limit > 100 {
		return nil, errors.New("limit must be between 1-100")
	}

	entries, err := s.repo.GetTopScores(limit)
	if err != nil {
		s.logger.Error("Failed to get top scores",
			zap.Int("limit", limit),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get leaderboard: %w", err)
	}

	// Add ranking position
	for i := range entries {
		entries[i].Rank = int64(i + 1)
	}

	s.logger.Debug("Retrieved top players",
		zap.Int("count", len(entries)),
	)
	return entries, nil
}

func (s *leaderboardService) GetHighestScore() (int64, error) {
	highestScore, err := s.repo.GetHighestScore()
	if err != nil {
		s.logger.Error("Failed to get highest score",
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to get highest score: %w", err)
	}

	s.logger.Debug("Retrieved highest score",
		zap.Int64("highestScore", highestScore),
	)
	return highestScore, nil
}