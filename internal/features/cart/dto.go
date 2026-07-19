package cart

// ItemRequest is a single line in a ReplaceRequest.
type ItemRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1,max=999"`
}

// ReplaceRequest replaces the caller's cart with the supplied items.
type ReplaceRequest struct {
	Items []ItemRequest `json:"items" validate:"required,dive"`
}

// Summary is the cart projection returned to clients.
type Summary struct {
	Items         []Item `json:"items"`
	TotalQuantity int    `json:"total_quantity"`
}
