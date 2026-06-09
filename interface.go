package review_sdk

import "context"

type Client interface {
	CreateReview(ctx context.Context, input NewReview) (*Review, error)
	GetCityReviews(ctx context.Context, cityID int) ([]Review, error)
}
