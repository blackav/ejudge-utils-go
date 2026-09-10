package ejclient

import (
	"context"
)

type Client interface {
	ListPendingReviews(ctx context.Context, req *ListPendingReviewsRequest) (int, *Reply[ListPendingReviewsResult], error)
	StartReview(ctx context.Context, req *StartReviewRequest) (int, *Reply[StartReviewResult], error)
	HeartbeatReview(ctx context.Context, req *HeartbeatReviewRequest) (int, *Reply[HeartbeatReviewResult], error)
	FinishReview(ctx context.Context, req *FinishReviewRequest) (int, *Reply[FinishReviewResult], error)
	ListActiveReviews(ctx context.Context, req *ListActiveReviewsRequest) (int, *Reply[ListActiveReviewsResult], error)
	GetActiveReview(ctx context.Context, req *GetActiveReviewRequest) (int, *Reply[GetActiveReviewResult], error)
}
