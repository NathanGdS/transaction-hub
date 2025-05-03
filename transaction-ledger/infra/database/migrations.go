package database

import (
	"github.com/NathanGdS/transaction-hub/transaction-ledger/domain"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&domain.Transaction{})
}
