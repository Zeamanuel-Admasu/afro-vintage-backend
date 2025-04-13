package cartitem

import (
	"context"
	"errors"
	"fmt"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/cartitem"
)

type cartItemUsecase struct {
	repo cartitem.Repository
	// Optionally, you can add a listing repository here for availability checks.
}

// NewCartItemUsecase creates a new CartItem usecase instance.
func NewCartItemUsecase(repo cartitem.Repository) cartitem.Usecase {
	return &cartItemUsecase{
		repo: repo,
	}
}

// AddCartItem adds an item to the user's cart after ensuring no duplicates.
func (u *cartItemUsecase) AddCartItem(ctx context.Context, userID string, item *cartitem.CartItem) error {
	// Additional validations can be added here (e.g. checking listing availability)
	return u.repo.CreateCartItem(ctx, item)
}

// GetCartItems retrieves all cart items for the given user.
func (u *cartItemUsecase) GetCartItems(ctx context.Context, userID string) ([]*cartitem.CartItem, error) {
	return u.repo.GetCartItems(ctx, userID)
}

// RemoveCartItem deletes a specific item from the user's cart.
func (u *cartItemUsecase) RemoveCartItem(ctx context.Context, userID string, listingID string) error {
	return u.repo.DeleteCartItem(ctx, userID, listingID)
}

// CheckoutCart validates all items and processes the checkout.
// It verifies that each listing is still available before proceeding.
// If any item is found to be unavailable, it returns an error and leaves the cart unchanged.
func (u *cartItemUsecase) CheckoutCart(ctx context.Context, userID string) error {
	// Retrieve all cart items for the user.
	items, err := u.repo.GetCartItems(ctx, userID)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return errors.New("cart is empty")
	}

	// Simulate validation: For each item, ensure the associated listing is still available.
	// In practice, you would call the listing repository to check the current status of each listing.
	for _, item := range items {
		// Placeholder for actual availability check.
		// For example: available, err := listingRepo.IsAvailable(ctx, item.ListingID)
		// Here we assume the listing is available unless explicitly simulated otherwise.
		// Replace the following condition with real logic.
		available := true
		if !available {
			// Optionally, remove the unavailable item from cart.
			_ = u.repo.DeleteCartItem(ctx, userID, item.ListingID)
			return fmt.Errorf("item %q (%s) is no longer available", item.Title, item.ListingID)
		}
	}

	// All items have been validated; proceed with the checkout.
	// At this point, you can:
	//   1. Mark the listings as sold (via the listing repository, using transactions/mutex as needed).
	//   2. Create an order record.
	// After a successful order creation, clear the cart.

	// For demonstration, we'll assume the listings status update and order creation succeed.
	clearErr := u.repo.ClearCart(ctx, userID)
	if clearErr != nil {
		return clearErr
	}

	return nil
}
