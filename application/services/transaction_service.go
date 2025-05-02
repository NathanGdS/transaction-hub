package services

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand/v2"
	"time"

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

func (s *TransactionService) ProcessTransaction(ctx context.Context, transaction *domain.Transaction) error {
	s.logger.Info("processando transação",
		zap.Any("transaction", transaction),
	)

	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()

	errChan := make(chan error, 1)
	done := make(chan struct{})

	go func() {
		defer close(done)

		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
			// Simulação de processamento
			time.Sleep(time.Duration(rand.IntN(5)) * time.Second)

			if ctx.Err() != nil {
				errChan <- ctx.Err()
				return
			}

			go s.publishResult(transaction, dto.TransactionStatusProcessed, "")

			errChan <- nil
		}
	}()

	select {
	case err := <-errChan:
		<-done // Aguarda a goroutine finalizar
		if err == nil {
			s.logger.Info("transação processada com sucesso",
				zap.Any("transaction", transaction),
			)
			return nil
		}

		errorMsg := err.Error()
		if err == context.DeadlineExceeded {
			errorMsg = "timeout ao processar transação"
		}

		s.logger.Error(errorMsg,
			zap.Any("transaction", transaction),
		)

		if publishErr := s.publishResult(transaction, dto.TransactionStatusFailed, errorMsg); publishErr != nil {
			return errors.Join(err, publishErr)
		}

		return err

	case <-ctx.Done():
		<-done // Aguarda a goroutine finalizar
		errorMsg := "timeout ao processar transação"
		s.logger.Error(errorMsg,
			zap.Any("transaction", transaction),
		)

		if publishErr := s.publishResult(transaction, dto.TransactionStatusFailed, errorMsg); publishErr != nil {
			return errors.Join(ctx.Err(), publishErr)
		}

		return ctx.Err()
	}
}

func (s *TransactionService) publishResult(transaction *domain.Transaction, status string, errorMsg string) error {
	processTransactionDto := dto.ProcessTransactionDto{
		TransactionID: transaction.ID,
		Status:        status,
		ErrorMessage:  errorMsg,
	}

	jsonData, err := json.Marshal(processTransactionDto)
	if err != nil {
		s.logger.Error("erro ao serializar DTO",
			zap.Error(err),
		)
		return err
	}

	if err := s.kafkaBroker.Publish("transaction-process-return", jsonData); err != nil {
		s.logger.Error("erro ao publicar mensagem no Kafka",
			zap.Error(err),
		)
		return err
	}

	return nil
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
