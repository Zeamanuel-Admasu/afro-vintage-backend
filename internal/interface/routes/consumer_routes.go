package routes

import (
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterConsumerRoutes(r *gin.Engine, ctrl *controllers.ConsumerController, jwtSvc auth.JWTService) {
	consumerGroup := r.Group("/orders")
	consumerGroup.Use(middlewares.AuthMiddleware(jwtSvc))

	consumerGroup.GET("/history", middlewares.AuthorizeRoles("consumer"), ctrl.GetOrderHistory) 
}
