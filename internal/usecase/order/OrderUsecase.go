package OrderUsecase

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/payment"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func NewOrderUsecase(bRepo bundle.Repository, oRepo order.Repository, wRepo warehouse.Repository, pRepo payment.Repository, uRepo user.Repository) *orderUseCaseImpl {
	return &orderUseCaseImpl{
		bundleRepo:    bRepo,
		orderRepo:     oRepo,
		warehouseRepo: wRepo,
		paymentRepo:   pRepo,
		userRepo:      uRepo,
	}
}

func simulateStripePayment(total float64) (string, error) {
	time.Sleep(500 * time.Millisecond)
	return fmt.Sprintf("ch_%d", rand.Intn(1000000)), nil
}

func processPayment(total float64) (fee float64, net float64, err error) {
	fee = total * 0.02
	net = total - fee
	_, err = simulateStripePayment(total)
	return
}

func (uc *orderUseCaseImpl) PurchaseBundle(ctx context.Context, bundleID, resellerID string) (*order.Order, *payment.Payment, *warehouse.WarehouseItem, error) {
	b, err := uc.bundleRepo.GetBundleByID(ctx, bundleID)
	if err != nil {
		return nil, nil, nil, err
	}

	availables, err := uc.bundleRepo.ListAvailableBundles(ctx)
	if err != nil {
		return nil, nil, nil, err
	}
	found := false
	for _, available := range availables {
		if available.ID == b.ID {
			found = true
			break
		}
	}
	if !found {
		return nil, nil, nil, errors.New("bundle not available")
	}

	// prevent self-purchase
	if b.SupplierID == resellerID {
		return nil, nil, nil, errors.New("reseller cannot purchase their own bundle")
	}

	// process payment
	fee, net, err := processPayment(b.Price)
	if err != nil {
		return nil, nil, nil, err
	}

	order := &order.Order{
		ID:          primitive.NewObjectID().Hex(),
		BundleID:    b.ID,
		ResellerID:  resellerID,
		SupplierID:  b.SupplierID,
		TotalPrice:  b.Price,
		PlatformFee: fee,
		Status:      order.OrderStatusProcessing,
		CreatedAt:   time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
	}
	if err := uc.orderRepo.CreateOrder(ctx, order); err != nil {
		return nil, nil, nil, err
	}

	payment := &payment.Payment{
		FromUserID:    resellerID,
		ToUserID:      b.SupplierID,
		Amount:        b.Price,
		PlatformFee:   fee,
		SellerEarning: net,
		Status:        "Paid",
		ReferenceID:   b.ID,
		Type:          payment.B2B,
		CreatedAt:     time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
	}
	if err := uc.paymentRepo.RecordPayment(ctx, payment); err != nil {
		return nil, nil, nil, err
	}

	if err := uc.bundleRepo.MarkAsPurchased(ctx, b.ID, resellerID); err != nil {
		return nil, nil, nil, err
	}

	warehouseItem := &warehouse.WarehouseItem{
		ID:         primitive.NewObjectID().Hex(),
		BundleID:   b.ID,
		ResellerID: resellerID,
		Status:     "pending",
	}
	if err := uc.warehouseRepo.AddItem(ctx, warehouseItem); err != nil {
		return nil, nil, nil, err
	}

	if err := uc.orderRepo.UpdateOrderStatus(ctx, order.ID, "completed"); err != nil {
		return nil, nil, nil, err
	}

	go func(itemID string) {
		time.Sleep(3 * time.Minute)
		if err := uc.warehouseRepo.MarkItemAsListed(context.Background(), itemID); err != nil {
			fmt.Println("Failed to mark warehouse item as listed:", err)
		}
	}(warehouseItem.ID)

	return order, payment, warehouseItem, nil
}

func (uc *orderUseCaseImpl) GetDashboardMetrics(ctx context.Context, supplierID string) (*order.DashboardMetrics, error) {
	bundles, err := uc.bundleRepo.ListBundles(ctx, supplierID)
	if err != nil {
		return nil, err
	}
	totalSales := 0.0
	activeCount := 0
	soldCount := 0
	var activeBundles []*bundle.Bundle
	bestSelling := 0.0

	// Get user trust score
	userData, err := uc.userRepo.GetByID(ctx, supplierID)
	if err != nil {
		return nil, err
	}

	for _, b := range bundles {
		if b.Status == "purchased" {
			totalSales += b.Price
			soldCount++
			// Track the highest selling bundle
			if b.Price > bestSelling {
				bestSelling = b.Price
			}
		} else if b.Status == "available" {
			activeCount++
			activeBundles = append(activeBundles, b)
		}
	}

	sort.Slice(activeBundles, func(i, j int) bool {
		return activeBundles[i].DateListed.After(activeBundles[j].DateListed)
	})

	performanceMetrice := order.PerformanceMetrics{
		TotalBundlesListed: len(bundles),
		ActiveCount:        activeCount,
		SoldCount:          soldCount,
	}

	metrics := &order.DashboardMetrics{
		TotalSales:         totalSales,
		ActiveBundles:      activeBundles,
		PerformanceMetrics: performanceMetrice,
		Rating:             userData.TrustScore,
		BestSelling:        bestSelling,
	}

	return metrics, nil
}

func (uc *orderUseCaseImpl) GetOrderByID(ctx context.Context, orderID string) (*order.Order, error) {
	return uc.orderRepo.GetOrderByID(ctx, orderID)
}
func (uc *orderUseCaseImpl) GetSoldBundleHistory(ctx context.Context, supplierID string) ([]*order.Order, error) {
    orders, err := uc.orderRepo.GetOrdersBySupplier(ctx, supplierID)
    if err != nil {
        return nil, err
    }

    var soldBundleOrders []*order.Order
    for _, order := range orders {
        if order.BundleID != "" && len(order.ProductIDs) == 0 {
            soldBundleOrders = append(soldBundleOrders, order)
        }
    }

    return soldBundleOrders, nil
}

func (uc *orderUseCaseImpl) GetResellerMetrics(ctx context.Context, resellerID string) (*order.ResellerMetrics, error) {
	// Get bought bundles
	bundles, err := uc.bundleRepo.ListPurchasedByReseller(ctx, resellerID)
	if err != nil {
		return nil, err
	}

	// Get user trust score
	userData, err := uc.userRepo.GetByID(ctx, resellerID)
	if err != nil {
		return nil, err
	}

	totalBoughtBundles := len(bundles)
	bestSelling := 0.0

	// Get all orders for this reseller to count actual sold items
	orders, err := uc.orderRepo.GetOrdersByReseller(ctx, resellerID)
	if err != nil {
		return nil, err
	}

	// Count actual sold items from orders
	totalItemsSold := 0
	for _, order := range orders {
		if order.Status == "completed" {
			// For product purchases, count the number of products
			if len(order.ProductIDs) > 0 {
				totalItemsSold += len(order.ProductIDs)
			}
			// For bundle purchases, get the bundle quantity
			if order.BundleID != "" {
				bundle, err := uc.bundleRepo.GetBundleByID(ctx, order.BundleID)
				if err == nil && bundle != nil {
					totalItemsSold += bundle.Quantity
				}
			}
		}
	}

	// Calculate best selling from bundles
	for _, b := range bundles {
		if b.Price > bestSelling {
			bestSelling = b.Price
		}
	}

	metrics := &order.ResellerMetrics{
		TotalBoughtBundles: totalBoughtBundles,
		TotalItemsSold:     totalItemsSold,
		Rating:             userData.TrustScore,
		BestSelling:        bestSelling,
		BoughtBundles:      bundles,
	}

	return metrics, nil
}