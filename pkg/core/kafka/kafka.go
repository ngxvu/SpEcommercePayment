package kafka

import (
	"github.com/segmentio/kafka-go"
	"payment/pkg/core/kafka/payment"
	"strings"
)

type App struct {
	Producer *payment.Producer
	Consumer *Consumer
}

func NewWriter(brokers, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(strings.Split(brokers, ",")...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func Close(writer *kafka.Writer, reader *kafka.Reader) {
	if writer != nil {
		writer.Close()
	}
	if reader != nil {
		reader.Close()
	}
}
