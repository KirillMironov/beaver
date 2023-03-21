package auth

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/KirillMironov/beaver/internal/aes"
	"github.com/KirillMironov/beaver/internal/log"
)

const (
	beaverFilename = "beaver.key"
	authMessage    = "beaver"
)

var (
	ErrInvalidMasterKey = errors.New("invalid master key")
	ErrUserExists       = errors.New("user already exists")
	ErrUserNotFound     = errors.New("user not found")
)

type Service struct {
	dataDir string
	logger  log.Logger
}

func NewService(dataDir string, logger log.Logger) (*Service, error) {
	service := &Service{
		dataDir: dataDir,
		logger:  logger,
	}

	return service, service.generateMasterKeyIfNotExists()
}

func (s Service) AddUser(passphrase, masterKey string) (User, error) {
	if err := s.verifyMasterKey(masterKey); err != nil {
		return User{}, err
	}

	ciphertext, err := aes.Encrypt([]byte(authMessage), []byte(passphrase))
	if err != nil {
		return User{}, err
	}

	userID := string(ciphertext)

	userDataDir := filepath.Join(s.dataDir, userID)

	if _, err = os.Stat(userDataDir); err == nil {
		return User{}, ErrUserExists
	}

	if err = os.MkdirAll(userDataDir, 0700); err != nil {
		return User{}, err
	}

	return User{
		ID:      userID,
		DataDir: userDataDir,
	}, nil
}

func (s Service) Authenticate(passphrase string) (User, error) {
	ciphertext, err := aes.Encrypt([]byte(authMessage), []byte(passphrase))
	if err != nil {
		return User{}, err
	}

	userID := string(ciphertext)

	userDataDir := filepath.Join(s.dataDir, userID)

	if _, err = os.Stat(userDataDir); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	return User{
		ID:      userID,
		DataDir: userDataDir,
	}, nil
}

func (s Service) verifyMasterKey(masterKey string) error {
	file, err := os.Open(filepath.Join(s.dataDir, beaverFilename))
	if err != nil {
		return err
	}
	defer file.Close()

	fileCiphertext, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	ciphertext, err := aes.Encrypt([]byte(authMessage), []byte(masterKey))
	if err != nil {
		return err
	}

	if !bytes.Equal(fileCiphertext, ciphertext) {
		return ErrInvalidMasterKey
	}

	return nil
}

func (s Service) generateMasterKeyIfNotExists() error {
	dirEntries, err := os.ReadDir(s.dataDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if len(dirEntries) > 0 {
		for _, entry := range dirEntries {
			if entry.Name() == beaverFilename {
				return nil
			}
		}

		return fmt.Errorf("data directory %q is not empty", s.dataDir)
	}

	if err = os.MkdirAll(s.dataDir, 0700); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(s.dataDir, beaverFilename))
	if err != nil {
		return err
	}
	defer file.Close()

	masterKey, err := aes.GenerateKey(aes.KeyLength)
	if err != nil {
		return err
	}

	ciphertext, err := aes.Encrypt([]byte(authMessage), masterKey)
	if err != nil {
		return err
	}

	if _, err = file.Write(ciphertext); err != nil {
		return err
	}

	s.logger.Infof("master key: %q", masterKey)

	return err
}
