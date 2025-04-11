package controllers

import (
	"net/http"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
	userUC user.Usecase
}

func NewAdminController(userUC user.Usecase) *AdminController {
	return &AdminController{userUC: userUC}
}

// GET /api/admin/users
func (a *AdminController) GetAllUsers(c *gin.Context) {
	roleParam := c.Query("role")

	// If no role is passed, return all users by calling ListByRole for each role
	if roleParam == "" {
		allRoles := []user.Role{
			user.RoleSupplier,
			user.RoleReseller,
			user.RoleConsumer,
			user.RoleAdmin,
		}

		var allUsers []*user.User
		for _, role := range allRoles {
			usersByRole, err := a.userUC.ListByRole(c.Request.Context(), role)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users by role"})
				return
			}
			allUsers = append(allUsers, usersByRole...)
		}

		c.JSON(http.StatusOK, allUsers)
		return
	}

	// Otherwise, get users by specified role
	role := user.Role(roleParam)
	users, err := a.userUC.ListByRole(c.Request.Context(), role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, users)
}
