package user

import "github.com/gin-gonic/gin"

// Register mounts user routes on the given router group (base /api/v1).
// All routes are protected by the auth middleware.
func Register(rg *gin.RouterGroup, h *Handler, auth gin.HandlerFunc) {
	g := rg.Group("/users")
	g.Use(auth)
	g.GET("/me", h.Me)
	g.PUT("/me", h.MeUpdate)
	g.GET("", h.List)
	g.GET("/:id", h.Get)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
	g.POST("", h.Create)
}
