package bundleorder

import "context"

type Usecase interface {
    CreateOrder(ctx context.Context, bundleID string, resellerID string) (*BundleOrder, error)
    GetOrdersBySellerID(ctx context.Context, sellerID string) ([]*BundleOrder, error) // New method
}