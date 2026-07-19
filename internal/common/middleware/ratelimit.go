package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
	ginlimiter "github.com/ulule/limiter/v3/drivers/middleware/gin"
)

// RateLimit returns a gin middleware limiting per-IP requests to perMin per minute.
// Redis is required for distributed state; the limiter store will panic if Redis is unreachable at startup.
func RateLimit(rdb *redis.Client, perMin int) gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  int64(perMin),
	}
	store, err := redisstore.NewStore(rdb)
	if err != nil {
		panic(err)
	}
	instance := limiter.New(store, rate)
	return ginlimiter.NewMiddleware(instance)
}
