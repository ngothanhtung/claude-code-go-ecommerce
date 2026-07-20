package reviews

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	ListForProduct(ctx context.Context, productID string) (ListResponse, error)
	GetUserReview(ctx context.Context, userID uuid.UUID, productID string) (*Review, error)
	Upsert(ctx context.Context, userID uuid.UUID, userName, userEmail, productID string, req UpsertRequest) (*Review, error)
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) ListForProduct(ctx context.Context, productID string) (ListResponse, error) {
	items, err := s.repo.ListForProduct(ctx, productID)
	if err != nil {
		return ListResponse{}, err
	}
	if items == nil {
		items = []Review{}
	}
	return ListResponse{Items: items, Total: int64(len(items))}, nil
}

// GetUserReview returns (nil, nil) when the user has not reviewed this
// product yet; only DB errors surface here.
func (s *service) GetUserReview(ctx context.Context, userID uuid.UUID, productID string) (*Review, error) {
	return s.repo.GetForUserAndProduct(ctx, userID, productID)
}

func (s *service) Upsert(ctx context.Context, userID uuid.UUID, userName, userEmail, productID string, req UpsertRequest) (*Review, error) {
	review := &Review{
		UserID:    userID,
		UserName:  userName,
		UserEmail: userEmail,
		ProductID: productID,
		OrderID:   "",
		Rating:    req.Rating,
		Comment:   req.Comment,
	}
	if err := s.repo.Upsert(ctx, review); err != nil {
		return nil, err
	}
	return review, nil
}
