package notifications

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
)

type Repository interface {
	ListForUser(ctx context.Context, userID uuid.UUID) ([]Notification, error)
	Create(ctx context.Context, n *Notification) error
	UpdateRead(ctx context.Context, id, userID uuid.UUID, isRead bool) error
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
}

type repo struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) ListForUser(ctx context.Context, userID uuid.UUID) ([]Notification, error) {
	var rows []Notification
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&rows).Error; err != nil {
		return nil, apperr.NewInternal("list notifications", err)
	}
	return rows, nil
}

func (r *repo) Create(ctx context.Context, n *Notification) error {
	if err := r.db.WithContext(ctx).Create(n).Error; err != nil {
		return apperr.NewInternal("create notification", err)
	}
	return nil
}

func (r *repo) UpdateRead(ctx context.Context, id, userID uuid.UUID, isRead bool) error {
	var row Notification
	findErr := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		First(&row).Error
	if errors.Is(findErr, gorm.ErrRecordNotFound) {
		return apperr.NewNotFound("notification not found")
	}
	if findErr != nil {
		return apperr.NewInternal("lookup notification", findErr)
	}
	if err := r.db.WithContext(ctx).
		Model(&Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", isRead).Error; err != nil {
		return apperr.NewInternal("update notification", err)
	}
	return nil
}

func (r *repo) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Model(&Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error; err != nil {
		return apperr.NewInternal("mark notifications read", err)
	}
	return nil
}
