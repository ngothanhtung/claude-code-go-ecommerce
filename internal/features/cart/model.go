package cart

import (
	"time"

	"github.com/google/uuid"
)

// Item is a single line in a user's cart (table cart_items).
type Item struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey" json:"-"`
	ProductID string    `gorm:"size:100;primaryKey" json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

func (Item) TableName() string { return "cart_items" }
