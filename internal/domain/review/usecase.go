package review

import "context"


type Usecase interface {
	SubmitReview(ctx context.Context, r *Review) error
}
