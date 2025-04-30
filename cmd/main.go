package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/NathanGdS/cali-challenge/models"
	akafka "github.com/NathanGdS/cali-challenge/pkg/broker"
	"github.com/segmentio/kafka-go"
)

func main() {
	// publish message
	go func() {
		for {
			randomNumber := rand.Intn(2)
			paymentMethod := "PIX"
			currencyCode := "BRL"
			if randomNumber == 1 {
				paymentMethod = "CREDIT_CARD"
				currencyCode = "USD"
			}

			transactionDto := &models.TransactionRequestDto{
				Amount:        rand.Float64() * 100,
				PaymentMethod: paymentMethod,
				CurrencyCode:  currencyCode,
				Description:   "Teste" + time.Now().Format("2006-01-02 15:04:05"),
			}

			jsonData, err := json.Marshal(transactionDto)
			if err != nil {
				log.Printf("Erro ao converter para JSON: %v", err)
				continue
			}

			fmt.Printf("Publicando mensagem: %s\n", string(jsonData))

			publishErr := akafka.Publish("test", "host.docker.internal:9094", jsonData)
			if publishErr != nil {
				log.Printf("Erro ao publicar mensagem: %v", publishErr)
			}
			time.Sleep(time.Second * 3)
		}
	}()

	msgChan := make(chan *kafka.Message)
	go akafka.Consume([]string{"test"}, "host.docker.internal:9094", msgChan)

	for msg := range msgChan {
		fmt.Printf("Consumindo mensagem: %s\n", string(msg.Value))

		var transactionDto models.TransactionRequestDto
		err := json.Unmarshal(msg.Value, &transactionDto)
		if err != nil {
			log.Printf("Erro ao converter para JSON: %v", err)
			continue
		}

		transaction, errs := models.TransactionFromDto(&transactionDto)
		if len(errs) > 0 {
			log.Printf("Erro ao converter DTO para Transaction: %v", errs)
			continue
		}

		fmt.Printf("Transação recebida: %+v\n", transaction)
	}
}
