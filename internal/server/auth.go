package server

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
	errInvalidMasterKey  = errors.New("invalid master key")
	errInvalidPassphrase = errors.New("invalid passphrase")
	errUserAlreadyExists = errors.New("user already exists")
	errUserNotFound      = errors.New("user not found")
	errEmptyUsername     = errors.New("username cannot be empty")
	errEmptyPassphrase   = errors.New("passphrase cannot be empty")
)

type User struct {
	Username string
	DataDir  string
	key      []byte
}

func (u User) Key() []byte {
	key := make([]byte, len(u.key))
	copy(key, u.key)
	return key
}

type Authenticator struct {
	dataDir string
	logger  log.Logger
}

func NewAuthenticator(dataDir string, logger log.Logger) (*Authenticator, error) {
	authenticator := &Authenticator{
		dataDir: dataDir,
		logger:  logger,
	}

	return authenticator, authenticator.generateMasterKeyIfNotExists()
}

func (a Authenticator) AddUser(username, passphrase, masterKey string) (User, error) {
	if username == "" {
		return User{}, errEmptyUsername
	}

	if passphrase == "" {
		return User{}, errEmptyPassphrase
	}

	userDataDir := filepath.Join(a.dataDir, username)

	if _, err := os.Stat(userDataDir); err == nil {
		return User{}, errUserAlreadyExists
	}

	if err := a.verifyMasterKey(masterKey); err != nil {
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
		key:      key,
	}, nil
}

func (a Authenticator) Authenticate(username, passphrase string) (User, error) {
	if username == "" {
		return User{}, errEmptyUsername
	}

	if passphrase == "" {
		return User{}, errEmptyPassphrase
	}

	userDataDir := filepath.Join(a.dataDir, username)

	fileCiphertext, err := os.ReadFile(filepath.Join(userDataDir, "."+username))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return User{}, errUserNotFound
		}
		return User{}, err
	}

	key := deriveKey(passphrase, username)

	plaintext, err := aes.Decrypt(fileCiphertext, key)
	if err != nil {
		return User{}, err
	}

	if !bytes.Equal(plaintext, []byte(authMessage)) {
		return User{}, errInvalidPassphrase
	}

	return User{
		Username: username,
		DataDir:  userDataDir,
		key:      key,
	}, nil
}

func (a Authenticator) verifyMasterKey(masterKey string) error {
	if len(masterKey) != aes.KeyLength {
		return errInvalidMasterKey
	}

	path := filepath.Join(a.dataDir, beaverFilename)

	ciphertext, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	plaintext, err := aes.Decrypt(ciphertext, []byte(masterKey))
	if err != nil {
		return err
	}

	if !bytes.Equal(plaintext, []byte(authMessage)) {
		return errInvalidMasterKey
	}

	return nil
}

func (a Authenticator) generateMasterKeyIfNotExists() error {
	dirEntries, err := os.ReadDir(a.dataDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	if len(dirEntries) > 0 {
		for _, entry := range dirEntries {
			if entry.Name() == beaverFilename {
				return nil
			}
		}

		return fmt.Errorf("data directory %q is not empty", a.dataDir)
	}

	if err = os.MkdirAll(a.dataDir, 0700); err != nil {
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

	if err = os.WriteFile(filepath.Join(a.dataDir, beaverFilename), ciphertext, 0400); err != nil {
		_ = os.Remove(a.dataDir)
		return err
	}

	a.logger.Infof("master key: %q", masterKey)

	return err
}

func deriveKey(passphrase, salt string) []byte {
	return pbkdf2.Key([]byte(passphrase), []byte(salt), 10000, aes.KeyLength, sha256.New)
}
