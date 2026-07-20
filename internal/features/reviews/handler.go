package reviews

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
	"github.com/ngothanhtung/go-tutorials/internal/common/middleware"
	"github.com/ngothanhtung/go-tutorials/internal/common/response"
	val "github.com/ngothanhtung/go-tutorials/pkg/validator"
)

type Handler struct {
	svc Service
	// LookUpUserName resolves a user ID to their display name.
	// Set by router.go to avoid circular imports.
	// Falls back to "Member" if unset or lookup fails.
	LookUpUserName func(ctx context.Context, userID uuid.UUID) string
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func userIDFrom(c *gin.Context) (uuid.UUID, *apperr.AppError) {
	cl := middleware.Claims(c)
	if cl == nil {
		return uuid.Nil, apperr.NewUnauthorized("unauthenticated")
	}
	id, err := uuid.Parse(cl.UserID)
	if err != nil {
		return uuid.Nil, apperr.NewValidation("invalid user id")
	}
	return id, nil
}

// List is public: anyone may read the reviews for a product.
func (h *Handler) List(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		response.Error(c, apperr.NewValidation("invalid product id"))
		return
	}
	res, err := h.svc.ListForProduct(c.Request.Context(), productID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

// GetMine returns the caller's own review for a product (auth required).
func (h *Handler) GetMine(c *gin.Context) {
	uid, aerr := userIDFrom(c)
	if aerr != nil {
		response.Error(c, aerr)
		return
	}
	productID := c.Param("id")
	if productID == "" {
		response.Error(c, apperr.NewValidation("invalid product id"))
		return
	}
	review, err := h.svc.GetUserReview(c.Request.Context(), uid, productID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, review)
}

// UpsertMe creates or updates the caller's review for a product (auth required).
func (h *Handler) UpsertMe(c *gin.Context) {
	uid, aerr := userIDFrom(c)
	if aerr != nil {
		response.Error(c, aerr)
		return
	}
	productID := c.Param("id")
	if productID == "" {
		response.Error(c, apperr.NewValidation("invalid product id"))
		return
	}
	var req UpsertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, val.WrapBind(err))
		return
	}
	if err := val.Struct(req); err != nil {
		response.Error(c, err)
		return
	}

	cl := middleware.Claims(c)
	userName, userEmail := "Member", ""
	if cl != nil {
		userEmail = cl.Email
	}
	if h.LookUpUserName != nil {
		if resolved := h.LookUpUserName(c.Request.Context(), uid); resolved != "" {
			userName = resolved
		}
	}

	// Determine whether this is a create or an update to return the right status.
	existing, err := h.svc.GetUserReview(c.Request.Context(), uid, productID)
	if err != nil {
		response.Error(c, err)
		return
	}
	isNew := existing == nil

	review, err := h.svc.Upsert(c.Request.Context(), uid, userName, userEmail, productID, req)
	if err != nil {
		response.Error(c, err)
		return
	}
	if isNew {
		response.Created(c, review)
		return
	}
	response.OK(c, review)
}
