package services

import (
	"context"
	"math"

	"github.com/NathanGdS/transaction-hub/pkg/akafka"
	"github.com/NathanGdS/transaction-hub/pkg/logger"
	"github.com/NathanGdS/transaction-hub/transaction-ledger/domain"
	"github.com/NathanGdS/transaction-hub/transaction-ledger/domain/dto"
	dRepo "github.com/NathanGdS/transaction-hub/transaction-ledger/domain/repository"
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

func (s *TransactionService) UpdateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	defer s.logger.Info("transação atualizada",
		zap.Any("transaction", transaction),
	)

	transaction.Mu.Lock()
	err := s.repository.Update(transaction)
	transaction.Mu.Unlock()

	return err
}

func (s *TransactionService) FindByID(ctx context.Context, id string) (*domain.Transaction, error) {
	return s.repository.FindByID(id)
}

func (s *TransactionService) FindPaginated(ctx context.Context, page, pageSize int) (*dto.PaginatedTransactionsResponseDto, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 50
	}

	transactions, total, err := s.repository.FindPaginated(page, pageSize)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &dto.PaginatedTransactionsResponseDto{
		Data:       transactions,
		Page:       page,
		PageSize:   pageSize,
		TotalItems: total,
		TotalPages: totalPages,
	}, nil
}
