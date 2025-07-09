package handlers

import (
	"ecom/server/services/products"
	"net/http"
)

type Handlers struct {
	ProductService *products.ProductService
}

func NewHandlers(productSvc *products.ProductService) *Handlers {
	return &Handlers{ProductService: productSvc}
}

func (h *Handlers) HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the E-commerce API!"))
}
