package products

import (
	"context"
	"ecom/server/customErrors"
	"ecom/server/repos"
	repoProducts "ecom/server/repos/products"
	"ecom/server/types"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type ProductService struct {
	Repo repos.IProductRepo
}

func NewService(repo repos.IProductRepo) *ProductService {
	return &ProductService{Repo: repo}
}

func (svc *ProductService) Get(ctx context.Context, productID int64) (types.Product, error) {
	p, err := svc.Repo.Get(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return p, customErrors.NotFound
		}
		return types.Product{}, err
	}
	return p, nil
}

func (svc *ProductService) GetAll(ctx context.Context, options repoProducts.GetAllOptions) (repoProducts.GetAllResult, error) {
	res, err := svc.Repo.GetAll(ctx, options)
	if err != nil {
		// The service layer can add more context or logic here.
		return repoProducts.GetAllResult{}, fmt.Errorf("failed to get all products: %w", err)
	}
	return res, nil
}
