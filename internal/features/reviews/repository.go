package reviews

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
)

type Repository interface {
	// ListForProduct returns all reviews for a product ordered by created_at desc.
	ListForProduct(ctx context.Context, productID string) ([]Review, error)
	// GetForUserAndProduct returns the review for the given user+product.
	// Returns (nil, nil) when no review exists.
	GetForUserAndProduct(ctx context.Context, userID uuid.UUID, productID string) (*Review, error)
	// Upsert inserts a new review or updates rating/comment/updated_at of an
	// existing one atomically.
	Upsert(ctx context.Context, review *Review) error
}

type repo struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) ListForProduct(ctx context.Context, productID string) ([]Review, error) {
	var rows []Review
	if err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("created_at DESC").
		Find(&rows).Error; err != nil {
		return nil, apperr.NewInternal("list product reviews", err)
	}
	return rows, nil
}

func (r *repo) GetForUserAndProduct(ctx context.Context, userID uuid.UUID, productID string) (*Review, error) {
	var row Review
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND product_id = ?", userID, productID).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, apperr.NewInternal("get user review", err)
	}
	return &row, nil
}

func (r *repo) Upsert(ctx context.Context, review *Review) error {
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing Review
		findErr := tx.Where("user_id = ? AND product_id = ?", review.UserID, review.ProductID).
			First(&existing).Error
		if errors.Is(findErr, gorm.ErrRecordNotFound) {
			if cErr := tx.Create(review).Error; cErr != nil {
				return apperr.NewInternal("create review", cErr)
			}
			return nil
		}
		if findErr != nil {
			return apperr.NewInternal("lookup review", findErr)
		}

		now := time.Now()
		if uErr := tx.Model(&existing).
			Updates(map[string]interface{}{
				"rating":     review.Rating,
				"comment":    review.Comment,
				"updated_at": now,
			}).Error; uErr != nil {
			return apperr.NewInternal("update review", uErr)
		}
		review.UpdatedAt = &now
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
