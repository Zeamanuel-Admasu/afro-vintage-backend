package routes

import (
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/middlewares"
    "github.com/gin-gonic/gin"
)

func RegisterBundleOrderRoutes(r *gin.Engine, ctrl *controllers.BundleOrderController, jwtSvc auth.JWTService) {
    orders := r.Group("/orders")
    orders.Use(middlewares.AuthMiddleware(jwtSvc))
    orders.Use(middlewares.AuthorizeRoles("reseller"))
    {
        orders.POST("", ctrl.CreateOrder)
    }
}