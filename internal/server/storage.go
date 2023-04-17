package server

import (
	"io"
	"os"
	"path/filepath"

	"github.com/KirillMironov/beaver/internal/aes"
)

type Storage struct {
	authenticator authenticator
}

type authenticator interface {
	Authenticate(Credentials) (User, error)
}

func NewStorage(authenticator authenticator) *Storage {
	return &Storage{authenticator: authenticator}
}

func (s Storage) Upload(credential Credentials, filename string, src io.Reader) error {
	user, err := s.authenticator.Authenticate(credential)
	if err != nil {
		return err
	}

	path := filepath.Join(user.DataDir, filename)

	dst, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer dst.Close()

	encrypter := aes.NewEncrypter(src, dst)

	return encrypter.Encrypt(user.Key())
}

func (s Storage) Download(credential Credentials, filename string, dst io.Writer) error {
	user, err := s.authenticator.Authenticate(credential)
	if err != nil {
		return err
	}

	path := filepath.Join(user.DataDir, filename)

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decrypter := aes.NewDecrypter(file, dst)

	return decrypter.Decrypt(user.Key())
}

func (s Storage) List(credential Credentials) ([]string, error) {
	user, err := s.authenticator.Authenticate(credential)
	if err != nil {
		return nil, err
	}

	dirEntries, err := os.ReadDir(user.DataDir)
	if err != nil {
		return nil, err
	}

	filenames := make([]string, 0, len(dirEntries))

	for _, entry := range dirEntries {
		if !entry.IsDir() {
			filenames = append(filenames, entry.Name())
		}
	}

	return filenames, nil
}
