package catalog

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ngothanhtung/go-tutorials/internal/common/paging"
	"github.com/ngothanhtung/go-tutorials/internal/common/response"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) ListCategories(c *gin.Context) {
	cats, err := h.svc.GetCategories(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, cats)
}

func (h *Handler) ListProducts(c *gin.Context) {
	q := paging.Parse(c)
	catID := c.Query("category_id")
	search := c.Query("q")
	res, err := h.svc.ListProducts(c.Request.Context(), q, catID, search)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, res)
}

func (h *Handler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	p, err := h.svc.GetProduct(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, p)
}

func (h *Handler) GetRelated(c *gin.Context) {
	id := c.Param("id")
	limit := 4
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = n
		}
	}
	products, err := h.svc.GetRelatedProducts(c.Request.Context(), id, limit)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, products)
}

func (h *Handler) ListPromos(c *gin.Context) {
	promos, err := h.svc.GetPromos(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.OK(c, promos)
}
