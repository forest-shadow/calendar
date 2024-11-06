package jobs

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/forest-shadow/calendar/internal/config"
	"github.com/forest-shadow/calendar/internal/events"
	"github.com/forest-shadow/calendar/internal/logger"
)

type eventService interface {
	GetRecentEvents(ctx context.Context, from time.Time) (*[]events.Event, error)
	DeleteOldEvents(ctx context.Context, date time.Time) (bool, error)
	MarkEventsAsNotified(ctx context.Context, ids []uuid.UUID, notifiedAt time.Time) error
}

type EventNotifierJob struct {
	ctx           context.Context
	cfg           *config.Config
	eventsService eventService
	logger        logger.Logger
}

func NewEventNotifier(ctx context.Context, cfg *config.Config, eventsService eventService, logger logger.Logger) *EventNotifierJob {
	return &EventNotifierJob{
		cfg:           cfg,
		ctx:           ctx,
		eventsService: eventsService,
		logger:        logger.With("component", "notifier_job"),
	}
}

func (j *EventNotifierJob) Start(errCh chan error) {
	j.logger.Info("started")
	notifierLogPath := j.cfg.Notifier.LogPath
	tickerInterval, err := time.ParseDuration(j.cfg.Notifier.Interval)
	if err != nil {
		errCh <- fmt.Errorf("invalid interval: %w", err)
		return
	}

	ticker := time.NewTicker(tickerInterval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				j.logger.Info("tick")
				events, err := j.eventsService.GetRecentEvents(j.ctx, time.Now())
				if err != nil {
					errCh <- fmt.Errorf("getting events for notification: %w", err)
					return
				}

				for _, event := range *events {
					j.logger.Infof("Event \"%s\" for user {%s} will be in %s", event.Title, event.UserID, *event.NotifyBefore)
					j.logger.With(
						"event_id", event.ID,
						"event_title", event.Title,
						"event_user_id", event.UserID,
						"event_notify_before", *event.NotifyBefore,
					).Info("notification")

					// write to file
					_, err := os.Stat(notifierLogPath)
					if os.IsNotExist(err) {
						_, err := os.Create(notifierLogPath)
						if err != nil {
							errCh <- fmt.Errorf("creating file: %v", err)
							return
						}
					}
					f, err := os.OpenFile(notifierLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0o644))
					if err != nil {
						errCh <- fmt.Errorf("opening file: %v", err)
						return
					}
					defer f.Close()
					_, err = f.WriteString(fmt.Sprintf("TS: %s. Event \"%s\" for user {%s} will be in %s\n", time.Now().UTC().Format(time.RFC3339), event.Title, event.UserID, *event.NotifyBefore))
					if err != nil {
						errCh <- fmt.Errorf("Error writing to file: %v", err)
						return
					}

					err = j.eventsService.MarkEventsAsNotified(j.ctx, []uuid.UUID{event.ID}, time.Now())
					if err != nil {
						errCh <- fmt.Errorf("marking event as notified: %w", err)
						return
					}
				}
			case <-j.ctx.Done():
				errCh <- fmt.Errorf("Notifier job stopped due to context cancellation")
				return
			}
		}
	}()
}
