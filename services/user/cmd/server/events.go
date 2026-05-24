package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type eventProducer struct {
	writer *kafka.Writer
}

func newEventProducer(cfg config) *eventProducer {
	if !cfg.KafkaEnabled || len(cfg.KafkaBrokers) == 0 {
		return &eventProducer{}
	}
	return &eventProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.KafkaBrokers...),
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
			Async:        true,
		},
	}
}

func (p *eventProducer) Publish(ctx context.Context, topic string, key string, payload any) {
	if p == nil || p.writer == nil {
		return
	}

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error al serializar evento: %v", err)
		return
	}
	if err := p.writer.WriteMessages(ctx, kafka.Message{Topic: topic, Key: []byte(key), Value: body, Time: time.Now()}); err != nil {
		log.Printf("Error al publicar mensaje en topic=%s: %v", topic, err)
	}
}

func (p *eventProducer) Close() {
	if p != nil && p.writer != nil {
		_ = p.writer.Close()
	}
}
