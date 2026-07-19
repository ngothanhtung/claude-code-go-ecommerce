package app

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ngothanhtung/go-tutorials/internal/config"
	"github.com/ngothanhtung/go-tutorials/internal/db"
	"github.com/ngothanhtung/go-tutorials/pkg/jwt"
	"github.com/ngothanhtung/go-tutorials/pkg/logger"
)

// App holds all shared dependencies.
type App struct {
	Cfg    *config.Config
	Log    *zap.Logger
	DB     *gorm.DB
	Redis  *redis.Client
	JWT    *jwt.Manager
	Engine *gin.Engine
}

// New builds the application: config, logging, db, redis, jwt, router.
func New(cfg *config.Config) (*App, error) {
	log := logger.New(cfg.App.Env)
	gdb, err := db.NewPostgres(cfg.DB)
	if err != nil {
		return nil, err
	}
	rdb := db.NewRedis(cfg.Redis)
	jm := jwt.New(cfg.JWT.Secret, cfg.JWT.AccessTTLMin, cfg.JWT.RefreshTTLHour)

	a := &App{Cfg: cfg, Log: log, DB: gdb, Redis: rdb, JWT: jm}
	a.Engine = a.buildRouter()
	return a, nil
}

// Run starts the HTTP server.
func (a *App) Run() error {
	addr := ":" + strconv.Itoa(a.Cfg.App.Port)
	a.Log.Info("server starting", zap.String("addr", addr))
	return a.Engine.Run(addr)
}
