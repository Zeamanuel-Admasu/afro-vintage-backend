package cartitem

import (
	"context"
	"errors"
	"fmt"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/cartitem"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
)

// UnavailableItem represents a cart item that failed validation.

type cartItemUsecase struct {
	repo        cartitem.Repository
	productRepo product.Repository // Used to fetch product details
}

// NewCartItemUsecase creates a new CartItem usecase instance.
// Note: productRepo is used for product lookup and validation during checkout.
func NewCartItemUsecase(repo cartitem.Repository, productRepo product.Repository) cartitem.Usecase {
	return &cartItemUsecase{
		repo:        repo,
		productRepo: productRepo,
	}
}

// AddCartItem adds an item to the user's cart after ensuring no duplicates.
func (u *cartItemUsecase) AddCartItem(ctx context.Context, userID string, item *cartitem.CartItem) error {
	// Additional validations can be added here (e.g. checking listing availability)
	prod, err := u.productRepo.GetProductByID(ctx, item.ListingID)
	if err != nil {
		return fmt.Errorf("failed to fetch product: %w", err)
	}
	if prod == nil {
		return errors.New("product not found")
	}
	// Check that the product is available.
	if prod.Status != "available" {
		return fmt.Errorf("product %s is not available", item.ListingID)
	}
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

// CheckoutCart validates all items in the cart by fetching each associated product.
// It ensures that the product exists and that its status is "available".
// If any product is not available, it returns a CheckoutValidationError containing details.
// If all validations pass, it simulates a successful checkout by clearing the cart.
func (u *cartItemUsecase) CheckoutCart(ctx context.Context, userID string) error {
	// Retrieve all cart items for the user.
	items, err := u.repo.GetCartItems(ctx, userID)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return errors.New("cart is empty")
	}

	var unavailableItems []cartitem.UnavailableItem

	// Loop through each cart item and validate by fetching the corresponding product.
	for _, item := range items {
		// Treat product.ID as the listingId.
		prod, err := u.productRepo.GetProductByID(ctx, item.ListingID)
		if err != nil || prod == nil {
			return fmt.Errorf("product with listing ID %s not found", item.ListingID)
		}
		// Ensure product.Status is "available" (add this field in your product struct if it does not exist).
		if prod.Status != "available" {
			unavailableItems = append(unavailableItems, cartitem.UnavailableItem{
				ListingID: prod.ID,
				Title:     item.Title,
			})
		}
	}

	// If any items are unavailable, return a validation error.
	if len(unavailableItems) > 0 {
		return &cartitem.CheckoutValidationError{
			Message:          "Some items are no longer available",
			UnavailableItems: unavailableItems,
		}
	}

	// Simulate a successful checkout process (e.g. order creation, marking as sold, etc.).
	// For now, we simply clear the cart.
	if err := u.repo.ClearCart(ctx, userID); err != nil {
		return err
	}

	// At this point, checkout validation passed.
	// The controller can provide a success message like:
	// "Checkout validation passed. Proceed to payment!"
	return nil
}
