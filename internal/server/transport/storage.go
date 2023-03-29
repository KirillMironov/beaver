package transport

import (
	"io"

	"github.com/KirillMironov/beaver/internal/server/transport/proto"
)

type StorageService struct {
	storage Storage
}

type Storage interface {
	Upload(username, passphrase string, filename string, src io.Reader) error
	Download(username, passphrase, filename string, dst io.Writer) error
}

func NewStorageService(storage Storage) *StorageService {
	return &StorageService{storage: storage}
}

func (ss StorageService) Upload(server proto.Storage_UploadServer) error {
	return nil
}

func (ss StorageService) Download(request *proto.FileRequest, server proto.Storage_DownloadServer) error {
	return nil
}
