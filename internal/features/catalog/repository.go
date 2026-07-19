package catalog

import (
	"context"
	"strings"

	"gorm.io/gorm"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
	paging "github.com/ngothanhtung/go-tutorials/internal/common/paging"
)

type Repository interface {
	GetCategories(ctx context.Context) ([]Category, error)
	GetProductByID(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, q paging.Query, categoryID, search string) ([]Product, int64, error)
	GetRelatedProducts(ctx context.Context, productID string, limit int) ([]Product, error)
	GetPromos(ctx context.Context) ([]Promo, error)
}

type repo struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) GetCategories(ctx context.Context) ([]Category, error) {
	var c []Category
	err := r.db.WithContext(ctx).Order("created_at ASC").Find(&c).Error
	if err != nil {
		return nil, apperr.NewInternal("get categories", err)
	}
	return c, nil
}

func (r *repo) GetProductByID(ctx context.Context, id string) (*Product, error) {
	var p Product
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return nil, apperr.NewNotFound("product not found")
		}
		return nil, apperr.NewInternal("get product", err)
	}
	return &p, nil
}

func (r *repo) ListProducts(ctx context.Context, q paging.Query, categoryID, search string) ([]Product, int64, error) {
	var products []Product
	var total int64
	q2 := r.db.WithContext(ctx).Model(&Product{})
	if categoryID != "" {
		q2 = q2.Where("category_id = ?", categoryID)
	}
	if search != "" {
		q2 = q2.Where("name ILIKE ?", "%"+search+"%")
	}
	if err := q2.Count(&total).Error; err != nil {
		return nil, 0, apperr.NewInternal("count products", err)
	}
	offset := (q.Page - 1) * q.PageSize
	if err := q2.Order(q.OrderBy()).Limit(q.PageSize).Offset(offset).Find(&products).Error; err != nil {
		return nil, 0, apperr.NewInternal("list products", err)
	}
	return products, total, nil
}

func (r *repo) GetRelatedProducts(ctx context.Context, productID string, limit int) ([]Product, error) {
	p, err := r.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}
	var products []Product
	err = r.db.WithContext(ctx).
		Where("category_id = ? AND id != ?", p.CategoryID, productID).
		Limit(limit).
		Find(&products).Error
	if err != nil {
		return nil, apperr.NewInternal("get related products", err)
	}
	return products, nil
}

func (r *repo) GetPromos(ctx context.Context) ([]Promo, error) {
	var promos []Promo
	err := r.db.WithContext(ctx).Find(&promos).Error
	if err != nil {
		return nil, apperr.NewInternal("get promos", err)
	}
	return promos, nil
}
