package consumers

import (
	"encoding/json"

	"github.com/NathanGdS/cali-challenge/models"
	"github.com/NathanGdS/cali-challenge/pkg/akafka"
	"github.com/NathanGdS/cali-challenge/pkg/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type TransactionConsumer struct {
	kafkaBroker *akafka.KafkaBroker
	logger      *zap.Logger
}

func NewTransactionConsumer(broker *akafka.KafkaBroker) *TransactionConsumer {
	return &TransactionConsumer{
		kafkaBroker: broker,
		logger:      logger.Log,
	}
}

func (c *TransactionConsumer) Start() {
	msgChan := make(chan *kafka.Message)
	go c.kafkaBroker.Consume([]string{"test"}, msgChan)

	for msg := range msgChan {
		c.processMessage(msg)
	}
}

func (c *TransactionConsumer) processMessage(msg *kafka.Message) {
	c.logger.Info("consumindo mensagem",
		zap.String("message", string(msg.Value)),
	)

	var transaction models.Transaction
	err := json.Unmarshal(msg.Value, &transaction)
	if err != nil {
		c.logger.Error("erro ao converter para JSON",
			zap.Error(err),
		)
		return
	}

	c.logger.Info("transação recebida",
		zap.Any("transaction", transaction),
	)
}
