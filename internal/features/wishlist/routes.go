package wishlist

import "github.com/gin-gonic/gin"

// Register mounts wishlist routes on the given router group. All routes are
// protected by the auth middleware.
func Register(rg *gin.RouterGroup, h *Handler, auth gin.HandlerFunc) {
	g := rg.Group("/wishlist")
	g.Use(auth)
	g.GET("", h.Get)
	g.PUT("", h.Replace)
	g.POST("/:productId", h.Add)
	g.DELETE("/:productId", h.Remove)
}
