package payment

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewPaymentProducer(brokers []string, topic string) *Producer {
	log.Printf("Creating Kafka producer for brokers: %v, topic: %s", brokers, topic)
	return &Producer{
		Writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    topic,
			Balancer: &kafka.Hash{},
		},
	}
}

func (p *Producer) SendMessage(ctx context.Context, key, value string) error {

	log.Printf("Sending message to topic %s: key=%s, value=%s", p.Writer.Topic, key, value)

	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}

	return p.Writer.WriteMessages(ctx, msg)
}

func (p *Producer) Close() error {
	if p.Writer == nil {
		return nil
	}
	return p.Writer.Close()
}
