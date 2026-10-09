package users

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Router struct {
	handler *Handler
	db      *gorm.DB
}

func NewRouter(handler *Handler, db *gorm.DB) *Router {
	return &Router{
		handler: handler,
		db:      db,
	}
}

func (r *Router) Setup(rg *gin.RouterGroup) {
	if err := r.db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("Gagal migrasi table user: %v", err)
	}

	userRouter := rg.Group("/user")
	{
		userRouter.POST("", r.handler.Create)
		userRouter.GET("", r.handler.GetAll)
		userRouter.PUT("/:id", r.handler.Update)
		userRouter.GET("/:id", r.handler.GetById)
		userRouter.DELETE("/:id", r.handler.Delete)
		userRouter.PATCH("/:id/status", r.handler.ChangeStatus)
	}
}
