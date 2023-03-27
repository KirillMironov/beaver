package storage

import (
	"io"
	"os"
	"path/filepath"

	"github.com/KirillMironov/beaver/internal/aes"
	"github.com/KirillMironov/beaver/internal/log"
	"github.com/KirillMironov/beaver/internal/server/auth"
)

type Storage struct {
	authenticator Authenticator
	logger        log.Logger
}

type (
	Authenticator interface {
		Authenticate(username, passphrase string) (auth.User, error)
	}

	File interface {
		io.Reader
		Name() string
	}
)

func NewStorage(authenticator Authenticator, logger log.Logger) *Storage {
	return &Storage{
		authenticator: authenticator,
		logger:        logger,
	}
}

func (s Storage) Upload(username, passphrase string, file File) error {
	user, err := s.authenticator.Authenticate(username, passphrase)
	if err != nil {
		return err
	}

	path := filepath.Join(user.DataDir, file.Name())

	dst, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer dst.Close()

	encrypter := aes.NewEncrypter(file, dst)

	return encrypter.Encrypt(user.Key)
}

func (s Storage) Download(username, passphrase, filename string, dst io.Writer) error {
	user, err := s.authenticator.Authenticate(username, passphrase)
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

	return decrypter.Decrypt(user.Key)
}

func (s Storage) List(username, passphrase string) ([]string, error) {
	user, err := s.authenticator.Authenticate(username, passphrase)
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
