package reviews

import (
	"time"

	"github.com/google/uuid"
)

// Review is a single (user, product) review row in product_reviews.
// The primary key is the composite (user_id, product_id), so a user may
// leave at most one review per product.
type Review struct {
	UserID    uuid.UUID  `gorm:"type:uuid;primaryKey" json:"-"`
	UserName  string     `gorm:"size:255;not null" json:"user_name"`
	UserEmail string     `gorm:"size:255;not null" json:"user_email"`
	ProductID string     `gorm:"size:100;primaryKey" json:"product_id"`
	OrderID   string     `gorm:"size:100;not null" json:"order_id"`
	Rating    int        `gorm:"not null" json:"rating"`
	Comment   string     `gorm:"type:text;not null;default:''" json:"comment"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (Review) TableName() string { return "product_reviews" }
