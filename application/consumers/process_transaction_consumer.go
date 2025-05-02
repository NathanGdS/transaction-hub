package consumers

import (
	"context"
	"encoding/json"

	"github.com/NathanGdS/cali-challenge/application/services"
	"github.com/NathanGdS/cali-challenge/domain/dto"
	dRepo "github.com/NathanGdS/cali-challenge/domain/repository"
	"github.com/NathanGdS/cali-challenge/infra/akafka"
	"github.com/NathanGdS/cali-challenge/infra/logger"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type ProcessTransactionConsumer struct {
	kafkaBroker akafka.KafkaBroker
	logger      *zap.Logger
	service     *services.TransactionService
}

func NewProcessTransactionConsumer(kafkaBroker *akafka.KafkaBroker, repository dRepo.TransactionRepository) *ProcessTransactionConsumer {
	return &ProcessTransactionConsumer{
		kafkaBroker: *kafkaBroker,
		logger:      logger.Log,
		service:     services.NewTransactionService(*kafkaBroker, repository),
	}
}

func (c *ProcessTransactionConsumer) Start() {
	msgChan := make(chan *kafka.Message)
	go c.kafkaBroker.Consume([]string{"transaction-process-return"}, msgChan)

	for msg := range msgChan {
		go c.processMessage(msg)
	}
}

func (c *ProcessTransactionConsumer) processMessage(msg *kafka.Message) {
	c.logger.Info("consumindo mensagem no tópico transaction-process-return",
		zap.String("message", string(msg.Value)),
	)

	var processTransactionDto dto.ProcessTransactionDto
	err := json.Unmarshal(msg.Value, &processTransactionDto)
	if err != nil {
		c.logger.Error("erro ao converter para JSON",
			zap.Error(err),
		)
	}

	transaction, err := c.service.FindByID(context.Background(), processTransactionDto.TransactionID)
	if err != nil {
		c.logger.Error("erro ao buscar transação",
			zap.Error(err),
		)
	}

	if processTransactionDto.Status == dto.TransactionStatusProcessed {
		transaction.TransactionProcessed()
	} else {
		transaction.ErrorProcessingTransaction(processTransactionDto.ErrorMessage)
	}
	go c.service.UpdateTransaction(context.Background(), transaction)
}
