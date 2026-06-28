package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"higolang/server/config"
	"higolang/server/internal/model"
	"higolang/server/internal/router"
	"higolang/server/internal/scheduler"
	"higolang/server/internal/service"
)

func main() {
	// Load configuration.
	cfg := config.Load()

	// Initialize logger.
	var logger *zap.Logger
	var err error
	if cfg.Server.Mode == "release" {
		logger, err = zap.NewProduction()
	} else {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Ensure data directory exists.
	if err := os.MkdirAll("data", 0755); err != nil {
		logger.Fatal("Failed to create data directory", zap.Error(err))
	}

	// Initialize database.
	db := model.InitDB(cfg.Database, cfg.Server.Mode)
	logger.Info("database initialized", zap.String("path", cfg.Database.Path))

	// Seed default data.
	if cfg.Seed.Enabled {
		seedData(db, logger)
	}

	// Start scheduler.
	feedSvc := service.NewFeedService(db)
	if cfg.Scheduler.Enabled {
		sched := scheduler.New(db, feedSvc, logger)
		if err := sched.Start(); err != nil {
			logger.Error("scheduler failed to start", zap.Error(err))
		}
		defer sched.Stop()
	}

	// Setup router and start server.
	r := router.Setup(cfg, db)
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("starting server", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
}

// seedData creates default admin, categories, and feed sources if the database is empty.
func seedData(db *gorm.DB, logger *zap.Logger) {
	// Check if admin already exists.
	var adminCount int64
	db.Model(&model.Admin{}).Count(&adminCount)
	if adminCount > 0 {
		return
	}

	logger.Info("seeding default data")

	// Create default admin: admin / admin123.
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := model.Admin{
		Username: "admin",
		Password: string(hashedPassword),
	}
	db.Create(&admin)

	// Create default categories.
	categories := []model.Category{
		{Name: "版本更新", Slug: "releases", SortOrder: 1},
		{Name: "官方博客", Slug: "official-blog", SortOrder: 2},
		{Name: "社区动态", Slug: "community", SortOrder: 3},
		{Name: "技术教程", Slug: "tutorials", SortOrder: 4},
		{Name: "行业资讯", Slug: "industry-news", SortOrder: 5},
	}
	for _, cat := range categories {
		db.Create(&cat)
	}

	// Find category IDs for feed sources.
	var releasesCat, blogCat, communityCat model.Category
	db.Where("slug = ?", "releases").First(&releasesCat)
	db.Where("slug = ?", "official-blog").First(&blogCat)
	db.Where("slug = ?", "community").First(&communityCat)

	// Create default feed sources.
	sources := []model.FeedSource{
		{
			Name:          "Go Official Blog",
			URL:           "https://go.dev/blog/feed.atom",
			FeedType:      "atom",
			CategoryID:    blogCat.ID,
			FetchInterval: 60,
			IsEnabled:     true,
		},
		{
			Name:          "Go GitHub Releases",
			URL:           "https://github.com/golang/go/releases.atom",
			FeedType:      "atom",
			CategoryID:    releasesCat.ID,
			FetchInterval: 120,
			IsEnabled:     true,
		},
		{
			Name:          "Go Release JSON",
			URL:           "https://go.dev/dl/?mode=json",
			FeedType:      "json",
			CategoryID:    releasesCat.ID,
			FetchInterval: 360,
			IsEnabled:     true,
		},
		{
			Name:          "Reddit r/golang",
			URL:           "https://www.reddit.com/r/golang/.rss",
			FeedType:      "rss",
			CategoryID:    communityCat.ID,
			FetchInterval: 60,
			IsEnabled:     true,
		},
	}
	for _, src := range sources {
		db.Create(&src)
	}

	// Create default site settings.
	settings := []model.SiteSetting{
		{Key: "site_name", Value: "HiGolang"},
		{Key: "description", Value: "Go 语言最新动态资讯聚合"},
		{Key: "logo_url", Value: ""},
		{Key: "meta_description", Value: "HiGolang - 关注 Go 语言最新版本、官方博客、社区动态和技术教程"},
		{Key: "meta_keywords", Value: "Go,Golang,Go语言,Go教程,Go新闻"},
	}
	for _, s := range settings {
		db.Create(&s)
	}

	// Seed a welcome article.
	now := time.Now()
	welcome := model.Article{
		Title:      "欢迎来到 HiGolang",
		Slug:       service.GenerateSlug("welcome-to-higolang"),
		Summary:    "HiGolang 是一个专注于 Go 语言最新动态的资讯聚合站点。我们自动收集来自 Go 官方博客、GitHub 发布、社区讨论等来源的最新内容。",
		Content:    "# 欢迎！\n\nHiGolang 自动聚合以下来源的 Go 语言资讯：\n\n- Go 官方博客 (go.dev/blog)\n- Go GitHub Releases\n- Reddit r/golang\n- Go 社区周刊\n\n您可以在后台管理面板中配置更多数据源。",
		Status:     "published",
		IsPinned:   true,
		PublishedAt: &now,
	}
	db.Create(&welcome)

	logger.Info("seed data created: admin/admin123, 5 categories, 4 feed sources")
}
