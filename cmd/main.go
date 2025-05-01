package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/NathanGdS/cali-challenge/handlers"
	"github.com/NathanGdS/cali-challenge/pkg/akafka"
	"github.com/NathanGdS/cali-challenge/pkg/consumers"
	"github.com/NathanGdS/cali-challenge/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Configs do kafka
	kafkaBroker := akafka.NewKafkaBroker("host.docker.internal:9094")
	defer kafkaBroker.Close()

	transactionConsumer := consumers.NewTransactionConsumer(kafkaBroker)
	go transactionConsumer.Start()

	// Configs Gin
	router := gin.Default()

	transactionHandler := handlers.NewTransactionHandler(kafkaBroker)
	router.POST("/transaction", transactionHandler.CreateTransaction)

	// Graceful shutdown config
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := router.Run(":8080"); err != nil {
			logger.Log.Fatal("erro ao iniciar o servidor",
				zap.Error(err),
			)
		}
	}()

	<-quit
	logger.Log.Info("encerrando o servidor...")
}
