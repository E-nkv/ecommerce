package products

import (
	"context"
	"ecom/server/repos"
	"ecom/server/types"

	"github.com/jackc/pgx/v5"
)

type ProductRepo struct {
	repos.IProductRepo
	DB *pgx.Conn
}

func NewProductRepo(db *pgx.Conn) *ProductRepo {
	return &ProductRepo{DB: db}
}

func (repo *ProductRepo) Get(productID int64) (types.Product, error) {
	var p types.Product
	sql := `
		SELECT
			p.id,
			p.name,
			p.description,
			p.price,
			c.id AS category_id,
			c.name AS category_name,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT('url', pi.url, 'alt_text', pi.alt_text)
				) FILTER (WHERE pi.url IS NOT NULL),
				'[]'
				) AS images
		FROM products p
		LEFT JOIN product_images pi ON p.id = pi.product_id
		LEFT JOIN categories c ON c.id = p.category_id
		WHERE p.id = $1
		GROUP BY p.id, p.name, c.id, c.name;

	`
	r := repo.DB.QueryRow(context.Background(), sql, productID)

	args := []any{&p.ID, &p.Name, &p.Description, &p.Price, &p.CategoryData.ID, &p.CategoryData.Name, &p.Images}
	if err := r.Scan(args...); err != nil {
		return p, err
	}
	return p, nil
}
