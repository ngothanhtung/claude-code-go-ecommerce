package rbac

import (
	"github.com/gin-gonic/gin"

	"github.com/ngothanhtung/go-tutorials/internal/common/middleware"
)

// AdminGuard is an alias to the shared role middleware for readability.
func AdminGuard() gin.HandlerFunc {
	return middleware.RequireRole("admin")
}
