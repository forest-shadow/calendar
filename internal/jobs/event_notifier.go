package jobs

import (
	"context"
	"fmt"
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

func NewEventNotifier(ctx context.Context, cfg *config.Config, eventsService eventService) (*EventNotifierJob, error) {
	logger, err := logger.NewNotifierLogger(cfg)
	if err != nil {
		return nil, fmt.Errorf("create logger: %w", err)
	}
	return &EventNotifierJob{
		cfg:           cfg,
		ctx:           ctx,
		eventsService: eventsService,
		logger:        logger,
	}, nil
}

func (j *EventNotifierJob) Start(errCh chan error) {
	j.logger.Info("started")
	tickerInterval, err := time.ParseDuration(j.cfg.Notifier.Interval)
	if err != nil {
		errCh <- fmt.Errorf("invalid interval: %w", err)
		return
	}

	ticker := time.NewTicker(tickerInterval)

	eventProducer, err := NewEventProducer(j.ctx, j.cfg)
	if err != nil {
		errCh <- fmt.Errorf("create event producer: %w", err)
		return
	}

	eventConsumer, err := NewEventConsumer(j.ctx, j.cfg)
	if err != nil {
		errCh <- fmt.Errorf("create event consumer: %w", err)
		return
	}

	go func() {
		err = eventConsumer.Consume(func(event events.Event) {
			j.logger.With(
				"event_id", event.ID,
				"event_title", event.Title,
				"event_user_id", event.UserID,
				"event_notify_before", *event.NotifyBefore,
				"notification_msg", fmt.Sprintf("Event \"%s\" for user {%s} will be in %s", event.Title, event.UserID, *event.NotifyBefore),
			).Info("notification")

			err = j.eventsService.MarkEventsAsNotified(j.ctx, []uuid.UUID{event.ID}, time.Now())
			if err != nil {
				errCh <- fmt.Errorf("marking event as notified: %w", err)
				return
			}
		})
		if err != nil {
			errCh <- fmt.Errorf("consuming events: %w", err)
			return
		}
	}()

	go func() {
		defer ticker.Stop()
		defer eventConsumer.Close()
		defer eventProducer.Close()
		for {
			select {
			case <-ticker.C:
				j.logger.Info("tick")
				recentEvents, err := j.eventsService.GetRecentEvents(j.ctx, time.Now())
				if err != nil {
					errCh <- fmt.Errorf("getting events for notification: %w", err)
					return
				}

				for _, event := range *recentEvents {
					if err := eventProducer.SendEvent(j.ctx, &event); err != nil {
						errCh <- fmt.Errorf("sending event: %w", err)
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
