package notifications

import "github.com/gin-gonic/gin"

// Register mounts notification routes on the given router group. All routes
// are protected by the auth middleware.
func Register(rg *gin.RouterGroup, h *Handler, auth gin.HandlerFunc) {
	g := rg.Group("/notifications")
	g.Use(auth)
	g.GET("", h.List)
	g.PATCH("/:id", h.Patch)
	g.POST("/mark-all-read", h.MarkAllRead)
}
