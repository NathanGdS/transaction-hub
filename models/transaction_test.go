package models_test

import (
	"testing"

	"github.com/NathanGdS/cali-challenge/dto"
	"github.com/NathanGdS/cali-challenge/models"
	"github.com/stretchr/testify/assert"
)

func TestTransaction_NewTransactionFromDto(t *testing.T) {
	transaction, err := dto.ToTransaction(&dto.TransactionRequestDto{
		Amount:        100,
		PaymentMethod: "PIX",
		CurrencyCode:  "BRL",
		Description:   "Teste",
	})

	assert.Empty(t, err)
	assert.NotNil(t, transaction)
}

func TestTransaction_Validate_InvalidAmount(t *testing.T) {
	transaction, err := dto.ToTransaction(&dto.TransactionRequestDto{
		Amount:        0,
		PaymentMethod: "PIX",
		CurrencyCode:  "BRL",
		Description:   "Teste",
	})
	assert.NotEmpty(t, err)
	assert.Nil(t, transaction)

	assert.Equal(t, err, []error{models.ErrorInvalidAmount})
}

func TestTransaction_Validate_InvalidPaymentMethod(t *testing.T) {
	transaction, err := dto.ToTransaction(&dto.TransactionRequestDto{
		Amount:        100,
		PaymentMethod: "INVALID",
		CurrencyCode:  "BRL",
		Description:   "Teste",
	})

	assert.NotEmpty(t, err)
	assert.Nil(t, transaction)

	assert.Equal(t, err, []error{models.ErrorInvalidPaymentMethod})
}

func TestTransaction_Validate_InvalidCurrencyCode(t *testing.T) {
	transaction, err := dto.ToTransaction(&dto.TransactionRequestDto{
		Amount:        100,
		PaymentMethod: "PIX",
		CurrencyCode:  "INVALID",
		Description:   "Teste",
	})

	assert.NotEmpty(t, err)
	assert.Nil(t, transaction)

	assert.Equal(t, err, []error{models.ErrorInvalidCurrencyCode})
}

func TestTransaction_Validate_InvalidDescription(t *testing.T) {
	transaction, err := dto.ToTransaction(&dto.TransactionRequestDto{
		Amount:        100,
		PaymentMethod: "PIX",
		CurrencyCode:  "BRL",
		Description:   "",
	})

	assert.NotEmpty(t, err)
	assert.Nil(t, transaction)

	assert.Equal(t, err, []error{models.ErrorInvalidDescription})
}
