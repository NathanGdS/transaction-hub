package dto

const (
	TransactionStatusProcessed = "PROCESSED"
	TransactionStatusFailed    = "FAILED"
)

type ProcessTransactionDto struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	ErrorMessage  string `json:"error_message"`
}
