package warehouse

import "context"

type Repository interface {
	AddItem(ctx context.Context, item *WarehouseItem) error
	GetItemsByReseller(ctx context.Context, resellerID string) ([]*WarehouseItem, error)
	GetItemsByBundle(ctx context.Context, bundleID string) ([]*WarehouseItem, error)
	MarkItemAsListed(ctx context.Context, itemID string) error
	MarkItemAsSkipped(ctx context.Context, itemID string) error
	DeleteItem(ctx context.Context, itemID string) error
}
