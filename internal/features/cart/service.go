package cart

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	Get(ctx context.Context, userID uuid.UUID) (Summary, error)
	Replace(ctx context.Context, userID uuid.UUID, req ReplaceRequest) (Summary, error)
	Clear(ctx context.Context, userID uuid.UUID) error
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Get(ctx context.Context, userID uuid.UUID) (Summary, error) {
	items, err := s.repo.List(ctx, userID)
	if err != nil {
		return Summary{}, err
	}
	if items == nil {
		items = []Item{}
	}
	total := 0
	for _, it := range items {
		total += it.Quantity
	}
	return Summary{Items: items, TotalQuantity: total}, nil
}

func (s *service) Replace(ctx context.Context, userID uuid.UUID, req ReplaceRequest) (Summary, error) {
	items, err := s.repo.Replace(ctx, userID, req)
	if err != nil {
		return Summary{}, err
	}
	if items == nil {
		items = []Item{}
	}
	total := 0
	for _, it := range items {
		total += it.Quantity
	}
	return Summary{Items: items, TotalQuantity: total}, nil
}

func (s *service) Clear(ctx context.Context, userID uuid.UUID) error {
	return s.repo.Clear(ctx, userID)
}
