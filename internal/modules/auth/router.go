package auth

import "github.com/gin-gonic/gin"

type Router struct {
	handler *Handler
}

func NewRouter(handler *Handler) *Router {
	return &Router{handler: handler}
}

func (r *Router) SetupPublic(rg *gin.RouterGroup) {
	authRouter := rg.Group("/auth")
	{
		authRouter.POST("/login", r.handler.Login)
		authRouter.POST("/refresh", r.handler.RefreshToken)
	}
}

func (r *Router) SetupProtected(rg *gin.RouterGroup) {
	authRouter := rg.Group("/auth")
	{
		authRouter.POST("/logout", r.handler.Logout)
	}
}

