package notifications

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type Service interface {
	ListForUser(ctx context.Context, userID uuid.UUID) (ListResponse, error)
	Create(ctx context.Context, n *Notification) error
	UpdateRead(ctx context.Context, id, userID uuid.UUID, isRead bool) error
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
	// CreateOrderNotification builds and persists an "order" notification for
	// the user. It is the callback target wired up by router.go into the
	// orders handler so that an order creation triggers an async
	// notification without orders importing notifications.
	CreateOrderNotification(ctx context.Context, userID uuid.UUID, orderID string, total float64) error
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) ListForUser(ctx context.Context, userID uuid.UUID) (ListResponse, error) {
	items, err := s.repo.ListForUser(ctx, userID)
	if err != nil {
		return ListResponse{}, err
	}
	if items == nil {
		items = []Notification{}
	}
	return ListResponse{Items: items, Total: int64(len(items))}, nil
}

func (s *service) Create(ctx context.Context, n *Notification) error {
	return s.repo.Create(ctx, n)
}

func (s *service) UpdateRead(ctx context.Context, id, userID uuid.UUID, isRead bool) error {
	return s.repo.UpdateRead(ctx, id, userID, isRead)
}

func (s *service) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	return s.repo.MarkAllRead(ctx, userID)
}

func (s *service) CreateOrderNotification(ctx context.Context, userID uuid.UUID, orderID string, total float64) error {
	n := &Notification{
		UserID:   userID,
		Title:    "Order placed",
		Subtitle: fmt.Sprintf("Your order %s for $%.2f was placed successfully.", orderID, total),
		Type:     "order",
		IsRead:   false,
	}
	return s.repo.Create(ctx, n)
}
