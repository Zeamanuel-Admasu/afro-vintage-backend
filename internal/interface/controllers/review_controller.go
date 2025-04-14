package controllers

import (
	"net/http"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/review"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/models"
	"github.com/gin-gonic/gin"
)

type ReviewController struct {
	usecase review.Usecase
}

func NewReviewController(usecase review.Usecase) *ReviewController {
	return &ReviewController{usecase: usecase}
}

func (ctrl *ReviewController) SubmitReview(c *gin.Context) {
	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	r := &review.Review{
		OrderID:   req.OrderID,
		ProductID: req.ProductID,
		UserID:    userID,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	if err := ctrl.usecase.SubmitReview(c.Request.Context(), r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "review submitted"})
}
