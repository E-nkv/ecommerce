package repos

import (
	"context"
	"ecom/server/repos/products"
	"ecom/server/types"
)

type IProductRepo interface {
	Get(ctx context.Context, productID int64) (types.Product, error)
	GetAll(ctx context.Context, options products.GetAllOptions) (products.GetAllResult, error)
}
