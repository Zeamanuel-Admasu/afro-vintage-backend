package main

import (
	"github.com/gin-gonic/gin"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/config"
	authinfra "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/infrastructure/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/infrastructure/mongo"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/routes"

	authusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/auth"
	cartitemusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/cartitem"

	bundleusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/bundle"
	orderusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/order"
	productusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/product"
	reviewusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/review"
	trustusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/trust"
	userusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/user"
	warehouse_usecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/warehouse"
)

func main() {
	// Load .env variables
	config.LoadEnv()

	// Set Gin to release mode
	gin.SetMode(gin.DebugMode)

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
	orderRepo := mongo.NewMongoOrderRepository(db) // Add order repository
	cartItemRepo := mongo.NewCartItemRepository(db)
	reviewRepo := mongo.NewReviewRepository(db)            // Add review repository
	warehouseRepo := mongo.NewMongoWarehouseRepository(db) // Add warehouse repository
	paymentRepo := mongo.NewMongoPaymentRepository(db)     // Add payment repository

	// Init Usecases
	userUC := userusecase.NewUserUsecase(userRepo)
	authUC := authusecase.NewAuthUsecase(userRepo, passSvc, jwtSvc)
	productUC := productusecase.NewProductUsecase(productRepo, bundleRepo)
	bundleUC := bundleusecase.NewBundleUsecase(bundleRepo)
	trustUC := trustusecase.NewTrustUsecase(productRepo, bundleRepo, userRepo)
	cartItemUC := cartitemusecase.NewCartItemUsecase(cartItemRepo, productRepo)

	reviewUC := reviewusecase.NewReviewUsecase(reviewRepo, orderRepo)                                     // Add review usecase
	orderSvc := orderusecase.NewOrderUsecase(bundleRepo, orderRepo, warehouseRepo, paymentRepo, userRepo) // Add order service
	warehouseSvc := warehouse_usecase.NewWarehouseUseCase(warehouseRepo)

	// Init Controllers
	authCtrl := controllers.NewAuthController(authUC)
	adminCtrl := controllers.NewAdminController(userUC, orderSvc)
	productCtrl := controllers.NewProductController(productUC, trustUC, bundleUC, warehouseRepo)
	bundleCtrl := controllers.NewBundleController(bundleUC, userUC)
	consumerCtrl := controllers.NewConsumerController(orderRepo)
	supplierCtrl := controllers.NewSupplierController(orderSvc) // Add consumer controller
	cartItemCtrl := controllers.NewCartItemController(cartItemUC)
	reviewCtrl := controllers.NewReviewController(reviewUC) // Add review controller
	warehouseCtrl := controllers.NewWarehouseController(warehouseSvc)
	orderCtrl := controllers.NewOrderController(orderSvc) // Add order controller

	// Init Gin Engine and Routes
	r := gin.Default()

	routes.RegisterAuthRoutes(r, authCtrl)
	routes.RegisterProductRoutes(r, productCtrl, jwtSvc, reviewCtrl) // Register product routes with review controller
	routes.RegisterAdminRoutes(r, adminCtrl, jwtSvc)
	routes.RegisterBundleRoutes(r, bundleCtrl, jwtSvc)
	routes.RegisterCartItemRoutes(r, cartItemCtrl, jwtSvc) // Register cart item routes

	routes.RegisterOrderRoutes(r, orderCtrl, consumerCtrl, jwtSvc) // Register order routes
	routes.RegisterSupplierRoutes(r, supplierCtrl, jwtSvc)
	routes.RegisterWarehouseRoutes(r, warehouseCtrl, jwtSvc)
	routes.RegisterResellerRoutes(r, supplierCtrl, jwtSvc)

	// Run server
	r.Run(":8080")
}
