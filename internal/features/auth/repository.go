package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
	usermodel "github.com/ngothanhtung/go-tutorials/internal/features/user"
)

type Repository interface {
	CreateUser(ctx context.Context, u *usermodel.User) error
	GetUserByEmail(ctx context.Context, email string) (*usermodel.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*usermodel.User, error)
	UpdateUser(ctx context.Context, u *usermodel.User) error
	SaveRefresh(ctx context.Context, r RefreshTokenRecord) error
	GetRefresh(ctx context.Context, id string) (*RefreshTokenRecord, error)
	RevokeRefresh(ctx context.Context, id string) error
	RevokeAllForUser(ctx context.Context, userID string) error
	SaveReset(ctx context.Context, r PasswordResetRecord) error
	GetReset(ctx context.Context, id string) (*PasswordResetRecord, error)
	MarkResetUsed(ctx context.Context, id string) error
}

type repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository { return &repo{db: db} }

func (r *repo) CreateUser(ctx context.Context, u *usermodel.User) error {
	if err := r.db.WithContext(ctx).Create(u).Error; err != nil {
		if isUniqueViolation(err) {
			return apperr.NewEmailExists()
		}
		return apperr.NewInternal("create user", err)
	}
	return nil
}

func (r *repo) GetUserByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	var u usermodel.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.NewNotFound("user not found")
	}
	if err != nil {
		return nil, apperr.NewInternal("get user", err)
	}
	return &u, nil
}

func (r *repo) GetUserByID(ctx context.Context, id uuid.UUID) (*usermodel.User, error) {
	var u usermodel.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.NewNotFound("user not found")
	}
	if err != nil {
		return nil, apperr.NewInternal("get user", err)
	}
	return &u, nil
}

func (r *repo) UpdateUser(ctx context.Context, u *usermodel.User) error {
	if err := r.db.WithContext(ctx).Model(u).Update("password_hash", u.PasswordHash).Error; err != nil {
		return apperr.NewInternal("update user", err)
	}
	return nil
}

func (r *repo) SaveRefresh(ctx context.Context, rec RefreshTokenRecord) error {
	return r.db.WithContext(ctx).Create(&rec).Error
}

func (r *repo) GetRefresh(ctx context.Context, id string) (*RefreshTokenRecord, error) {
	var rec RefreshTokenRecord
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&rec).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, apperr.NewInternal("get refresh", err)
	}
	return &rec, nil
}

func (r *repo) RevokeRefresh(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&RefreshTokenRecord{}).Where("id = ?", id).Update("revoked", true).Error
}

func (r *repo) RevokeAllForUser(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Model(&RefreshTokenRecord{}).Where("user_id = ?", userID).Update("revoked", true).Error
}

func (r *repo) SaveReset(ctx context.Context, rec PasswordResetRecord) error {
	return r.db.WithContext(ctx).Create(&rec).Error
}

func (r *repo) GetReset(ctx context.Context, id string) (*PasswordResetRecord, error) {
	var rec PasswordResetRecord
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&rec).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, apperr.NewInternal("get reset", err)
	}
	return &rec, nil
}

func (r *repo) MarkResetUsed(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&PasswordResetRecord{}).Where("id = ?", id).Update("used", true).Error
}

func isUniqueViolation(err error) bool {
	return err != nil && (contains(err.Error(), "duplicate key") || contains(err.Error(), "UNIQUE"))
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
