package server

import (
	"go.uber.org/fx"

	"github.com/KirillMironov/beaver/internal/log"
	"github.com/KirillMironov/beaver/internal/server/config"
)

func Module() fx.Option {
	return fx.Module(
		"server",
		fx.Provide(
			fx.Annotate(
				func(cfg config.Config, logger log.Logger) (*Authenticator, error) {
					return NewAuthenticator(cfg.DataDir, logger)
				},
				fx.As(new(authenticator)),
			),
			NewStorage,
		),
	)
}
