package main

import (
	"github.com/geraldiadityo/go-backend/internal/config"
	"github.com/geraldiadityo/go-backend/internal/modules/pegawai"
	"github.com/geraldiadityo/go-backend/internal/modules/role"
	"github.com/geraldiadityo/go-backend/internal/modules/users"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	PegawaiRouter *pegawai.Router
	RoleRouter    *role.Router
	UserRouter    *users.Router
	Config        *config.AppConfig
	DB            *gorm.DB
}

func NewServer(
	pr *pegawai.Router,
	rr *role.Router,
	ur *users.Router,
	cfg *config.AppConfig,
	db *gorm.DB,
) *Server {
	return &Server{
		PegawaiRouter: pr,
		RoleRouter:    rr,
		UserRouter:    ur,
		Config:        cfg,
		DB:            db,
	}
}

func (s *Server) SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	s.PegawaiRouter.Setup(api)
	s.RoleRouter.Setup(api)
	s.UserRouter.Setup(api)
}
