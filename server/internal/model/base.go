package model

import (
	"log"
	"time"

	"github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"higolang/server/config"
)

// DB is the global database instance.
var DB *gorm.DB

// BaseModel provides common fields for all models.
type BaseModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// InitDB opens the SQLite database and runs auto-migration.
func InitDB(cfg config.DatabaseConfig, ginMode string) *gorm.DB {
	// Register the SQLite driver with foreign-key support.
	_ = sqlite3.SQLiteDriver{}
	gormsqlite := sqlite.New(sqlite.Config{
		DriverName: "sqlite3",
		DSN:        cfg.Path,
	})

	logLevel := logger.Warn
	if ginMode == "debug" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(gormsqlite, &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Enable WAL mode and foreign keys for SQLite.
	sqlDB, _ := db.DB()
	_, _ = sqlDB.Exec("PRAGMA journal_mode=WAL")
	_, _ = sqlDB.Exec("PRAGMA foreign_keys=ON")

	// Auto-migrate all models.
	if err := db.AutoMigrate(
		&Admin{},
		&Category{},
		&Tag{},
		&Article{},
		&FeedSource{},
		&FetchLog{},
		&SiteSetting{},
	); err != nil {
		log.Fatalf("Failed to auto-migrate: %v", err)
	}

	DB = db
	return db
}

// SiteSetting stores key-value configuration in the database.
type SiteSetting struct {
	Key       string    `gorm:"primaryKey" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	UpdatedAt time.Time `json:"updated_at"`
}
