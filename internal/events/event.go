package events

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID           uuid.UUID  `json:"id"`
	Title        string     `json:"title"`
	StartTime    time.Time  `json:"start_time"`
	EndTime      time.Time  `json:"end_time"`
	Description  *string    `json:"description,omitempty"`
	UserID       uuid.UUID  `json:"user_id"`
	NotifyBefore *string    `json:"notify_before,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	NotifiedAt   *time.Time `json:"notified_at,omitempty"`
}

type UpdateEvent = map[string]any
