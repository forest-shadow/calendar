package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/forest-shadow/calendar/internal/config"
	"github.com/forest-shadow/calendar/internal/events"

	"github.com/twmb/franz-go/pkg/kgo"
)

type EventConsumer struct {
	ctx context.Context
	c   *kgo.Client
}

func NewEventConsumer(ctx context.Context, cfg *config.Config) (*EventConsumer, error) {
	seeds := []string{cfg.Kafka.BootstrapServers}
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(seeds...),
		kgo.ConsumeTopics(cfg.Kafka.Topic),
	)
	if err != nil {
		return nil, fmt.Errorf("create kafka event consumer client: %w", err)
	}

	return &EventConsumer{
		ctx: ctx,
		c:   cl,
	}, nil
}

func (ec *EventConsumer) Consume(itemHandler func(event events.Event)) error {
	for {
		fetches := ec.c.PollFetches(ec.ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			return fmt.Errorf("poll events fetches: %v", errs)
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			var event events.Event
			if err := json.Unmarshal(record.Value, &event); err != nil {
				return fmt.Errorf("unmarshal event: %w", err)
			}
			if itemHandler != nil {
				itemHandler(event)
			}
		}
	}
}

func (ec *EventConsumer) Close() {
	ec.c.Close()
}
