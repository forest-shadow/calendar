package events

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type repository interface {
	Create(ctx context.Context, event Event) error
	Delete(ctx context.Context, id string) error
	Get(ctx context.Context, id uuid.UUID) (Event, error)
	GetList(ctx context.Context, from time.Time, to time.Time) ([]Event, error)
	Update(ctx context.Context, id uuid.UUID, event EventUpdateDTO) error
}

type EventService struct {
	repo repository
}

func NewEventService(repo repository) *EventService {
	return &EventService{repo: repo}
}

func (s *EventService) Create(ctx context.Context, event Event) error {
	return s.repo.Create(ctx, event)
}

func (s *EventService) Get(ctx context.Context, id uuid.UUID) (Event, error) {
	return s.repo.Get(ctx, id)
}

func (s *EventService) GetList(ctx context.Context, from time.Time, to time.Time) ([]Event, error) {
	return s.repo.GetList(ctx, from, to)
}

func (s *EventService) Update(ctx context.Context, id uuid.UUID, event EventUpdateDTO) error {
	return s.repo.Update(ctx, id, event)
}

func (s *EventService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

