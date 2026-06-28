package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"higolang/server/internal/model"
)

// TagHandler handles tag-related requests.
type TagHandler struct {
	db *gorm.DB
}

// NewTagHandler creates a new TagHandler.
func NewTagHandler(db *gorm.DB) *TagHandler {
	return &TagHandler{db: db}
}

type tagWithCount struct {
	model.Tag
	ArticleCount int64 `json:"article_count"`
}

// ListTags returns all tags with article counts.
func (h *TagHandler) ListTags(c *gin.Context) {
	var tags []model.Tag
	if err := h.db.Order("id DESC").Find(&tags).Error; err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	var result []tagWithCount
	for _, tag := range tags {
		var count int64
		h.db.Table("article_tags").Where("tag_id = ?", tag.ID).Count(&count)
		result = append(result, tagWithCount{Tag: tag, ArticleCount: count})
	}

	OK(c, result)
}

type tagInput struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

// CreateTag creates a new tag (admin).
func (h *TagHandler) CreateTag(c *gin.Context) {
	var input tagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	tag := model.Tag{Name: input.Name, Slug: input.Slug}
	if err := h.db.Create(&tag).Error; err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	OK(c, tag)
}

// UpdateTag updates a tag (admin).
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid tag ID")
		return
	}

	var input tagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	var tag model.Tag
	if err := h.db.First(&tag, id).Error; err != nil {
		Fail(c, http.StatusNotFound, "tag not found")
		return
	}

	h.db.Model(&tag).Updates(map[string]interface{}{
		"name": input.Name,
		"slug": input.Slug,
	})

	OK(c, tag)
}

// DeleteTag deletes a tag (admin).
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid tag ID")
		return
	}

	if err := h.db.Delete(&model.Tag{}, id).Error; err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	OK(c, nil)
}
