package bundle

import "context"

type Repository interface {
	CreateBundle(ctx context.Context, b *Bundle) error
	GetBundleByID(ctx context.Context, id string) (*Bundle, error)
	ListBundles(ctx context.Context, supplierID string) ([]*Bundle, error)
	ListAvailableBundles(ctx context.Context) ([]*Bundle, error)
	ListPurchasedByReseller(ctx context.Context, resellerID string) ([]*Bundle, error)
	UpdateBundleStatus(ctx context.Context, id string, status string) error
	MarkAsPurchased(ctx context.Context, bundleID string, resellerID string) error
}
