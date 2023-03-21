package auth

import (
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
			name:    "valid data dir",
			dataDir: t.TempDir(),
			wantErr: false,
		},
		{
			name:    "invalid data dir",
			dataDir: "",
			wantErr: true,
		},
		{
			name: "not empty data dir",
			dataDir: func() string {
				t.Helper()

				path := t.TempDir()

				file, err := os.Create(filepath.Join(path, "file"))
				if err != nil {
					t.Fatal(err)
				}
				defer file.Close()

				return path
			}(),
			wantErr: true,
		},
		{
			name: "data dir is not writable",
			dataDir: func() string {
				t.Helper()

				path := filepath.Join(t.TempDir(), "data")

				if err := os.MkdirAll(path, 0000); err != nil {
					t.Fatal(err)
				}

				return path
			}(),
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

			log := logger.First()
			lastIndex := strings.LastIndex(log, " ")
			masterKey := strings.Trim(log[lastIndex+1:], `"`)

			ciphertext, err := os.ReadFile(filepath.Join(tc.dataDir, beaverFilename))
			if err != nil {
				t.Fatal(err)
			}

			plaintext, err := aes.Decrypt(ciphertext, []byte(masterKey))
			if err != nil {
				t.Fatal(err)
			}

			if string(plaintext) != authMessage {
				t.Fatal("master key is invalid")
			}
		})
	}
}
