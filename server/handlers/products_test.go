package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	repoProducts "ecom/server/repos/products"
	productSvc "ecom/server/services/products"
	"ecom/server/types"

	"github.com/go-chi/chi/v4"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testServer *httptest.Server

// TestMain boots a test server for all E2E tests in this package.
func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
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

	repo := repoProducts.NewProductRepo(db)
	service := productSvc.NewService(repo)
	handler := NewHandlers(service)

	router := chi.NewRouter()
	router.Get("/products/{id}", handler.HandleGetProduct)
	router.Get("/products", handler.HandleGetProducts)

	testServer = httptest.NewServer(router)
	defer testServer.Close()

	code := m.Run()
	os.Exit(code)
}

// TestGetProductE2E tests the single product endpoint.
func TestGetProductE2E(t *testing.T) {
	t.Run("Success - Get existing product", func(t *testing.T) {

		resp, err := http.Get(testServer.URL + "/products/1")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var product types.Product
		err = json.NewDecoder(resp.Body).Decode(&product)
		require.NoError(t, err)

		assert.Equal(t, int64(1), product.ID)

	})

	t.Run("Failure - Product not found", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/products/9999")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Failure - Invalid product ID format", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/products/abc")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	})
}

// TestGetProductsE2E tests the product list endpoint with various queries.
func TestGetProductsE2E(t *testing.T) {
	t.Run("Success - No params, default response", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/products")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result repoProducts.GetAllResult
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, 100, result.TotalCount)
		assert.Len(t, result.Products, 20) // Default limit
		// Default sort is created_at desc, so newest (id 100) is first.
		assert.Equal(t, int64(100), result.Products[0].ID)
	})

	t.Run("Success - Sorting by price ascending", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/products?sortBy=price&order=asc")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result repoProducts.GetAllResult
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Len(t, result.Products, 20)
		// From seed data, cheapest product is "Cereal" at 1.99
		assert.Equal(t, "Cereal", result.Products[0].Name)
		assert.Equal(t, 1.99, result.Products[0].Price)
		assert.LessOrEqual(t, result.Products[0].Price, result.Products[1].Price)
	})

	t.Run("Success - Filtering by price range and search", func(t *testing.T) {
		url := fmt.Sprintf("%s/products?search=laptop&priceMin=1000&priceMax=2000", testServer.URL)
		resp, err := http.Get(url)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result repoProducts.GetAllResult
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, 1, result.TotalCount)
		require.Len(t, result.Products, 1)
		assert.Equal(t, "QuantumLeap X1 Laptop", result.Products[0].Name)
	})

	t.Run("Failure - Invalid query parameter value", func(t *testing.T) {
		resp, err := http.Get(testServer.URL + "/products?priceMin=abc")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Contains(t, string(body), "invalid 'priceMin' value")
	})

	t.Run("Failure - Validation rule failed", func(t *testing.T) {
		// sortBy must be one of 'price' or 'created_at'
		resp, err := http.Get(testServer.URL + "/products?sortBy=name")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Contains(t, string(body), "validation failed")
	})
}
