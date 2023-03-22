package auth

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/crypto/pbkdf2"

	"github.com/KirillMironov/beaver/internal/aes"
	"github.com/KirillMironov/beaver/internal/log"
	"github.com/KirillMironov/beaver/internal/rand"
)

const (
	beaverFilename = "beaver.key"
	authMessage    = "beaver"
)

var (
	ErrInvalidMasterKey  = errors.New("invalid master key")
	ErrInvalidPassphrase = errors.New("invalid passphrase")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	errEmptyUsername     = errors.New("username cannot be empty")
	errEmptyPassphrase   = errors.New("passphrase cannot be empty")
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

func (s Service) AddUser(username, passphrase, masterKey string) (User, error) {
	if username == "" {
		return User{}, errEmptyUsername
	}

	if passphrase == "" {
		return User{}, errEmptyPassphrase
	}

	userDataDir := filepath.Join(s.dataDir, username)

	if _, err := os.Stat(userDataDir); err == nil {
		return User{}, ErrUserAlreadyExists
	}

	if err := s.verifyMasterKey(masterKey); err != nil {
		return User{}, err
	}

	key := deriveKey(passphrase, username)

	ciphertext, err := aes.Encrypt([]byte(authMessage), key)
	if err != nil {
		return User{}, err
	}

	if err = os.MkdirAll(userDataDir, 0700); err != nil {
		return User{}, err
	}

	if err = os.WriteFile(filepath.Join(userDataDir, "."+username), ciphertext, 0400); err != nil {
		_ = os.Remove(userDataDir)
		return User{}, err
	}

	return User{
		Username: username,
		DataDir:  userDataDir,
	}, nil
}

func (s Service) Authenticate(username, passphrase string) (User, error) {
	if username == "" {
		return User{}, errEmptyUsername
	}

	if passphrase == "" {
		return User{}, errEmptyPassphrase
	}

	userDataDir := filepath.Join(s.dataDir, username)

	fileCiphertext, err := os.ReadFile(filepath.Join(userDataDir, "."+username))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	key := deriveKey(passphrase, username)

	plaintext, err := aes.Decrypt(fileCiphertext, key)
	if err != nil {
		return User{}, err
	}

	if !bytes.Equal(plaintext, []byte(authMessage)) {
		return User{}, ErrInvalidPassphrase
	}

	return User{
		Username: username,
		DataDir:  userDataDir,
	}, nil
}

func (s Service) verifyMasterKey(masterKey string) error {
	if len(masterKey) != aes.KeyLength {
		return ErrInvalidMasterKey
	}

	path := filepath.Join(s.dataDir, beaverFilename)

	ciphertext, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	plaintext, err := aes.Decrypt(ciphertext, []byte(masterKey))
	if err != nil {
		return err
	}

	if !bytes.Equal(plaintext, []byte(authMessage)) {
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

	masterKey, err := rand.Key(aes.KeyLength)
	if err != nil {
		return err
	}

	ciphertext, err := aes.Encrypt([]byte(authMessage), masterKey)
	if err != nil {
		return err
	}

	if err = os.WriteFile(filepath.Join(s.dataDir, beaverFilename), ciphertext, 0400); err != nil {
		_ = os.Remove(s.dataDir)
		return err
	}

	s.logger.Infof("master key: %q", masterKey)

	return err
}

func deriveKey(passphrase, salt string) []byte {
	return pbkdf2.Key([]byte(passphrase), []byte(salt), 10000, aes.KeyLength, sha256.New)
}
