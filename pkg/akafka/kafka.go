package akafka

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/NathanGdS/transaction-hub/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// KafkaBroker é a interface que define os métodos necessários para um broker Kafka
type KafkaBroker interface {
	Publish(topic string, message []byte) error
	Close() error
	Consume(topics []string, msgChan chan *kafka.Message)
	CreateTopicsIfNotExists(topics []string) error
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
			BatchSize:    100,
			BatchTimeout: 5,
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
		Brokers:           []string{k.brokerURL},
		GroupID:           "transaction-group",
		Topic:             topics[0],
		MaxWait:           1 * time.Second,
		HeartbeatInterval: 5 * time.Second,
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

func (k *KafkaBrokerImpl) CreateTopicsIfNotExists(topics []string) error {
	conn, err := kafka.Dial("tcp", k.brokerURL)
	if err != nil {
		return fmt.Errorf("erro ao conectar com kafka: %v", err)
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return fmt.Errorf("erro ao obter controller: %v", err)
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return fmt.Errorf("erro ao conectar com controller: %v", err)
	}
	defer controllerConn.Close()

	for _, topic := range topics {
		topicConfigs := []kafka.TopicConfig{
			{
				Topic:             topic,
				NumPartitions:     1,
				ReplicationFactor: 1,
			},
		}

		err = controllerConn.CreateTopics(topicConfigs...)
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			// ignore error if topic already exists
			k.logger.Info("topico já existe",
				zap.String("topic", topic),
			)
		}
	}

	return nil
}
