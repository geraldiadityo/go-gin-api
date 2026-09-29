package pegawai

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(rg *gin.RouterGroup, db *gorm.DB) {
	if err := db.AutoMigrate(&Pegawai{}); err != nil {
		log.Fatalf("Gagal migrasi table pegawai: %v", err)
	}

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	pegawaiRouter := rg.Group("/pegawai")
	{
		pegawaiRouter.GET("", handler.GetAll)
		pegawaiRouter.POST("", handler.Create)
		pegawaiRouter.GET("/:id", handler.GetById)
		pegawaiRouter.PUT("/:id", handler.Update)
		pegawaiRouter.DELETE("/:id", handler.Delete)
	}
}
