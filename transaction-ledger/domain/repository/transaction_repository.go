package repository

import "github.com/NathanGdS/transaction-hub/transaction-ledger/domain"

type TransactionRepository interface {
	Create(transaction *domain.Transaction) error
	FindByID(id string) (*domain.Transaction, error)
	Update(transaction *domain.Transaction) error
	Delete(id string) error
	FindAll() ([]*domain.Transaction, error)
	FindPaginated(page, pageSize int) ([]domain.Transaction, int64, error)
}
