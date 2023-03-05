package auth

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/KirillMironov/beaver/pkg/aes"
	"github.com/KirillMironov/beaver/pkg/log"
)

const (
	beaverFilename = "beaver.key"
	secretMessage  = "beaver"
)

type Service struct {
	dataDir string
	logger  log.Logger
}

func NewService(dataDir string, logger log.Logger) (*Service, error) {
	entries, err := os.ReadDir(dataDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	if len(entries) > 0 {
		return nil, fmt.Errorf("data directory %q is not empty", dataDir)
	}

	if err = os.MkdirAll(dataDir, 0700); err != nil {
		return nil, err
	}

	file, err := os.Create(filepath.Join(dataDir, beaverFilename))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	key, err := aes.GenerateKey(aes.KeyLength)
	if err != nil {
		return nil, err
	}

	ciphertext, err := aes.Encrypt([]byte(secretMessage), key)
	if err != nil {
		return nil, err
	}

	if _, err = file.Write(ciphertext); err != nil {
		return nil, err
	}

	logger.Infof("secret key: %s", key)

	return &Service{
		dataDir: dataDir,
		logger:  logger,
	}, nil
}

func (s Service) Create(passphrase, key string) (User, error) {
	if err := s.validateKey(key); err != nil {
		return User{}, err
	}

	now := time.Now().UnixNano()
	id := strconv.FormatInt(now, 10)

	user := User{
		ID:         id,
		Passphrase: passphrase,
	}

	return user, nil
}

func (s Service) validateKey(key string) error {
	file, err := os.Open(filepath.Join(s.dataDir, beaverFilename))
	if err != nil {
		return err
	}
	defer file.Close()

	fileCiphertext, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	ciphertext, err := aes.Encrypt([]byte(secretMessage), []byte(key))
	if err != nil {
		return err
	}

	if !bytes.Equal(fileCiphertext, ciphertext) {
		return errors.New("invalid secret key")
	}

	return nil
}
