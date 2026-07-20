package catalog

type ProductListResponse struct {
	Items    []ProductPublic `json:"items"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}
