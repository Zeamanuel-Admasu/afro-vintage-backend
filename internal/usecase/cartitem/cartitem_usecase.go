package cartitem

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/cartitem"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models"
	"github.com/google/uuid"
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

func (u *cartItemUsecase) AddCartItem(ctx context.Context, userID string, listingID string) error {
	// Fetch product details using the provided productID (listingID).
	prod, err := u.productRepo.GetProductByID(ctx, listingID)
	if err != nil {
		return fmt.Errorf("failed to fetch product: %w", err)
	}
	if prod == nil {
		return errors.New("product not found")
	}
	// Check that the product is available.
	if prod.Status != "available" {
		return fmt.Errorf("product %s is not available", listingID)
	}

	// Build a new CartItem using the product details.
	cartItem := &cartitem.CartItem{
		ID:        uuid.NewString(),
		UserID:    userID,
		ListingID: prod.ID, // Using product.ID as the listing id.
		Title:     prod.Title,
		Price:     prod.Price,
		ImageURL:  prod.ImageURL,
		Grade:     prod.Grade,
		CreatedAt: time.Now(),
	}

	return u.repo.CreateCartItem(ctx, cartItem)
}

// GetCartItems retrieves all cart items for the given user.
func (u *cartItemUsecase) GetCartItems(ctx context.Context, userID string) ([]*cartitem.CartItem, error) {
	return u.repo.GetCartItems(ctx, userID)
}

// RemoveCartItem deletes a specific item from the user's cart.
func (u *cartItemUsecase) RemoveCartItem(ctx context.Context, userID string, listingID string) error {
	return u.repo.DeleteCartItem(ctx, userID, listingID)
}

// CheckoutCart processes a full cart checkout.
func (u *cartItemUsecase) CheckoutCart(ctx context.Context, userID string) (*models.CheckoutResponse, error) {
	// Retrieve all cart items.
	items, err := u.repo.GetCartItems(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New("cart is empty")
	}

	var total float64
	var checkoutItems []models.CheckoutItemResponse

	// Validate each item.
	for _, item := range items {
		prod, err := u.productRepo.GetProductByID(ctx, item.ListingID)
		if err != nil || prod == nil {
			return nil, fmt.Errorf("product with ListingID %s not found", item.ListingID)
		}
		if prod.Status != "available" {
			return nil, fmt.Errorf("item %q is no longer available", prod.Title)
		}

		// Accumulate for checkout.
		total += prod.Price
		checkoutItems = append(checkoutItems, models.CheckoutItemResponse{
			ListingID: prod.ID,
			Title:     prod.Title,
			Price:     prod.Price,
			SellerID:  prod.ResellerID.Hex(),
			Status:    prod.Status,
		})
		// Mark product as sold (simulate update; ideally via product repository with a transaction).
		prod.Status = "sold"
		// (Update product in DB here if needed)
	}

	// Calculate fees.
	platformFee := total * 0.02
	netPayable := total - platformFee

	// Create PaymentRecord and Order here (simulate payment via Stripe).
	// (This code would call another repository to save PaymentRecord)

	// Clear the purchased items from the cart.
	if err := u.repo.ClearCart(ctx, userID); err != nil {
		return nil, err
	}

	// Launch a goroutine to simulate order delivery update after 3 minutes.
	go func() {
		time.Sleep(3 * time.Minute)
		// Update order status to "Delivered".
		// (Call order repository or similar here)
	}()

	return &models.CheckoutResponse{
		TotalAmount: total,
		Items:       checkoutItems,
		PlatformFee: platformFee,
		NetPayable:  netPayable,
	}, nil
}

// CheckoutSingleItem processes checkout for a single cart item.
func (u *cartItemUsecase) CheckoutSingleItem(ctx context.Context, userID, listingID string) (*models.CheckoutResponse, error) {
	// Fetch the specific cart item.
	// (Option 1: Filter from GetCartItems; Option 2: Add a method to repo to get single item)
	items, err := u.repo.GetCartItems(ctx, userID)
	if err != nil {
		return nil, err
	}
	var targetItem *cartitem.CartItem
	for _, i := range items {
		if i.ListingID == listingID {
			targetItem = i
			break
		}
	}
	if targetItem == nil {
		return nil, errors.New("item not found in cart")
	}

	prod, err := u.productRepo.GetProductByID(ctx, listingID)
	if err != nil || prod == nil {
		return nil, fmt.Errorf("product with ListingID %s not found", listingID)
	}
	if prod.Status != "available" {
		return nil, fmt.Errorf("item %q is no longer available", prod.Title)
	}

	total := prod.Price
	platformFee := total * 0.02
	netPayable := total - platformFee

	checkoutItem := models.CheckoutItemResponse{
		ListingID: prod.ID,
		Title:     prod.Title,
		Price:     prod.Price,
		SellerID:  prod.ResellerID.Hex(),
		Status:    prod.Status,
	}

	// Mark product as sold and update db if required.
	prod.Status = "sold"

	// Simulate payment record creation & update.
	// Remove the single item from cart.
	if err := u.repo.DeleteCartItem(ctx, userID, listingID); err != nil {
		return nil, err
	}

	go func() {
		time.Sleep(3 * time.Minute)
		// Update order status to "Delivered".
	}()

	return &models.CheckoutResponse{
		TotalAmount: total,
		Items:       []models.CheckoutItemResponse{checkoutItem},
		PlatformFee: platformFee,
		NetPayable:  netPayable,
	}, nil
}
