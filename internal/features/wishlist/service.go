package wishlist

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Get(ctx context.Context, userID uuid.UUID) (Summary, error)
	Replace(ctx context.Context, userID uuid.UUID, req ReplaceRequest) error
	Add(ctx context.Context, userID uuid.UUID, productID string) error
	Remove(ctx context.Context, userID uuid.UUID, productID string) error
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Get(ctx context.Context, userID uuid.UUID) (Summary, error) {
	items, err := s.repo.List(ctx, userID)
	if err != nil {
		return Summary{}, err
	}
	ids := make([]string, 0, len(items))
	for _, it := range items {
		ids = append(ids, it.ProductID)
	}
	return Summary{ProductIDs: ids, Total: len(ids)}, nil
}

func (s *service) Replace(ctx context.Context, userID uuid.UUID, req ReplaceRequest) error {
	return s.repo.Replace(ctx, userID, req.ProductIDs)
}

func (s *service) Add(ctx context.Context, userID uuid.UUID, productID string) error {
	return s.repo.Add(ctx, userID, productID)
}

func (s *service) Remove(ctx context.Context, userID uuid.UUID, productID string) error {
	return s.repo.Remove(ctx, userID, productID)
}
