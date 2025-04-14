package order

import (
	"context"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/payment"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
)

type Usecase interface {
	PurchaseBundle(ctx context.Context, bundleID, resellerID string) (*Order, *payment.Payment, *warehouse.WarehouseItem, error)
	GetDashboardMetrics(ctx context.Context, supplierID string) (*DashboardMetrics, error)
	GetOrderByID(ctx context.Context, orderID string) (*Order, error)
}
