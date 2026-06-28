package model

import "time"

// FeedSource represents an RSS/Atom/JSON feed to fetch articles from.
type FeedSource struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Name          string     `gorm:"size:200;not null" json:"name"`
	URL           string     `gorm:"size:1000;not null" json:"url"`
	FeedType      string     `gorm:"size:50;not null" json:"feed_type"` // rss, atom, json
	CategoryID    uint       `json:"category_id"`
	FetchInterval int        `gorm:"default:30" json:"fetch_interval"` // minutes
	IsEnabled     bool       `gorm:"default:true" json:"is_enabled"`
	LastFetchedAt *time.Time `json:"last_fetched_at"`
	CreatedAt     time.Time  `json:"created_at"`
}
