package repository

import (
	"github.com/NathanGdS/transaction-hub/transaction-ledger/domain"
	"gorm.io/gorm"
)

type TransactionRepositoryGorm struct {
	db *gorm.DB
}

func NewTransactionRepositoryGorm(db *gorm.DB) *TransactionRepositoryGorm {
	return &TransactionRepositoryGorm{
		db: db,
	}
}

func (r *TransactionRepositoryGorm) Create(transaction *domain.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *TransactionRepositoryGorm) FindByID(id string) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := r.db.First(&transaction, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepositoryGorm) Update(transaction *domain.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *TransactionRepositoryGorm) Delete(id string) error {
	return r.db.Delete(&domain.Transaction{}, "id = ?", id).Error
}

func (r *TransactionRepositoryGorm) FindAll() ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction
	err := r.db.Find(&transactions).Error
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepositoryGorm) FindPaginated(page, pageSize int) ([]domain.Transaction, int64, error) {
	var transactions []domain.Transaction
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&domain.Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(pageSize).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}
