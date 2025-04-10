package main

import (
	// ...

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/config"
	productinfra "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/infrastructure/mongo"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/routes"
	productusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/product"
	"github.com/gin-gonic/gin"
)

func main() {
	db := config.ConnectMongo("mongodb://localhost:27017", "afro_vintage")

	productRepo := productinfra.NewMongoProductRepository(db)
	productUC := productusecase.NewProductUsecase(productRepo)
	productCtrl := controllers.NewProductController(productUC)
	r := gin.Default()
	routes.RegisterProductRoutes(r, productCtrl)
}
