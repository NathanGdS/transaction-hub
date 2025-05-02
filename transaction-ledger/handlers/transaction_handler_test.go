package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NathanGdS/cali-challenge/pkg/akafka"
	"github.com/NathanGdS/cali-challenge/pkg/logger"
	"github.com/NathanGdS/cali-challenge/transaction-ledger/application/services"
	"github.com/NathanGdS/cali-challenge/transaction-ledger/domain"
	"github.com/NathanGdS/cali-challenge/transaction-ledger/domain/dto"
	dRepo "github.com/NathanGdS/cali-challenge/transaction-ledger/domain/repository"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockKafkaBroker struct {
	mock.Mock
}

func (m *MockKafkaBroker) Publish(topic string, message []byte) error {
	args := m.Called(topic, message)
	return args.Error(0)
}

func (m *MockKafkaBroker) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockKafkaBroker) Consume(topics []string, msgChan chan *kafka.Message) {
	m.Called(topics, msgChan)
}

func (m *MockKafkaBroker) CreateTopicsIfNotExists(topics []string) error {
	args := m.Called(topics)
	return args.Error(0)
}

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(transaction *domain.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) FindByID(id string) (*domain.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(transaction *domain.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTransactionRepository) FindAll() ([]*domain.Transaction, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

var _ akafka.KafkaBroker = (*MockKafkaBroker)(nil)
var _ dRepo.TransactionRepository = (*MockTransactionRepository)(nil)

// Mock logger
var mockLogger *zap.Logger

func init() {
	mockLogger = zap.NewNop()
}

func TestNewTransactionHandler(t *testing.T) {
	// Arrange
	mockKafka := new(MockKafkaBroker)
	mockRepo := new(MockTransactionRepository)
	service := services.NewTransactionService(mockKafka, mockRepo)

	// Act
	handler := NewTransactionHandler(service)

	// Assert
	assert.NotNil(t, handler)
	assert.Equal(t, service, handler.transactionService)
	assert.NotNil(t, handler.logger)
}

func TestCreateTransaction(t *testing.T) {
	// Configurando o Gin para modo de teste
	gin.SetMode(gin.TestMode)
	logger.Log = mockLogger

	t.Run("Should create a transaction with success", func(t *testing.T) {
		// Arrange
		mockKafka := new(MockKafkaBroker)
		mockRepo := new(MockTransactionRepository)
		service := services.NewTransactionService(mockKafka, mockRepo)
		handler := NewTransactionHandler(service)

		requestDto := dto.TransactionRequestDto{
			Amount:        100.0,
			PaymentMethod: domain.PaymentMethodPIX,
			CurrencyCode:  "BRL",
			Description:   "Test transaction",
		}

		mockRepo.On("Create", mock.AnythingOfType("*domain.Transaction")).Return(nil)
		mockKafka.On("Publish", "process-transaction", mock.Anything).Return(nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(requestDto)
		c.Request = httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		// Act
		handler.CreateTransaction(c)

		// Assert
		assert.Equal(t, http.StatusCreated, w.Code)

		var response dto.TransactionResponseDto
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.ID)
		mockKafka.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Should return error when the JSON is invalid", func(t *testing.T) {
		// Arrange
		mockKafka := new(MockKafkaBroker)
		mockRepo := new(MockTransactionRepository)
		service := services.NewTransactionService(mockKafka, mockRepo)
		handler := NewTransactionHandler(service)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		invalidJSON := []byte(`{"amount": "invalid"}`)
		c.Request = httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(invalidJSON))
		c.Request.Header.Set("Content-Type", "application/json")

		// Act
		handler.CreateTransaction(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Should return error when the Kafka fails", func(t *testing.T) {
		// Arrange
		mockKafka := new(MockKafkaBroker)
		mockRepo := new(MockTransactionRepository)
		service := services.NewTransactionService(mockKafka, mockRepo)
		handler := NewTransactionHandler(service)

		requestDto := dto.TransactionRequestDto{
			Amount:        100.0,
			PaymentMethod: domain.PaymentMethodPIX,
			CurrencyCode:  "BRL",
			Description:   "Test transaction",
		}

		mockRepo.On("Create", mock.AnythingOfType("*domain.Transaction")).Return(nil)
		mockKafka.On("Publish", "process-transaction", mock.Anything).Return(errors.New("erro ao publicar no kafka"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		jsonData, _ := json.Marshal(requestDto)
		c.Request = httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(jsonData))
		c.Request.Header.Set("Content-Type", "application/json")

		// Act
		handler.CreateTransaction(c)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockKafka.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}
