package wishlist

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
)

type Repository interface {
	List(ctx context.Context, userID uuid.UUID) ([]Item, error)
	Replace(ctx context.Context, userID uuid.UUID, ids []string) error
	Add(ctx context.Context, userID uuid.UUID, productID string) error
	Remove(ctx context.Context, userID uuid.UUID, productID string) error
}

type repo struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) List(ctx context.Context, userID uuid.UUID) ([]Item, error) {
	var items []Item
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&items).Error; err != nil {
		return nil, apperr.NewInternal("list wishlist", err)
	}
	return items, nil
}

func (r *repo) Replace(ctx context.Context, userID uuid.UUID, ids []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&Item{}).Error; err != nil {
			return apperr.NewInternal("clear wishlist", err)
		}
		if len(ids) == 0 {
			return nil
		}
		rows := make([]Item, 0, len(ids))
		for _, pid := range ids {
			rows = append(rows, Item{UserID: userID, ProductID: pid})
		}
		if err := tx.Create(&rows).Error; err != nil {
			if strings.Contains(err.Error(), "violates foreign key") {
				return apperr.NewNotFound("product not found")
			}
			return apperr.NewInternal("insert wishlist", err)
		}
		return nil
	})
}

func (r *repo) Add(ctx context.Context, userID uuid.UUID, productID string) error {
	row := &Item{UserID: userID, ProductID: productID}
	err := r.db.WithContext(ctx).Where("user_id = ? AND product_id = ?", userID, productID).FirstOrCreate(row).Error
	if err != nil {
		if strings.Contains(err.Error(), "violates foreign key") {
			return apperr.NewNotFound("product not found")
		}
		return apperr.NewInternal("add wishlist item", err)
	}
	return nil
}

func (r *repo) Remove(ctx context.Context, userID uuid.UUID, productID string) error {
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Delete(&Item{}).Error; err != nil {
		return apperr.NewInternal("remove wishlist item", err)
	}
	return nil
}
