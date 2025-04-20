package cartitem

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/cartitem"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- Mocks ---

type MockCartItemRepository struct {
	mock.Mock
}

func (m *MockCartItemRepository) CreateCartItem(ctx context.Context, item *cartitem.CartItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}

func (m *MockCartItemRepository) GetCartItems(ctx context.Context, userID string) ([]*cartitem.CartItem, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*cartitem.CartItem), args.Error(1)
}

func (m *MockCartItemRepository) DeleteCartItem(ctx context.Context, userID, listingID string) error {
	args := m.Called(ctx, userID, listingID)
	return args.Error(0)
}

func (m *MockCartItemRepository) ClearCart(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*product.Product), args.Error(1)
}

// Added dummy implementation so that it satisfies product.Repository.
func (m *MockProductRepository) AddProduct(ctx context.Context, prod *product.Product) error {
	return nil
}

// Added dummy DeleteProduct method as required by product.Repository.
func (m *MockProductRepository) DeleteProduct(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProductRepository) ListProductsByReseller(ctx context.Context, resellerID string, page, limit int) ([]*product.Product, error) {
	args := m.Called(ctx, resellerID, page, limit)
	return args.Get(0).([]*product.Product), args.Error(1)
}

func (m *MockProductRepository) ListAvailableProducts(ctx context.Context, page, limit int) ([]*product.Product, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).([]*product.Product), args.Error(1)
}

func (m *MockProductRepository) UpdateProduct(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockProductRepository) GetProductsByBundleID(ctx context.Context, bundleID string) ([]*product.Product, error) {
	args := m.Called(ctx, bundleID)
	return args.Get(0).([]*product.Product), args.Error(1)
}

// --- Test Suite ---

type CartItemUsecaseTestSuite struct {
	suite.Suite
	ctx             context.Context
	usecase         cartitem.Usecase
	mockCartRepo    *MockCartItemRepository
	mockProductRepo *MockProductRepository
	userID          string
}

func (suite *CartItemUsecaseTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.mockCartRepo = new(MockCartItemRepository)
	suite.mockProductRepo = new(MockProductRepository)
	suite.usecase = NewCartItemUsecase(suite.mockCartRepo, suite.mockProductRepo)
	suite.userID = "user123"
}

// --- Helper: create a dummy product ---
func createTestProduct(id string, price float64, status string, title string) *product.Product {
	// Use primitive.NewObjectID() to assign the ResellerID.
	return &product.Product{
		ID:         id,
		Title:      title,
		Price:      price,
		ImageURL:   "image.jpg",
		Grade:      "A",
		Status:     status,
		ResellerID: primitive.NewObjectID(),
	}
}

// --- Tests for AddCartItem ---

func (suite *CartItemUsecaseTestSuite) TestAddCartItem_Success() {
	testListingID := "prod123"
	prod := createTestProduct(testListingID, 100.0, "available", "Test Product")
	// Expect product lookup.
	suite.mockProductRepo.On("GetProductByID", suite.ctx, testListingID).Return(prod, nil).Once()
	// Expect repository CreateCartItem call.
	suite.mockCartRepo.On("CreateCartItem", suite.ctx, mock.MatchedBy(func(item *cartitem.CartItem) bool {
		return item.UserID == suite.userID && item.ListingID == prod.ID && item.Title == prod.Title
	})).Return(nil).Once()

	err := suite.usecase.AddCartItem(suite.ctx, suite.userID, testListingID)
	assert.NoError(suite.T(), err)
	suite.mockProductRepo.AssertExpectations(suite.T())
	suite.mockCartRepo.AssertExpectations(suite.T())
}

func (suite *CartItemUsecaseTestSuite) TestAddCartItem_ProductNotFound() {
	testListingID := "prodNotFound"
	suite.mockProductRepo.On("GetProductByID", suite.ctx, testListingID).Return(nil, nil).Once()

	err := suite.usecase.AddCartItem(suite.ctx, suite.userID, testListingID)
	assert.EqualError(suite.T(), err, "product not found")
	suite.mockProductRepo.AssertExpectations(suite.T())
}

func (suite *CartItemUsecaseTestSuite) TestAddCartItem_ProductNotAvailable() {
	testListingID := "prod123"
	prod := createTestProduct(testListingID, 100.0, "sold", "Test Product")
	suite.mockProductRepo.On("GetProductByID", suite.ctx, testListingID).Return(prod, nil).Once()

	err := suite.usecase.AddCartItem(suite.ctx, suite.userID, testListingID)
	expectedErr := fmt.Sprintf("product %s is not available", testListingID)
	assert.EqualError(suite.T(), err, expectedErr)
	suite.mockProductRepo.AssertExpectations(suite.T())
}

// --- Tests for GetCartItems ---

func (suite *CartItemUsecaseTestSuite) TestGetCartItems_Success() {
	cartItems := []*cartitem.CartItem{
		{
			ID:        "item1",
			UserID:    suite.userID,
			ListingID: "prod123",
			Title:     "Test Product",
			Price:     100.0,
			ImageURL:  "img.jpg",
			Grade:     "A",
			CreatedAt: time.Now(),
		},
	}
	suite.mockCartRepo.On("GetCartItems", suite.ctx, suite.userID).Return(cartItems, nil).Once()

	items, err := suite.usecase.GetCartItems(suite.ctx, suite.userID)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), items, 1)
	suite.mockCartRepo.AssertExpectations(suite.T())
}

// --- Tests for RemoveCartItem ---

func (suite *CartItemUsecaseTestSuite) TestRemoveCartItem_Success() {
	listingID := "prod123"
	suite.mockCartRepo.On("DeleteCartItem", suite.ctx, suite.userID, listingID).Return(nil).Once()

	err := suite.usecase.RemoveCartItem(suite.ctx, suite.userID, listingID)
	assert.NoError(suite.T(), err)
	suite.mockCartRepo.AssertExpectations(suite.T())
}

// --- Tests for CheckoutCart ---

func (suite *CartItemUsecaseTestSuite) TestCheckoutCart_Success() {
	// Set up two cart items.
	now := time.Now()
	cartItems := []*cartitem.CartItem{
		{
			ID:        "item1",
			UserID:    suite.userID,
			ListingID: "prod1",
			Title:     "Test Product 1",
			Price:     100.0,
			ImageURL:  "img1.jpg",
			Grade:     "A",
			CreatedAt: now,
		},
		{
			ID:        "item2",
			UserID:    suite.userID,
			ListingID: "prod2",
			Title:     "Test Product 2",
			Price:     200.0,
			ImageURL:  "img2.jpg",
			Grade:     "B",
			CreatedAt: now,
		},
	}
	suite.mockCartRepo.On("GetCartItems", suite.ctx, suite.userID).Return(cartItems, nil).Once()
	// For each cart item, define product lookup.
	prod1 := createTestProduct("prod1", 100.0, "available", "Test Product 1")
	prod2 := createTestProduct("prod2", 200.0, "available", "Test Product 2")
	suite.mockProductRepo.On("GetProductByID", suite.ctx, "prod1").Return(prod1, nil).Once()
	suite.mockProductRepo.On("GetProductByID", suite.ctx, "prod2").Return(prod2, nil).Once()
	// Expect ClearCart call.
	suite.mockCartRepo.On("ClearCart", suite.ctx, suite.userID).Return(nil).Once()

	resp, err := suite.usecase.CheckoutCart(suite.ctx, suite.userID)
	assert.NoError(suite.T(), err)
	// Total amount should be 300.0, fee = 6, net = 294.
	assert.Equal(suite.T(), 300.0, resp.TotalAmount)
	assert.Equal(suite.T(), 6.0, resp.PlatformFee)
	assert.Equal(suite.T(), 294.0, resp.NetPayable)
	assert.Len(suite.T(), resp.Items, 2)
	suite.mockCartRepo.AssertExpectations(suite.T())
	suite.mockProductRepo.AssertExpectations(suite.T())
}

func (suite *CartItemUsecaseTestSuite) TestCheckoutCart_EmptyCart() {
	suite.mockCartRepo.On("GetCartItems", suite.ctx, suite.userID).Return([]*cartitem.CartItem{}, nil).Once()

	resp, err := suite.usecase.CheckoutCart(suite.ctx, suite.userID)
	assert.Nil(suite.T(), resp)
	assert.EqualError(suite.T(), err, "cart is empty")
	suite.mockCartRepo.AssertExpectations(suite.T())
}

// --- Tests for CheckoutSingleItem ---

func (suite *CartItemUsecaseTestSuite) TestCheckoutSingleItem_Success() {
	// Set up a cart with one item.
	now := time.Now()
	cartItems := []*cartitem.CartItem{
		{
			ID:        "item1",
			UserID:    suite.userID,
			ListingID: "prod1",
			Title:     "Test Product 1",
			Price:     100.0,
			ImageURL:  "img1.jpg",
			Grade:     "A",
			CreatedAt: now,
		},
	}
	suite.mockCartRepo.On("GetCartItems", suite.ctx, suite.userID).Return(cartItems, nil).Once()
	prod1 := createTestProduct("prod1", 100.0, "available", "Test Product 1")
	suite.mockProductRepo.On("GetProductByID", suite.ctx, "prod1").Return(prod1, nil).Once()
	// Expect deletion of the single item from cart.
	suite.mockCartRepo.On("DeleteCartItem", suite.ctx, suite.userID, "prod1").Return(nil).Once()

	resp, err := suite.usecase.CheckoutSingleItem(suite.ctx, suite.userID, "prod1")
	assert.NoError(suite.T(), err)
	// Total should be 100, fee=2, net=98.
	assert.Equal(suite.T(), 100.0, resp.TotalAmount)
	assert.Equal(suite.T(), 2.0, resp.PlatformFee)
	assert.Equal(suite.T(), 98.0, resp.NetPayable)
	assert.Len(suite.T(), resp.Items, 1)
	suite.mockCartRepo.AssertExpectations(suite.T())
	suite.mockProductRepo.AssertExpectations(suite.T())
}

func (suite *CartItemUsecaseTestSuite) TestCheckoutSingleItem_ItemNotFoundInCart() {
	// Empty cart scenario.
	suite.mockCartRepo.On("GetCartItems", suite.ctx, suite.userID).Return([]*cartitem.CartItem{}, nil).Once()

	resp, err := suite.usecase.CheckoutSingleItem(suite.ctx, suite.userID, "prod1")
	assert.Nil(suite.T(), resp)
	assert.EqualError(suite.T(), err, "item not found in cart")
	suite.mockCartRepo.AssertExpectations(suite.T())
}

func (suite *CartItemUsecaseTestSuite) TestCheckoutSingleItem_ProductNotAvailable() {
	// Set up a cart with one item.
	now := time.Now()
	cartItems := []*cartitem.CartItem{
		{
			ID:        "item1",
			UserID:    suite.userID,
			ListingID: "prod1",
			Title:     "Test Product 1",
			Price:     100.0,
			ImageURL:  "img1.jpg",
			Grade:     "A",
			CreatedAt: now,
		},
	}
	suite.mockCartRepo.On("GetCartItems", suite.ctx, suite.userID).Return(cartItems, nil).Once()
	prod1 := createTestProduct("prod1", 100.0, "sold", "Test Product 1")
	suite.mockProductRepo.On("GetProductByID", suite.ctx, "prod1").Return(prod1, nil).Once()

	resp, err := suite.usecase.CheckoutSingleItem(suite.ctx, suite.userID, "prod1")
	assert.Nil(suite.T(), resp)
	expectedErr := fmt.Sprintf("item %q is no longer available", prod1.Title)
	assert.EqualError(suite.T(), err, expectedErr)
	suite.mockCartRepo.AssertExpectations(suite.T())
	suite.mockProductRepo.AssertExpectations(suite.T())
}

func TestCartItemUsecaseTestSuite(t *testing.T) {
	suite.Run(t, new(CartItemUsecaseTestSuite))
}
