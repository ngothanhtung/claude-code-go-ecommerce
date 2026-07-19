package user

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
	paging "github.com/ngothanhtung/go-tutorials/internal/common/paging"
)

type Service interface {
	Create(ctx context.Context, req CreateRequest) (PublicUser, error)
	GetByID(ctx context.Context, id uuid.UUID) (PublicUser, error)
	GetByEmail(ctx context.Context, email string) (PublicUser, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (PublicUser, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, q paging.Query) (ListResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service { return &service{repo: repo} }

func (s *service) Create(ctx context.Context, req CreateRequest) (PublicUser, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return PublicUser{}, apperr.NewInternal("hash password", err)
	}
	u := &User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hash),
	}
	if req.RoleID != "" {
		if rid, err := uuid.Parse(req.RoleID); err == nil {
			u.RoleID = &rid
		}
	}
	if err := s.repo.Create(ctx, u); err != nil {
		return PublicUser{}, err
	}
	return u.ToPublic(), nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (PublicUser, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return PublicUser{}, err
	}
	return u.ToPublic(), nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (PublicUser, error) {
	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return PublicUser{}, err
	}
	return u.ToPublic(), nil
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req UpdateRequest) (PublicUser, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return PublicUser{}, err
	}
	if req.Name != "" {
		u.Name = req.Name
	}
	if req.Email != "" {
		u.Email = req.Email
	}
	if err := s.repo.Update(ctx, u); err != nil {
		return PublicUser{}, err
	}
	return u.ToPublic(), nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *service) List(ctx context.Context, q paging.Query) (ListResponse, error) {
	users, total, err := s.repo.List(ctx, q)
	if err != nil {
		return ListResponse{}, err
	}
	out := make([]PublicUser, 0, len(users))
	for _, u := range users {
		out = append(out, u.ToPublic())
	}
	return ListResponse{Items: out, Total: total, Page: q.Page, PageSize: q.PageSize}, nil
}
