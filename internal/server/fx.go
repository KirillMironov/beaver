package server

import (
	"go.uber.org/fx"

	"github.com/KirillMironov/beaver/internal/server/config"
)

var Module = fx.Module(
	"server",
	fx.Provide(
		func(cfg config.Config) string {
			return cfg.DataDir
		},
		fx.Annotate(NewAuthenticator, fx.As(new(authenticator))),
		NewStorage,
	),
)
