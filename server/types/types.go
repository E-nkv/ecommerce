package types

import (
	"encoding/json"
	"time"
)

type User struct {
}

type MiniProduct struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	AvgRating    float32 `json:"average_rating"` // Average rating out of 5
	CategoryData struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"category_data"`
	Image     json.RawMessage `json:"image"` // A single image object {url:string, alt_text:string}
	CreatedAt time.Time       `json:"created_at"`
}

type Product struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	AvgRating    float64 `json:"average_rating"` // Average rating out of 5
	CategoryData struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"category_data"`
	Images    json.RawMessage `json:"images"` // JSON array of image objects {url:string, alt_text:string}
	CreatedAt time.Time       `json:"created_at"`
}

// GetProductsRequest defines query params for the product list endpoint.
// Pointers are used for optional fields.
type GetProductsRequest struct {
	SortBy       string   `validate:"omitempty,oneof=price created_at"`
	Order        string   `validate:"omitempty,oneof=asc desc"`
	PriceMin     *float64 `validate:"omitempty,gte=0"`
	PriceMax     *float64 `validate:"omitempty,gte=0,gtfield=PriceMin"`
	SearchString *string  `validate:"omitempty,min=1,max=100"`
	MinScore     *int     `validate:"omitempty,gte=1,lte=5"`
	PageNum      int      `validate:"omitempty,gte=1,lte=100"`
	Cursor       []string `validate:"omitempty,len=2"`
}
