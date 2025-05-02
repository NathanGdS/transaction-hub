package consumers

import (
	"context"
	"encoding/json"

	"github.com/NathanGdS/cali-challenge/application/services"
	"github.com/NathanGdS/cali-challenge/domain"
	dRepo "github.com/NathanGdS/cali-challenge/domain/repository"
	"github.com/NathanGdS/cali-challenge/infra/akafka"
	"github.com/NathanGdS/cali-challenge/infra/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type TransactionConsumer struct {
	kafkaBroker akafka.KafkaBroker
	logger      *zap.Logger
	service     *services.TransactionService
}

func NewTransactionConsumer(broker *akafka.KafkaBroker, repository dRepo.TransactionRepository) *TransactionConsumer {
	return &TransactionConsumer{
		kafkaBroker: *broker,
		logger:      logger.Log,
		service:     services.NewTransactionService(*broker, repository),
	}
}

func (c *TransactionConsumer) Start() {
	msgChan := make(chan *kafka.Message)
	go c.kafkaBroker.Consume([]string{"process-transaction"}, msgChan)

	for msg := range msgChan {
		go c.processMessage(msg)
	}
}

func (c *TransactionConsumer) processMessage(msg *kafka.Message) {
	c.logger.Info("consumindo mensagem",
		zap.String("message", string(msg.Value)),
	)

	var transaction domain.Transaction
	err := json.Unmarshal(msg.Value, &transaction)
	if err != nil {
		c.logger.Error("erro ao converter para JSON",
			zap.Error(err),
		)
		return
	}

	c.service.ProcessTransaction(context.Background(), &transaction)
}
