package repos

import (
	"ecom/server/types"
)

type IProductRepo interface {
	Get(productID int64) (types.Product, error)
}
