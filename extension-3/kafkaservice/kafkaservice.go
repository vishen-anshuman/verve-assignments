package kafkaservice

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaService struct {
	Writer *kafka.Writer
}

func InitKafkaService(brokers []string, topic string) *KafkaService {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaService{
		Writer: writer,
	}
}

func (ks *KafkaService) WriteLog(message string) error {
	err := ks.Writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: []byte(message),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to write message to Kafka: %v", err)
	}
	log.Printf("Message sent to Kafka: Value=%s", message)
	return nil
}
