package controllers

import (
	"net/http"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/cartitem"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models"
	"github.com/gin-gonic/gin"
)

type CartItemController struct {
	usecase cartitem.Usecase
}

// NewCartItemController creates a new CartItemController.
func NewCartItemController(usecase cartitem.Usecase) *CartItemController {
	return &CartItemController{
		usecase: usecase,
	}
}

func (ctr *CartItemController) AddCartItem(c *gin.Context) {
	var req models.CreateCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real application, extract the userID from the authentication context.
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Call the usecase with only the product (listing) ID
	if err := ctr.usecase.AddCartItem(c.Request.Context(), userID, req.ListingID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "item added to cart"})
}

// GetCartItems handles GET /api/cart
func (ctr *CartItemController) GetCartItems(c *gin.Context) {
	// Extract userID from context (e.g., middleware)
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	items, err := ctr.usecase.GetCartItems(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert domain items to response models.
	var responses []models.CartItemResponse
	for _, item := range items {
		responses = append(responses, models.CartItemResponse{
			ID:        item.ID,
			ListingID: item.ListingID,
			Title:     item.Title,
			Price:     item.Price,
			ImageURL:  item.ImageURL,
			Grade:     item.Grade,
			CreatedAt: item.CreatedAt.Format(time.RFC3339),
		})
	}
	c.JSON(http.StatusOK, responses)
}

// RemoveCartItem handles DELETE /api/cart/items/:listingID
func (ctr *CartItemController) RemoveCartItem(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	listingID := c.Param("listingID")
	if listingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "listingID is required"})
		return
	}

	err := ctr.usecase.RemoveCartItem(c.Request.Context(), userID, listingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed from cart"})
}

// Add this new method for single item checkout.
func (ctr *CartItemController) CheckoutSingleItem(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	listingID := c.Param("listingId")
	if listingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "listingId parameter is required"})
		return
	}

	resp, err := ctr.usecase.CheckoutSingleItem(c.Request.Context(), userID, listingID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment successful. Order confirmed.",
		"data":    resp,
	})
}

// Update CheckoutCart (full cart checkout) method to return response.
func (ctr *CartItemController) CheckoutCart(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	resp, err := ctr.usecase.CheckoutCart(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Payment successful. Order confirmed.",
		"data":    resp,
	})
}
