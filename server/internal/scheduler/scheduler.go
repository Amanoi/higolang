package scheduler

import (
	"time"

	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"higolang/server/internal/model"
	"higolang/server/internal/service"
)

// Scheduler manages periodic feed fetching.
type Scheduler struct {
	sched   gocron.Scheduler
	feedSvc *service.FeedService
	db      *gorm.DB
	logger  *zap.Logger
}

// New creates a new Scheduler.
func New(db *gorm.DB, feedSvc *service.FeedService, logger *zap.Logger) *Scheduler {
	return &Scheduler{db: db, feedSvc: feedSvc, logger: logger}
}

// Start initializes and starts the scheduler.
func (s *Scheduler) Start() error {
	sched, err := gocron.NewScheduler()
	if err != nil {
		return err
	}
	s.sched = sched

	// Schedule fetching for each enabled source.
	var sources []model.FeedSource
	s.db.Where("is_enabled = ?", true).Find(&sources)

	for _, src := range sources {
		interval := time.Duration(src.FetchInterval) * time.Minute
		if interval <= 0 {
			interval = 30 * time.Minute
		}

		sourceID := src.ID
		sourceName := src.Name

		_, err := sched.NewJob(
			gocron.DurationJob(interval),
			gocron.NewTask(func() {
				s.logger.Info("fetching feed", zap.Uint("source_id", sourceID), zap.String("name", sourceName))
				var source model.FeedSource
				if err := s.db.First(&source, sourceID).Error; err != nil {
					s.logger.Error("source not found", zap.Uint("source_id", sourceID))
					return
				}
				count, err := s.feedSvc.FetchFeedSource(source)
				if err != nil {
					s.logger.Error("fetch failed", zap.Uint("source_id", sourceID), zap.Error(err))
				} else {
					s.logger.Info("fetch complete", zap.Uint("source_id", sourceID), zap.Int("new_articles", count))
				}
			}),
		)
		if err != nil {
			s.logger.Error("failed to schedule job", zap.Error(err))
		}
	}

	// Also schedule a global fetch-all every 30 minutes as fallback.
	_, err = sched.NewJob(
		gocron.DurationJob(30*time.Minute),
		gocron.NewTask(func() {
			s.logger.Info("running global fetch-all")
			s.feedSvc.FetchAllEnabled()
		}),
	)
	if err != nil {
		s.logger.Error("failed to schedule global fetch", zap.Error(err))
	}

	s.sched.Start()
	s.logger.Info("scheduler started", zap.Int("sources", len(sources)))
	return nil
}

// Stop gracefully stops the scheduler.
func (s *Scheduler) Stop() {
	if s.sched != nil {
		_ = s.sched.Shutdown()
		s.logger.Info("scheduler stopped")
	}
}
