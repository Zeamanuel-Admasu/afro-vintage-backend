package routes

import (
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterSupplierRoutes(r *gin.Engine, ctrl *controllers.SupplierController, jwtSvc auth.JWTService) {
	supplierGroup := r.Group("/supplier")
	supplierGroup.Use(middlewares.AuthMiddleware(jwtSvc))

	supplierGroup.GET("/dashboard", middlewares.AuthorizeRoles("supplier"), ctrl.GetDashboardMetrics)
    supplierGroup.GET("/sold-history", ctrl.ListSoldBundles) 
   
}
