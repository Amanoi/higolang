package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"higolang/server/internal/model"
	"higolang/server/internal/service"
)

// FeedHandler handles feed source management and manual fetching.
type FeedHandler struct {
	db        *gorm.DB
	feedSvc   *service.FeedService
}

// NewFeedHandler creates a new FeedHandler.
func NewFeedHandler(db *gorm.DB, feedSvc *service.FeedService) *FeedHandler {
	return &FeedHandler{db: db, feedSvc: feedSvc}
}

// ListFeedSources returns all feed sources.
func (h *FeedHandler) ListFeedSources(c *gin.Context) {
	var sources []model.FeedSource
	if err := h.db.Order("id ASC").Find(&sources).Error; err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}
	OK(c, sources)
}

type feedSourceInput struct {
	Name          string `json:"name" binding:"required"`
	URL           string `json:"url" binding:"required"`
	FeedType      string `json:"feed_type" binding:"required"`
	CategoryID    uint   `json:"category_id"`
	FetchInterval int    `json:"fetch_interval"`
	IsEnabled     bool   `json:"is_enabled"`
}

// CreateFeedSource creates a new feed source (admin).
func (h *FeedHandler) CreateFeedSource(c *gin.Context) {
	var input feedSourceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	if input.FetchInterval <= 0 {
		input.FetchInterval = 30
	}

	source := model.FeedSource{
		Name:          input.Name,
		URL:           input.URL,
		FeedType:      input.FeedType,
		CategoryID:    input.CategoryID,
		FetchInterval: input.FetchInterval,
		IsEnabled:     input.IsEnabled,
	}

	if err := h.db.Create(&source).Error; err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	OK(c, source)
}

// UpdateFeedSource updates a feed source (admin).
func (h *FeedHandler) UpdateFeedSource(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid feed source ID")
		return
	}

	var input feedSourceInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	var source model.FeedSource
	if err := h.db.First(&source, id).Error; err != nil {
		Fail(c, http.StatusNotFound, "feed source not found")
		return
	}

	h.db.Model(&source).Updates(map[string]interface{}{
		"name":           input.Name,
		"url":            input.URL,
		"feed_type":      input.FeedType,
		"category_id":    input.CategoryID,
		"fetch_interval": input.FetchInterval,
		"is_enabled":     input.IsEnabled,
	})

	OK(c, source)
}

// DeleteFeedSource deletes a feed source (admin).
func (h *FeedHandler) DeleteFeedSource(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid feed source ID")
		return
	}

	if err := h.db.Delete(&model.FeedSource{}, id).Error; err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	OK(c, nil)
}

// ManualFetch triggers a manual fetch for a single feed source (admin).
func (h *FeedHandler) ManualFetch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid feed source ID")
		return
	}

	var source model.FeedSource
	if err := h.db.First(&source, id).Error; err != nil {
		Fail(c, http.StatusNotFound, "feed source not found")
		return
	}

	count, fetchErr := h.feedSvc.FetchFeedSource(source)
	if fetchErr != nil {
		OK(c, gin.H{"fetched": count, "error": fetchErr.Error()})
		return
	}

	OK(c, gin.H{"fetched": count})
}

// ListFetchLogs returns paginated fetch logs.
func (h *FeedHandler) ListFetchLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := h.db.Model(&model.FetchLog{})

	if sourceID := c.Query("source_id"); sourceID != "" {
		sid, err := strconv.ParseUint(sourceID, 10, 64)
		if err == nil {
			query = query.Where("source_id = ?", sid)
		}
	}

	var total int64
	query.Count(&total)

	var logs []model.FetchLog
	offset := (page - 1) * pageSize
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs)

	OKList(c, logs, total, page, pageSize)
}
