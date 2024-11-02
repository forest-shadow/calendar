package events

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/forest-shadow/calendar/internal/logger"
)

type eventRepository interface {
	Create(ctx context.Context, event Event) error
	Get(ctx context.Context, id uuid.UUID) (Event, error)
	GetList(ctx context.Context, from time.Time, to time.Time) ([]Event, error)
	Update(ctx context.Context, id uuid.UUID, event UpdateEvent) error
	Delete(ctx context.Context, id string) error
	GetRecentEvents(ctx context.Context, now time.Time) (*[]Event, error)
	DeleteOldEvents(ctx context.Context, now time.Time) (bool, error)
}

type Service struct {
	repo   eventRepository
	logger logger.Logger
}

func NewEventsService(repo eventRepository, logger logger.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}

func (s *Service) Create(ctx context.Context, event Event) error {
	return s.repo.Create(ctx, event)
}

func (s *Service) Get(ctx context.Context, id uuid.UUID) (Event, error) {
	return s.repo.Get(ctx, id)
}

func (s *Service) GetList(ctx context.Context, from time.Time, to time.Time) ([]Event, error) {
	return s.repo.GetList(ctx, from, to)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, event UpdateEvent) error {
	return s.repo.Update(ctx, id, event)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetRecentEvents(ctx context.Context, from time.Time) (*[]Event, error) {
	events, err := s.repo.GetRecentEvents(ctx, from)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Service) DeleteOldEvents(ctx context.Context, date time.Time) (bool, error) {
	return s.repo.DeleteOldEvents(ctx, date)
}
