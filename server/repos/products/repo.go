package products

import (
	"context"
	"ecom/server/types"
	"fmt"
	"math"
	"strings"

	"github.com/jackc/pgx/v5"
)

type ProductRepo struct {
	DB *pgx.Conn
}

func NewProductRepo(db *pgx.Conn) *ProductRepo {
	return &ProductRepo{DB: db}
}

type PaginationOptions struct {
	Cursor  []string // [sort_key, tie_breaker_id] for keyset pagination.
	PageNum int
}
type SortOptions struct {
	SortBy string // "price" or "created_at"
	Order  string // "asc" or "desc"
}
type FiltersOptions struct {
	PriceMin     *float64
	PriceMax     *float64
	SearchString *string
	MinScore     *int
}
type GetAllOptions struct {
	Filters    FiltersOptions
	Pagination PaginationOptions
	Sort       SortOptions
}

type GetAllResult struct {
	Products   []types.MiniProduct
	TotalPages int
	TotalCount int
}

// MapRequestToGetAllOptions converts the request to repo options, with defaults.
func MapRequestToGetAllOptions(req *types.GetProductsRequest) GetAllOptions {
	sortBy := "created_at"
	if req.SortBy != "" {
		sortBy = req.SortBy
	}
	order := "desc"
	if req.Order != "" {
		order = req.Order
	}

	return GetAllOptions{
		Filters: FiltersOptions{
			PriceMin:     req.PriceMin,
			PriceMax:     req.PriceMax,
			SearchString: req.SearchString,
			MinScore:     req.MinScore,
		},
		Pagination: PaginationOptions{
			PageNum: req.PageNum, // Repo applies a default if this is 0
			Cursor:  req.Cursor,
		},
		Sort: SortOptions{
			SortBy: sortBy,
			Order:  order,
		},
	}
}
func (repo *ProductRepo) Get(ctx context.Context, productID int64) (types.Product, error) {
	var p types.Product
	sql := `
		SELECT
			p.id,
			p.name,
			p.description,
			p.price,
			p.created_at,
			c.id AS category_id,
			c.name AS category_name,
			COALESCE(
				JSON_AGG(
					JSON_BUILD_OBJECT('url', pi.url, 'alt_text', pi.alt_text)
				) FILTER (WHERE pi.url IS NOT NULL),
				'[]'
				) AS images,
			COALESCE(AVG(r.score), 0) AS average_rating
		FROM products p
		LEFT JOIN product_images pi ON p.id = pi.product_id
		LEFT JOIN ratings r ON p.id = r.product_id
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
		GROUP BY p.id, c.id;

	`
	r := repo.DB.QueryRow(ctx, sql, productID)

	args := []any{&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.CategoryData.ID, &p.CategoryData.Name, &p.Images, &p.AvgRating}
	if err := r.Scan(args...); err != nil {
		return p, err
	}
	return p, nil
}

func (repo *ProductRepo) GetAll(ctx context.Context, options GetAllOptions) (GetAllResult, error) {
	res := GetAllResult{
		Products: make([]types.MiniProduct, 0),
	}

	var args, countArgs []any
	var mainWhereClauses, filterWhereClauses []string

	// These filter clauses apply to both the main query and the count query.
	if options.Filters.PriceMin != nil {
		args = append(args, *options.Filters.PriceMin)
		filterWhereClauses = append(filterWhereClauses, fmt.Sprintf("p.price >= $%d", len(args)))
	}
	if options.Filters.PriceMax != nil {
		args = append(args, *options.Filters.PriceMax)
		filterWhereClauses = append(filterWhereClauses, fmt.Sprintf("p.price <= $%d", len(args)))
	}
	if options.Filters.SearchString != nil {
		searchArg := "%" + *options.Filters.SearchString + "%"
		args = append(args, searchArg)
		filterWhereClauses = append(filterWhereClauses, fmt.Sprintf("(p.name ILIKE $%d OR p.description ILIKE $%d)", len(args), len(args)))
	}

	// The count query uses only the filter arguments.
	countArgs = append(countArgs, args...)
	countWhereSQL := ""
	if len(filterWhereClauses) > 0 {
		countWhereSQL = "WHERE " + strings.Join(filterWhereClauses, " AND ")
	}

	// The main query's WHERE clause starts with the filters and may have the cursor condition added.
	mainWhereClauses = filterWhereClauses

	if len(options.Pagination.Cursor) == 2 { // Cursor pagination is active.
		sortValue := options.Pagination.Cursor[0]
		idValue := options.Pagination.Cursor[1]

		sortByField := "p.created_at"
		if options.Sort.SortBy == "price" {
			sortByField = "p.price"
		}

		operator := ">"
		if strings.ToLower(options.Sort.Order) == "desc" {
			operator = "<"
		}

		// Append cursor args to the main query args, then create the clause.
		args = append(args, sortValue, idValue)
		mainWhereClauses = append(mainWhereClauses, fmt.Sprintf("(%s, p.id) %s ($%d, $%d)", sortByField, operator, len(args)-1, len(args)))
	}

	whereSQL := ""
	if len(mainWhereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(mainWhereClauses, " AND ")
	}

	havingSQL := ""
	countHavingSQL := ""
	if options.Filters.MinScore != nil {
		scoreArg := *options.Filters.MinScore

		args = append(args, scoreArg)
		havingSQL = fmt.Sprintf("HAVING COALESCE(AVG(r.score), 0) >= $%d", len(args))

		countArgs = append(countArgs, scoreArg)
		countHavingSQL = fmt.Sprintf("HAVING COALESCE(AVG(r.score), 0) >= $%d", len(countArgs))
	}

	sortBy := "p.created_at"
	if options.Sort.SortBy == "price" {
		sortBy = "p.price"
	}
	order := "DESC"
	if strings.ToLower(options.Sort.Order) == "asc" {
		order = "ASC"
	}
	orderSQL := fmt.Sprintf("ORDER BY %s %s, p.id %s", sortBy, order, order)

	limit := 20 // Default limit
	if options.Pagination.PageNum > 0 {
		limit = options.Pagination.PageNum
	}
	args = append(args, limit)
	limitSQL := fmt.Sprintf("LIMIT $%d", len(args))

	mainQuerySQL := fmt.Sprintf(`
		SELECT
			p.id,
			p.name,
			p.price,
			p.created_at,
			c.id AS category_id,
			c.name AS category_name,
			COALESCE(
				(SELECT json_build_object('url', pi.url, 'alt_text', pi.alt_text)
				 FROM product_images pi
				 WHERE pi.product_id = p.id
				 ORDER BY pi.product_id ASC
				 LIMIT 1),
				'{}'
			) AS image,
			COALESCE(AVG(r.score), 0) AS average_rating
		FROM products p
		LEFT JOIN product_images pi ON p.id = pi.product_id
		LEFT JOIN categories c ON c.id = p.category_id
		LEFT JOIN ratings r ON r.product_id = p.id
		%s
		GROUP BY p.id, c.id
		%s
		%s
		%s
	`, whereSQL, havingSQL, orderSQL, limitSQL)

	rows, err := repo.DB.Query(ctx, mainQuerySQL, args...)
	if err != nil {
		return res, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p types.MiniProduct
		var avgRating float64

		scanArgs := []any{&p.ID, &p.Name, &p.Price, &p.CreatedAt, &p.CategoryData.ID, &p.CategoryData.Name, &p.Image, &avgRating}
		if err := rows.Scan(scanArgs...); err != nil {
			return res, fmt.Errorf("failed to scan product row: %w", err)
		}

		p.AvgRating = float32(math.Round(avgRating*100) / 100)
		res.Products = append(res.Products, p)
	}
	if err := rows.Err(); err != nil {
		return res, fmt.Errorf("error iterating product rows: %w", err)
	}

	// The total count query respects filters but ignores pagination (cursor/limit).
	var countSQL string
	if options.Filters.MinScore != nil {
		countSQL = fmt.Sprintf(`
			SELECT COUNT(*) FROM (
				SELECT 1
				FROM products p
				LEFT JOIN ratings r ON r.product_id = p.id
				%s
				GROUP BY p.id
				%s
			) AS sub
		`, countWhereSQL, countHavingSQL)
	} else {
		countSQL = fmt.Sprintf("SELECT COUNT(p.id) FROM products p %s", countWhereSQL)
	}

	err = repo.DB.QueryRow(ctx, countSQL, countArgs...).Scan(&res.TotalCount)
	if err != nil {
		return res, fmt.Errorf("failed to count products: %w", err)
	}

	if limit > 0 {
		res.TotalPages = int(math.Ceil(float64(res.TotalCount) / float64(limit)))
	}

	return res, nil
}
