package jobs

import (
	"context"
	"fmt"
	"time"

	"github.com/forest-shadow/calendar/internal/config"
	"github.com/forest-shadow/calendar/internal/logger"
)

type EventCleanerJob struct {
	ctx          context.Context
	cfg          *config.Config
	eventService eventService
	logger       logger.Logger
}

func NewEventCleanerJob(ctx context.Context, cfg *config.Config, eventService eventService, logger logger.Logger) *EventCleanerJob {
	return &EventCleanerJob{
		ctx:          ctx,
		cfg:          cfg,
		eventService: eventService,
		logger:       logger.With("component", "cleaner_job"),
	}
}

func (j *EventCleanerJob) Start(errCh chan error) {
	j.logger.Info("started")
	tickerInterval, err := time.ParseDuration(j.cfg.Cleaner.Interval)
	if err != nil {
		j.logger.Error("invalid interval", "error", err)
		return
	}

	timer := time.NewTicker(tickerInterval)

	go func() {
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				j.logger.Info("ticked")
				rowsAffected, err := j.eventService.DeleteOldEvents(j.ctx, time.Now())
				if rowsAffected {
					j.logger.Info("old events deleted")
				}
				if err != nil {
					errCh <- fmt.Errorf("error deleting old events: %w", err)
					return
				}
			case <-j.ctx.Done():
				errCh <- fmt.Errorf("scheduler job stopped due to context cancellation")
				return
			}
		}
	}()
}
