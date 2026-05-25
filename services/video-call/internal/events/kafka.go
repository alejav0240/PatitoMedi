package events

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writers map[string]*kafka.Writer
}

func New(brokers []string) *Producer {
	topics := []string{"call-started", "call-ended", "call-failed"}
	writers := make(map[string]*kafka.Writer, len(topics))
	for _, t := range topics {
		writers[t] = &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        t,
			Balancer:     &kafka.LeastBytes{},
			WriteTimeout: 5 * time.Second,
		}
	}
	return &Producer{writers: writers}
}

func (p *Producer) publish(ctx context.Context, topic, appointmentID string, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("kafka marshal %s: %v", topic, err)
		return
	}
	w, ok := p.writers[topic]
	if !ok {
		return
	}
	if err := w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(appointmentID),
		Value: data,
	}); err != nil {
		log.Printf("kafka write %s: %v", topic, err)
	}
}

func (p *Producer) CallStarted(ctx context.Context, appointmentID string, participants []string) {
	p.publish(ctx, "call-started", appointmentID, map[string]any{
		"appointmentId": appointmentID,
		"participants":  participants,
		"startedAt":     time.Now().UTC().Format(time.RFC3339),
	})
}

func (p *Producer) CallEnded(ctx context.Context, appointmentID, initiatorID string) {
	p.publish(ctx, "call-ended", appointmentID, map[string]any{
		"appointmentId": appointmentID,
		"initiatorId":   initiatorID,
		"endedAt":       time.Now().UTC().Format(time.RFC3339),
	})
}

func (p *Producer) CallFailed(ctx context.Context, appointmentID, reason string) {
	p.publish(ctx, "call-failed", appointmentID, map[string]any{
		"appointmentId": appointmentID,
		"reason":        reason,
		"failedAt":      time.Now().UTC().Format(time.RFC3339),
	})
}

func (p *Producer) Close() {
	for _, w := range p.writers {
		w.Close()
	}
}
