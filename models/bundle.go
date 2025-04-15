package models

type CreateBundleRequest struct {
	Title              string         `json:"title" binding:"required"`
	SampleImage        string         `json:"sample_image"`
	NumberOfItems      int            `json:"number_of_items" binding:"required"`
	Grade              string         `json:"grade" binding:"required"`
	Price              float64        `json:"price" binding:"required"`
	Description        string         `json:"description"`
	SizeRange          string         `json:"size_range"`
	ClothingTypes      []string       `json:"clothing_types"`
	Type               string         `json:"type" binding:"required"`
	EstimatedBreakdown map[string]int `json:"estimated_breakdown"`
	DeclaredRating     int            `json:"declared_rating" binding:"required"`
}

type BundleResponse struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Grade  string  `json:"grade"`
	Price  float64 `json:"price"`
	Type   string  `json:"type"`
	Status string  `json:"status"`
}
