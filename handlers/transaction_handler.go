package handlers

import (
	"net/http"

	"github.com/NathanGdS/cali-challenge/application/services"
	"github.com/NathanGdS/cali-challenge/domain/dto"
	"github.com/NathanGdS/cali-challenge/infra/logger"
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
	h.logger.Info("Solicitação recebida para criar uma transação")

	var transactionDto dto.TransactionRequestDto
	if err := c.ShouldBindJSON(&transactionDto); err != nil {
		h.logger.Error("erro ao validar JSON",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, errs := h.transactionService.CreateTransaction(c.Request.Context(), &transactionDto)
	if len(errs) > 0 {
		h.logger.Error("erro ao criar transação",
			zap.Any("errors", errs),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": errs})
		return
	}

	response := dto.FromTransaction(&transaction)

	c.JSON(http.StatusCreated, response)
}
