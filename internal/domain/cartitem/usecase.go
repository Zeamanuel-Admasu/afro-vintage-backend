package cartitem

import "context"

type Usecase interface {
	// AddCartItem adds an item to the user's cart.
	// It should perform validations such as duplicate checks and item availability.
	AddCartItem(ctx context.Context, userID string, listingID string) error

	// GetCartItems retrieves all items from the user's cart.
	GetCartItems(ctx context.Context, userID string) ([]*CartItem, error)

	// RemoveCartItem deletes a specific item from the user's cart.
	RemoveCartItem(ctx context.Context, userID string, listingID string) error

	// CheckoutCart processes a checkout by validating item availability,
	// marking items as sold, creating an order, and then clearing the cart.
	CheckoutCart(ctx context.Context, userID string) error
}
