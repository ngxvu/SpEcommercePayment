package bootstrap

import (
	"context"
	"log"
	"payment/pkg/core/configloader"
	"payment/pkg/core/kafka"
	"payment/pkg/core/kafka/payment"
)

func InitKafka(parent context.Context) (*kafka.App, func()) {
	cfg := configloader.GetConfig()

	producer := payment.NewPaymentProducer(cfg.KafkaBrokers, cfg.KafkaTopicPaymentAuthorized)
	log.Println("Kafka producer created:", cfg.KafkaBrokers, "topic:", cfg.KafkaTopicPaymentAuthorized)
	//consumer := kafka.NewConsumer(cfg.KafkaBrokers, cfg.KafkaTopicOrder, "order_group")
	//
	//ctx, cancel := context.WithCancel(parent)
	//
	//go consumer.Listen(ctx, func(data []byte) {
	//	log.Printf("Received: %s", string(data))
	//})

	stop := func() {
		//cancel()
		if err := producer.Close(); err != nil {
			log.Printf("producer close error: %v", err)
		}
		//if err := consumer.Reader.Close(); err != nil {
		//	log.Printf("consumer close error: %v", err)
		//}
	}

	return &kafka.App{
		Producer: producer,
		//Consumer: consumer,
	}, stop
}
