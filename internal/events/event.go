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

type EventUpdateDTO struct {
	Title        *string    `json:"title,omitempty"`
	StartTime    *time.Time `json:"start_time,omitempty"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	Description  *string    `json:"description,omitempty"`
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	NotifyBefore *string    `json:"notify_before,omitempty"`
}

type changesBuilder EventUpdateDTO

func (e changesBuilder) ToMap() map[string]any {
	result := make(map[string]any)
	if e.Title != nil {
		result["title"] = *e.Title
	}
	if e.StartTime != nil {
		result["start_time"] = *e.StartTime
	}
	if e.EndTime != nil {
		result["end_time"] = *e.EndTime
	}
	if e.Description != nil {
		result["description"] = *e.Description
	}
	if e.UserID != nil {
		result["user_id"] = *e.UserID
	}
	if e.NotifyBefore != nil {
		result["notify_before"] = *e.NotifyBefore
	}
	return result
}

type EventsListDTO struct {
	From time.Time `query:"from"`
	To   time.Time `query:"to"`
}