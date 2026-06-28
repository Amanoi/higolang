package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"higolang/server/internal/service"
)

// ArticleHandler handles article-related requests.
type ArticleHandler struct {
	articleSvc *service.ArticleService
}

// NewArticleHandler creates a new ArticleHandler.
func NewArticleHandler(articleSvc *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleSvc: articleSvc}
}

// ListArticles returns a paginated list of published articles (public).
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	params := service.ListParams{
		Page:     page,
		PageSize: pageSize,
		Category: c.Query("category"),
		Tag:      c.Query("tag"),
		Search:   c.Query("search"),
	}

	result, err := h.articleSvc.List(params)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	OKList(c, result.Articles, result.Total, page, pageSize)
}

// GetArticle returns a single article by slug (public).
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	slug := c.Param("slug")
	article, err := h.articleSvc.GetBySlug(slug)
	if err != nil {
		Fail(c, http.StatusNotFound, "article not found")
		return
	}

	// Increment view count asynchronously.
	go h.articleSvc.IncrementViewCount(article.ID)

	OK(c, article)
}

// AdminListArticles returns a paginated list of all articles (admin).
func (h *ArticleHandler) AdminListArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	params := service.ListParams{
		Page:     page,
		PageSize: pageSize,
		Category: c.Query("category"),
		Tag:      c.Query("tag"),
		Search:   c.Query("search"),
		Status:   c.Query("status"), // empty = all for admin
	}

	result, err := h.articleSvc.List(params)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	OKList(c, result.Articles, result.Total, page, pageSize)
}

// AdminGetArticle returns a single article by ID (admin).
func (h *ArticleHandler) AdminGetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid article ID")
		return
	}

	article, err := h.articleSvc.GetByID(uint(id))
	if err != nil {
		Fail(c, http.StatusNotFound, "article not found")
		return
	}

	OK(c, article)
}

// AdminCreateArticle creates a new article (admin).
func (h *ArticleHandler) AdminCreateArticle(c *gin.Context) {
	var input service.CreateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	article, err := h.articleSvc.Create(input)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	OK(c, article)
}

// AdminUpdateArticle updates an existing article (admin).
func (h *ArticleHandler) AdminUpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid article ID")
		return
	}

	var input service.UpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	article, err := h.articleSvc.Update(uint(id), input)
	if err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	OK(c, article)
}

// AdminDeleteArticle deletes an article (admin).
func (h *ArticleHandler) AdminDeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		Fail(c, http.StatusBadRequest, "invalid article ID")
		return
	}

	if err := h.articleSvc.Delete(uint(id)); err != nil {
		Fail(c, http.StatusInternalServerError, err.Error())
		return
	}

	OK(c, nil)
}

// GetGoVersion returns the latest Go version from the database.
func (h *ArticleHandler) GetGoVersion(c *gin.Context) {
	result, err := h.articleSvc.List(service.ListParams{
		Page:     1,
		PageSize: 1,
		Category: "releases",
	})
	if err != nil || len(result.Articles) == 0 {
		OK(c, gin.H{"version": "unknown"})
		return
	}
	OK(c, gin.H{"version": result.Articles[0].Title, "url": result.Articles[0].SourceURL})
}
