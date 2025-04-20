package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/cartitem"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockCartItemUsecase struct {
	mock.Mock
}

func (m *MockCartItemUsecase) AddCartItem(ctx context.Context, userID, listingID string) error {
	args := m.Called(ctx, userID, listingID)
	return args.Error(0)
}

func (m *MockCartItemUsecase) GetCartItems(ctx context.Context, userID string) ([]*cartitem.CartItem, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*cartitem.CartItem), args.Error(1)
}

func (m *MockCartItemUsecase) RemoveCartItem(ctx context.Context, userID, listingID string) error {
	args := m.Called(ctx, userID, listingID)
	return args.Error(0)
}

// Change signature to return *models.CheckoutResponse instead of interface{}
func (m *MockCartItemUsecase) CheckoutCart(ctx context.Context, userID string) (*models.CheckoutResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CheckoutResponse), args.Error(1)
}

// Change signature to return *models.CheckoutResponse instead of interface{}
func (m *MockCartItemUsecase) CheckoutSingleItem(ctx context.Context, userID, listingID string) (*models.CheckoutResponse, error) {
	args := m.Called(ctx, userID, listingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CheckoutResponse), args.Error(1)
}

type CartItemControllerTestSuite struct {
	suite.Suite
	controller *CartItemController
	mockUC     *MockCartItemUsecase
	router     *gin.Engine
	userID     string
}

func (suite *CartItemControllerTestSuite) SetupTest() {
	suite.mockUC = new(MockCartItemUsecase)
	suite.controller = NewCartItemController(suite.mockUC)
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
	suite.userID = "user123"

	// Add a simple auth middleware that sets userID in the context.
	suite.router.Use(func(c *gin.Context) {
		c.Set("userID", suite.userID)
	})
}

func (suite *CartItemControllerTestSuite) TestAddCartItem_Success() {
	// Setup
	req := models.CreateCartItemRequest{
		ListingID: "listing123",
	}
	suite.mockUC.On("AddCartItem", mock.Anything, suite.userID, req.ListingID).Return(nil)

	// Execute
	jsonData, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/cart", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	suite.router.POST("/api/cart", suite.controller.AddCartItem)
	suite.router.ServeHTTP(w, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "item added to cart", response["message"])
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *CartItemControllerTestSuite) TestAddCartItem_Unauthorized() {
	// Setup
	req := models.CreateCartItemRequest{
		ListingID: "listing123",
	}

	// Execute
	jsonData, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/cart", bytes.NewBuffer(jsonData))
	httpReq.Header.Set("Content-Type", "application/json")
	// Create new router without auth middleware.
	router := gin.Default()
	router.POST("/api/cart", suite.controller.AddCartItem)
	router.ServeHTTP(w, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *CartItemControllerTestSuite) TestAddCartItem_InvalidRequest() {
	// Setup
	invalidReq := "invalid json"

	// Execute
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/cart", bytes.NewBuffer([]byte(invalidReq)))
	httpReq.Header.Set("Content-Type", "application/json")
	suite.router.POST("/api/cart", suite.controller.AddCartItem)
	suite.router.ServeHTTP(w, httpReq)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *CartItemControllerTestSuite) TestGetCartItems_Success() {
	// Setup
	now := time.Now()
	items := []*cartitem.CartItem{
		{
			ID:        "item1",
			ListingID: "listing1",
			Title:     "Test Item 1",
			Price:     100.0,
			ImageURL:  "image1.jpg",
			Grade:     "A",
			CreatedAt: now,
		},
	}
	suite.mockUC.On("GetCartItems", mock.Anything, suite.userID).Return(items, nil)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/cart", nil)
	suite.router.GET("/api/cart", suite.controller.GetCartItems)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var response []models.CartItemResponse
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Len(suite.T(), response, 1)
	assert.Equal(suite.T(), items[0].ID, response[0].ID)
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *CartItemControllerTestSuite) TestGetCartItems_Unauthorized() {
	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/cart", nil)
	// Create new router without auth middleware.
	router := gin.Default()
	router.GET("/api/cart", suite.controller.GetCartItems)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *CartItemControllerTestSuite) TestRemoveCartItem_Success() {
	// Setup
	listingID := "listing123"
	suite.mockUC.On("RemoveCartItem", mock.Anything, suite.userID, listingID).Return(nil)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/cart/items/"+listingID, nil)
	suite.router.DELETE("/api/cart/items/:listingID", suite.controller.RemoveCartItem)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "item removed from cart", response["message"])
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *CartItemControllerTestSuite) TestRemoveCartItem_Unauthorized() {
	// Setup
	listingID := "listing123"

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/cart/items/"+listingID, nil)
	// Create new router without auth middleware.
	router := gin.Default()
	router.DELETE("/api/cart/items/:listingID", suite.controller.RemoveCartItem)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *CartItemControllerTestSuite) TestCheckoutCart_Success() {
	// Setup
	dummyResp := &models.CheckoutResponse{
		TotalAmount: 100.0,
		PlatformFee: 2.0,
		NetPayable:  98.0,
		Items: []models.CheckoutItemResponse{
			{
				ListingID: "item1",
				Title:     "Test Item 1",
				Price:     100.0,
				SellerID:  "seller1",
				Status:    "available",
			},
		},
	}
	suite.mockUC.On("CheckoutCart", mock.Anything, suite.userID).Return(dummyResp, nil)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/checkout", nil)
	suite.router.POST("/api/checkout", suite.controller.CheckoutCart)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), true, response["success"])
	assert.Equal(suite.T(), "Payment successful. Order confirmed.", response["message"])

	// Verify that data matches dummyResp
	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), dummyResp.TotalAmount, data["totalAmount"])
	assert.Equal(suite.T(), dummyResp.PlatformFee, data["platformFee"])
	assert.Equal(suite.T(), dummyResp.NetPayable, data["netPayable"])
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *CartItemControllerTestSuite) TestCheckoutCart_ValidationError() {
	// Setup
	suite.mockUC.On("CheckoutCart", mock.Anything, suite.userID).Return(nil, errors.New("some items are unavailable"))

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/checkout", nil)
	suite.router.POST("/api/checkout", suite.controller.CheckoutCart)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), "some items are unavailable", response["error"])
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *CartItemControllerTestSuite) TestCheckoutCart_Unauthorized() {
	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/checkout", nil)
	// Create new router without auth middleware.
	router := gin.Default()
	router.POST("/api/checkout", suite.controller.CheckoutCart)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *CartItemControllerTestSuite) TestCheckoutSingleItem_Success() {
	// Setup
	dummyResp := &models.CheckoutResponse{
		TotalAmount: 100.0,
		PlatformFee: 2.0,
		NetPayable:  98.0,
		Items: []models.CheckoutItemResponse{
			{
				ListingID: "item1",
				Title:     "Test Item 1",
				Price:     100.0,
				SellerID:  "seller1",
				Status:    "available",
			},
		},
	}
	listingID := "listing123"
	suite.mockUC.On("CheckoutSingleItem", mock.Anything, suite.userID, listingID).Return(dummyResp, nil)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/checkout/"+listingID, nil)
	suite.router.POST("/api/checkout/:listingId", suite.controller.CheckoutSingleItem)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(suite.T(), true, response["success"])
	assert.Equal(suite.T(), "Payment successful. Order confirmed.", response["message"])

	// Verify that data matches dummyResp
	data := response["data"].(map[string]interface{})
	assert.Equal(suite.T(), dummyResp.TotalAmount, data["totalAmount"])
	assert.Equal(suite.T(), dummyResp.PlatformFee, data["platformFee"])
	assert.Equal(suite.T(), dummyResp.NetPayable, data["netPayable"])
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *CartItemControllerTestSuite) TestCheckoutSingleItem_Unauthorized() {
	// Setup
	listingID := "listing123"
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/checkout/"+listingID, nil)
	// Create new router without auth middleware.
	router := gin.Default()
	router.POST("/api/checkout/:listingId", suite.controller.CheckoutSingleItem)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func TestCartItemControllerSuite(t *testing.T) {
	suite.Run(t, new(CartItemControllerTestSuite))
}
