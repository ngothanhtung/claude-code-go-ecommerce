package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
	"github.com/ngothanhtung/go-tutorials/internal/common/response"
	"github.com/ngothanhtung/go-tutorials/pkg/jwt"
)

// Auth validates the bearer token and stores claims in context.
func Auth(jm *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Error(c, apperr.NewUnauthorized("missing bearer token"))
			c.Abort()
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := jm.Parse(token)
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

// Claims returns the parsed JWT claims from context, or nil if missing.
func Claims(c *gin.Context) *jwt.Claims {
	v, ok := c.Get("claims")
	if !ok {
		return nil
	}
	cl, ok := v.(*jwt.Claims)
	if !ok {
		return nil
	}
	return cl
}

// RequireRole aborts with 403 if the claims lack the required role.
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cl := Claims(c)
		if cl == nil {
			response.Error(c, apperr.NewUnauthorized("unauthenticated"))
			c.Abort()
			return
		}
		for _, r := range cl.Roles {
			if r == role {
				c.Next()
				return
			}
		}
		response.Error(c, apperr.NewForbidden("insufficient role"))
		c.Abort()
	}
}
