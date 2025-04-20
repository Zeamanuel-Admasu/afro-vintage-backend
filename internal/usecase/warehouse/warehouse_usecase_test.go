package warehouse_usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository is a mock implementation of the warehouse.Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) AddItem(ctx context.Context, item *warehouse.WarehouseItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockRepository) GetItemsByReseller(ctx context.Context, resellerID string) ([]*warehouse.WarehouseItem, error) {
	args := m.Called(ctx, resellerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*warehouse.WarehouseItem), args.Error(1)
}

func (m *MockRepository) GetItemsByBundle(ctx context.Context, bundleID string) ([]*warehouse.WarehouseItem, error) {
	args := m.Called(ctx, bundleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*warehouse.WarehouseItem), args.Error(1)
}

func (m *MockRepository) MarkItemAsListed(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockRepository) MarkItemAsSkipped(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockRepository) DeleteItem(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockRepository) HasResellerReceivedBundle(ctx context.Context, resellerID string, bundleID string) (bool, error) {
	args := m.Called(ctx, resellerID, bundleID)
	return args.Bool(0), args.Error(1)
}
func (m *MockRepository) CountByStatus(ctx context.Context, status string) (int, error) {
	args := m.Called(ctx, status)
	return args.Int(0), args.Error(1)
}

func TestNewWarehouseUseCase(t *testing.T) {
	// Arrange
	mockRepo := new(MockRepository)

	// Act
	useCase := NewWarehouseUseCase(mockRepo)

	// Assert
	assert.NotNil(t, useCase)
}

func TestGetWarehouseItems(t *testing.T) {
	tests := []struct {
		name        string
		resellerID  string
		mockItems   []*warehouse.WarehouseItem
		mockError   error
		expectError bool
	}{
		{
			name:       "Success - Returns items",
			resellerID: "reseller1",
			mockItems: []*warehouse.WarehouseItem{
				{
					ID:         "item1",
					ResellerID: "reseller1",
					BundleID:   "bundle1",
					ProductID:  "product1",
					Status:     "pending",
					CreatedAt:  time.Now().Format(time.RFC3339),
				},
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "Success - Returns empty list",
			resellerID:  "reseller2",
			mockItems:   []*warehouse.WarehouseItem{},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "Error - Repository error",
			resellerID:  "reseller3",
			mockItems:   nil,
			mockError:   errors.New("database error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockRepository)
			useCase := NewWarehouseUseCase(mockRepo)
			ctx := context.Background()

			mockRepo.On("GetItemsByReseller", ctx, tt.resellerID).Return(tt.mockItems, tt.mockError)

			// Act
			items, err := useCase.GetWarehouseItems(ctx, tt.resellerID)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, items)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockItems, items)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
