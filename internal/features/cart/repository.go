package cart

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
)

type Repository interface {
	List(ctx context.Context, userID uuid.UUID) ([]Item, error)
	Replace(ctx context.Context, userID uuid.UUID, req ReplaceRequest) ([]Item, error)
	Clear(ctx context.Context, userID uuid.UUID) error
}

type repo struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) List(ctx context.Context, userID uuid.UUID) ([]Item, error) {
	var items []Item
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&items).Error; err != nil {
		return nil, apperr.NewInternal("list cart items", err)
	}
	return items, nil
}

func (r *repo) Replace(ctx context.Context, userID uuid.UUID, req ReplaceRequest) ([]Item, error) {
	var items []Item
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&Item{}).Error; err != nil {
			return apperr.NewInternal("clear cart", err)
		}
		if len(req.Items) == 0 {
			return nil
		}
		rows := make([]Item, 0, len(req.Items))
		for _, it := range req.Items {
			rows = append(rows, Item{
				UserID:    userID,
				ProductID: it.ProductID,
				Quantity:  it.Quantity,
			})
		}
		if err := tx.Create(&rows).Error; err != nil {
			if strings.Contains(err.Error(), "violates foreign key") {
				return apperr.NewNotFound("product not found")
			}
			return apperr.NewInternal("insert cart items", err)
		}
		items = rows
		return nil
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repo) Clear(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&Item{}).Error; err != nil {
		return apperr.NewInternal("clear cart", err)
	}
	return nil
}
