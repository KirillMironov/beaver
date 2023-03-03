package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/KirillMironov/beaver/pkg/log"
)

const (
	beaverKeyFile = "beaver.key"
	message       = "beaver"
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

	file, err := os.Create(filepath.Join(dataDir, beaverKeyFile))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	secretKey, err := writeSecretKey(file)
	if err != nil {
		return nil, err
	}

	logger.Infof("secret key: %s", secretKey)

	return &Service{
		dataDir: dataDir,
		logger:  logger,
	}, nil
}

func (s Service) Create(passphrase, secretKey string) (User, error) {
	if err := s.validateSecretKey(secretKey); err != nil {
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

func (s Service) validateSecretKey(key string) error {
	file, err := os.Open(filepath.Join(s.dataDir, beaverKeyFile))
	if err != nil {
		return err
	}
	defer file.Close()

	secretKey, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	data, err := encrypt([]byte(message), []byte(key))
	if err != nil {
		return err
	}

	if data != string(secretKey) {
		return errors.New("invalid secret key")
	}

	return nil
}

func writeSecretKey(dst io.Writer) (string, error) {
	key := make([]byte, 32)
	plaintext := []byte(message)

	if _, err := rand.Read(key); err != nil {
		return "", err
	}

	data, err := encrypt(plaintext, key)
	if err != nil {
		return "", err
	}

	if _, err = dst.Write([]byte(data)); err != nil {
		return "", err
	}

	return data, err
}

func encrypt(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]

	if _, err = rand.Read(iv); err != nil {
		return "", err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return hex.EncodeToString(ciphertext), nil
}
