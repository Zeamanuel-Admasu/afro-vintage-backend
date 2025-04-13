package bundleorder

import (
    "context"
    "errors"
    "fmt"
    "time"

    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundleorder"
)

type BundleOrderUsecase struct {
    orderRepo  bundleorder.Repository
    bundleRepo bundle.Repository
}

func NewBundleOrderUsecase(orderRepo bundleorder.Repository, bundleRepo bundle.Repository) *BundleOrderUsecase {
    return &BundleOrderUsecase{
        orderRepo:  orderRepo,
        bundleRepo: bundleRepo,
    }
}

func (uc *BundleOrderUsecase) CreateOrder(ctx context.Context, bundleID string, resellerID string) (*bundleorder.BundleOrder, error) {
    // Fetch the bundle
    bundle, err := uc.bundleRepo.GetBundleByID(ctx, bundleID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch bundle: %w", err)
    }

    // Validate bundle status
    if bundle.Status != "available" {
        return nil, errors.New("bundle not available for trade")
    }

    // Validate reseller isn't the supplier
    if bundle.SupplierID == resellerID {
        return nil, errors.New("cannot purchase your own bundle")
    }

    // Create the order
    orderID := fmt.Sprintf("order_%s_%d", resellerID, time.Now().Unix())
    newOrder := &bundleorder.BundleOrder{
        ID:             orderID,
        BuyerID:        resellerID,
        SellerID:       bundle.SupplierID,
        BundleID:       bundleID,
        WarehouseStatus: "unpacked",
        CreatedAt:      time.Now(),
    }
    if err := uc.orderRepo.CreateOrder(ctx, newOrder); err != nil {
        return nil, fmt.Errorf("failed to create order: %w", err)
    }

    // Mark the bundle as purchased
    if err := uc.bundleRepo.MarkAsPurchased(ctx, bundleID, resellerID); err != nil {
        return nil, fmt.Errorf("failed to mark bundle as purchased: %w", err)
    }

    return newOrder, nil
}

func (uc *BundleOrderUsecase) GetOrdersBySellerID(ctx context.Context, sellerID string) ([]*bundleorder.BundleOrder, error) {
    orders, err := uc.orderRepo.GetOrdersBySellerID(ctx, sellerID)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch orders for seller: %w", err)
    }
    return orders, nil
}