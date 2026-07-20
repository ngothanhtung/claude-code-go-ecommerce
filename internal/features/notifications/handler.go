package notifications

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
	"github.com/ngothanhtung/go-tutorials/internal/common/middleware"
	"github.com/ngothanhtung/go-tutorials/internal/common/response"
	val "github.com/ngothanhtung/go-tutorials/pkg/validator"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

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

func (h *Handler) List(c *gin.Context) {
	uid, aerr := userIDFrom(c)
	if aerr != nil {
		response.Error(c, aerr)
		return
	}
	res, err := h.svc.ListForUser(c.Request.Context(), uid)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Patch(c *gin.Context) {
	uid, aerr := userIDFrom(c)
	if aerr != nil {
		response.Error(c, aerr)
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperr.NewValidation("invalid notification id"))
		return
	}
	var req UpdateReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, val.WrapBind(err))
		return
	}
	if err := val.Struct(req); err != nil {
		response.Error(c, err)
		return
	}
	if err := h.svc.UpdateRead(c.Request.Context(), id, uid, req.IsRead); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) MarkAllRead(c *gin.Context) {
	uid, aerr := userIDFrom(c)
	if aerr != nil {
		response.Error(c, aerr)
		return
	}
	if err := h.svc.MarkAllRead(c.Request.Context(), uid); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
