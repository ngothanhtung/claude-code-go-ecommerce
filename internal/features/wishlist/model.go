package wishlist

import (
	"time"

	"github.com/google/uuid"
)

// Item is a single product saved to a user's wishlist (table wishlists).
type Item struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"-"`
	ProductID string    `gorm:"size:100;primaryKey" json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (Item) TableName() string { return "wishlists" }
