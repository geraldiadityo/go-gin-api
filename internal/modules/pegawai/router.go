package pegawai

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
	if err := r.db.AutoMigrate(&Pegawai{}); err != nil {
		log.Fatalf("Gagal Migrasi table pegawai: %v", err)
	}

	pegawaiRouter := rg.Group("/pegawai")
	{
		pegawaiRouter.GET("", r.handler.GetAll)
		pegawaiRouter.POST("", r.handler.Create)
		pegawaiRouter.GET("/:id", r.handler.GetById)
		pegawaiRouter.PUT("/:id", r.handler.Update)
		pegawaiRouter.DELETE("/:id", r.handler.Delete)
	}
}
