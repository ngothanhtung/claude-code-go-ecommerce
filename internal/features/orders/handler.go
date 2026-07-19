package orders

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

func userUUIDFromClaims(c *gin.Context) (uuid.UUID, *apperr.AppError) {
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
	uid, aerr := userUUIDFromClaims(c)
	if aerr != nil {
		response.Error(c, aerr)
		return
	}
	res, err := h.svc.ListForUser(c.Request.Context(), uid)
	if err != nil {
		response.Error(c, err)
		return
	}
	if res == nil {
		res = []Order{}
	}
	response.OK(c, res)
}

func (h *Handler) Get(c *gin.Context) {
	uid, aerr := userUUIDFromClaims(c)
	if aerr != nil {
		response.Error(c, aerr)
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperr.NewValidation("invalid order id"))
		return
	}
	res, err := h.svc.GetForUser(c.Request.Context(), uid, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) Create(c *gin.Context) {
	uid, aerr := userUUIDFromClaims(c)
	if aerr != nil {
		response.Error(c, aerr)
		return
	}
	var req CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, val.WrapBind(err))
		return
	}
	if err := val.Struct(req); err != nil {
		response.Error(c, err)
		return
	}
	res, err := h.svc.Create(c.Request.Context(), uid, req)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Created(c, res)
}

func (h *Handler) UpdateStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, apperr.NewValidation("invalid order id"))
		return
	}
	var req StatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, val.WrapBind(err))
		return
	}
	if err := val.Struct(req); err != nil {
		response.Error(c, err)
		return
	}
	res, err := h.svc.UpdateStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}
