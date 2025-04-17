package OrderUsecase

import (
	"context"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/payment"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
)

type OrderUseCase interface {
	PurchaseBundle(ctx context.Context, bundleID, resellerID string) (*order.Order, *payment.Payment, *warehouse.WarehouseItem, error)
	GetDashboardMetrics(ctx context.Context, supplierID string) (*order.DashboardMetrics, error)
	GetOrderByID(ctx context.Context, orderID string)
}

type orderUseCaseImpl struct {
	bundleRepo    bundle.Repository
	orderRepo     order.Repository
	warehouseRepo warehouse.Repository
	paymentRepo   payment.Repository
	userRepo      user.Repository
}
