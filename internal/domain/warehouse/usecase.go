package warehouse

import (
	"context"
)

type Usecase interface {
	GetWarehouseItems(ctx context.Context, resellerID string) ([]*WarehouseItem, error)
}
