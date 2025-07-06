package models

import "time"

type LeaderboardEntry struct {
	UserID        string    `json:"userId"`
	Username      string    `json:"username"`
	Score         float64   `json:"score"`
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`
}