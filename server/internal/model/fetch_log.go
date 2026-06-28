package model

import "time"

// FetchLog records the result of a feed fetch operation.
type FetchLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	SourceID     uint      `gorm:"index" json:"source_id"`
	Status       string    `gorm:"size:50" json:"status"` // success, error
	FetchedCount int       `json:"fetched_count"`
	Message      string    `gorm:"type:text" json:"message"`
	DurationMs   int       `json:"duration_ms"`
	CreatedAt    time.Time `json:"created_at"`
}
