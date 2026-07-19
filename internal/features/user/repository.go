package user

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
	paging "github.com/ngothanhtung/go-tutorials/internal/common/paging"
)

type Repository interface {
	Create(ctx context.Context, u *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, u *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, q paging.Query) ([]User, int64, error)
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) Create(ctx context.Context, u *User) error {
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		if isUniqueViolation(err) {
			return apperr.NewEmailExists()
		}
		return apperr.NewInternal("create user", err)
	}
	return nil
}

func (r *repo) GetByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.NewNotFound("user not found")
	}
	if err != nil {
		return nil, apperr.NewInternal("get user", err)
	}
	return &u, nil
}

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.NewNotFound("user not found")
	}
	if err != nil {
		return nil, apperr.NewInternal("get user", err)
	}
	return &u, nil
}

func (r *repo) Update(ctx context.Context, u *User) error {
	if err := r.db.WithContext(ctx).Model(u).Updates(map[string]interface{}{
		"name":  u.Name,
		"email": u.Email,
	}).Error; err != nil {
		return apperr.NewInternal("update user", err)
	}
	return nil
}

func (r *repo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Delete(&User{}, "id = ?", id).Error; err != nil {
		return apperr.NewInternal("delete user", err)
	}
	return nil
}

func (r *repo) List(ctx context.Context, q paging.Query) ([]User, int64, error) {
	var users []User
	var total int64
	if err := r.db.WithContext(ctx).Model(&User{}).Count(&total).Error; err != nil {
		return nil, 0, apperr.NewInternal("count users", err)
	}
	offset := (q.Page - 1) * q.PageSize
	if err := r.db.WithContext(ctx).
		Order(q.OrderBy()).
		Limit(q.PageSize).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, apperr.NewInternal("list users", err)
	}
	return users, total, nil
}

func isUniqueViolation(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "UNIQUE"))
}
