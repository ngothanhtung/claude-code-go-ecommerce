package catalog

import (
	"context"

	"github.com/ngothanhtung/go-tutorials/internal/common/paging"
)

type Service interface {
	GetCategories(ctx context.Context) ([]Category, error)
	GetProduct(ctx context.Context, id string) (ProductPublic, error)
	ListProducts(ctx context.Context, q paging.Query, categoryID, search string) (ProductListResponse, error)
	GetRelatedProducts(ctx context.Context, productID string, limit int) ([]ProductPublic, error)
	GetPromos(ctx context.Context) ([]Promo, error)
}

type service struct{ repo Repository }

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) GetCategories(ctx context.Context) ([]Category, error) {
	return s.repo.GetCategories(ctx)
}

func (s *service) GetProduct(ctx context.Context, id string) (ProductPublic, error) {
	p, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return ProductPublic{}, err
	}
	return p.ToPublic(), nil
}

func (s *service) ListProducts(ctx context.Context, q paging.Query, categoryID, search string) (ProductListResponse, error) {
	products, total, err := s.repo.ListProducts(ctx, q, categoryID, search)
	if err != nil {
		return ProductListResponse{}, err
	}
	publics := make([]ProductPublic, 0, len(products))
	for _, p := range products {
		publics = append(publics, p.ToPublic())
	}
	return ProductListResponse{Items: publics, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}

func (s *service) GetRelatedProducts(ctx context.Context, productID string, limit int) ([]ProductPublic, error) {
	if limit <= 0 {
		limit = 4
	}
	products, err := s.repo.GetRelatedProducts(ctx, productID, limit)
	if err != nil {
		return nil, err
	}
	publics := make([]ProductPublic, 0, len(products))
	for _, p := range products {
		publics = append(publics, p.ToPublic())
	}
	return publics, nil
}

func (s *service) GetPromos(ctx context.Context) ([]Promo, error) {
	return s.repo.GetPromos(ctx)
}
