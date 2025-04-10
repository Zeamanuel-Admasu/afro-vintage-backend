package routes

import (
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r *gin.Engine, ctrl *controllers.ProductController) {
	productGroup := r.Group("/products")
	{
		productGroup.POST("", ctrl.Create)
		productGroup.GET("/:id", ctrl.GetByID)
		productGroup.GET("", ctrl.ListAvailable)
		productGroup.GET("/reseller/:id", ctrl.ListByReseller)
		productGroup.PUT("/:id", ctrl.Update)
		productGroup.DELETE("/:id", ctrl.Delete)
	}
}
