package reviewusecase

import (
	"context"
	"errors"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/review"
	"github.com/google/uuid"
)

type reviewUsecase struct {
	reviewRepo review.Repository
	orderRepo  order.Repository
}

func NewReviewUsecase(reviewRepo review.Repository, orderRepo order.Repository) review.Usecase {
	return &reviewUsecase{
		reviewRepo: reviewRepo,
		orderRepo:  orderRepo,
	}
}

func (u *reviewUsecase) SubmitReview(ctx context.Context, r *review.Review) error {
	// Check if the order exists and is delivered
	order, err := u.orderRepo.GetOrderByID(ctx, r.OrderID)
	if err != nil || order == nil {
		return errors.New("order not found")
	}
	if order.Status != "Delivered" { 
		return errors.New("cannot review before delivery")
	}

	// Check if the user already reviewed this product
	existingReview, err := u.reviewRepo.GetReviewByUserAndProduct(ctx, r.UserID, r.ProductID)
	if err != nil {
		return err
	}
	if existingReview != nil {
		return errors.New("you already reviewed this item")
	}

	// Save the review
	r.ID = uuid.NewString()
	r.CreatedAt = time.Now().Format(time.RFC3339)
	return u.reviewRepo.CreateReview(ctx, r)
}
