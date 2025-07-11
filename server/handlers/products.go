package handlers

import (
	"ecom/server/customErrors"
	"ecom/server/handlers/validations"
	repoProducts "ecom/server/repos/products"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v4"
)

func (h *Handlers) HandleGetProduct(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid product ID format")
		return
	}

	product, err := h.ProductService.Get(r.Context(), productID)

	if err != nil {
		if errors.Is(err, customErrors.NotFound) {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, product)
}

func (h *Handlers) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	req, err := validations.ParseAndValidateGetProducts(r.URL.Query())
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	options := repoProducts.MapRequestToGetAllOptions(req)

	products, err := h.ProductService.GetAll(r.Context(), options)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}

	writeJSON(w, http.StatusOK, products)
}
func (h *Handlers) HandleRateProduct(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotImplemented, "unimplemented")
}
