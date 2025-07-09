package types

type User struct {
}

type Product struct {
	ID           int64   `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	CategoryData struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"category_data"`
	Images []struct {
		URL     string `json:"url"`
		AltText string `json:"alt_text"`
	} `json:"images"`
}
