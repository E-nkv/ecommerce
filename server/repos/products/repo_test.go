package products

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testRepo *ProductRepo

// TestMain sets up the db connection for all repo tests.
func TestMain(m *testing.M) {
	if err := godotenv.Load("../../../.env"); err != nil {
		log.Println("Could not load .env file, will rely on environment variables.")
	}

	dburl := os.Getenv("DB_URL")
	if dburl == "" {
		log.Fatal("DB_URL is not set. Please provide it via .env file or environment variable.")
	}

	db, err := pgx.Connect(context.Background(), dburl)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close(context.Background())

	testRepo = NewProductRepo(db)

	// Make sure the DB is seeded. We need our 100 test products.
	var count int
	err = db.QueryRow(context.Background(), "SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil || count < 100 {
		log.Fatalf("Database is not seeded correctly or is empty. Please run the seed migration. Expected 100 products, found %d.", count)
	}

	code := m.Run()

	os.Exit(code)
}

func TestProductRepo_GetAll(t *testing.T) {
	// Helpers to create pointers from literals for test options.
	float64Ptr := func(f float64) *float64 { return &f }
	stringPtr := func(s string) *string { return &s }
	intPtr := func(i int) *int { return &i }

	testCases := []struct {
		name               string
		options            GetAllOptions
		expectedCount      int
		expectedTotalCount int
		asserter           func(t *testing.T, res GetAllResult)
	}{
		{
			name:               "No options, default sort (created_at desc)",
			options:            GetAllOptions{},
			expectedCount:      20, // Default limit
			expectedTotalCount: 100,
			asserter: func(t *testing.T, res GetAllResult) {
				require.NotEmpty(t, res.Products)
				// Default sort is created_at DESC, id DESC. Since created_at is similar from seed, we check ID.
				assert.Greater(t, res.Products[0].ID, res.Products[1].ID)
				assert.Equal(t, int64(100), res.Products[0].ID, "Expected the last inserted product (ID 100) first")
			},
		},
		{
			name: "Sort by price ascending",
			options: GetAllOptions{
				Sort: SortOptions{SortBy: "price", Order: "asc"},
			},
			expectedCount:      20,
			expectedTotalCount: 100,
			asserter: func(t *testing.T, res GetAllResult) {
				require.GreaterOrEqual(t, len(res.Products), 2)
				assert.LessOrEqual(t, res.Products[0].Price, res.Products[1].Price)
				assert.Equal(t, "Cereal", res.Products[0].Name, "Expected cheapest product first") // From seed data
				assert.Equal(t, 1.99, res.Products[0].Price)
			},
		},
		{
			name: "Sort by price descending",
			options: GetAllOptions{
				Sort: SortOptions{SortBy: "price", Order: "desc"},
			},
			expectedCount:      20,
			expectedTotalCount: 100,
			asserter: func(t *testing.T, res GetAllResult) {
				require.GreaterOrEqual(t, len(res.Products), 2)
				assert.GreaterOrEqual(t, res.Products[0].Price, res.Products[1].Price)
				assert.Equal(t, "QuantumLeap X1 Laptop", res.Products[0].Name, "Expected most expensive product first") // From seed data
				assert.Equal(t, 1499.99, res.Products[0].Price)
			},
		},
		{
			name: "Sort by date ascending",
			options: GetAllOptions{
				Sort: SortOptions{SortBy: "created_at", Order: "asc"},
			},
			expectedCount:      20,
			expectedTotalCount: 100,
			asserter: func(t *testing.T, res GetAllResult) {
				require.NotEmpty(t, res.Products)
				// Sort is created_at ASC, id ASC.
				assert.Less(t, res.Products[0].ID, res.Products[1].ID)
				assert.Equal(t, int64(1), res.Products[0].ID, "Expected the first inserted product (ID 1) first")
			},
		},
		{
			name: "Filter by Price Range",
			options: GetAllOptions{
				Filters: FiltersOptions{PriceMin: float64Ptr(50.0), PriceMax: float64Ptr(60.0)},
				Sort:    SortOptions{SortBy: "price", Order: "asc"},
			},
			expectedCount:      7, // Manually counted from seed data
			expectedTotalCount: 7,
			asserter: func(t *testing.T, res GetAllResult) {
				for _, p := range res.Products {
					assert.GreaterOrEqual(t, p.Price, 50.0)
					assert.LessOrEqual(t, p.Price, 60.0)
				}
			},
		},
		{
			name: "Filter by Search String",
			options: GetAllOptions{
				Filters: FiltersOptions{SearchString: stringPtr("drone")},
			},
			expectedCount:      1,
			expectedTotalCount: 1,
			asserter: func(t *testing.T, res GetAllResult) {
				require.Len(t, res.Products, 1)
				assert.Equal(t, "Stealth Drone Pro", res.Products[0].Name)
			},
		},
		{
			name: "Filter by Min Rating (HAVING clause)",
			options: GetAllOptions{
				Filters: FiltersOptions{MinScore: intPtr(5)}, // Only products with a perfect average score
			},
			// Rating tests are non-deterministic due to random seed data. We just check that results match the criteria and the total count is correct.
			asserter: func(t *testing.T, res GetAllResult) {
				for _, p := range res.Products {
					assert.GreaterOrEqual(t, p.AvgRating, float32(5.0))
				}
				// Verify TotalCount is correct by running a separate query
				var countWithMinScore int
				err := testRepo.DB.QueryRow(context.Background(), `
					SELECT COUNT(*) FROM (
						SELECT 1 FROM products p
						LEFT JOIN ratings r ON r.product_id = p.id
						GROUP BY p.id
						HAVING COALESCE(AVG(r.score), 0) >= 5
					) as sub
				`).Scan(&countWithMinScore)
				require.NoError(t, err)
				assert.Equal(t, countWithMinScore, res.TotalCount, "TotalCount should match a direct query")
				assert.Len(t, res.Products, min(20, countWithMinScore))
			},
		},
		{
			name: "Pagination with custom limit",
			options: GetAllOptions{
				Pagination: PaginationOptions{PageNum: 7},
			},
			expectedCount:      7,
			expectedTotalCount: 100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := testRepo.GetAll(context.Background(), tc.options)
			require.NoError(t, err)

			// For non-deterministic tests (like random ratings), we don't assert a fixed count.
			if tc.expectedCount > 0 {
				assert.Len(t, res.Products, tc.expectedCount)
			}
			if tc.expectedTotalCount > 0 {
				assert.Equal(t, tc.expectedTotalCount, res.TotalCount)
			}

			if tc.asserter != nil {
				tc.asserter(t, res)
			}
		})
	}

	t.Run("Cursor Pagination", func(t *testing.T) {
		optsPage1 := GetAllOptions{
			Sort:       SortOptions{SortBy: "price", Order: "asc"},
			Pagination: PaginationOptions{PageNum: 5},
		}
		resPage1, err := testRepo.GetAll(context.Background(), optsPage1)
		require.NoError(t, err)
		require.Len(t, resPage1.Products, 5, "Page 1 should have 5 products")
		require.Equal(t, 100, resPage1.TotalCount, "Total count should be 100 regardless of pagination")

		lastItemPage1 := resPage1.Products[4]
		cursorPrice := strconv.FormatFloat(lastItemPage1.Price, 'f', -1, 64)
		cursorID := strconv.FormatInt(lastItemPage1.ID, 10)

		optsPage2 := GetAllOptions{
			Sort: SortOptions{SortBy: "price", Order: "asc"},
			Pagination: PaginationOptions{
				PageNum: 5,
				Cursor:  []string{cursorPrice, cursorID},
			},
		}
		resPage2, err := testRepo.GetAll(context.Background(), optsPage2)
		require.NoError(t, err)
		require.Len(t, resPage2.Products, 5, "Page 2 should have 5 products")

		firstItemPage2 := resPage2.Products[0]

		assert.GreaterOrEqual(t, firstItemPage2.Price, lastItemPage1.Price, "Price of page 2's first item should be >= page 1's last item")
		if firstItemPage2.Price == lastItemPage1.Price {
			assert.Greater(t, firstItemPage2.ID, lastItemPage1.ID, "If prices are equal, ID of page 2's first item must be > page 1's last item")
		}

		page1IDs := make(map[int64]bool)
		for _, p := range resPage1.Products {
			page1IDs[p.ID] = true
		}
		for _, p := range resPage2.Products {
			assert.False(t, page1IDs[p.ID], "Product with ID %d from page 1 should not be in page 2", p.ID)
		}
	})
}
