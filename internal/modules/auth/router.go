package auth

import (
	"time"

	"github.com/geraldiadityo/go-backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Router struct {
	handler *Handler
	rdb     *redis.Client
}

func NewRouter(handler *Handler, rdb *redis.Client) *Router {
	return &Router{
		handler: handler,
		rdb:     rdb,
	}
}

func (r *Router) SetupPublic(rg *gin.RouterGroup) {
	authRouter := rg.Group("/auth")
	{
		authRouter.POST("/login", middleware.LoginRateLimiter(r.rdb, 10, 1*time.Minute), r.handler.Login)
		authRouter.POST("/refresh", r.handler.RefreshToken)
	}
}

func (r *Router) SetupProtected(rg *gin.RouterGroup) {
	authRouter := rg.Group("/auth")
	{
		authRouter.POST("/logout", r.handler.Logout)
	}
}
