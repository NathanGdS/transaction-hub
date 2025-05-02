package handlers

import (
	"net/http"

	"github.com/NathanGdS/cali-challenge/pkg/logger"
	"github.com/NathanGdS/cali-challenge/transaction-ledger/application/services"
	"github.com/NathanGdS/cali-challenge/transaction-ledger/domain/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
	logger             *zap.Logger
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		logger:             logger.Log,
	}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var transactionDto dto.TransactionRequestDto
	if err := c.ShouldBindJSON(&transactionDto); err != nil {
		h.logger.Error("erro ao validar JSON",
			zap.Error(err),
		)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": []string{err.Error()}})
		return
	}

	transaction, errs := h.transactionService.CreateTransaction(c.Request.Context(), &transactionDto)
	if len(errs) > 0 {
		h.logger.Error("erro ao criar transação",
			zap.Any("errors", errs),
		)

		// Pré-aloca o slice de erros
		errorMessages := make([]string, 0, len(errs))
		for _, err := range errs {
			errorMessages = append(errorMessages, err.Error())
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errorMessages})
		return
	}

	c.JSON(http.StatusCreated, dto.FromTransaction(&transaction))
}
