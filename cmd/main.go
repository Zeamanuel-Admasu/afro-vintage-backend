package main

import (
	"github.com/gin-gonic/gin"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/config"
	authinfra "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/infrastructure/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/infrastructure/mongo"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/routes"

	authusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/auth"
	bundleusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/bundle"
	bundleorderusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/bundleorder"
	productusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/product"
	userusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/user"
)

func main() {
	// Load .env variables
	config.LoadEnv()

	// Load grouped app config
	appConfig := config.LoadAppConfig()

	// Connect to MongoDB
	db := config.ConnectMongo(appConfig.DBURI, appConfig.DBName)

	// Init shared services
	jwtSvc := authinfra.NewJWTService(appConfig.JWTSecret)
	passSvc := authinfra.NewPasswordService()

	// Init Repositories
	userRepo := mongo.NewMongoUserRepository(db)
	productRepo := mongo.NewMongoProductRepository(db)
	bundleRepo := mongo.NewBundleRepository(db)
	bundleOrderRepo := mongo.NewBundleOrderRepository(db)

	// Init Usecases
	userUC := userusecase.NewUserUsecase(userRepo)
	authUC := authusecase.NewAuthUsecase(userRepo, passSvc, jwtSvc)
	productUC := productusecase.NewProductUsecase(productRepo)
	bundleUC := bundleusecase.NewBundleUsecase(bundleRepo)
	bundleOrderUC := bundleorderusecase.NewBundleOrderUsecase(bundleOrderRepo, bundleRepo)

	// Init Controllers
	authCtrl := controllers.NewAuthController(authUC)
	adminCtrl := controllers.NewAdminController(userUC)
	productCtrl := controllers.NewProductController(productUC)
	bundleCtrl := controllers.NewBundleController(bundleUC)
	bundleOrderCtrl := controllers.NewBundleOrderController(bundleOrderUC)
	supplierCtrl := controllers.NewSupplierController(bundleOrderUC, bundleUC) // Added

	// Init Gin Engine and Routes
	r := gin.Default()

	routes.RegisterAuthRoutes(r, authCtrl)
	routes.RegisterAdminRoutes(r, adminCtrl, jwtSvc)
	routes.RegisterProductRoutes(r, productCtrl, jwtSvc)
	routes.RegisterBundleRoutes(r, bundleCtrl, jwtSvc)
	routes.RegisterBundleOrderRoutes(r, bundleOrderCtrl, jwtSvc)
	routes.RegisterSupplierRoutes(r, supplierCtrl, jwtSvc) // Added

	// Run server
	r.Run(":8080")
}