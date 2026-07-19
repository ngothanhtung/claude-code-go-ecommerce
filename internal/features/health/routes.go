package health

import "github.com/gin-gonic/gin"

// Register mounts health endpoints. Accepts either an Engine or a RouterGroup.
func Register(rg gin.IRouter, h *Handler) {
	rg.GET("/healthz", h.Live)
	rg.GET("/readyz", h.Ready)
}
