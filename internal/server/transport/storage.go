package transport

import (
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/KirillMironov/beaver/internal/grpcutil"
	"github.com/KirillMironov/beaver/internal/log"
	"github.com/KirillMironov/beaver/internal/server/transport/proto"
)

type StorageService struct {
	storage Storage
	logger  log.Logger
}

type Storage interface {
	Upload(username, passphrase, filename string, src io.Reader) error
	Download(username, passphrase, filename string, dst io.Writer) error
}

func NewStorageService(storage Storage, logger log.Logger) *StorageService {
	return &StorageService{
		storage: storage,
		logger:  logger,
	}
}

func (s StorageService) Upload(request *proto.UploadRequest, stream proto.Storage_UploadServer) error {
	user := request.GetUser()

	reader := grpcutil.StreamToReader[*proto.File](stream.Context(), stream)

	if err := s.storage.Upload(user.GetUsername(), user.GetPassphrase(), request.GetFilename(), reader); err != nil {
		s.logger.Errorf("failed to upload file: %v", err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}

func (s StorageService) Download(request *proto.DownloadRequest, stream proto.Storage_DownloadServer) error {
	user := request.GetUser()

	writer := grpcutil.StreamToWriter[*proto.File](stream.Context(), stream)

	if err := s.storage.Download(user.GetUsername(), user.GetPassphrase(), request.GetFilename(), writer); err != nil {
		s.logger.Errorf("failed to download file: %v", err)
		return status.Error(codes.Internal, err.Error())
	}

	return nil
}
