package controllers

import (
    "net/http"
    "time"

    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundleorder"
    "github.com/gin-gonic/gin"
)

type SupplierController struct {
    orderUsecase  bundleorder.Usecase
    bundleUsecase bundle.Usecase
}

func NewSupplierController(orderUsecase bundleorder.Usecase, bundleUsecase bundle.Usecase) *SupplierController {
    return &SupplierController{
        orderUsecase:  orderUsecase,
        bundleUsecase: bundleUsecase,
    }
}

type SoldBundle struct {
    BundleID    string  `json:"bundle_id"`
    Title       string  `json:"title"`
    BuyerID     string  `json:"buyer_id"`
    OrderID     string  `json:"order_id"`
    Price       float64 `json:"price"`
    OrderDate   string  `json:"order_date"`
    OrderStatus string  `json:"order_status"`
}

func (ctrl *SupplierController) GetSoldBundles(c *gin.Context) {
    // Extract userID from context (set by AuthMiddleware)
    userID, exists := c.Get("userID")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    userIDStr, ok := userID.(string)
    if !ok {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }

    // Fetch orders for the supplier
    orders, err := ctrl.orderUsecase.GetOrdersBySellerID(c.Request.Context(), userIDStr)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch orders: " + err.Error()})
        return
    }

    if len(orders) == 0 {
        c.JSON(http.StatusOK, gin.H{
            "message": "No sold bundles found",
            "data":    []SoldBundle{},
        })
        return
    }

    // Fetch bundle details one by one using GetBundleByID
    var soldBundles []SoldBundle
    for _, order := range orders {
        // Fetch the bundle for this order
        bundle, err := ctrl.bundleUsecase.GetBundleByID(c.Request.Context(), userIDStr, order.BundleID)
        if err != nil {
            // If the bundle doesn't exist or isn't owned by the supplier, skip this order
            continue
        }

        // Add the sold bundle to the result
        soldBundles = append(soldBundles, SoldBundle{
            BundleID:    order.BundleID,
            Title:       bundle.Title,
            BuyerID:     order.BuyerID,
            OrderID:     order.ID,
            Price:       bundle.Price,
            OrderDate:   order.CreatedAt.Format(time.RFC3339),
            OrderStatus: order.WarehouseStatus,
        })
    }

    if len(soldBundles) == 0 {
        c.JSON(http.StatusOK, gin.H{
            "message": "No sold bundles found",
            "data":    []SoldBundle{},
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Sold bundles retrieved successfully",
        "data":    soldBundles,
    })
}