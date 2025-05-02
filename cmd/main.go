package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/NathanGdS/cali-challenge/application/consumers"
	"github.com/NathanGdS/cali-challenge/application/services"
	"github.com/NathanGdS/cali-challenge/handlers"
	"github.com/NathanGdS/cali-challenge/infra/akafka"
	"github.com/NathanGdS/cali-challenge/infra/database"
	"github.com/NathanGdS/cali-challenge/infra/logger"
	"github.com/NathanGdS/cali-challenge/infra/repository"
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

	transactionConsumer := consumers.NewTransactionConsumer(&kafkaBroker, txRepository)
	go transactionConsumer.Start()

	processTransactionConsumer := consumers.NewProcessTransactionConsumer(&kafkaBroker, txRepository)
	go processTransactionConsumer.Start()

	// Configs Gin
	router := gin.Default()

	transactionService := services.NewTransactionService(kafkaBroker, txRepository)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
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
