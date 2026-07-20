package notifications

import (
	"time"

	"github.com/google/uuid"
)

// Notification is a row in the notifications table.
type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	Title     string    `gorm:"size:255;not null" json:"title"`
	Subtitle  string    `gorm:"type:text;not null;default:''" json:"subtitle"`
	Type      string    `gorm:"size:30;not null" json:"type"`
	IsRead    bool      `gorm:"not null;default:false" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

func (Notification) TableName() string { return "notifications" }
