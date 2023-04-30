package main

import (
	"context"
	"net"

	"go.uber.org/fx"
	"google.golang.org/grpc"

	"github.com/KirillMironov/beaver/internal/jwt"
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
		fx.Provide(
			config.Load,
			fx.Annotate(
				func(cfg config.Config) *jwt.Manager[server.User] {
					return jwt.NewManager[server.User](cfg.JWT.Secret, cfg.JWT.TokenTTL)
				},
				fx.As(new(jwt.TokenManager[server.User])),
			),
			fx.Annotate(log.New, fx.As(new(log.Logger))),
			fx.Annotate(server.NewStorage, fx.As(new(transport.Storage))),
			fx.Annotate(
				func(cfg config.Config, logger log.Logger, tokenManager jwt.TokenManager[server.User]) (*server.Authenticator, error) {
					return server.NewAuthenticator(cfg.DataDir, logger, tokenManager)
				},
				fx.As(new(transport.Authenticator)),
			),
			fx.Annotate(transport.NewStorageService, fx.As(new(proto.StorageServer))),
			fx.Annotate(transport.NewAuthenticatorService, fx.As(new(proto.AuthenticatorServer))),
		),
		fx.Invoke(
			startServer,
		),
	)
}

func startServer(lifecycle fx.Lifecycle, cfg config.Config, logger log.Logger, storage proto.StorageServer, authenticator proto.AuthenticatorServer) error {
	listener, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	proto.RegisterStorageServer(grpcServer, storage)
	proto.RegisterAuthenticatorServer(grpcServer, authenticator)

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err = grpcServer.Serve(listener); err != nil {
					logger.Errorf("failed to serve: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})

	return nil
}
