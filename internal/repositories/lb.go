package repositories

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/whotterre/tiermaster/internal/models"
	"go.uber.org/zap"
)

type LeaderboardRepo interface {
	AddScore(entry models.LeaderboardEntry) error
	GetTopScores(limit int) ([]EntryWithRank, error)
	GetUserRank(userID string) (int64, error)
}

type leaderBoardRepo struct {
	client *redis.Client
	logger *zap.Logger
}

func NewLeaderBoardRepo(client *redis.Client, logger *zap.Logger) LeaderboardRepo {
	return &leaderBoardRepo{
		client: client,
		logger: logger,
	}
}

func (r *leaderBoardRepo) AddScore(entry models.LeaderboardEntry) error {
	// Serialize the full entry data as JSON
	entryData, err := json.Marshal(entry)
	if err != nil {
		r.logger.Error("Failed to marshal leaderboard entry")
		return fmt.Errorf("Failed to marshal leaderboard entry: %w", err)
	}

	// Execute both commands in pipeline
	_, err = r.client.Pipelined(func(pipe redis.Pipeliner) error {
		zaddCmd := pipe.ZAdd("leaderboard:rankings", redis.Z{
			Score:  entry.Score,
			Member: entry.UserID,
		})
		if zaddCmd.Err() != nil {
			r.logger.Error("Failed to add score to rankings")
			return fmt.Errorf("Failed to add score to rankings: %w", zaddCmd.Err())
		}
		hsetCmd := pipe.HSet("leaderboard:entries", entry.UserID, entryData)
		if hsetCmd.Err() != nil {
			r.logger.Error("Failed to set leaderboard entry")
			return fmt.Errorf("Failed to set leaderboard entry: %w", hsetCmd.Err())
		}
		return nil
	})

	if err != nil {
		r.logger.Error("Redis pipeline failed")
		return fmt.Errorf("Redis pipeline failed: %w", err)
	}

	return nil
}

type EntryWithRank struct {
	models.LeaderboardEntry
	Rank int64
}

func (r *leaderBoardRepo) GetTopScores(limit int) ([]EntryWithRank, error) {
	// Get top N entries with scores
	results, err := r.client.ZRevRangeWithScores("leaderboard:rankings", 0, int64(limit-1)).Result()
	if err != nil {
		r.logger.Error("Failed to get rankings", zap.Error(err))
		return nil, fmt.Errorf("failed to get rankings: %w", err)
	}

	entries := make([]EntryWithRank, 0, len(results))
	for i, res := range results {
		userID, ok := res.Member.(string)
		if !ok {
			r.logger.Warn("Invalid userID type in leaderboard",
				zap.Any("member", res.Member),
			)
			continue
		}

		// Get full user data
		entryJSON, err := r.client.HGet("leaderboard:entries", userID).Result()
		if err != nil {
			r.logger.Error("Failed to get user entry",
				zap.String("userID", userID),
				zap.Error(err),
			)
			continue 
		}

		var entry EntryWithRank
		if err := json.Unmarshal([]byte(entryJSON), &entry); err != nil {
			r.logger.Error("Failed to unmarshal entry",
				zap.String("userID", userID),
				zap.Error(err),
			)
			continue
		}

		// Add rank (1-based index)
		entry.Rank = int64(i + 1)
		entry.Score = res.Score // Ensure score matches ranking

		entries = append(entries, entry)
	}

	return entries, nil
}

func (r *leaderBoardRepo) GetUserRank(userID string) (int64, error) {
	rank, err := r.client.ZRevRank("leaderboard:rankings", userID).Result()
	if err != nil {
		if err == redis.Nil {
			return -1, nil 
		}
		return 0, fmt.Errorf("failed to get user rank: %w", err)
	}
	return rank + 1, nil 
}