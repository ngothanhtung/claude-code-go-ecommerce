package health

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ngothanhtung/go-tutorials/internal/common/response"
)

type Handler struct {
	pingDB    func() error
	pingRedis func() error
}

func NewHandler(pingDB, pingRedis func() error) *Handler {
	return &Handler{pingDB: pingDB, pingRedis: pingRedis}
}

func (h *Handler) Live(c *gin.Context) {
	response.OK(c, gin.H{"status": "alive"})
}

func (h *Handler) Ready(c *gin.Context) {
	if err := h.pingDB(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "db down", "error": err.Error()})
		return
	}
	if err := h.pingRedis(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "redis down", "error": err.Error()})
		return
	}
	response.OK(c, gin.H{"status": "ready"})
}
