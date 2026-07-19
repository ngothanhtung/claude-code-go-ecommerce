package cart

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

func (h *Handler) Get(c *gin.Context) {
	uid, aerr := userIDFrom(c)
	if aerr != nil {
		response.Error(c, aerr)
		return
	}
	res, err := h.svc.Get(c.Request.Context(), uid)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Replace(c *gin.Context) {
	uid, aerr := userIDFrom(c)
	if aerr != nil {
		response.Error(c, aerr)
		return
	}
	var req ReplaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, val.WrapBind(err))
		return
	}
	if err := val.Struct(req); err != nil {
		response.Error(c, err)
		return
	}
	res, err := h.svc.Replace(c.Request.Context(), uid, req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Clear(c *gin.Context) {
	uid, aerr := userIDFrom(c)
	if aerr != nil {
		response.Error(c, aerr)
		return
	}
	if err := h.svc.Clear(c.Request.Context(), uid); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}
