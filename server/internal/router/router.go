package router

import (
	"github.com/gin-gonic/gin"

	"higolang/server/config"
	"higolang/server/internal/handler"
	"higolang/server/internal/middleware"
	"higolang/server/internal/service"
	"gorm.io/gorm"
)

// Setup registers all routes and returns the Gin engine.
func Setup(cfg config.Config, db *gorm.DB) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	// Middleware.
	r.Use(middleware.CORS())

	// Services.
	authSvc := service.NewAuthService(db, cfg.JWT)
	articleSvc := service.NewArticleService(db)
	feedSvc := service.NewFeedService(db)

	// Handlers.
	authHandler := handler.NewAuthHandler(authSvc)
	articleHandler := handler.NewArticleHandler(articleSvc)
	categoryHandler := handler.NewCategoryHandler(db)
	tagHandler := handler.NewTagHandler(db)
	feedHandler := handler.NewFeedHandler(db, feedSvc)
	dashboardHandler := handler.NewDashboardHandler(db)
	settingHandler := handler.NewSettingHandler(db)

	// Public API.
	v1 := r.Group("/api/v1")
	{
		v1.GET("/articles", articleHandler.ListArticles)
		v1.GET("/articles/search", articleHandler.ListArticles) // search uses same handler with query param
		v1.GET("/articles/:slug", articleHandler.GetArticle)
		v1.GET("/categories", categoryHandler.ListCategories)
		v1.GET("/tags", tagHandler.ListTags)
		v1.GET("/go-version", articleHandler.GetGoVersion)
		v1.GET("/settings/public", settingHandler.GetPublicSettings)
	}

	// Admin API.
	admin := v1.Group("/admin")
	{
		admin.POST("/login", authHandler.Login)

		// Auth-protected routes.
		protected := admin.Group("")
		protected.Use(middleware.Auth(authSvc))
		{
			protected.GET("/dashboard", dashboardHandler.GetDashboard)

			// Articles.
			protected.GET("/articles", articleHandler.AdminListArticles)
			protected.GET("/articles/:id", articleHandler.AdminGetArticle)
			protected.POST("/articles", articleHandler.AdminCreateArticle)
			protected.PUT("/articles/:id", articleHandler.AdminUpdateArticle)
			protected.DELETE("/articles/:id", articleHandler.AdminDeleteArticle)

			// Categories.
			protected.GET("/categories", categoryHandler.ListCategories)
			protected.POST("/categories", categoryHandler.CreateCategory)
			protected.PUT("/categories/:id", categoryHandler.UpdateCategory)
			protected.DELETE("/categories/:id", categoryHandler.DeleteCategory)

			// Tags.
			protected.GET("/tags", tagHandler.ListTags)
			protected.POST("/tags", tagHandler.CreateTag)
			protected.PUT("/tags/:id", tagHandler.UpdateTag)
			protected.DELETE("/tags/:id", tagHandler.DeleteTag)

			// Feed sources.
			protected.GET("/feed-sources", feedHandler.ListFeedSources)
			protected.POST("/feed-sources", feedHandler.CreateFeedSource)
			protected.PUT("/feed-sources/:id", feedHandler.UpdateFeedSource)
			protected.DELETE("/feed-sources/:id", feedHandler.DeleteFeedSource)
			protected.POST("/feed-sources/:id/fetch", feedHandler.ManualFetch)

			// Fetch logs.
			protected.GET("/fetch-logs", feedHandler.ListFetchLogs)

			// Settings.
			protected.GET("/settings", settingHandler.AdminGetSettings)
			protected.PUT("/settings", settingHandler.AdminUpdateSettings)
		}
	}

	return r
}
