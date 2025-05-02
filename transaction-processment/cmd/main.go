package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/NathanGdS/cali-challenge/pkg/akafka"
	"github.com/NathanGdS/cali-challenge/pkg/logger"
	"github.com/NathanGdS/cali-challenge/transaction-processment/application/consumers"
)

func main() {
	// Configs do kafka
	kafkaBroker := akafka.NewKafkaBroker("host.docker.internal:9094")
	defer kafkaBroker.Close()

	kafkaBroker.CreateTopicsIfNotExists([]string{"process-transaction", "transaction-process-return"})

	transactionConsumer := consumers.NewProcessTransactionConsumer(&kafkaBroker)
	// usando como loop infinito da aplicação
	transactionConsumer.Start()

	// Graceful shutdown config
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Log.Info("encerrando o servidor...")
}
