package transport

import (
	"context"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/KirillMironov/beaver/internal/grpcutil"
	"github.com/KirillMironov/beaver/internal/log"
	"github.com/KirillMironov/beaver/internal/server"
	"github.com/KirillMironov/beaver/internal/server/transport/proto"
)

const authorizationHeader = "authorization"

type StorageService struct {
	authenticator Authenticator
	storage       Storage
	logger        log.Logger
}

type Storage interface {
	Upload(user server.User, filename string, src io.Reader) error
	Download(user server.User, filename string, dst io.Writer) error
	List(user server.User) ([]string, error)
}

func NewStorageService(authenticator Authenticator, storage Storage, logger log.Logger) *StorageService {
	return &StorageService{
		authenticator: authenticator,
		storage:       storage,
		logger:        logger,
	}
}

func (s StorageService) Upload(request *proto.FileRequest, stream proto.Storage_UploadServer) error {
	user, err := s.authenticate(stream.Context())
	if err != nil {
		return err
	}

	reader := grpcutil.StreamToReader[*proto.File](stream.Context(), stream)

	if err = s.storage.Upload(user, request.GetFilename(), reader); err != nil {
		s.logger.Errorf("failed to upload file: %v", err)
		return status.Error(codes.Internal, "")
	}

	return nil
}

func (s StorageService) Download(request *proto.FileRequest, stream proto.Storage_DownloadServer) error {
	user, err := s.authenticate(stream.Context())
	if err != nil {
		return err
	}

	writer := grpcutil.StreamToWriter[*proto.File](stream.Context(), stream)

	if err = s.storage.Download(user, request.GetFilename(), writer); err != nil {
		s.logger.Errorf("failed to download file: %v", err)
		return status.Error(codes.Internal, "")
	}

	return nil
}

func (s StorageService) List(ctx context.Context, _ *emptypb.Empty) (*proto.ListResponse, error) {
	user, err := s.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	filenames, err := s.storage.List(user)
	if err != nil {
		s.logger.Errorf("failed to list files: %v", err)
		return nil, status.Error(codes.Internal, "")
	}

	return &proto.ListResponse{Filenames: filenames}, nil
}

func (s StorageService) authenticate(ctx context.Context) (server.User, error) {
	token := grpcutil.HeaderFromContext(ctx, authorizationHeader)
	if token == "" {
		return server.User{}, status.Error(codes.Unauthenticated, `provide token in "authorization" header`)
	}

	user, err := s.authenticator.ValidateToken(token)
	if err != nil {
		return server.User{}, status.Error(codes.Unauthenticated, "invalid token")
	}

	return user, nil
}
