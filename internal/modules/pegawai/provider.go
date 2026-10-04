package pegawai

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewRepository,
	NewService,
	NewHandler,
	NewRouter,
)
