package main

import (
	"go.uber.org/fx"

	"github.com/KirillMironov/beaver/internal/log"
	"github.com/KirillMironov/beaver/internal/server"
	"github.com/KirillMironov/beaver/internal/server/config"
	"github.com/KirillMironov/beaver/internal/server/transport"
	"github.com/KirillMironov/beaver/internal/server/transport/proto"
)

func main() {
	fx.New(options()).Run()
}

func options() fx.Option {
	return fx.Options(
		server.Module,
		fx.Provide(
			config.Load,
			fx.Annotate(log.New, fx.As(new(log.Logger))),
			fx.Annotate(server.NewStorage, fx.As(new(transport.Storage))),
			fx.Annotate(transport.NewStorageService, fx.As(new(proto.StorageServer))),
		),
		fx.Invoke(
			startServer,
		),
	)
}

func startServer(server proto.StorageServer) {
	println(server)
}
