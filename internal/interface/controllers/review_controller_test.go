package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/review"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockReviewUsecase struct {
	mock.Mock
}

func (m *MockReviewUsecase) SubmitReview(ctx context.Context, r *review.Review) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

type ReviewControllerTestSuite struct {
	suite.Suite
	usecase    *MockReviewUsecase
	controller *ReviewController
	router     *gin.Engine
}

func (suite *ReviewControllerTestSuite) SetupTest() {
	suite.usecase = new(MockReviewUsecase)
	suite.controller = NewReviewController(suite.usecase)
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
}

func TestReviewControllerTestSuite(t *testing.T) {
	suite.Run(t, new(ReviewControllerTestSuite))
}

func (suite *ReviewControllerTestSuite) TestSubmitReview_Success() {
	// Setup
	req := models.CreateReviewRequest{
		OrderID:   "order123",
		ProductID: "product123",
		Rating:    4,
		Comment:   "Great product!",
	}

	suite.usecase.On("SubmitReview", mock.Anything, mock.Anything).
		Return(nil)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user123")

	body, _ := json.Marshal(req)
	request := httptest.NewRequest("POST", "/reviews", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	c.Request = request

	// Execute
	suite.controller.SubmitReview(c)

	// Assert
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	suite.usecase.AssertExpectations(suite.T())
}

func (suite *ReviewControllerTestSuite) TestSubmitReview_InvalidPayload() {
	// Setup
	invalidPayload := "invalid json"

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user123")

	request := httptest.NewRequest("POST", "/reviews", bytes.NewBufferString(invalidPayload))
	request.Header.Set("Content-Type", "application/json")
	c.Request = request

	// Execute
	suite.controller.SubmitReview(c)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	suite.usecase.AssertNotCalled(suite.T(), "SubmitReview")
}

func (suite *ReviewControllerTestSuite) TestSubmitReview_Unauthorized() {
	// Setup
	req := models.CreateReviewRequest{
		OrderID:   "order123",
		ProductID: "product123",
		Rating:    4,
		Comment:   "Great product!",
	}

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Note: No userID set

	body, _ := json.Marshal(req)
	request := httptest.NewRequest("POST", "/reviews", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	c.Request = request

	// Execute
	suite.controller.SubmitReview(c)

	// Assert
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
	suite.usecase.AssertNotCalled(suite.T(), "SubmitReview")
}

func (suite *ReviewControllerTestSuite) TestSubmitReview_UseCaseError() {
	// Setup
	req := models.CreateReviewRequest{
		OrderID:   "order123",
		ProductID: "product123",
		Rating:    4,
		Comment:   "Great product!",
	}

	expectedError := errors.New("review already exists")
	suite.usecase.On("SubmitReview", mock.Anything, mock.Anything).
		Return(expectedError)

	// Create test request
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", "user123")

	body, _ := json.Marshal(req)
	request := httptest.NewRequest("POST", "/reviews", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	c.Request = request

	// Execute
	suite.controller.SubmitReview(c)

	// Assert
	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
	suite.usecase.AssertExpectations(suite.T())
}
