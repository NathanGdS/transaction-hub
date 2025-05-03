package handlers

import (
	"net/http"
	"strconv"

	"github.com/NathanGdS/transaction-hub/pkg/logger"
	"github.com/NathanGdS/transaction-hub/transaction-ledger/application/services"
	"github.com/NathanGdS/transaction-hub/transaction-ledger/domain/dto"
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

func (h *TransactionHandler) GetTransactionsPaginated(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "50")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "página inválida"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tamanho de página inválido"})
		return
	}

	result, err := h.transactionService.FindPaginated(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao buscar transações"})
		return
	}

	c.JSON(http.StatusOK, result)
}
