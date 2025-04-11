package controllers

import (
    "net/http"
    "time"
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundle"
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/models"
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/models/common"
    "github.com/gin-gonic/gin"
)

type BundleController struct {
    bundleUsecase bundle.Usecase
}

func NewBundleController(bundleUsecase bundle.Usecase) *BundleController {
    return &BundleController{
        bundleUsecase: bundleUsecase,
    }
}

func (c *BundleController) CreateBundle(ctx *gin.Context) {
    // Extract Supplier ID from JWT
    supplierID, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, common.APIResponse{
            Success: false,
            Message: "user ID not found in context",
        })
        return
    }

    // Validate supplierID is a non-empty string
    supplierIDStr, ok := supplierID.(string)
    if !ok || supplierIDStr == "" {
        ctx.JSON(http.StatusUnauthorized, common.APIResponse{
            Success: false,
            Message: "invalid or empty user ID in context",
        })
        return
    }

    // Validate role
    role, exists := ctx.Get("role")
    if !exists || role != "supplier" {
        ctx.JSON(http.StatusForbidden, common.APIResponse{
            Success: false,
            Message: "only suppliers can create bundles",
        })
        return
    }

    // Parse request body
    var req models.CreateBundleRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, common.APIResponse{
            Success: false,
            Message: "invalid request: " + err.Error(),
        })
        return
    }

    // Map DTO to domain entity
    b := &bundle.Bundle{
        ID:                 "bundle_" + supplierIDStr + "_" + time.Now().String(),
        SupplierID:         supplierIDStr,
        Title:              req.Title,
        Description:        req.Description,
        SampleImage:        req.SampleImage,
        Quantity:           req.NumberOfItems,
        Grade:              req.Grade,
        SortingLevel:       bundle.SortingLevel(req.Type),
        EstimatedBreakdown: req.EstimatedBreakdown,
        Type:               req.ClothingTypes[0], // Assuming single type for now
        Price:              req.Price,
        Status:             "available",
        CreatedAt:          time.Now().Format(time.RFC3339),
    }

    if err := c.bundleUsecase.CreateBundle(ctx, supplierIDStr, b); err != nil {
        ctx.JSON(http.StatusBadRequest, common.APIResponse{
            Success: false,
            Message: err.Error(),
        })
        return
    }

    // Map domain entity to response DTO
    resp := models.BundleResponse{
        ID:     b.ID,
        Title:  b.Title,
        Grade:  b.Grade,
        Price:  b.Price,
        Type:   string(b.SortingLevel),
        Status: b.Status,
    }

    ctx.JSON(http.StatusOK, common.APIResponse{
        Success: true,
        Message: "Bundle successfully created and listed!",
        Data:    resp,
    })
}

func (c *BundleController) ListBundles(ctx *gin.Context) {
    // Extract Supplier ID from JWT
    supplierID, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, common.APIResponse{
            Success: false,
            Message: "user ID not found in context",
        })
        return
    }

    // Validate supplierID is a non-empty string
    supplierIDStr, ok := supplierID.(string)
    if !ok || supplierIDStr == "" {
        ctx.JSON(http.StatusUnauthorized, common.APIResponse{
            Success: false,
            Message: "invalid or empty user ID in context",
        })
        return
    }

    // Validate role
    role, exists := ctx.Get("role")
    if !exists || role != "supplier" {
        ctx.JSON(http.StatusForbidden, common.APIResponse{
            Success: false,
            Message: "only suppliers can list bundles",
        })
        return
    }

    // Fetch bundles for the supplier
    bundles, err := c.bundleUsecase.ListBundles(ctx, supplierIDStr)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, common.APIResponse{
            Success: false,
            Message: "failed to retrieve bundles: " + err.Error(),
        })
        return
    }

    // Map domain entities to response DTOs
    var resp []models.BundleResponse
    for _, b := range bundles {
        bundleResp := models.BundleResponse{
            ID:     b.ID,
            Title:  b.Title,
            Grade:  b.Grade,
            Price:  b.Price,
            Type:   string(b.SortingLevel),
            Status: b.Status,
        }
        resp = append(resp, bundleResp)
    }

    ctx.JSON(http.StatusOK, common.APIResponse{
        Success: true,
        Message: "Bundles retrieved successfully",
        Data:    resp,
    })
}

func (c *BundleController) DeleteBundle(ctx *gin.Context) {
    // Extract Supplier ID from JWT
    supplierID, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, common.APIResponse{
            Success: false,
            Message: "user ID not found in context",
        })
        return
    }

    // Validate supplierID is a non-empty string
    supplierIDStr, ok := supplierID.(string)
    if !ok || supplierIDStr == "" {
        ctx.JSON(http.StatusUnauthorized, common.APIResponse{
            Success: false,
            Message: "invalid or empty user ID in context",
        })
        return
    }

    // Validate role
    role, exists := ctx.Get("role")
    if !exists || role != "supplier" {
        ctx.JSON(http.StatusForbidden, common.APIResponse{
            Success: false,
            Message: "only suppliers can delete bundles",
        })
        return
    }

    // Extract bundle ID from URL parameter
    bundleID := ctx.Param("id")
    if bundleID == "" {
        ctx.JSON(http.StatusBadRequest, common.APIResponse{
            Success: false,
            Message: "bundle ID is required",
        })
        return
    }

    // Delete the bundle
    err := c.bundleUsecase.DeleteBundle(ctx, supplierIDStr, bundleID)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, common.APIResponse{
            Success: false,
            Message: err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, common.APIResponse{
        Success: true,
        Message: "Bundle successfully deactivated",
        Data:    nil,
    })
}