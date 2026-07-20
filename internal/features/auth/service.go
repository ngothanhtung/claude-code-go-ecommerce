package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
	"github.com/ngothanhtung/go-tutorials/internal/config"
	usermodel "github.com/ngothanhtung/go-tutorials/internal/features/user"
	"github.com/ngothanhtung/go-tutorials/pkg/jwt"
)

type Service interface {
	Register(ctx context.Context, req RegisterRequest) (AuthResult, error)
	Login(ctx context.Context, req LoginRequest) (TokenPair, error)
	Refresh(ctx context.Context, req RefreshRequest) (TokenPair, error)
	Logout(ctx context.Context, req LogoutRequest) error
	ForgotPassword(ctx context.Context, req ForgotRequest) (string, error)
	ResetPassword(ctx context.Context, req ResetRequest) error
}

type service struct {
	repo          Repository
	jm            *jwt.Manager
	rdb           *redis.Client
	cfg           config.JWTConfig
	accessSeconds int
}

func NewService(repo Repository, jm *jwt.Manager, rdb *redis.Client, cfg config.JWTConfig) Service {
	return &service{
		repo:          repo,
		jm:            jm,
		rdb:           rdb,
		cfg:           cfg,
		accessSeconds: cfg.AccessTTLMin * 60,
	}
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (AuthResult, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResult{}, apperr.NewInternal("hash password", err)
	}
	u := &usermodel.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hash),
	}
	if req.RoleID != "" {
		if rid, e := uuid.Parse(req.RoleID); e == nil {
			u.RoleID = &rid
		}
	}
	if err := s.repo.CreateUser(ctx, u); err != nil {
		return AuthResult{}, err
	}
	pair, err := s.issueTokens(ctx, u)
	if err != nil {
		return AuthResult{}, err
	}
	return AuthResult{User: toUserView(u, nil), TokenPair: pair}, nil
}

func (s *service) Login(ctx context.Context, req LoginRequest) (TokenPair, error) {
	u, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return TokenPair{}, apperr.NewInvalidCredentials()
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		return TokenPair{}, apperr.NewInvalidCredentials()
	}
	return s.issueTokens(ctx, u)
}

func (s *service) Refresh(ctx context.Context, req RefreshRequest) (TokenPair, error) {
	claims, err := s.jm.Parse(req.RefreshToken)
	if err != nil {
		return TokenPair{}, err
	}
	rec, err := s.repo.GetRefresh(ctx, claims.RegisteredClaims.ID)
	if err != nil {
		return TokenPair{}, err
	}
	if rec == nil || rec.Revoked {
		return TokenPair{}, apperr.NewTokenInvalid()
	}
	if time.Now().After(rec.ExpiresAt) {
		return TokenPair{}, apperr.NewTokenExpired()
	}
	u, err := s.repo.GetUserByID(ctx, uuid.MustParse(claims.UserID))
	if err != nil {
		return TokenPair{}, err
	}
	// Rotate: revoke old, issue new
	_ = s.repo.RevokeRefresh(ctx, rec.ID)
	_ = s.rdb
	return s.issueTokens(ctx, u)
}

func (s *service) Logout(ctx context.Context, req LogoutRequest) error {
	if req.RefreshToken == "" {
		return nil
	}
	claims, err := s.jm.Parse(req.RefreshToken)
	if err != nil {
		return nil // best-effort logout
	}
	return s.repo.RevokeRefresh(ctx, claims.RegisteredClaims.ID)
}

func (s *service) ForgotPassword(ctx context.Context, req ForgotRequest) (string, error) {
	u, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		// do not leak existence
		return "", nil
	}
	raw := randomToken(32)
	hash := sha256.Sum256([]byte(raw))
	rec := PasswordResetRecord{
		ID:        uuid.NewString(),
		UserID:    u.ID.String(),
		TokenHash: hex.EncodeToString(hash[:]),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
	if err := s.repo.SaveReset(ctx, rec); err != nil {
		return "", apperr.NewInternal("save reset", err)
	}
	return rec.ID + ":" + raw, nil
}

func (s *service) ResetPassword(ctx context.Context, req ResetRequest) error {
	raw := extractRawToken(req.Token)
	id := extractIDToken(req.Token)
	if id == "" || raw == "" {
		return apperr.NewTokenInvalid()
	}
	rec, err := s.repo.GetReset(ctx, id)
	if err != nil {
		return err
	}
	if rec == nil || rec.Used || time.Now().After(rec.ExpiresAt) {
		return apperr.NewTokenInvalid()
	}
	expected := sha256.Sum256([]byte(raw))
	if hex.EncodeToString(expected[:]) != rec.TokenHash {
		return apperr.NewTokenInvalid()
	}
	u, err := s.repo.GetUserByID(ctx, uuid.MustParse(rec.UserID))
	if err != nil {
		return err
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return apperr.NewInternal("hash password", err)
	}
	u.PasswordHash = string(newHash)
	if err := s.repo.UpdateUser(ctx, u); err != nil {
		return err
	}
	if err := s.repo.MarkResetUsed(ctx, id); err != nil {
		return apperr.NewInternal("mark reset used", err)
	}
	return s.repo.RevokeAllForUser(ctx, u.ID.String())
}

func (s *service) issueTokens(ctx context.Context, u *usermodel.User) (TokenPair, error) {
	roles := []string{"user"}
	if u.RoleID != nil {
		role, err := s.repo.GetRoleName(ctx, *u.RoleID)
		if err != nil {
			return TokenPair{}, err
		}
		if role != "" {
			roles = []string{role}
		}
	}
	refreshID := uuid.NewString()
	refreshClaims := jwt.Claims{
		UserID: u.ID.String(),
		Email:  u.Email,
		Roles:  roles,
	}
	access, _, err := s.jm.GenerateAccessToken(u.ID.String(), u.Email, roles)
	if err != nil {
		return TokenPair{}, apperr.NewInternal("sign access", err)
	}
	refresh, ttl, err := s.jm.GenerateRefreshTokenWithID(refreshID, refreshClaims)
	if err != nil {
		return TokenPair{}, apperr.NewInternal("sign refresh", err)
	}
	rec := RefreshTokenRecord{
		ID:        refreshID,
		UserID:    u.ID.String(),
		ExpiresAt: time.Now().Add(ttl),
	}
	if err := s.repo.SaveRefresh(ctx, rec); err != nil {
		return TokenPair{}, apperr.NewInternal("save refresh", err)
	}
	return TokenPair{AccessToken: access, RefreshToken: refresh, ExpiresIn: s.accessSeconds}, nil
}

func toUserView(u *usermodel.User, roles []string) UserView {
	if roles == nil {
		roles = []string{"user"}
	}
	return UserView{ID: u.ID.String(), Email: u.Email, Name: u.Name, Roles: roles}
}

func randomToken(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Token format: "id:raw" — split on first ':'
func extractRawToken(token string) string {
	for i := 0; i < len(token); i++ {
		if token[i] == ':' {
			return token[i+1:]
		}
	}
	return ""
}

func extractIDToken(token string) string {
	for i := 0; i < len(token); i++ {
		if token[i] == ':' {
			return token[:i]
		}
	}
	return ""
}
