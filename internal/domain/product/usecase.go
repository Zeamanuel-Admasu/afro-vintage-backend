package product

import "context"

type Usecase interface {
	AddProduct(ctx context.Context, p *Product) error
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProductsByReseller(ctx context.Context, resellerID string) ([]*Product, error)
	ListAvailableProducts(ctx context.Context) ([]*Product, error)
	DeleteProduct(ctx context.Context, id string) error
	UpdateProduct(ctx context.Context, id string, updates map[string]interface{}) error
}
