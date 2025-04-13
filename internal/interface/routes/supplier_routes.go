package routes

import (
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/middlewares"
    "github.com/gin-gonic/gin"
)

func RegisterSupplierRoutes(r *gin.Engine, ctrl *controllers.SupplierController, jwtSvc auth.JWTService) {
    supplier := r.Group("/supplier")
    supplier.Use(middlewares.AuthMiddleware(jwtSvc))
    supplier.Use(middlewares.AuthorizeRoles("supplier"))
    {
        supplier.GET("/sold-bundles", ctrl.GetSoldBundles)
    }
}
