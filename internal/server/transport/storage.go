package transport

import (
	"context"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/KirillMironov/beaver/internal/grpcutil"
	"github.com/KirillMironov/beaver/internal/log"
	"github.com/KirillMironov/beaver/internal/server"
	"github.com/KirillMironov/beaver/internal/server/transport/proto"
)

type StorageService struct {
	storage Storage
	logger  log.Logger
}

type Storage interface {
	Upload(credentials server.Credentials, filename string, src io.Reader) error
	Download(credentials server.Credentials, filename string, dst io.Writer) error
	List(credentials server.Credentials) ([]string, error)
}

func NewStorageService(storage Storage, logger log.Logger) *StorageService {
	return &StorageService{
		storage: storage,
		logger:  logger,
	}
}

func (s StorageService) Upload(request *proto.FileRequest, stream proto.Storage_UploadServer) error {
	reader := grpcutil.StreamToReader[*proto.File](stream.Context(), stream)

	credentials := convertCredentials(request.GetCredentials())

	if err := s.storage.Upload(credentials, request.GetFilename(), reader); err != nil {
		s.logger.Errorf("failed to upload file: %v", err)
		return status.Error(codes.Internal, "")
	}

	return nil
}

func (s StorageService) Download(request *proto.FileRequest, stream proto.Storage_DownloadServer) error {
	writer := grpcutil.StreamToWriter[*proto.File](stream.Context(), stream)

	credentials := convertCredentials(request.GetCredentials())

	if err := s.storage.Download(credentials, request.GetFilename(), writer); err != nil {
		s.logger.Errorf("failed to download file: %v", err)
		return status.Error(codes.Internal, "")
	}

	return nil
}

func (s StorageService) List(_ context.Context, request *proto.Credentials) (*proto.ListResponse, error) {
	credentials := convertCredentials(request)

	filenames, err := s.storage.List(credentials)
	if err != nil {
		s.logger.Errorf("failed to list files: %v", err)
		return nil, status.Error(codes.Internal, "")
	}

	return &proto.ListResponse{Filenames: filenames}, nil
}

func convertCredentials(credentials *proto.Credentials) server.Credentials {
	return server.Credentials{
		Username:   credentials.GetUsername(),
		Passphrase: credentials.GetPassphrase(),
	}
}
