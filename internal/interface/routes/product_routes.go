package routes

import (
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.Engine, ctrl *controllers.ProductController, jwtSvc auth.JWTService, reviewCtrl *controllers.ReviewController) {
	productGroup := r.Group("/products")
	productGroup.Use(middlewares.AuthMiddleware(jwtSvc))

	productGroup.POST("", middlewares.AuthorizeRoles("reseller", "admin"), ctrl.Create)
	productGroup.GET("/:id", middlewares.AuthorizeRoles("consumer", "reseller", "admin"), ctrl.GetByID)
	productGroup.GET("", middlewares.AuthorizeRoles("consumer", "reseller", "admin"), ctrl.ListAvailable)
	productGroup.GET("/reseller/:id", middlewares.AuthorizeRoles("admin", "reseller"), ctrl.ListByReseller)
	productGroup.PUT("/:id", middlewares.AuthorizeRoles("reseller", "admin"), ctrl.Update)
	productGroup.DELETE("/:id", middlewares.AuthorizeRoles("reseller", "admin"), ctrl.Delete)
	productGroup.POST("/reviews", middlewares.AuthorizeRoles("consumer"), reviewCtrl.SubmitReview)
}
