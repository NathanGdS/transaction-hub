package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/NathanGdS/cali-challenge/models"
	akafka "github.com/NathanGdS/cali-challenge/pkg/akafka"
	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	kafkaBroker *akafka.KafkaBroker
}

func NewTransactionHandler(broker *akafka.KafkaBroker) *TransactionHandler {
	return &TransactionHandler{
		kafkaBroker: broker,
	}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var transactionDto models.TransactionRequestDto
	if err := c.ShouldBindJSON(&transactionDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonData, err := json.Marshal(transactionDto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar a transação"})
		return
	}

	if err := h.kafkaBroker.Publish("test", jsonData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao processar a transação"})
		return
	}

	c.Status(http.StatusCreated)
}
