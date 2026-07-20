package notifications

// UpdateReadRequest flips the is_read flag on a single notification.
type UpdateReadRequest struct {
	IsRead bool `json:"is_read"`
}

// ListResponse is the standard paginated projection.
type ListResponse struct {
	Items []Notification `json:"items"`
	Total int64          `json:"total"`
}
