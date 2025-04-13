package controllers

import (
    "net/http"

    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/bundleorder"
    "github.com/gin-gonic/gin"
)

type BundleOrderController struct {
    orderUsecase *bundleorder.BundleOrderUsecase
}

func NewBundleOrderController(orderUsecase *bundleorder.BundleOrderUsecase) *BundleOrderController {
    return &BundleOrderController{
        orderUsecase: orderUsecase,
    }
}

func (ctrl *BundleOrderController) CreateOrder(c *gin.Context) {
    // Extract userID from context (set by AuthMiddleware)
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

    // Parse request body
    var req struct {
        BundleID string `json:"bundle_id" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    // Create the order
    _, err := ctrl.orderUsecase.CreateOrder(c.Request.Context(), req.BundleID, userIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Bundle purchased and added to warehouse",
    })
}