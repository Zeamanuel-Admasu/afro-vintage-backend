package controllers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateOrder(ctx context.Context, o *order.Order) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrdersByConsumer(ctx context.Context, consumerID string) ([]*order.Order, error) {
	args := m.Called(ctx, consumerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*order.Order), args.Error(1)
}

func (m *MockOrderRepository) GetOrderByID(ctx context.Context, orderID string) (*order.Order, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.Order), args.Error(1)
}

func (m *MockOrderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status order.OrderStatus) error {
	args := m.Called(ctx, orderID, status)
	return args.Error(0)
}

func (m *MockOrderRepository) DeleteOrder(ctx context.Context, orderID string) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}

type ConsumerControllerTestSuite struct {
	suite.Suite
	mockRepo    *MockOrderRepository
	controller  *ConsumerController
	testContext *gin.Context
	recorder    *httptest.ResponseRecorder
}

func (suite *ConsumerControllerTestSuite) SetupTest() {
	suite.mockRepo = new(MockOrderRepository)
	suite.controller = NewConsumerController(suite.mockRepo)

	// Set up a Gin context and response recorder for testing
	gin.SetMode(gin.TestMode)
	suite.recorder = httptest.NewRecorder()
	suite.testContext, _ = gin.CreateTestContext(suite.recorder)
}

func (suite *ConsumerControllerTestSuite) TestGetOrderHistory_Success() {
	orders := []*order.Order{
		{
			ID:         "order1",
			ConsumerID: "consumer1",
			ProductIDs: []string{"product1"},
			TotalPrice: 100.0,
			Status:     order.Pending,
			CreatedAt:  time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
		},
	}

	suite.mockRepo.On("GetOrdersByConsumer", mock.Anything, "consumer1").Return(orders, nil)

	suite.testContext.Set("userID", "consumer1")
	suite.controller.GetOrderHistory(suite.testContext)

	assert.Equal(suite.T(), http.StatusOK, suite.recorder.Code)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ConsumerControllerTestSuite) TestGetOrderHistory_NoOrders() {
	suite.mockRepo.On("GetOrdersByConsumer", mock.Anything, "consumer1").Return([]*order.Order{}, nil)

	suite.testContext.Set("userID", "consumer1")
	suite.controller.GetOrderHistory(suite.testContext)

	assert.Equal(suite.T(), http.StatusOK, suite.recorder.Code)
	assert.Contains(suite.T(), suite.recorder.Body.String(), "No orders yet.")
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ConsumerControllerTestSuite) TestGetOrderHistory_ErrorFetchingOrders() {
	suite.mockRepo.On("GetOrdersByConsumer", mock.Anything, "consumer1").Return(nil, errors.New("database error"))

	suite.testContext.Set("userID", "consumer1")
	suite.controller.GetOrderHistory(suite.testContext)

	assert.Equal(suite.T(), http.StatusInternalServerError, suite.recorder.Code)
	assert.Contains(suite.T(), suite.recorder.Body.String(), "Failed to fetch orders")
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestConsumerControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ConsumerControllerTestSuite))
}
