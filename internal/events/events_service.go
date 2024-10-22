package events

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/forest-shadow/calendar/internal/logger"
)

type EventService interface {
	repository

	StartEventJobs(ctx context.Context) chan error
	StartSchedulerJob(ctx context.Context, errCh chan error)
	StartNotifierJob(ctx context.Context, errCh chan error)
}

type EventManger struct {
	repo   EventRepository
	logger logger.Logger
}

func NewEventService(repo EventRepository, logger logger.Logger) EventService {
	return &EventManger{repo: repo, logger: logger}
}

func (s *EventManger) Create(ctx context.Context, event Event) error {
	return s.repo.Create(ctx, event)
}

func (s *EventManger) Get(ctx context.Context, id uuid.UUID) (Event, error) {
	return s.repo.Get(ctx, id)
}

func (s *EventManger) GetList(ctx context.Context, from time.Time, to time.Time) ([]Event, error) {
	return s.repo.GetList(ctx, from, to)
}

func (s *EventManger) Update(ctx context.Context, id uuid.UUID, event EventUpdateDTO) error {
	return s.repo.Update(ctx, id, event)
}

func (s *EventManger) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *EventManger) StartEventJobs(ctx context.Context) chan error {
	errCh := make(chan error, 2)

	s.StartSchedulerJob(ctx, errCh)
	s.StartNotifierJob(ctx, errCh)

	return errCh
}

func (s *EventManger) StartSchedulerJob(ctx context.Context, errCh chan error) {
	jobName := "Scheduler Job"
	s.logger.Info(fmt.Sprintf("%s: started", jobName))
	tickerInterval := 30 * time.Second
	// tickerInterval := 24 * time.Hour

	timer := time.NewTicker(tickerInterval)

	go func() {
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				s.logger.Info(fmt.Sprintf("%s: job ticked", jobName))
				rowsAffected, err := s.repo.DeleteOldEvents(ctx, time.Now())
				if rowsAffected {
					s.logger.Info(fmt.Sprintf("%s: old events deleted", jobName))
				}
				if err != nil {
					errCh <- fmt.Errorf("error deleting old events: %w", err)
					return
				}
			case <-ctx.Done():
				errCh <- fmt.Errorf("scheduler job stopped due to context cancellation")
				return
			}
		}
	}()
}

func (s *EventManger) StartNotifierJob(ctx context.Context, errCh chan error) {
	jobName := "Notifier Job"

	s.logger.Info(fmt.Sprintf("%s: started", jobName))
	notifierLogPath := "tmp/events.log"
	tickerInterval := 5 * time.Second

	ticker := time.NewTicker(tickerInterval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.logger.Info(fmt.Sprintf("%s: job ticked", jobName))
				events, err := s.repo.GetRecentEvents(ctx, time.Now())
				if err != nil {
					errCh <- fmt.Errorf("error getting events for notification: %w", err)
					return
				}

				for _, event := range *events {
					s.logger.Infof("%s: Event \"%s\" for user {%s} will be in %s", jobName, event.Title, event.UserID, *event.NotifyBefore)

					// write to file
					_, err := os.Stat(notifierLogPath)
					if os.IsNotExist(err) {
						_, err := os.Create(notifierLogPath)
						if err != nil {
							errCh <- fmt.Errorf("%s: Error creating file: %v", jobName, err)
							return
						}
					}
					f, err := os.OpenFile(notifierLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0o644))
					if err != nil {
						errCh <- fmt.Errorf("%s: Error opening file: %v", jobName, err)
						return
					}
					defer f.Close()
					_, err = f.WriteString(fmt.Sprintf("TS: %s. %s: Event \"%s\" for user {%s} will be in %s\n", time.Now().UTC().Format(time.RFC3339), jobName, event.Title, event.UserID, *event.NotifyBefore))
					if err != nil {
						errCh <- fmt.Errorf("%s: Error writing to file: %v", jobName, err)
						return
					}
				}
			case <-ctx.Done():
				errCh <- fmt.Errorf("Notifier job stopped due to context cancellation")
				return
			}
		}
	}()
}
