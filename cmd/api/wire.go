//go:build wireinject
// +build wireinject

package main

import (
	"github.com/geraldiadityo/go-backend/internal/config"
	"github.com/geraldiadityo/go-backend/internal/modules/auth"
	"github.com/geraldiadityo/go-backend/internal/modules/pegawai"
	"github.com/geraldiadityo/go-backend/internal/modules/role"
	"github.com/geraldiadityo/go-backend/internal/modules/users"
	"github.com/google/wire"
)

func InitializeApp() (*Server, error) {
	wire.Build(
		config.ProviderSet,
		pegawai.ProviderSet,
		role.ProviderSet,
		users.ProviderSet,
		auth.ProviderSet,
		NewServer,
	)

	return nil, nil
}
