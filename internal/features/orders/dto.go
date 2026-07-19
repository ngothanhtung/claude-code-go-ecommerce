package orders

// CreateItemRequest is a single line in a CreateRequest.
type CreateItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1,max=999"`
}

// CreateRequest creates an order from the supplied items.
type CreateRequest struct {
	Items           []CreateItemRequest `json:"items" validate:"required,min=1,dive"`
	ShippingAddress string              `json:"shipping_address" validate:"required,min=5,max=500"`
	PaymentMethod   string              `json:"payment_method" validate:"required,oneof=cod demoCard"`
}

// StatusRequest updates an order's status (admin only).
type StatusRequest struct {
	Status string `json:"status" validate:"required,oneof=processing shipped delivered cancelled"`
}
