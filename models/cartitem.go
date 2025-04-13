package models

// CreateCartItemRequest is the request model for adding a new item to a cart.
// Now it only accepts the product (listing) ID from the client.
type CreateCartItemRequest struct {
	ListingID string `json:"listing_id" binding:"required"`
}

// CartItemResponse remains the same.
type CartItemResponse struct {
	ID        string  `json:"id"`
	ListingID string  `json:"listing_id"`
	Title     string  `json:"title"`
	Price     float64 `json:"price"`
	ImageURL  string  `json:"image_url"`
	Grade     string  `json:"grade"`
	CreatedAt string  `json:"created_at"`
}
