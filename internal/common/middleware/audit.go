package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Audit writes request metadata to audit_logs unless disabled or path skipped.
func Audit(log *zap.Logger, db *gorm.DB, enabled bool) gin.HandlerFunc {
	skip := map[string]bool{"/healthz": true, "/readyz": true, "/metrics": true}
	return func(c *gin.Context) {
		if !enabled || skip[c.Request.URL.Path] {
			c.Next()
			return
		}
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		cl := Claims(c)
		var uid *string
		if cl != nil && cl.UserID != "" {
			uid = &cl.UserID
		}
		rid, _ := c.Get(RequestIDKey)
		record := map[string]interface{}{
			"user_id":    uid,
			"action":     c.Request.Method,
			"resource":   c.Request.URL.Path,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"latency_ms": latency.Milliseconds(),
			"ip":         c.ClientIP(),
			"request_id": toStr(rid),
			"created_at": time.Now(),
		}
		if db != nil {
			if err := db.Table("audit_logs").Create(record).Error; err != nil {
				log.Warn("audit write failed", zap.Error(err))
			}
		}
	}
}

