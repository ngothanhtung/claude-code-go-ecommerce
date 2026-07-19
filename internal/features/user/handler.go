package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
	"github.com/ngothanhtung/go-tutorials/internal/common/middleware"
	"github.com/ngothanhtung/go-tutorials/internal/common/paging"
	"github.com/ngothanhtung/go-tutorials/internal/common/response"
	val "github.com/ngothanhtung/go-tutorials/pkg/validator"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Create(c *gin.Context) {
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, val.WrapBind(err))
		return
	}
	if err := val.Struct(req); err != nil {
		response.Error(c, err)
		return
	}
	pu, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, pu)
}

func (h *Handler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperr.NewValidation("invalid id"))
		return
	}
	pu, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, pu)
}

func (h *Handler) List(c *gin.Context) {
	q := paging.Parse(c)
	res, err := h.svc.List(c.Request.Context(), q)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperr.NewValidation("invalid id"))
		return
	}
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, val.WrapBind(err))
		return
	}
	if err := val.Struct(req); err != nil {
		response.Error(c, err)
		return
	}
	pu, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, pu)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperr.NewValidation("invalid id"))
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, nil)
}

func (h *Handler) Me(c *gin.Context) {
	cl := middleware.Claims(c)
	if cl == nil {
		response.Error(c, apperr.NewUnauthorized("unauthenticated"))
		return
	}
	id, err := uuid.Parse(cl.UserID)
	if err != nil {
		response.Error(c, apperr.NewValidation("invalid id"))
		return
	}
	pu, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, pu)
}
