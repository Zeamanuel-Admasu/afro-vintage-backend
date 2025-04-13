package models

// CreateCartItemRequest is the request model for adding a new item to a cart.
type CreateCartItemRequest struct {
	ListingID string  `json:"listing_id" binding:"required"`
	Title     string  `json:"title" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
	ImageURL  string  `json:"image_url" binding:"required"`
	Grade     int     `json:"grade" binding:"required"`
}

// CartItemResponse is the response model for a cart item.
type CartItemResponse struct {
	ID        string  `json:"id"`
	ListingID string  `json:"listing_id"`
	Title     string  `json:"title"`
	Price     float64 `json:"price"`
	ImageURL  string  `json:"image_url"`
	Grade     int     `json:"grade"`
	CreatedAt string  `json:"created_at"`
}
