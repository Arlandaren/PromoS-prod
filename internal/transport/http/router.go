package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"solution/internal/transport/api/v1/b2b"
	"solution/internal/transport/api/v1/b2c"
)

type Router interface {
	RouteInit()
	SetContext(ctx context.Context)
	ContextMiddleware(c *gin.Context)
}

type MainRouter struct {
	router     *gin.Engine
	ctx        context.Context
	b2bHandler b2b.BusinessHandler
	b2cHandler b2c.UserHandler
}

func NewRouter() *MainRouter {
	router := &MainRouter{
		router:     gin.Default(),
		b2bHandler: b2b.NewHandler(),
		b2cHandler: b2c.NewHandler(),
	}

	return router
}

func (r *MainRouter) RouteInit() {
	r.router.Use(r.ContextMiddleware)

	r.router.GET("api/ping", func(c *gin.Context) { c.String(200, "pong") })

	r.b2bHandler.Route(r.router)
	r.b2cHandler.Route(r.router)
}

func (r *MainRouter) SetContext(ctx context.Context) {
	r.ctx = ctx
}

func (r *MainRouter) ContextMiddleware(c *gin.Context) {
	c.Set("context", r.ctx)
	c.Next()
}
