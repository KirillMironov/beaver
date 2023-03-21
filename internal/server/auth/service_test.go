package auth

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/KirillMironov/beaver/internal/aes"
	"github.com/KirillMironov/beaver/internal/observer"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		dataDir string
		wantErr bool
	}{
		{
			name:    "success",
			dataDir: t.TempDir(),
			wantErr: false,
		},
		{
			name:    "error",
			dataDir: "",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			logger := observer.New()

			_, err := NewService(tc.dataDir, logger)
			if err != nil != tc.wantErr {
				t.Fatalf("NewService() error = %v, wantErr %v", err, tc.wantErr)
			}

			if tc.wantErr {
				return
			}

			words := strings.Split(logger.First(), " ")

			key := words[len(words)-1][1 : len(words[len(words)-1])-1]

			secretKey := []byte(key)

			err = verifySecret(filepath.Join(tc.dataDir, beaverFilename), secretKey)
			if err != nil != tc.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func verifySecret(path string, key []byte) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	ciphertext, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	plaintext, err := aes.Decrypt(ciphertext, key)
	if err != nil {
		return err
	}

	if string(plaintext) != authMessage {
		return errors.New("key is invalid")
	}

	return nil
}
