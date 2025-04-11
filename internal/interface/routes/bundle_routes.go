package routes

import (
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/middlewares"
    "github.com/gin-gonic/gin"
)

func RegisterBundleRoutes(r *gin.Engine, ctrl *controllers.BundleController, jwtSvc auth.JWTService) {
    bundleGroup := r.Group("/bundles")
    bundleGroup.Use(middlewares.AuthMiddleware(jwtSvc)) // All routes require valid token

    bundleGroup.POST("", middlewares.AuthorizeRoles("supplier"), ctrl.CreateBundle)
    bundleGroup.GET("", middlewares.AuthorizeRoles("supplier"), ctrl.ListBundles)
    bundleGroup.DELETE("/:id", middlewares.AuthorizeRoles("supplier"), ctrl.DeleteBundle) // Added
}