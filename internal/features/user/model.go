package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	Email        string     `gorm:"uniqueIndex;size:255" json:"email"`
	Name         string     `gorm:"size:255" json:"name"`
	PasswordHash string     `gorm:"size:255" json:"-"`
	PhotoURL     *string    `gorm:"type:text" json:"photo_url,omitempty"`
	RoleID       *uuid.UUID `gorm:"type:uuid" json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `gorm:"index" json:"-"`
}

func (User) TableName() string { return "users" }

// PublicUser is the safe projection returned to clients.
type PublicUser struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	PhotoURL  *string   `json:"photo_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func (u User) ToPublic() PublicUser {
	return PublicUser{ID: u.ID, Email: u.Email, Name: u.Name, PhotoURL: u.PhotoURL, CreatedAt: u.CreatedAt}
}
