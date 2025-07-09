package products

import (
	"database/sql"
	"ecom/server/customErrors"
	"ecom/server/repos"
	"ecom/server/types"
	"fmt"
	"strconv"
)

type ProductService struct {
	Repo repos.IProductRepo
}

func NewService(repo repos.IProductRepo) *ProductService {
	return &ProductService{Repo: repo}
}

func (svc *ProductService) Get(productID string) (types.Product, error) {
	v, err := strconv.ParseInt(productID, 10, 64)
	if err != nil {
		return types.Product{}, fmt.Errorf("expected productID to be an int64")
	}
	p, err := svc.Repo.Get(v)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return p, customErrors.NotFound
		}
		return types.Product{}, err
	}
	return p, nil
}
