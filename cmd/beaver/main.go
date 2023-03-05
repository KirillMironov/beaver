package main

import (
	"go.uber.org/fx"

	"github.com/KirillMironov/beaver/internal/auth"
	"github.com/KirillMironov/beaver/pkg/log"
)

func main() {
	fx.New(options()).Run()
}

func options() fx.Option {
	return fx.Options(
		fx.Provide(
			log.New,
			auth.NewService,
		),
		fx.NopLogger,
	)
}
