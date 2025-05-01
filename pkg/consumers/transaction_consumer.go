package consumers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/NathanGdS/cali-challenge/models"
	akafka "github.com/NathanGdS/cali-challenge/pkg/akafka"
	"github.com/segmentio/kafka-go"
)

type TransactionConsumer struct {
	kafkaBroker *akafka.KafkaBroker
}

func NewTransactionConsumer(broker *akafka.KafkaBroker) *TransactionConsumer {
	return &TransactionConsumer{
		kafkaBroker: broker,
	}
}

func (c *TransactionConsumer) Start() {
	msgChan := make(chan *kafka.Message)
	go c.kafkaBroker.Consume([]string{"test"}, msgChan)

	for msg := range msgChan {
		c.processMessage(msg)
	}
}

func (c *TransactionConsumer) processMessage(msg *kafka.Message) {
	fmt.Printf("Consumindo mensagem: %s\n", string(msg.Value))

	var transactionDto models.TransactionRequestDto
	err := json.Unmarshal(msg.Value, &transactionDto)
	if err != nil {
		log.Printf("Erro ao converter para JSON: %v", err)
		return
	}

	transaction, errs := models.TransactionFromDto(&transactionDto)
	if len(errs) > 0 {
		log.Printf("Erro ao converter DTO para Transaction: %v", errs)
		return
	}

	fmt.Printf("Transação recebida: %+v\n", transaction)
}
