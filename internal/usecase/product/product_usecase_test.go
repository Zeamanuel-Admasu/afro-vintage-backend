package productusecase

import (
	"context"
	"errors"
	"testing"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) AddProduct(ctx context.Context, p *product.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockRepository) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

func (m *MockRepository) ListProductsByReseller(ctx context.Context, resellerID string, page, limit int) ([]*product.Product, error) {
	args := m.Called(ctx, resellerID, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*product.Product), args.Error(1)
}

func (m *MockRepository) ListAvailableProducts(ctx context.Context, page, limit int) ([]*product.Product, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*product.Product), args.Error(1)
}

func (m *MockRepository) DeleteProduct(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) UpdateProduct(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

type ProductUsecaseTestSuite struct {
	suite.Suite
	mockRepo *MockRepository
	usecase  product.Usecase
}

func (suite *ProductUsecaseTestSuite) SetupTest() {
	suite.mockRepo = new(MockRepository)
	suite.usecase = NewProductUsecase(suite.mockRepo)
}

func (suite *ProductUsecaseTestSuite) TestAddProduct() {
	product := &product.Product{
		ID:    "test-id",
		Title: "Test Product",
	}

	suite.mockRepo.On("AddProduct", mock.Anything, product).Return(nil)

	err := suite.usecase.AddProduct(context.Background(), product)
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestGetProductByID() {
	expectedProduct := &product.Product{
		ID:    "test-id",
		Title: "Test Product",
	}

	suite.mockRepo.On("GetProductByID", mock.Anything, "test-id").Return(expectedProduct, nil)

	result, err := suite.usecase.GetProductByID(context.Background(), "test-id")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedProduct, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestListProductsByReseller() {
	expectedProducts := []*product.Product{
		{ID: "test-id-1", Title: "Product 1"},
		{ID: "test-id-2", Title: "Product 2"},
	}

	suite.mockRepo.On("ListProductsByReseller", mock.Anything, "reseller-id", 1, 10).Return(expectedProducts, nil)

	result, err := suite.usecase.ListProductsByReseller(context.Background(), "reseller-id", 1, 10)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedProducts, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestListProductsByReseller_EmptyState() {
	suite.mockRepo.On("ListProductsByReseller", mock.Anything, "reseller-id", 1, 10).Return([]*product.Product{}, nil)

	result, err := suite.usecase.ListProductsByReseller(context.Background(), "reseller-id", 1, 10)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestListProductsByReseller_LoadFailure() {
	suite.mockRepo.On("ListProductsByReseller", mock.Anything, "reseller-id", 1, 10).Return(nil, errors.New("database error"))

	result, err := suite.usecase.ListProductsByReseller(context.Background(), "reseller-id", 1, 10)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "database error", err.Error())
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestListAvailableProducts() {
	expectedProducts := []*product.Product{
		{ID: "test-id-1", Title: "Product 1"},
		{ID: "test-id-2", Title: "Product 2"},
	}

	suite.mockRepo.On("ListAvailableProducts", mock.Anything, 1, 10).Return(expectedProducts, nil)

	result, err := suite.usecase.ListAvailableProducts(context.Background(), 1, 10)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedProducts, result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestListAvailableProducts_EmptyState() {
	suite.mockRepo.On("ListAvailableProducts", mock.Anything, 1, 10).Return([]*product.Product{}, nil)

	result, err := suite.usecase.ListAvailableProducts(context.Background(), 1, 10)
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestListAvailableProducts_LoadFailure() {
	suite.mockRepo.On("ListAvailableProducts", mock.Anything, 1, 10).Return(nil, errors.New("database error"))

	result, err := suite.usecase.ListAvailableProducts(context.Background(), 1, 10)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), result)
	assert.Equal(suite.T(), "database error", err.Error())
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestDeleteProduct() {
	suite.mockRepo.On("DeleteProduct", mock.Anything, "test-id").Return(nil)

	err := suite.usecase.DeleteProduct(context.Background(), "test-id")
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ProductUsecaseTestSuite) TestUpdateProduct() {
	updates := map[string]interface{}{
		"Title": "Updated Title",
	}

	suite.mockRepo.On("UpdateProduct", mock.Anything, "test-id", updates).Return(nil)

	err := suite.usecase.UpdateProduct(context.Background(), "test-id", updates)
	assert.NoError(suite.T(), err)
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestProductUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(ProductUsecaseTestSuite))
}
