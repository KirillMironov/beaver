package main

import (
	"go.uber.org/fx"

	"github.com/KirillMironov/beaver/pkg/logger"
)

func main() {
	fx.New(options()).Run()
}

func options() fx.Option {
	return fx.Options(
		fx.Provide(
			logger.New,
		),
	)
}
