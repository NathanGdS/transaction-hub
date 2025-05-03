package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/NathanGdS/transaction-hub/pkg/akafka"
	"github.com/NathanGdS/transaction-hub/pkg/logger"
	"github.com/NathanGdS/transaction-hub/transaction-ledger/application/consumers"
	"github.com/NathanGdS/transaction-hub/transaction-ledger/application/services"
	"github.com/NathanGdS/transaction-hub/transaction-ledger/handlers"
	"github.com/NathanGdS/transaction-hub/transaction-ledger/infra/database"
	"github.com/NathanGdS/transaction-hub/transaction-ledger/infra/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// Configs do kafka
	kafkaBroker := akafka.NewKafkaBroker("host.docker.internal:9094")
	defer kafkaBroker.Close()

	kafkaBroker.CreateTopicsIfNotExists([]string{"process-transaction", "transaction-process-return"})

	db, err := database.NewPostgresConnection()
	if err != nil {
		logger.Log.Fatal("erro ao conectar ao banco de dados",
			zap.Error(err),
		)
	}
	database.RunMigrations(db)

	txRepository := repository.NewTransactionRepositoryGorm(db)

	processTransactionConsumer := consumers.NewProcessTransactionConsumer(&kafkaBroker, txRepository)
	go processTransactionConsumer.Start()

	// Configs Gin
	router := gin.Default()

	transactionService := services.NewTransactionService(kafkaBroker, txRepository)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	router.POST("/transaction", transactionHandler.CreateTransaction)
	router.GET("/transactions", transactionHandler.GetTransactionsPaginated)

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
