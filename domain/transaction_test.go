package domain_test

import (
	"testing"

	"github.com/NathanGdS/cali-challenge/domain"
	"github.com/NathanGdS/cali-challenge/domain/dto"
	"github.com/stretchr/testify/assert"
)

func TestTransaction_NewTransactionFromDto(t *testing.T) {
	// Arrange
	transaction, err := dto.ToTransaction(&dto.TransactionRequestDto{
		Amount:        100,
		PaymentMethod: "PIX",
		CurrencyCode:  "BRL",
		Description:   "Teste",
	})

	// Assert
	assert.Empty(t, err)
	assert.NotNil(t, transaction)
}

func TestTransaction_Validate_InvalidAmount(t *testing.T) {
	// Arrange
	transaction, err := dto.ToTransaction(&dto.TransactionRequestDto{
		Amount:        0,
		PaymentMethod: "PIX",
		CurrencyCode:  "BRL",
		Description:   "Teste",
	})

	// Assert
	assert.NotEmpty(t, err)
	assert.Nil(t, transaction)

	assert.Equal(t, err, []error{domain.ErrorInvalidAmount})
}

func TestTransaction_Validate_InvalidPaymentMethod(t *testing.T) {
	// Arrange
	transaction, err := dto.ToTransaction(&dto.TransactionRequestDto{
		Amount:        100,
		PaymentMethod: "INVALID",
		CurrencyCode:  "BRL",
		Description:   "Teste",
	})

	// Assert
	assert.NotEmpty(t, err)
	assert.Nil(t, transaction)

	assert.Equal(t, err, []error{domain.ErrorInvalidPaymentMethod})
}

func TestTransaction_Validate_InvalidCurrencyCode(t *testing.T) {
	// Arrange
	transaction, err := dto.ToTransaction(&dto.TransactionRequestDto{
		Amount:        100,
		PaymentMethod: "PIX",
		CurrencyCode:  "INVALID",
		Description:   "Teste",
	})

	// Assert
	assert.NotEmpty(t, err)
	assert.Nil(t, transaction)

	assert.Equal(t, err, []error{domain.ErrorInvalidCurrencyCode})
}

func TestTransaction_Validate_InvalidDescription(t *testing.T) {
	// Arrange
	transaction, err := dto.ToTransaction(&dto.TransactionRequestDto{
		Amount:        100,
		PaymentMethod: "PIX",
		CurrencyCode:  "BRL",
		Description:   "",
	})

	// Assert
	assert.NotEmpty(t, err)
	assert.Nil(t, transaction)

	assert.Equal(t, err, []error{domain.ErrorInvalidDescription})
}
