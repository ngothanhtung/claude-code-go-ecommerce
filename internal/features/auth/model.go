package auth

import "time"

// RefreshTokenRecord maps to the refresh_tokens table.
type RefreshTokenRecord struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;index" json:"user_id"`
	TokenHash string    `gorm:"size:255" json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
}

func (RefreshTokenRecord) TableName() string { return "refresh_tokens" }

// PasswordResetRecord maps to password_reset_tokens.
type PasswordResetRecord struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;index" json:"user_id"`
	TokenHash string    `gorm:"size:255" json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

func (PasswordResetRecord) TableName() string { return "password_reset_tokens" }
