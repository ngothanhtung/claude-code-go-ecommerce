package orders

import (
	"time"

	"github.com/google/uuid"
)

// Order is a header row in store_orders.
type Order struct {
	ID              uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID          uuid.UUID   `gorm:"type:uuid;index" json:"user_id"`
	Total           float64     `gorm:"type:decimal(12,2)" json:"total"`
	Status          string      `gorm:"size:30" json:"status"`
	ShippingAddress string      `gorm:"type:text" json:"shipping_address"`
	PaymentMethod   string      `gorm:"size:30" json:"payment_method"`
	Items           []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

func (Order) TableName() string { return "store_orders" }

// OrderItem is a snapshot row in order_items.
type OrderItem struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OrderID   uuid.UUID `gorm:"type:uuid" json:"order_id"`
	ProductID string    `gorm:"size:100" json:"product_id"`
	Name      string    `gorm:"size:255" json:"name"`
	UnitPrice float64   `gorm:"type:decimal(10,2)" json:"unit_price"`
	Quantity  int       `json:"quantity"`
}

func (OrderItem) TableName() string { return "order_items" }
