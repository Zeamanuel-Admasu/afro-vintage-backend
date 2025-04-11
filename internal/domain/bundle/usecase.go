package bundle

import "context"

type Usecase interface {
    CreateBundle(ctx context.Context, supplierID string, bundle *Bundle) error
    ListBundles(ctx context.Context, supplierID string) ([]*Bundle, error)
    DeleteBundle(ctx context.Context, supplierID string, bundleID string) error // Added
}