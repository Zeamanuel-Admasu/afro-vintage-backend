package warehouse_usecase

import (
	"context"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/warehouse"
)

type WarehouseUseCase interface {
	GetWarehouseItems(ctx context.Context, resellerID string) ([]*warehouse.WarehouseItem, error)
}

type warehouseUseCaseImpl struct {
	warehouseRepo warehouse.Repository
}

func NewWarehouseUseCase(repo warehouse.Repository) WarehouseUseCase {
	return &warehouseUseCaseImpl{warehouseRepo: repo}
}

func (uc *warehouseUseCaseImpl) GetWarehouseItems(ctx context.Context, resellerID string) ([]*warehouse.WarehouseItem, error) {
	return uc.warehouseRepo.GetItemsByReseller(ctx, resellerID)
}
