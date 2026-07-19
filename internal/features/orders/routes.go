package orders

import "github.com/gin-gonic/gin"

// Register mounts order routes on the given router group. All routes under
// /orders are protected by auth; PATCH /:id/status additionally requires the
// admin middleware supplied. The group is already authenticated.
func Register(rg *gin.RouterGroup, h *Handler, auth, admin gin.HandlerFunc) {
	g := rg.Group("/orders")
	g.Use(auth)
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.POST("", h.Create)
	g.PATCH("/:id/status", admin, h.UpdateStatus)
}
