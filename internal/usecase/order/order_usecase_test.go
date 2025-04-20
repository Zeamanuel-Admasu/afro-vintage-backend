package OrderUsecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/payment"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Repositories
type MockBundleRepo struct {
	mock.Mock
}

func (m *MockBundleRepo) CreateBundle(ctx context.Context, b *bundle.Bundle) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockBundleRepo) GetBundleByID(ctx context.Context, id string) (*bundle.Bundle, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bundle.Bundle), args.Error(1)
}

func (m *MockBundleRepo) ListBundles(ctx context.Context, supplierID string) ([]*bundle.Bundle, error) {
	args := m.Called(ctx, supplierID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockBundleRepo) ListAvailableBundles(ctx context.Context) ([]*bundle.Bundle, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockBundleRepo) ListPurchasedByReseller(ctx context.Context, resellerID string) ([]*bundle.Bundle, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*bundle.Bundle), args.Error(1)
}

func (m *MockBundleRepo) UpdateBundleStatus(ctx context.Context, id string, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockBundleRepo) MarkAsPurchased(ctx context.Context, bundleID string, resellerID string) error {
	args := m.Called(ctx, bundleID, resellerID)
	return args.Error(0)
}

func (m *MockBundleRepo) DeleteBundle(ctx context.Context, bundleID string) error {
	args := m.Called(ctx, bundleID)
	return args.Error(0)
}

func (m *MockBundleRepo) UpdateBundle(ctx context.Context, id string, updatedData map[string]interface{}) error {
	args := m.Called(ctx, id, updatedData)
	return args.Error(0)
}

func (m *MockBundleRepo) DecreaseBundleQuantity(ctx context.Context, bundleID string) error {
	args := m.Called(ctx, bundleID)
	return args.Error(0)
}

type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) CreateOrder(ctx context.Context, o *order.Order) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *MockOrderRepo) GetOrdersByConsumer(ctx context.Context, consumerID string) ([]*order.Order, error) {
	args := m.Called(ctx, consumerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*order.Order), args.Error(1)
}

func (m *MockOrderRepo) GetOrderByID(ctx context.Context, orderID string) (*order.Order, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*order.Order), args.Error(1)
}

func (m *MockOrderRepo) UpdateOrderStatus(ctx context.Context, orderID string, status order.OrderStatus) error {
	args := m.Called(ctx, orderID, status)
	return args.Error(0)
}

func (m *MockOrderRepo) DeleteOrder(ctx context.Context, orderID string) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}

func (m *MockOrderRepo) GetOrdersBySupplier(ctx context.Context, supplierID string) ([]*order.Order, error) {
	args := m.Called(ctx, supplierID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*order.Order), args.Error(1)
}

func (m *MockOrderRepo) GetOrdersByReseller(ctx context.Context, resellerID string) ([]*order.Order, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*order.Order), args.Error(1)
}

type MockWarehouseRepo struct {
	mock.Mock
}

func (m *MockWarehouseRepo) AddItem(ctx context.Context, item *warehouse.WarehouseItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockWarehouseRepo) GetItemsByReseller(ctx context.Context, resellerID string) ([]*warehouse.WarehouseItem, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*warehouse.WarehouseItem), args.Error(1)
}

func (m *MockWarehouseRepo) GetItemsByBundle(ctx context.Context, bundleID string) ([]*warehouse.WarehouseItem, error) {
	args := m.Called(ctx, bundleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*warehouse.WarehouseItem), args.Error(1)
}

func (m *MockWarehouseRepo) MarkItemAsListed(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockWarehouseRepo) MarkItemAsSkipped(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}
func (m *MockWarehouseRepo) CountByStatus(ctx context.Context, status string) (int, error) {
	args := m.Called(ctx, status)
	return args.Int(0), args.Error(1)
}

func (m *MockWarehouseRepo) DeleteItem(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockWarehouseRepo) HasResellerReceivedBundle(ctx context.Context, resellerID string, bundleID string) (bool, error) {
	args := m.Called(ctx, resellerID, bundleID)
	return args.Bool(0), args.Error(1)
}

type MockPaymentRepo struct {
	mock.Mock
}

func (m *MockPaymentRepo) RecordPayment(ctx context.Context, p *payment.Payment) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockPaymentRepo) GetPaymentsByUser(ctx context.Context, userID string) ([]*payment.Payment, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*payment.Payment), args.Error(1)
}
func (m *MockPaymentRepo) GetAllPlatformFees(ctx context.Context) (float64, float64, error) {
	args := m.Called(ctx)
	return args.Get(0).(float64), args.Get(1).(float64), args.Error(2)
}

func (m *MockPaymentRepo) GetPaymentsByType(ctx context.Context, userID string, pType payment.PaymentType) ([]*payment.Payment, error) {
	args := m.Called(ctx, userID, pType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*payment.Payment), args.Error(1)
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepo) CreateUser(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}
func (m *MockUserRepo) CountActiveUsers(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepo) ListUsersByRole(ctx context.Context, role user.Role) ([]*user.User, error) {
	args := m.Called(ctx, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *MockUserRepo) UpdateUser(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepo) FindUserByUsername(ctx context.Context, username string) (*user.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserRepo) UpdateTrustData(ctx context.Context, user *user.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) GetBlacklistedUsers(ctx context.Context) ([]*user.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*user.User), args.Error(1)
}
func (m *MockBundleRepo) CountBundles(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

// Test Cases
func TestNewOrderUsecase(t *testing.T) {
	// Arrange
	mockBundleRepo := new(MockBundleRepo)
	mockOrderRepo := new(MockOrderRepo)
	mockWarehouseRepo := new(MockWarehouseRepo)
	mockPaymentRepo := new(MockPaymentRepo)
	mockUserRepo := new(MockUserRepo)

	// Act
	useCase := NewOrderUsecase(mockBundleRepo, mockOrderRepo, mockWarehouseRepo, mockPaymentRepo, mockUserRepo)

	// Assert
	assert.NotNil(t, useCase)
}

func TestPurchaseBundle(t *testing.T) {
	tests := []struct {
		name         string
		bundleID     string
		resellerID   string
		mockBundle   *bundle.Bundle
		mockError    error
		expectError  bool
		errorMessage string
	}{
		{
			name:       "Success - Valid purchase",
			bundleID:   "bundle1",
			resellerID: "reseller1",
			mockBundle: &bundle.Bundle{
				ID:         "bundle1",
				SupplierID: "supplier1",
				Price:      100.0,
				Status:     "available",
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:         "Error - Bundle not found",
			bundleID:     "nonexistent",
			resellerID:   "reseller1",
			mockBundle:   nil,
			mockError:    errors.New("bundle not found"),
			expectError:  true,
			errorMessage: "bundle not found",
		},
		{
			name:       "Error - Self purchase attempt",
			bundleID:   "bundle1",
			resellerID: "supplier1",
			mockBundle: &bundle.Bundle{
				ID:         "bundle1",
				SupplierID: "supplier1",
				Price:      100.0,
				Status:     "available",
			},
			mockError:    nil,
			expectError:  true,
			errorMessage: "reseller cannot purchase their own bundle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockBundleRepo := new(MockBundleRepo)
			mockOrderRepo := new(MockOrderRepo)
			mockWarehouseRepo := new(MockWarehouseRepo)
			mockPaymentRepo := new(MockPaymentRepo)
			mockUserRepo := new(MockUserRepo)
			useCase := NewOrderUsecase(mockBundleRepo, mockOrderRepo, mockWarehouseRepo, mockPaymentRepo, mockUserRepo)
			ctx := context.Background()

			mockBundleRepo.On("GetBundleByID", ctx, tt.bundleID).Return(tt.mockBundle, tt.mockError)
			if tt.mockBundle != nil {
				mockBundleRepo.On("ListAvailableBundles", ctx).Return([]*bundle.Bundle{tt.mockBundle}, nil)
			}

			if !tt.expectError {
				mockOrderRepo.On("CreateOrder", ctx, mock.AnythingOfType("*order.Order")).Return(nil)
				mockPaymentRepo.On("RecordPayment", ctx, mock.AnythingOfType("*payment.Payment")).Return(nil)
				mockBundleRepo.On("MarkAsPurchased", ctx, tt.bundleID, tt.resellerID).Return(nil)
				mockWarehouseRepo.On("AddItem", ctx, mock.AnythingOfType("*warehouse.WarehouseItem")).Return(nil)
				mockOrderRepo.On("UpdateOrderStatus", ctx, mock.AnythingOfType("string"), order.OrderStatus("completed")).Return(nil)
			}

			// Act
			order, payment, warehouseItem, err := useCase.PurchaseBundle(ctx, tt.bundleID, tt.resellerID)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMessage)
				assert.Nil(t, order)
				assert.Nil(t, payment)
				assert.Nil(t, warehouseItem)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.NotNil(t, payment)
				assert.NotNil(t, warehouseItem)
			}
			mockBundleRepo.AssertExpectations(t)
			mockOrderRepo.AssertExpectations(t)
			mockWarehouseRepo.AssertExpectations(t)
			mockPaymentRepo.AssertExpectations(t)
		})
	}
}

func TestGetDashboardMetrics(t *testing.T) {
	tests := []struct {
		name           string
		supplierID     string
		mockBundles    []*bundle.Bundle
		mockUser       *user.User
		mockError      error
		expectError    bool
		expectedSales  float64
		expectedCounts order.PerformanceMetrics
		expectedRating int
		expectedBest   float64
	}{
		{
			name:       "Success - With bundles",
			supplierID: "supplier1",
			mockBundles: []*bundle.Bundle{
				{
					ID:         "bundle1",
					Status:     "purchased",
					Price:      100.0,
					DateListed: time.Now(),
				},
				{
					ID:         "bundle2",
					Status:     "available",
					Price:      200.0,
					DateListed: time.Now(),
				},
			},
			mockUser: &user.User{
				ID:         "supplier1",
				TrustScore: 85,
			},
			mockError:     nil,
			expectError:   false,
			expectedSales: 100.0,
			expectedCounts: order.PerformanceMetrics{
				TotalBundlesListed: 2,
				ActiveCount:        1,
				SoldCount:          1,
			},
			expectedRating: 85,
			expectedBest:   100.0,
		},
		{
			name:        "Success - No bundles",
			supplierID:  "supplier2",
			mockBundles: []*bundle.Bundle{},
			mockUser: &user.User{
				ID:         "supplier2",
				TrustScore: 90,
			},
			mockError:     nil,
			expectError:   false,
			expectedSales: 0.0,
			expectedCounts: order.PerformanceMetrics{
				TotalBundlesListed: 0,
				ActiveCount:        0,
				SoldCount:          0,
			},
			expectedRating: 90,
			expectedBest:   0.0,
		},
		{
			name:           "Error - Repository error",
			supplierID:     "supplier3",
			mockBundles:    nil,
			mockUser:       nil,
			mockError:      errors.New("database error"),
			expectError:    true,
			expectedSales:  0.0,
			expectedCounts: order.PerformanceMetrics{},
			expectedRating: 0,
			expectedBest:   0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockBundleRepo := new(MockBundleRepo)
			mockOrderRepo := new(MockOrderRepo)
			mockWarehouseRepo := new(MockWarehouseRepo)
			mockPaymentRepo := new(MockPaymentRepo)
			mockUserRepo := new(MockUserRepo)
			useCase := NewOrderUsecase(mockBundleRepo, mockOrderRepo, mockWarehouseRepo, mockPaymentRepo, mockUserRepo)
			ctx := context.Background()

			mockBundleRepo.On("ListBundles", ctx, tt.supplierID).Return(tt.mockBundles, tt.mockError)
			if tt.mockUser != nil {
				mockUserRepo.On("GetByID", ctx, tt.supplierID).Return(tt.mockUser, nil)
			}

			// Act
			metrics, err := useCase.GetDashboardMetrics(ctx, tt.supplierID)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, metrics)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, metrics)
				assert.Equal(t, tt.expectedSales, metrics.TotalSales)
				assert.Equal(t, tt.expectedCounts, metrics.PerformanceMetrics)
				assert.Equal(t, tt.expectedRating, metrics.Rating)
				assert.Equal(t, tt.expectedBest, metrics.BestSelling)
			}
			mockBundleRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}

func TestGetOrderByID(t *testing.T) {
	tests := []struct {
		name        string
		orderID     string
		mockOrder   *order.Order
		mockError   error
		expectError bool
	}{
		{
			name:    "Success - Order found",
			orderID: "order1",
			mockOrder: &order.Order{
				ID:         "order1",
				ResellerID: "reseller1",
				Status:     order.OrderStatusCompleted,
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "Error - Order not found",
			orderID:     "nonexistent",
			mockOrder:   nil,
			mockError:   errors.New("order not found"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockBundleRepo := new(MockBundleRepo)
			mockOrderRepo := new(MockOrderRepo)
			mockWarehouseRepo := new(MockWarehouseRepo)
			mockPaymentRepo := new(MockPaymentRepo)
			mockUserRepo := new(MockUserRepo)
			useCase := NewOrderUsecase(mockBundleRepo, mockOrderRepo, mockWarehouseRepo, mockPaymentRepo, mockUserRepo)
			ctx := context.Background()

			mockOrderRepo.On("GetOrderByID", ctx, tt.orderID).Return(tt.mockOrder, tt.mockError)

			// Act
			order, err := useCase.GetOrderByID(ctx, tt.orderID)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, order)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.Equal(t, tt.mockOrder.ID, order.ID)
			}
			mockOrderRepo.AssertExpectations(t)
		})
	}
}

func TestGetSoldBundleHistory(t *testing.T) {
	tests := []struct {
		name          string
		supplierID    string
		mockOrders    []*order.Order
		mockError     error
		expectError   bool
		expectedCount int
	}{
		{
			name:       "Success - With sold bundles",
			supplierID: "supplier1",
			mockOrders: []*order.Order{
				{
					ID:         "order1",
					BundleID:   "bundle1",
					ProductIDs: []string{},
				},
				{
					ID:         "order2",
					BundleID:   "bundle2",
					ProductIDs: []string{},
				},
			},
			mockError:     nil,
			expectError:   false,
			expectedCount: 2,
		},
		{
			name:          "Success - No sold bundles",
			supplierID:    "supplier2",
			mockOrders:    []*order.Order{},
			mockError:     nil,
			expectError:   false,
			expectedCount: 0,
		},
		{
			name:          "Error - Repository error",
			supplierID:    "supplier3",
			mockOrders:    nil,
			mockError:     errors.New("database error"),
			expectError:   true,
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockBundleRepo := new(MockBundleRepo)
			mockOrderRepo := new(MockOrderRepo)
			mockWarehouseRepo := new(MockWarehouseRepo)
			mockPaymentRepo := new(MockPaymentRepo)
			mockUserRepo := new(MockUserRepo)
			useCase := NewOrderUsecase(mockBundleRepo, mockOrderRepo, mockWarehouseRepo, mockPaymentRepo, mockUserRepo)
			ctx := context.Background()

			mockOrderRepo.On("GetOrdersBySupplier", ctx, tt.supplierID).Return(tt.mockOrders, tt.mockError)

			// Act
			orders, err := useCase.GetSoldBundleHistory(ctx, tt.supplierID)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, orders)
			} else {
				assert.NoError(t, err)
				if tt.expectedCount == 0 {
					assert.Empty(t, orders)
				} else {
					assert.NotNil(t, orders)
					assert.Len(t, orders, tt.expectedCount)
				}
			}
			mockOrderRepo.AssertExpectations(t)
		})
	}
}
