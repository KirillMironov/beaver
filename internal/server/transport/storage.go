package transport

import (
	"io"

	"github.com/KirillMironov/beaver/internal/grpcutil"
	"github.com/KirillMironov/beaver/internal/server/transport/proto"
)

type StorageService struct {
	storage Storage
}

type Storage interface {
	Upload(username, passphrase, filename string, src io.Reader) error
	Download(username, passphrase, filename string, dst io.Writer) error
}

func NewStorageService(storage Storage) *StorageService {
	return &StorageService{storage: storage}
}

func (s StorageService) Upload(request *proto.UploadRequest, stream proto.Storage_UploadServer) error {
	user := request.User

	reader := grpcutil.StreamToReader[*proto.File](stream.Context(), stream)

	return s.storage.Upload(user.Username, user.Passphrase, request.Filename, reader)
}

func (s StorageService) Download(request *proto.DownloadRequest, stream proto.Storage_DownloadServer) error {
	user := request.User

	writer := grpcutil.StreamToWriter[*proto.File](stream.Context(), stream)

	return s.storage.Download(user.Username, user.Passphrase, request.Filename, writer)
}
