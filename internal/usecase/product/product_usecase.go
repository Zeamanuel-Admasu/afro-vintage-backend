// internal/usecase/product/product_usecase.go
package productusecase

import (
	"context"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
)

type productUsecase struct {
	repo product.Repository
}

func NewProductUsecase(repo product.Repository) product.Usecase {
	return &productUsecase{repo: repo}
}

func (uc *productUsecase) AddProduct(ctx context.Context, p *product.Product) error {
	return uc.repo.AddProduct(ctx, p)
}

func (uc *productUsecase) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	return uc.repo.GetProductByID(ctx, id)
}

func (uc *productUsecase) ListProductsByReseller(ctx context.Context, resellerID string, page, limit int) ([]*product.Product, error) {
	return uc.repo.ListProductsByReseller(ctx, resellerID, page, limit)
}

func (uc *productUsecase) ListAvailableProducts(ctx context.Context, page, limit int) ([]*product.Product, error) {
	return uc.repo.ListAvailableProducts(ctx, page, limit)
}

func (uc *productUsecase) DeleteProduct(ctx context.Context, id string) error {
	return uc.repo.DeleteProduct(ctx, id)
}

func (uc *productUsecase) UpdateProduct(ctx context.Context, id string, updates map[string]interface{}) error {
	return uc.repo.UpdateProduct(ctx, id, updates)
}
