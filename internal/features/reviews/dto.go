package reviews

// UpsertRequest creates or updates the caller's review for a product.
type UpsertRequest struct {
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment" validate:"max=500"`
}

// ListResponse is the public projection of the reviews for a product.
type ListResponse struct {
	Items []Review `json:"items"`
	Total int64    `json:"total"`
}
