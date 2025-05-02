package services

import (
	"context"

	"github.com/NathanGdS/cali-challenge/domain"
	"github.com/NathanGdS/cali-challenge/domain/dto"
	dRepo "github.com/NathanGdS/cali-challenge/domain/repository"
	"github.com/NathanGdS/cali-challenge/infra/akafka"
	"github.com/NathanGdS/cali-challenge/infra/logger"
	"go.uber.org/zap"
)

type TransactionService struct {
	kafkaBroker akafka.KafkaBroker
	logger      *zap.Logger
	repository  dRepo.TransactionRepository
}

func NewTransactionService(kafkaBroker akafka.KafkaBroker, repository dRepo.TransactionRepository) *TransactionService {
	return &TransactionService{kafkaBroker: kafkaBroker, logger: logger.Log, repository: repository}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, transactionDto *dto.TransactionRequestDto) (domain.Transaction, []error) {
	transaction, errs := domain.NewTransaction(transactionDto.Amount, transactionDto.PaymentMethod, transactionDto.CurrencyCode, transactionDto.Description)
	if len(errs) > 0 {
		return domain.Transaction{}, errs
	}

	jsonData, err := transaction.ToJson()
	if err != nil {
		return domain.Transaction{}, []error{err}
	}

	err = s.repository.Create(transaction)
	if err != nil {
		return domain.Transaction{}, []error{err}
	}

	if err := s.kafkaBroker.Publish("process-transaction", jsonData); err != nil {
		s.logger.Error("erro ao publicar no Kafka",
			zap.Error(err),
		)
		return domain.Transaction{}, []error{err}
	}

	s.logger.Info("transação publicada com sucesso",
		zap.Any("transaction", transactionDto),
	)

	return *transaction, nil
}
