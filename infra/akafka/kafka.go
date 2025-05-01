package akafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/NathanGdS/cali-challenge/infra/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// KafkaBroker é a interface que define os métodos necessários para um broker Kafka
type KafkaBroker interface {
	Publish(topic string, message []byte) error
	Close() error
	Consume(topics []string, msgChan chan *kafka.Message)
}

// KafkaBrokerImpl é a implementação concreta do KafkaBroker
type KafkaBrokerImpl struct {
	brokerURL string
	writer    *kafka.Writer
	mu        sync.Mutex
	logger    *zap.Logger
}

func NewKafkaBroker(brokerURL string) KafkaBroker {
	return &KafkaBrokerImpl{
		brokerURL: brokerURL,
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      []string{brokerURL},
			Async:        true,
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    1,
			BatchTimeout: 0,
		}),
		logger: logger.Log,
	}
}

func (k *KafkaBrokerImpl) Publish(topic string, message []byte) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.writer.Topic = topic
	err := k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: message,
		},
	)

	if err != nil {
		k.logger.Error("erro ao publicar mensagem",
			zap.Error(err),
			zap.String("topic", topic),
		)
		return fmt.Errorf("erro ao publicar mensagem: %v", err)
	}

	k.logger.Info("mensagem publicada com sucesso",
		zap.String("topic", topic),
	)
	return nil
}

func (k *KafkaBrokerImpl) Close() error {
	return k.writer.Close()
}

func (k *KafkaBrokerImpl) Consume(topics []string, msgChan chan *kafka.Message) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{k.brokerURL},
		GroupID: "transaction-group",
		Topic:   topics[0],
		MaxWait: 1 * time.Second,
	})

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			k.logger.Error("erro ao ler mensagem",
				zap.Error(err),
				zap.String("topic", topics[0]),
			)
			continue
		}
		msgChan <- &msg
	}
}
