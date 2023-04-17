package transport

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/KirillMironov/beaver/internal/log"
	"github.com/KirillMironov/beaver/internal/server"
	"github.com/KirillMironov/beaver/internal/server/transport/proto"
)

type AuthenticatorService struct {
	authenticator Authenticator
	logger        log.Logger
}

type Authenticator interface {
	AddUser(credentials server.Credentials, masterKey string) (server.User, error)
}

func NewAuthenticatorService(authenticator Authenticator, logger log.Logger) *AuthenticatorService {
	return &AuthenticatorService{
		authenticator: authenticator,
		logger:        logger,
	}
}

func (a AuthenticatorService) AddUser(_ context.Context, request *proto.AddUserRequest) (*proto.Response, error) {
	credentials := server.Credentials{
		Username:   request.GetUsername(),
		Passphrase: request.GetPassphrase(),
	}

	_, err := a.authenticator.AddUser(credentials, request.GetMasterKey())
	if err != nil {
		a.logger.Errorf("failed to add user: %v", err)
		return nil, status.Error(codes.Internal, "")
	}

	return &proto.Response{}, nil
}
