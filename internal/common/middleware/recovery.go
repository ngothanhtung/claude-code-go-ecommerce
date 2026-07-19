package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
	"github.com/ngothanhtung/go-tutorials/internal/common/response"
)

func Recovery(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic recovered", zap.Any("error", r))
				response.Error(c, apperr.NewInternal("panic recovered", nil))
				c.Abort()
			}
		}()
		c.Next()
	}
}
