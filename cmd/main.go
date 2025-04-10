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

	// Init Usecases
	userUC := userusecase.NewUserUsecase(userRepo)
	authUC := authusecase.NewAuthUsecase(userRepo, passSvc, jwtSvc)
	productUC := productusecase.NewProductUsecase(productRepo)
	bundleUC := bundleusecase.NewBundleUsecase(bundleRepo)

	// Init Controllers
	authCtrl := controllers.NewAuthController(authUC)
	adminCtrl := controllers.NewAdminController(userUC)
	productCtrl := controllers.NewProductController(productUC)
	bundleCtrl := controllers.NewBundleController(bundleUC)

	// Init Gin Engine and Routes
	r := gin.Default()

	routes.RegisterAuthRoutes(r, authCtrl)
	routes.RegisterProductRoutes(r, productCtrl, jwtSvc)
	routes.RegisterAdminRoutes(r, adminCtrl, jwtSvc)
	routes.RegisterBundleRoutes(r, bundleCtrl, jwtSvc)

	// Run server
	r.Run(":8080")
}
