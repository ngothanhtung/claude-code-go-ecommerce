package wishlist

// ReplaceRequest replaces the caller's wishlist with the supplied product IDs.
type ReplaceRequest struct {
	ProductIDs []string `json:"product_ids" validate:"required"`
}

// Summary is the wishlist projection returned to clients.
type Summary struct {
	ProductIDs []string `json:"product_ids"`
	Total      int      `json:"total"`
}
