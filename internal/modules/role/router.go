package role

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(rg *gin.RouterGroup, db *gorm.DB) {
	if err := db.AutoMigrate(&Role{}); err != nil {
		log.Fatalf("Gagal migrasi table role: %v", err)
	}

	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	roleRouter := rg.Group("/role")
	{
		roleRouter.GET("", handler.GetAll)
		roleRouter.POST("", handler.Create)
		roleRouter.GET("/:id", handler.GetById)
		roleRouter.PUT("/:id", handler.Update)
		roleRouter.DELETE("/:id", handler.Delete)
	}
}
