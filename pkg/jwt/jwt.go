package jwt

import (
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
)

type Claims struct {
	UserID string   `json:"uid"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	jwtv5.RegisteredClaims
}

type Manager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func New(secret string, accessTTLMin, refreshTTLHour int) *Manager {
	return &Manager{
		secret:     []byte(secret),
		accessTTL:  time.Duration(accessTTLMin) * time.Minute,
		refreshTTL: time.Duration(refreshTTLHour) * time.Hour,
	}
}

func (m *Manager) sign(claims Claims, ttl time.Duration) (string, time.Duration, error) {
	claims.RegisteredClaims = jwtv5.RegisteredClaims{
		ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(ttl)),
		IssuedAt:  jwtv5.NewNumericDate(time.Now()),
	}
	tok := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	s, err := tok.SignedString(m.secret)
	return s, ttl, err
}

func (m *Manager) GenerateAccessToken(userID, email string, roles []string) (string, time.Duration, error) {
	return m.sign(Claims{UserID: userID, Email: email, Roles: roles}, m.accessTTL)
}

func (m *Manager) GenerateRefreshToken(userID, email string, roles []string) (string, time.Duration, error) {
	return m.sign(Claims{UserID: userID, Email: email, Roles: roles}, m.refreshTTL)
}

// GenerateRefreshTokenWithID builds a refresh token whose RegisteredClaims.ID equals the supplied id.
// Used by the auth feature to persist a refresh token record keyed by id.
func (m *Manager) GenerateRefreshTokenWithID(id string, claims Claims) (string, time.Duration, error) {
	claims.RegisteredClaims = jwtv5.RegisteredClaims{
		ID:        id,
		ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(m.refreshTTL)),
		IssuedAt:  jwtv5.NewNumericDate(time.Now()),
	}
	tok := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)
	s, err := tok.SignedString(m.secret)
	return s, m.refreshTTL, err
}

func (m *Manager) Parse(tokenString string) (*Claims, error) {
	claims := &Claims{}
	tok, err := jwtv5.ParseWithClaims(tokenString, claims, func(t *jwtv5.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, apperr.NewTokenInvalid()
		}
		return m.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwtv5.ErrTokenExpired) {
			return nil, apperr.NewTokenExpired()
		}
		return nil, apperr.NewTokenInvalid()
	}
	if !tok.Valid {
		return nil, apperr.NewTokenInvalid()
	}
	return claims, nil
}
