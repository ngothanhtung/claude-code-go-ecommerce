package auth

import "github.com/gin-gonic/gin"

// Register mounts auth routes (base /api/v1).
// These routes are intentionally NOT protected by auth middleware.
func Register(rg *gin.RouterGroup, h *Handler) {
	g := rg.Group("/auth")
	g.POST("/register", h.Register)
	g.POST("/login", h.Login)
	g.POST("/refresh", h.Refresh)
	g.POST("/logout", h.Logout)
	g.POST("/forgot-password", h.Forgot)
	g.POST("/reset-password", h.Reset)
}
