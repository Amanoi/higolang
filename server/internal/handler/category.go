package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"higolang/server/internal/model"
)

// CategoryHandler handles category-related requests.
type CategoryHandler struct {
	db *gorm.DB
}

// NewCategoryHandler creates a new CategoryHandler.
func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

type categoryWithCount struct {
	model.Category
	ArticleCount int64 `json:"article_count"`
}

// ListCategories returns all categories with article counts.
func (h *CategoryHandler) ListCategories(c *gin.Context) {
	var categories []model.Category
	if err := h.db.Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	var result []categoryWithCount
	for _, cat := range categories {
		var count int64
		h.db.Model(&model.Article{}).Where("category_id = ? AND status = ?", cat.ID, "published").Count(&count)
		result = append(result, categoryWithCount{Category: cat, ArticleCount: count})
	}

	OK(c, result)
}

type categoryInput struct {
	Name      string `json:"name" binding:"required"`
	Slug      string `json:"slug" binding:"required"`
	SortOrder int    `json:"sort_order"`
}

// CreateCategory creates a new category (admin).
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var input categoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	cat := model.Category{
		Name:      input.Name,
		Slug:      input.Slug,
		SortOrder: input.SortOrder,
	}
	if err := h.db.Create(&cat).Error; err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	OK(c, cat)
}

// UpdateCategory updates a category (admin).
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid category ID")
		return
	}

	var input categoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	var cat model.Category
	if err := h.db.First(&cat, id).Error; err != nil {
		Fail(c, http.StatusNotFound, "category not found")
		return
	}

	h.db.Model(&cat).Updates(map[string]interface{}{
		"name":       input.Name,
		"slug":       input.Slug,
		"sort_order": input.SortOrder,
	})

	OK(c, cat)
}

// DeleteCategory deletes a category (admin).
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid category ID")
		return
	}

	if err := h.db.Delete(&model.Category{}, id).Error; err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	OK(c, nil)
}
