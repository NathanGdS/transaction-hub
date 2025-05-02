package services

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand/v2"
	"time"

	"github.com/NathanGdS/cali-challenge/pkg/akafka"
	"github.com/NathanGdS/cali-challenge/pkg/logger"
	"github.com/NathanGdS/cali-challenge/transaction-ledger/domain"
	"github.com/NathanGdS/cali-challenge/transaction-ledger/domain/dto"
	"go.uber.org/zap"
)

type ProcessTransactionService struct {
	kafkaBroker akafka.KafkaBroker
	logger      *zap.Logger
}

func NewProcessTransactionService(kafkaBroker akafka.KafkaBroker) *ProcessTransactionService {
	return &ProcessTransactionService{kafkaBroker: kafkaBroker, logger: logger.Log}
}

func (s *ProcessTransactionService) ProcessTransaction(ctx context.Context, transaction *domain.Transaction) error {
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

func (s *ProcessTransactionService) publishResult(transaction *domain.Transaction, status string, errorMsg string) error {
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
