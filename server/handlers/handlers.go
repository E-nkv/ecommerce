package handlers

import (
	"ecom/server/services"
	"net/http"
)

type Handlers struct {
	svc services.IService
}

func NewHandlers(svc services.IService) *Handlers {
	return &Handlers{svc: svc}
}

func (h *Handlers) HandleHome(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the E-commerce API!"))
}
