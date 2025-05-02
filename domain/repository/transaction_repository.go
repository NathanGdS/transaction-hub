package repository

import "github.com/NathanGdS/cali-challenge/domain"

type TransactionRepository interface {
	Create(transaction *domain.Transaction) error
	FindByID(id string) (*domain.Transaction, error)
	Update(transaction *domain.Transaction) error
	Delete(id string) error
	FindAll() ([]*domain.Transaction, error)
}
