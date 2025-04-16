package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) ListByRole(ctx context.Context, role user.Role) ([]*user.User, error) {
	args := m.Called(ctx, role)
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *MockUserUsecase) GetByID(ctx context.Context, id string) (*user.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*user.User), args.Error(1)
}

func (m *MockUserUsecase) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockUserUsecase) GetBlacklistedUsers(ctx context.Context) ([]*user.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*user.User), args.Error(1)
}

func (m *MockUserUsecase) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserUsecase) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*user.User), args.Error(1)
}

type AdminControllerTestSuite struct {
	suite.Suite
	controller *AdminController
	mockUC     *MockUserUsecase
	router     *gin.Engine
}

func (suite *AdminControllerTestSuite) SetupTest() {
	suite.mockUC = new(MockUserUsecase)
	suite.controller = NewAdminController(suite.mockUC)
	gin.SetMode(gin.TestMode)
	suite.router = gin.Default()
}

func (suite *AdminControllerTestSuite) TestGetAllUsers_NoRoleParam() {
	// Setup
	users := []*user.User{
		{ID: "1", Name: "User1", Role: string(user.RoleSupplier)},
		{ID: "2", Name: "User2", Role: string(user.RoleConsumer)},
	}

	suite.mockUC.On("ListByRole", mock.Anything, user.RoleSupplier).Return(users[:1], nil)
	suite.mockUC.On("ListByRole", mock.Anything, user.RoleReseller).Return([]*user.User{}, nil)
	suite.mockUC.On("ListByRole", mock.Anything, user.RoleConsumer).Return(users[1:], nil)
	suite.mockUC.On("ListByRole", mock.Anything, user.RoleAdmin).Return([]*user.User{}, nil)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/users", nil)
	suite.router.GET("/api/admin/users", suite.controller.GetAllUsers)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *AdminControllerTestSuite) TestGetAllUsers_WithRoleParam() {
	// Setup
	users := []*user.User{
		{ID: "1", Name: "User1", Role: string(user.RoleSupplier)},
	}
	suite.mockUC.On("ListByRole", mock.Anything, user.RoleSupplier).Return(users, nil)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/users?role=supplier", nil)
	suite.router.GET("/api/admin/users", suite.controller.GetAllUsers)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *AdminControllerTestSuite) TestDeleteUserIfBlacklisted_Success() {
	// Setup
	userID := "1"
	userData := &user.User{
		ID:         userID,
		Role:       string(user.RoleSupplier),
		TrustScore: 50,
	}
	suite.mockUC.On("GetByID", mock.Anything, userID).Return(userData, nil)
	suite.mockUC.On("Update", mock.Anything, userID, map[string]interface{}{"is_deleted": true}).Return(nil)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/admin/users/"+userID, nil)
	suite.router.DELETE("/api/admin/users/:userId", suite.controller.DeleteUserIfBlacklisted)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *AdminControllerTestSuite) TestDeleteUserIfBlacklisted_UserNotFound() {
	// Setup
	userID := "1"
	suite.mockUC.On("GetByID", mock.Anything, userID).Return((*user.User)(nil), nil)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/admin/users/"+userID, nil)
	suite.router.DELETE("/api/admin/users/:userId", suite.controller.DeleteUserIfBlacklisted)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *AdminControllerTestSuite) TestDeleteUserIfBlacklisted_InvalidRole() {
	// Setup
	userID := "1"
	userData := &user.User{
		ID:         userID,
		Role:       string(user.RoleConsumer),
		TrustScore: 50,
	}
	suite.mockUC.On("GetByID", mock.Anything, userID).Return(userData, nil)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/admin/users/"+userID, nil)
	suite.router.DELETE("/api/admin/users/:userId", suite.controller.DeleteUserIfBlacklisted)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *AdminControllerTestSuite) TestGetTrustScores_AllRoles() {
	// Setup
	users := []*user.User{
		{ID: "1", Name: "User1", Role: string(user.RoleSupplier), TrustScore: 70},
		{ID: "2", Name: "User2", Role: string(user.RoleReseller), TrustScore: 50},
	}
	suite.mockUC.On("ListByRole", mock.Anything, user.RoleSupplier).Return(users[:1], nil)
	suite.mockUC.On("ListByRole", mock.Anything, user.RoleReseller).Return(users[1:], nil)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/trust-scores", nil)
	suite.router.GET("/api/admin/trust-scores", suite.controller.GetTrustScores)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockUC.AssertExpectations(suite.T())
}

func (suite *AdminControllerTestSuite) TestGetBlacklistedUsers() {
	// Setup
	users := []*user.User{
		{ID: "1", Name: "User1", Role: string(user.RoleSupplier), TrustScore: 50},
		{ID: "2", Name: "User2", Role: string(user.RoleReseller), TrustScore: 40},
	}
	suite.mockUC.On("GetBlacklistedUsers", mock.Anything).Return(users, nil)

	// Execute
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/blacklisted-users", nil)
	suite.router.GET("/api/admin/blacklisted-users", suite.controller.GetBlacklistedUsers)
	suite.router.ServeHTTP(w, req)

	// Assert
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockUC.AssertExpectations(suite.T())
}

func TestAdminControllerSuite(t *testing.T) {
	suite.Run(t, new(AdminControllerTestSuite))
}
