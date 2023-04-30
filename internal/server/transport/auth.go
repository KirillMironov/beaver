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
	AddUser(username, passphrase, masterKey string) (token string, err error)
	Authenticate(username, passphrase string) (token string, err error)
	ValidateToken(token string) (server.User, error)
}

func NewAuthenticatorService(authenticator Authenticator, logger log.Logger) *AuthenticatorService {
	return &AuthenticatorService{
		authenticator: authenticator,
		logger:        logger,
	}
}

func (a AuthenticatorService) AddUser(_ context.Context, request *proto.AddUserRequest) (*proto.Token, error) {
	token, err := a.authenticator.AddUser(request.GetUsername(), request.GetPassphrase(), request.GetMasterKey())
	if err != nil {
		a.logger.Errorf("failed to add user: %v", err)
		return nil, status.Error(codes.Internal, "")
	}

	return &proto.Token{Token: token}, nil
}

func (a AuthenticatorService) Authenticate(_ context.Context, request *proto.AuthenticateRequest) (*proto.Token, error) {
	token, err := a.authenticator.Authenticate(request.GetUsername(), request.GetPassphrase())
	if err != nil {
		a.logger.Errorf("failed to authenticate user: %v", err)
		return nil, status.Error(codes.Internal, "")
	}

	return &proto.Token{Token: token}, nil
}
