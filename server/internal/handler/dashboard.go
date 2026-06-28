package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"higolang/server/internal/model"
)

// DashboardHandler handles dashboard statistics.
type DashboardHandler struct {
	db *gorm.DB
}

// NewDashboardHandler creates a new DashboardHandler.
func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// GetDashboard returns dashboard statistics.
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	// Total articles.
	var totalArticles int64
	h.db.Model(&model.Article{}).Count(&totalArticles)

	// This week's articles.
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday()))
	var weekCount int64
	h.db.Model(&model.Article{}).Where("created_at >= ?", weekStart).Count(&weekCount)

	// Category counts.
	var categories []model.Category
	h.db.Find(&categories)

	type catStat struct {
		Name  string `json:"name"`
		Count int64  `json:"count"`
	}
	var catStats []catStat
	for _, cat := range categories {
		var count int64
		h.db.Model(&model.Article{}).Where("category_id = ?", cat.ID).Count(&count)
		catStats = append(catStats, catStat{Name: cat.Name, Count: count})
	}

	// Recent fetch logs.
	var recentLogs []model.FetchLog
	h.db.Order("created_at DESC").Limit(10).Find(&recentLogs)

	// Source status.
	var enabledSources, disabledSources int64
	h.db.Model(&model.FeedSource{}).Where("is_enabled = ?", true).Count(&enabledSources)
	h.db.Model(&model.FeedSource{}).Where("is_enabled = ?", false).Count(&disabledSources)

	OK(c, gin.H{
		"total_articles":   totalArticles,
		"week_count":       weekCount,
		"category_stats":   catStats,
		"recent_logs":      recentLogs,
		"enabled_sources":  enabledSources,
		"disabled_sources": disabledSources,
	})
}
