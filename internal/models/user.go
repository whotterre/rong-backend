package models

import "time"

// Leaderboard entry model
type LeaderboardEntry struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Score    int64  `json:"score"`
	TimeStamp time.Time `json:"timestamp"`
}
