package user

import "github.com/ngothanhtung/go-tutorials/internal/common/paging"

type CreateRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Password string `json:"password" validate:"required,min=8,max=72"`
	RoleID   string `json:"role_id"`
}

type UpdateRequest struct {
	Name  string `json:"name" validate:"omitempty,min=2,max=100"`
	Email string `json:"email" validate:"omitempty,email"`
}

type UpdateMeRequest struct {
	Name     string `json:"name" validate:"omitempty,min=2,max=100"`
	PhotoURL string `json:"photo_url"`
}

type ListResponse struct {
	Items    []PublicUser `json:"items"`
	Total    int64        `json:"total"`
	Page     int          `json:"page"`
	PageSize int          `json:"page_size"`
}

// keep paging import non-unused even if not referenced
var _ = paging.Query{}
