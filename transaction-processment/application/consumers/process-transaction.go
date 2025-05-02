package consumers

import (
	"context"
	"encoding/json"

	"github.com/NathanGdS/cali-challenge/pkg/akafka"
	"github.com/NathanGdS/cali-challenge/pkg/logger"
	"github.com/NathanGdS/cali-challenge/transaction-ledger/domain"
	"github.com/NathanGdS/cali-challenge/transaction-processment/application/services"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type ProcessTransactionConsumer struct {
	kafkaBroker akafka.KafkaBroker
	logger      *zap.Logger
	service     *services.ProcessTransactionService
}

func NewProcessTransactionConsumer(broker *akafka.KafkaBroker) *ProcessTransactionConsumer {
	return &ProcessTransactionConsumer{
		kafkaBroker: *broker,
		logger:      logger.Log,
		service:     services.NewProcessTransactionService(*broker),
	}
}

func (c *ProcessTransactionConsumer) Start() {
	msgChan := make(chan *kafka.Message)
	go c.kafkaBroker.Consume([]string{"process-transaction"}, msgChan)

	for msg := range msgChan {
		go c.processMessage(msg)
	}
}

func (c *ProcessTransactionConsumer) processMessage(msg *kafka.Message) {
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
