package orders

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	ListForUser(ctx context.Context, userID uuid.UUID) ([]Order, error)
	GetForUser(ctx context.Context, userID, id uuid.UUID) (*Order, error)
	Create(ctx context.Context, userID uuid.UUID, req CreateRequest) (*Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*Order, error)
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) ListForUser(ctx context.Context, userID uuid.UUID) ([]Order, error) {
	return s.repo.ListForUser(ctx, userID)
}

func (s *service) GetForUser(ctx context.Context, userID, id uuid.UUID) (*Order, error) {
	return s.repo.GetForUser(ctx, userID, id)
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, req CreateRequest) (*Order, error) {
	return s.repo.Create(ctx, userID, req)
}

func (s *service) UpdateStatus(ctx context.Context, id uuid.UUID, status string) (*Order, error) {
	return s.repo.UpdateStatus(ctx, id, status)
}
