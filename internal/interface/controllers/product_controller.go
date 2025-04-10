package controllers

import (
	"net/http"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"github.com/gin-gonic/gin"
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
	products, err := h.Usecase.ListAvailableProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductController) ListByReseller(c *gin.Context) {
	resellerID := c.Param("id")
	products, err := h.Usecase.ListProductsByReseller(c.Request.Context(), resellerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
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
