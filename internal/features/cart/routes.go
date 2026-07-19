package cart

import "github.com/gin-gonic/gin"

// Register mounts cart routes on the given router group. All routes are
// protected by the auth middleware.
func Register(rg *gin.RouterGroup, h *Handler, auth gin.HandlerFunc) {
	g := rg.Group("/cart")
	g.Use(auth)
	g.GET("", h.Get)
	g.PUT("", h.Replace)
	g.DELETE("", h.Clear)
}
