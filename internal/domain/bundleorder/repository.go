package bundleorder

import "context"

type Repository interface {
    CreateOrder(ctx context.Context, order *BundleOrder) error
    GetOrdersByBuyerID(ctx context.Context, buyerID string) ([]*BundleOrder, error)
    GetOrderByBundleID(ctx context.Context, bundleID string) (*BundleOrder, error)
    GetOrdersBySellerID(ctx context.Context, sellerID string) ([]*BundleOrder, error)
}