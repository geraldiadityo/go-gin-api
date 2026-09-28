package pegawai

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(rg *gin.RouterGroup, db *gorm.DB) {
	if err := db.AutoMigrate(&Pegawai{}); err
}