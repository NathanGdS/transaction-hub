package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NathanGdS/cali-challenge/models"
	"github.com/NathanGdS/cali-challenge/pkg/akafka"
	"github.com/NathanGdS/cali-challenge/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TransactionHandler struct {
	kafkaBroker *akafka.KafkaBroker
	logger      *zap.Logger
}

func NewTransactionHandler(broker *akafka.KafkaBroker) *TransactionHandler {
	return &TransactionHandler{
		kafkaBroker: broker,
		logger:      logger.Log,
	}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	h.logger.Info("Solicitação recebida para criar uma transação")

	var transactionDto models.TransactionRequestDto
	if err := c.ShouldBindJSON(&transactionDto); err != nil {
		h.logger.Error("erro ao validar JSON",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonData, err := json.Marshal(transactionDto)
	if err != nil {
		h.logger.Error("erro ao converter para JSON",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar a transação"})
		return
	}

	if err := h.kafkaBroker.Publish("test", jsonData); err != nil {
		h.logger.Error("erro ao publicar no Kafka",
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar a transação"})
		return
	}

	h.logger.Info("transação publicada com sucesso",
		zap.Any("transaction", transactionDto),
	)
	c.Status(http.StatusCreated)
}
