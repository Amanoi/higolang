package model

import (
	"time"

	"gorm.io/gorm"
)

// Article represents a blog post or news article.
type Article struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:500;not null" json:"title"`
	Slug        string         `gorm:"uniqueIndex;size:500;not null" json:"slug"`
	Summary     string         `gorm:"type:text" json:"summary"`
	Content     string         `gorm:"type:text" json:"content"`
	CoverURL    string         `gorm:"size:1000" json:"cover_url"`
	SourceURL   string         `gorm:"size:1000;index" json:"source_url"`
	SourceType  string         `gorm:"size:50;default:manual" json:"source_type"`
	CategoryID  uint           `json:"category_id"`
	Category    *Category      `json:"category,omitempty"`
	Tags        []Tag          `gorm:"many2many:article_tags;" json:"tags"`
	IsPinned    bool           `gorm:"default:false" json:"is_pinned"`
	Status      string         `gorm:"size:50;default:published;index" json:"status"`
	ViewCount   int            `gorm:"default:0" json:"view_count"`
	PublishedAt *time.Time     `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
