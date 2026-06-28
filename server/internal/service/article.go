package service

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"higolang/server/internal/model"
)

// ArticleService handles article CRUD and querying.
type ArticleService struct {
	db *gorm.DB
}

// NewArticleService creates a new ArticleService.
func NewArticleService(db *gorm.DB) *ArticleService {
	return &ArticleService{db: db}
}

// ListParams holds query parameters for listing articles.
type ListParams struct {
	Page     int
	PageSize int
	Category string
	Tag      string
	Search   string
	Status   string // optional; empty means only "published"
}

// ListResult holds paginated article results.
type ListResult struct {
	Articles []model.Article `json:"list"`
	Total    int64           `json:"total"`
}

// List returns a paginated, filtered list of articles.
func (s *ArticleService) List(params ListParams) (*ListResult, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 10
	}

	query := s.db.Model(&model.Article{})

	// Default to published only unless an explicit status is given.
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	} else {
		query = query.Where("status = ?", "published")
	}

	if params.Category != "" {
		query = query.Joins("JOIN categories ON categories.id = articles.category_id").
			Where("categories.slug = ?", params.Category)
	}

	if params.Tag != "" {
		query = query.Joins("JOIN article_tags ON article_tags.article_id = articles.id").
			Joins("JOIN tags ON tags.id = article_tags.tag_id").
			Where("tags.slug = ?", params.Tag)
	}

	if params.Search != "" {
		keyword := "%" + params.Search + "%"
		query = query.Where("articles.title LIKE ? OR articles.content LIKE ?", keyword, keyword)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var articles []model.Article
	offset := (params.Page - 1) * params.PageSize
	if err := query.Preload("Category").Preload("Tags").
		Order("is_pinned DESC, published_at DESC, created_at DESC").
		Offset(offset).Limit(params.PageSize).
		Find(&articles).Error; err != nil {
		return nil, err
	}

	return &ListResult{Articles: articles, Total: total}, nil
}

// GetBySlug retrieves a single article by its slug.
func (s *ArticleService) GetBySlug(slug string) (*model.Article, error) {
	var article model.Article
	if err := s.db.Preload("Category").Preload("Tags").
		Where("slug = ?", slug).First(&article).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// GetByID retrieves a single article by ID.
func (s *ArticleService) GetByID(id uint) (*model.Article, error) {
	var article model.Article
	if err := s.db.Preload("Category").Preload("Tags").
		First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// CreateInput holds fields for creating a new article.
type CreateInput struct {
	Title      string   `json:"title" binding:"required"`
	Summary    string   `json:"summary"`
	Content    string   `json:"content"`
	CoverURL   string   `json:"cover_url"`
	SourceURL  string   `json:"source_url"`
	SourceType string   `json:"source_type"`
	CategoryID uint     `json:"category_id"`
	TagIDs     []uint   `json:"tag_ids"`
	IsPinned   bool     `json:"is_pinned"`
	Status     string   `json:"status"`
}

// Create creates a new article and associates tags.
func (s *ArticleService) Create(input CreateInput) (*model.Article, error) {
	slug := GenerateSlug(input.Title)

	now := time.Now()
	article := model.Article{
		Title:       input.Title,
		Slug:        slug,
		Summary:     input.Summary,
		Content:     input.Content,
		CoverURL:    input.CoverURL,
		SourceURL:   input.SourceURL,
		SourceType:  input.SourceType,
		CategoryID:  input.CategoryID,
		IsPinned:    input.IsPinned,
		Status:      input.Status,
		PublishedAt: &now,
	}

	if article.Status == "" {
		article.Status = "published"
	}
	if article.SourceType == "" {
		article.SourceType = "manual"
	}

	if err := s.db.Create(&article).Error; err != nil {
		return nil, err
	}

	if len(input.TagIDs) > 0 {
		var tags []model.Tag
		s.db.Where("id IN ?", input.TagIDs).Find(&tags)
		if err := s.db.Model(&article).Association("Tags").Replace(tags); err != nil {
			return nil, err
		}
	}

	return &article, nil
}

// UpdateInput holds fields for updating an article.
type UpdateInput struct {
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	Content    string   `json:"content"`
	CoverURL   string   `json:"cover_url"`
	SourceURL  string   `json:"source_url"`
	SourceType string   `json:"source_type"`
	CategoryID uint     `json:"category_id"`
	TagIDs     []uint   `json:"tag_ids"`
	IsPinned   bool     `json:"is_pinned"`
	Status     string   `json:"status"`
}

// Update modifies an existing article.
func (s *ArticleService) Update(id uint, input UpdateInput) (*model.Article, error) {
	var article model.Article
	if err := s.db.First(&article, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"title":       input.Title,
		"summary":     input.Summary,
		"content":     input.Content,
		"cover_url":   input.CoverURL,
		"source_url":  input.SourceURL,
		"source_type": input.SourceType,
		"category_id": input.CategoryID,
		"is_pinned":   input.IsPinned,
		"status":      input.Status,
	}
	if err := s.db.Model(&article).Updates(updates).Error; err != nil {
		return nil, err
	}

	if input.TagIDs != nil {
		var tags []model.Tag
		s.db.Where("id IN ?", input.TagIDs).Find(&tags)
		if err := s.db.Model(&article).Association("Tags").Replace(tags); err != nil {
			return nil, err
		}
	}

	return s.GetByID(id)
}

// Delete removes an article by ID.
func (s *ArticleService) Delete(id uint) error {
	return s.db.Select("Tags").Delete(&model.Article{}, id).Error
}

// IncrementViewCount atomically increments the view count for an article.
func (s *ArticleService) IncrementViewCount(id uint) {
	s.db.Model(&model.Article{}).Where("id = ?", id).
		UpdateColumn("view_count", gorm.Expr("view_count + 1"))
}

// GenerateSlug creates a URL-friendly slug from a title.
// For non-ASCII titles (e.g. Chinese), a UUID-based slug is used.
func GenerateSlug(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))

	// Check if the title is primarily ASCII.
	isASCII := true
	for _, r := range slug {
		if r > 127 {
			isASCII = false
			break
		}
	}

	if isASCII {
		var b strings.Builder
		for _, r := range slug {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
				b.WriteRune(r)
			} else if r == ' ' || r == '-' || r == '_' {
				b.WriteByte('-')
			}
		}
		result := b.String()
		// Collapse multiple hyphens.
		for strings.Contains(result, "--") {
			result = strings.ReplaceAll(result, "--", "-")
		}
		result = strings.Trim(result, "-")
		if result == "" {
			return uuid.New().String()[:8]
		}
		if len(result) > 200 {
			result = result[:200]
		}
		return result
	}

	// Non-ASCII: use UUID-based slug.
	return uuid.New().String()[:8]
}
