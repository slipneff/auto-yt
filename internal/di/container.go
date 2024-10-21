package di

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/slipneff/auto-yt/internal/pkg/router"
	"github.com/slipneff/auto-yt/internal/pkg/service/jwt"

	"github.com/slipneff/auto-yt/internal/utils/config"
)

type Container struct {
	cfg *config.Config
	ctx context.Context

	handler    *router.Handler
	httpServer *http.Server
	jwtService *jwt.Service
}

func New(ctx context.Context, cfg *config.Config) *Container {
	return &Container{cfg: cfg, ctx: ctx}
}

func (c *Container) GetHttpServer() *http.Server {
	return get(&c.httpServer, func() *http.Server {
		return &http.Server{
			Addr:           fmt.Sprintf("%s:%d", c.cfg.Host, c.cfg.Port),
			Handler:        c.GetHttpHandler().InitRoutes(),
			MaxHeaderBytes: 1 << 20,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
		}
	})
}

func (c *Container) GetJWTService() *jwt.Service {
	return get(&c.jwtService, func() *jwt.Service {
		return jwt.New(c.cfg)
	})
}
func (c *Container) GetHttpHandler() *router.Handler {
	return get(&c.handler, func() *router.Handler {
		return router.NewRouter(c.cfg)
	})
}

func get[T comparable](obj *T, builder func() T) T {
	if *obj != *new(T) {
		return *obj
	}

	*obj = builder()
	return *obj
}
