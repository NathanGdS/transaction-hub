package akafka

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

type KafkaBroker struct {
	brokerURL string
	writer    *kafka.Writer
	mu        sync.Mutex
}

func NewKafkaBroker(brokerURL string) *KafkaBroker {
	return &KafkaBroker{
		brokerURL: brokerURL,
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: []string{brokerURL},
			Async:   true,
		}),
	}
}

func (k *KafkaBroker) Publish(topic string, message []byte) error {
	k.mu.Lock()
	defer k.mu.Unlock()

	k.writer.Topic = topic
	err := k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: message,
		},
	)

	if err != nil {
		return fmt.Errorf("erro ao publicar mensagem: %v", err)
	}

	return nil
}

func (k *KafkaBroker) Close() error {
	return k.writer.Close()
}

func (k *KafkaBroker) Consume(topics []string, msgChan chan *kafka.Message) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{k.brokerURL},
		GroupID: "transaction-group",
		Topic:   topics[0],
	})

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Erro ao ler mensagem: %v", err)
			continue
		}
		msgChan <- &msg
	}
}
