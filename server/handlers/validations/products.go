package validations

import (
	"ecom/server/types"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// ParseAndValidateGetProducts pulls and validates query params for the product list.
func ParseAndValidateGetProducts(q url.Values) (*types.GetProductsRequest, error) {

	req := &types.GetProductsRequest{}

	req.SortBy = q.Get("sortBy")
	req.Order = q.Get("order")

	if val := q.Get("priceMin"); val != "" {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid 'priceMin' value: must be a number")
		}
		req.PriceMin = &f
	}

	if val := q.Get("priceMax"); val != "" {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid 'priceMax' value: must be a number")
		}
		req.PriceMax = &f
	}

	if val := q.Get("search"); val != "" {
		req.SearchString = &val
	}

	if val := q.Get("minScore"); val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("invalid 'minScore' value: must be an integer")
		}
		req.MinScore = &i
	}

	if val := q.Get("limit"); val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			return nil, fmt.Errorf("invalid 'limit' value: must be an integer")
		}
		req.PageNum = i
	}

	if q.Has("cursor") {
		// Assumes cursor is sent as a comma-separated string, e.g., ?cursor=value1,value2
		req.Cursor = strings.Split(q.Get("cursor"), ",")
	}

	if err := validate.Struct(req); err != nil {
		// TODO: format validation errors nicely for the client
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return req, nil
}
