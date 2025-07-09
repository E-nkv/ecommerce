package handlers

import (
	"ecom/server/customErrors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v4"
)

func (h *Handlers) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "id")
	fmt.Println("str pr id is: ", productID)
	product, err := h.ProductService.Get(productID)
	if err != nil {
		switch err {
		case customErrors.NotFound:
			writeError(w, http.StatusBadRequest, err.Error())
		}
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, product)
}
func (h *Handlers) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "unimplemented")
}
func (h *Handlers) HandleRateProduct(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "unimplemented")
}
