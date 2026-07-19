package app

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ngothanhtung/go-tutorials/internal/common/middleware"
	"github.com/ngothanhtung/go-tutorials/internal/common/storage"
	"github.com/ngothanhtung/go-tutorials/internal/db"
	"github.com/ngothanhtung/go-tutorials/internal/features/auth"
	"github.com/ngothanhtung/go-tutorials/internal/features/cart"
	"github.com/ngothanhtung/go-tutorials/internal/features/health"
	"github.com/ngothanhtung/go-tutorials/internal/features/orders"
	"github.com/ngothanhtung/go-tutorials/internal/features/rbac"
	"github.com/ngothanhtung/go-tutorials/internal/features/uploads"
	"github.com/ngothanhtung/go-tutorials/internal/features/user"
	"github.com/ngothanhtung/go-tutorials/internal/features/wishlist"
)

func (a *App) buildRouter() *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS(a.Cfg.CORS.Origins))
	r.Use(middleware.Logger(a.Log))
	r.Use(middleware.Recovery(a.Log))
	r.Use(middleware.RateLimit(a.Redis, a.Cfg.Rate.PerMin))
	r.Use(middleware.Audit(a.Log, a.DB, a.Cfg.Audit.Enabled))

	authM := middleware.Auth(a.JWT)
	adminM := rbac.AdminGuard()
	r.Static("/static", a.Cfg.Upload.Dir)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")
	// Health endpoints live at the root (per spec), not under /api/v1.
	// Each readyz call must use a fresh context (handlers run asynchronously
	// long after buildRouter returns).
	health.Register(r, health.NewHandler(
		func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			return db.Ping(ctx, a.DB)
		},
		func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			return db.PingRedis(ctx, a.Redis)
		},
	))

	store := storage.NewLocal(a.Cfg.Upload)
	authRepo := auth.NewRepository(a.DB)
	authSvc := auth.NewService(authRepo, a.JWT, a.Redis, a.Cfg.JWT)
	auth.Register(v1, auth.NewHandler(authSvc))

	userRepo := user.NewRepository(a.DB)
	userSvc := user.NewService(userRepo)
	user.Register(v1, user.NewHandler(userSvc), authM)

	uploads.Register(v1, uploads.NewHandler(store), authM)

	// Storefront (auth-protected): cart, wishlist, orders.
	cartRepo := cart.NewRepository(a.DB)
	cartSvc := cart.NewService(cartRepo)
	cart.Register(v1, cart.NewHandler(cartSvc), authM)

	wishlistRepo := wishlist.NewRepository(a.DB)
	wishlistSvc := wishlist.NewService(wishlistRepo)
	wishlist.Register(v1, wishlist.NewHandler(wishlistSvc), authM)

	ordersRepo := orders.NewRepository(a.DB)
	ordersSvc := orders.NewService(ordersRepo)
	orders.Register(v1, orders.NewHandler(ordersSvc), authM, adminM)

	return r
}
