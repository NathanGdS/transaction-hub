package main

import (
	"fmt"
	"log"
	"time"

	akafka "github.com/NathanGdS/cali-challenge/pkg/broker"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

func main() {
	// publish message
	go func() {
		for {
			msg := "Teste " + uuid.New().String()
			err := akafka.Publish("test", "host.docker.internal:9094", []byte(msg))
			if err != nil {
				log.Printf("Erro ao publicar mensagem: %v", err)
			}
			time.Sleep(time.Second * 10)
		}
	}()

	msgChan := make(chan *kafka.Message)
	go akafka.Consume([]string{"test"}, "host.docker.internal:9094", msgChan)

	for msg := range msgChan {
		fmt.Println(string(msg.Value))
	}
}
