package main

import (
	"github.com/geraldiadityo/go-backend/internal/config"
	"github.com/geraldiadityo/go-backend/internal/middleware"
	"github.com/geraldiadityo/go-backend/internal/modules/auth"
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
	AuthRouter    *auth.Router
	Config        *config.AppConfig
	DB            *gorm.DB
}

func NewServer(
	pr *pegawai.Router,
	rr *role.Router,
	ur *users.Router,
	ar *auth.Router,
	cfg *config.AppConfig,
	db *gorm.DB,
) *Server {
	return &Server{
		PegawaiRouter: pr,
		RoleRouter:    rr,
		UserRouter:    ur,
		AuthRouter:    ar,
		Config:        cfg,
		DB:            db,
	}
}

func (s *Server) SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	
	// Rute Publik
	s.AuthRouter.SetupPublic(api)

	// Rute Terproteksi
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(s.Config))
	{
		s.AuthRouter.SetupProtected(protected)
		s.PegawaiRouter.Setup(protected)
		s.RoleRouter.Setup(protected)
		s.UserRouter.Setup(protected)
	}
}
