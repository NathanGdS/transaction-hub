package dto

import tx "github.com/NathanGdS/transaction-hub/transaction-ledger/domain"

type TransactionRequestDto struct {
	Amount        float64 `json:"amount" validate:"required,min=0"`
	PaymentMethod string  `json:"paymentMethod" validate:"required,oneof=PIX CREDIT_CARD"`
	CurrencyCode  string  `json:"currencyCode" validate:"required,oneof=BRL USD"`
	Description   string  `json:"description" validate:"required"`
}

type TransactionResponseDto struct {
	ID string `json:"id"`
}

type PaginatedTransactionsResponseDto struct {
	Data       []tx.Transaction `json:"data"`
	Page       int              `json:"page"`
	PageSize   int              `json:"pageSize"`
	TotalItems int64            `json:"totalItems"`
	TotalPages int              `json:"totalPages"`
}

func ToTransaction(dto *TransactionRequestDto) (*tx.Transaction, []error) {
	transaction, err := tx.NewTransaction(dto.Amount, dto.PaymentMethod, dto.CurrencyCode, dto.Description)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func FromTransaction(model *tx.Transaction) TransactionResponseDto {
	return TransactionResponseDto{
		ID: model.ID,
	}
}
