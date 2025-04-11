package controllers

import (
	"net/http"
	"strconv"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductController struct {
	Usecase product.Usecase
}

func NewProductController(usecase product.Usecase) *ProductController {
	return &ProductController{Usecase: usecase}
}

func (h *ProductController) Create(c *gin.Context) {
	var p product.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	// Extract userID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Ensure userID is a string
	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Convert userID to primitive.ObjectID
	resellerID, err := primitive.ObjectIDFromHex(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID format"})
		return
	}

	// Set ResellerID
	p.ResellerID = resellerID

	// Generate a new ID for the product
	p.ID = p.GenerateID()

	// Add product to the repository
	if err := h.Usecase.AddProduct(c.Request.Context(), &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "product created"})
}

func (h *ProductController) GetByID(c *gin.Context) {
	id := c.Param("id")
	prod, err := h.Usecase.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, prod)
}

func (h *ProductController) ListAvailable(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := h.Usecase.ListAvailableProducts(c.Request.Context(), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load products", "details": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no products available", "products": []product.Product{}})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductController) ListByReseller(c *gin.Context) {
	resellerID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	products, err := h.Usecase.ListProductsByReseller(c.Request.Context(), resellerID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load reseller products", "details": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no products found for this reseller", "products": []product.Product{}})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductController) Update(c *gin.Context) {
	id := c.Param("id")
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid update payload"})
		return
	}
	if err := h.Usecase.UpdateProduct(c.Request.Context(), id, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product updated"})
}

func (h *ProductController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.Usecase.DeleteProduct(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product deleted"})
}
