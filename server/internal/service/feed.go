package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
	"gorm.io/gorm"

	"higolang/server/internal/model"
)

// FeedService handles fetching and parsing of RSS/Atom/JSON feeds.
type FeedService struct {
	db *gorm.DB
}

// NewFeedService creates a new FeedService.
func NewFeedService(db *gorm.DB) *FeedService {
	return &FeedService{db: db}
}

// FetchFeedSource fetches a single feed source and creates articles from new items.
// Returns the number of newly created articles.
func (s *FeedService) FetchFeedSource(source model.FeedSource) (int, error) {
	start := time.Now()
	count := 0
	var errMsg string

	switch source.FeedType {
	case "rss", "atom":
		c, err := s.fetchRSS(source)
		count = c
		if err != nil {
			errMsg = err.Error()
		}
	case "json":
		c, err := s.fetchJSON(source)
		count = c
		if err != nil {
			errMsg = err.Error()
		}
	default:
		errMsg = fmt.Sprintf("unsupported feed type: %s", source.FeedType)
	}

	// Update last_fetched_at.
	now := time.Now()
	s.db.Model(&source).Update("last_fetched_at", now)

	// Write a fetch log entry.
	status := "success"
	if errMsg != "" {
		status = "error"
	}
	fetchLog := model.FetchLog{
		SourceID:     source.ID,
		Status:       status,
		FetchedCount: count,
		Message:      errMsg,
		DurationMs:   int(time.Since(start).Milliseconds()),
	}
	s.db.Create(&fetchLog)

	if errMsg != "" {
		return count, fmt.Errorf("%s", errMsg)
	}
	return count, nil
}

// fetchRSS handles RSS and Atom feeds via gofeed.
func (s *FeedService) fetchRSS(source model.FeedSource) (int, error) {
	fp := gofeed.NewParser()
	fp.Client = &http.Client{Timeout: 30 * time.Second}

	feed, err := fp.ParseURL(source.URL)
	if err != nil {
		return 0, fmt.Errorf("parse feed error: %w", err)
	}

	count := 0
	for _, item := range feed.Items {
		// Skip duplicates by source_url.
		if item.Link != "" {
			var existing model.Article
			err := s.db.Where("source_url = ?", item.Link).First(&existing).Error
			if err == nil {
				continue // already exists
			}
		}

		slug := GenerateSlug(item.Title)
		summary := ""
		if item.Description != "" {
			summary = item.Description
		} else if item.Content != "" {
			summary = item.Content
		}
		if len(summary) > 500 {
			summary = summary[:500]
		}

		content := item.Content
		if content == "" {
			content = item.Description
		}

		now := time.Now()
		article := model.Article{
			Title:       item.Title,
			Slug:        slug,
			Summary:     summary,
			Content:     content,
			SourceURL:   item.Link,
			SourceType:  source.FeedType,
			CategoryID:  source.CategoryID,
			Status:      "published",
			PublishedAt: &now,
		}
		if item.PublishedParsed != nil {
			article.PublishedAt = item.PublishedParsed
		}

		if err := s.db.Create(&article).Error; err == nil {
			count++
		}
	}

	return count, nil
}

// goReleaseJSON models the top-level Go downloads JSON feed.
type goReleaseJSON struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	Files   []struct {
		Filename string `json:"filename"`
		OS       string `json:"os"`
		Arch     string `json:"arch"`
		Version  string `json:"version"`
		SHA256   string `json:"sha256"`
		Size     int64  `json:"size"`
		Kind     string `json:"kind"`
	} `json:"files"`
}

// fetchJSON handles JSON feeds (e.g. Go release announcements).
func (s *FeedService) fetchJSON(source model.FeedSource) (int, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(source.URL)
	if err != nil {
		return 0, fmt.Errorf("http get error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("read body error: %w", err)
	}

	// Try parsing as an array of releases first.
	var releases []goReleaseJSON
	if err := json.Unmarshal(body, &releases); err != nil {
		// Try as a single object.
		var single goReleaseJSON
		if err2 := json.Unmarshal(body, &single); err2 != nil {
			return 0, fmt.Errorf("json parse error: %w", err)
		}
		releases = []goReleaseJSON{single}
	}

	count := 0
	for _, rel := range releases {
		sourceURL := fmt.Sprintf("%s#%s", source.URL, rel.Version)

		var existing model.Article
		if err := s.db.Where("source_url = ?", sourceURL).First(&existing).Error; err == nil {
			continue
		}

		slug := GenerateSlug(rel.Version)
		summary := fmt.Sprintf("Go %s release", rel.Version)
		content := fmt.Sprintf("<h2>Go %s</h2><p>Stable: %v</p>", rel.Version, rel.Stable)

		now := time.Now()
		article := model.Article{
			Title:       fmt.Sprintf("Go %s Released", rel.Version),
			Slug:        slug,
			Summary:     summary,
			Content:     content,
			SourceURL:   sourceURL,
			SourceType:  "json",
			CategoryID:  source.CategoryID,
			Status:      "published",
			PublishedAt: &now,
		}

		if err := s.db.Create(&article).Error; err == nil {
			count++
		}
	}

	return count, nil
}

// FetchAllEnabled iterates all enabled feed sources and fetches each one.
func (s *FeedService) FetchAllEnabled() {
	var sources []model.FeedSource
	s.db.Where("is_enabled = ?", true).Find(&sources)
	for _, src := range sources {
		s.FetchFeedSource(src)
	}
}
