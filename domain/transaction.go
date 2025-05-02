package domain

import (
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrorInvalidPaymentMethod = errors.New("payment method must be PIX or CREDIT_CARD")
	ErrorInvalidCurrencyCode  = errors.New("currency code must be BRL or USD")
	ErrorInvalidAmount        = errors.New("amount must be greater than 0")
	ErrorInvalidDescription   = errors.New("description is required")
)

const (
	PaymentMethodPIX        = "PIX"
	PaymentMethodCreditCard = "CREDIT_CARD"
)

const (
	TransactionPending    = "PENDING"
	TransactionProcessing = "PROCESSING"
	TransactionFinished   = "FINISHED"
)

type Transaction struct {
	ID            string         `json:"id" gorm:"primaryKey;type:uuid"`
	Amount        float64        `json:"amount" gorm:"type:decimal(10,2);not null"`
	PaymentMethod string         `json:"paymentMethod" gorm:"type:varchar(20);not null"`
	CurrencyCode  string         `json:"currencyCode" gorm:"type:varchar(3);not null"`
	Description   string         `json:"description" gorm:"type:text;not null"`
	Status        string         `json:"status" gorm:"type:varchar(20);not null"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"type:timestamp;not null"`
	UpdatedAt     time.Time      `json:"updatedAt" gorm:"type:timestamp;not null"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index"`
}

func (t *Transaction) Validate() []error {
	var errors []error
	if t.Amount <= 0 {
		errors = append(errors, ErrorInvalidAmount)
	}

	if t.PaymentMethod != PaymentMethodPIX && t.PaymentMethod != PaymentMethodCreditCard {
		errors = append(errors, ErrorInvalidPaymentMethod)
	}

	if t.CurrencyCode != "BRL" && t.CurrencyCode != "USD" {
		errors = append(errors, ErrorInvalidCurrencyCode)
	}

	if t.Description == "" {
		errors = append(errors, ErrorInvalidDescription)
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}

func NewTransaction(amount float64, paymentMethod string, currencyCode string, description string) (*Transaction, []error) {
	transaction := &Transaction{
		ID:            uuid.New().String(),
		Amount:        math.Round(amount*100) / 100,
		PaymentMethod: paymentMethod,
		CurrencyCode:  currencyCode,
		Description:   description,
		Status:        TransactionPending,
	}
	err := transaction.Validate()
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (t *Transaction) ToJson() ([]byte, error) {
	return json.Marshal(t)
}
