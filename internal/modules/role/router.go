package role

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
	if err := r.db.AutoMigrate(&Role{}); err != nil {
		log.Fatalf("Gagal migrasi table role: %v", err)
	}

	roleRouter := rg.Group("/role")
	{
		roleRouter.GET("", r.handler.GetAll)
		roleRouter.POST("", r.handler.Create)
		roleRouter.GET("/:id", r.handler.GetById)
		roleRouter.PUT("/:id", r.handler.Update)
		roleRouter.DELETE("/:id", r.handler.Update)
	}
}
