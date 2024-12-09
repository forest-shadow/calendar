package jobs

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/forest-shadow/calendar/internal/config"
	"github.com/forest-shadow/calendar/internal/events"

	"github.com/twmb/franz-go/pkg/kgo"
)

type EventProducer struct {
	ctx context.Context
	p   *kgo.Client
}

func NewEventProducer(ctx context.Context, cfg *config.Config) (*EventProducer, error) {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.Kafka.BootstrapServers),
		kgo.DefaultProduceTopic(cfg.Kafka.Topic),
		kgo.AllowAutoTopicCreation(),
		kgo.RecordRetries(1),
		kgo.RequiredAcks(kgo.NoAck()),
		kgo.DisableIdempotentWrite(),
	)
	if err != nil {
		return nil, fmt.Errorf("create kafka client: %w", err)
	}

	return &EventProducer{ctx: ctx, p: cl}, nil
}

func (p *EventProducer) SendEvent(ctx context.Context, event *events.Event) error {
	value, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}
	fmt.Println(string(value))
	rec := &kgo.Record{
		Key:   []byte(event.ID.String()),
		Value: value,
	}
	if err := p.p.ProduceSync(ctx, rec).FirstErr(); err != nil {
		return fmt.Errorf("produce event: %w", err)
	}
	return nil
}
