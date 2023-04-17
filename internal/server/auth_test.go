package server

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/KirillMironov/beaver/internal/aes"
	"github.com/KirillMironov/beaver/internal/log/observer"
)

func TestNewAuthenticator(t *testing.T) {
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

				if _, err := os.CreateTemp(path, "file"); err == nil {
					t.Skip("skipping test for root user")
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

			_, err := NewAuthenticator(tc.dataDir, logger)
			if err != nil != tc.wantErr {
				t.Fatalf("NewAuthenticator() error = %v, wantErr %v", err, tc.wantErr)
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

func TestAuthenticator_AddUser(t *testing.T) {
	t.Parallel()

	authenticator, masterKey := newAuthenticator(t)

	tests := []struct {
		name        string
		credentials Credentials
		masterKey   string
		wantErr     error
	}{
		{
			name:        "invalid master key",
			credentials: Credentials{Username: "user", Passphrase: "passphrase"},
			masterKey:   "invalid",
			wantErr:     errInvalidMasterKey,
		},
		{
			name:        "valid user",
			credentials: Credentials{Username: "user", Passphrase: "passphrase"},
			masterKey:   masterKey,
			wantErr:     nil,
		},
		{
			name:        "user already exists",
			credentials: Credentials{Username: "user", Passphrase: "passphrase"},
			masterKey:   masterKey,
			wantErr:     errUserAlreadyExists,
		},
		{
			name:        "empty username",
			credentials: Credentials{Username: "", Passphrase: "passphrase"},
			masterKey:   masterKey,
			wantErr:     errEmptyUsername,
		},
		{
			name:        "empty passphrase",
			credentials: Credentials{Username: "user-2", Passphrase: ""},
			masterKey:   masterKey,
			wantErr:     errEmptyPassphrase,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			user, err := authenticator.AddUser(tc.credentials, tc.masterKey)
			if err != tc.wantErr {
				t.Fatalf("AddUser() error = %v, wantErr %v", err, tc.wantErr)
			}

			if tc.wantErr != nil {
				return
			}

			if user.Username != tc.credentials.Username {
				t.Fatalf("AddUser() username = %v, want %v", user.Username, tc.credentials.Username)
			}

			if _, err = os.Stat(user.DataDir); err != nil {
				t.Fatalf("AddUser() user data dir does not exist: %v", err)
			}
		})
	}
}

func TestAuthenticator_Authenticate(t *testing.T) {
	t.Parallel()

	authenticator, masterKey := newAuthenticator(t)

	_, err := authenticator.AddUser(Credentials{Username: "user", Passphrase: "passphrase"}, masterKey)
	if err != nil {
		t.Fatalf("AddUser() error = %v", err)
	}

	tests := []struct {
		name        string
		credentials Credentials
		wantErr     bool
	}{
		{
			name:        "user exists",
			credentials: Credentials{Username: "user", Passphrase: "passphrase"},
			wantErr:     false,
		},
		{
			name:        "user not found",
			credentials: Credentials{Username: "user-2", Passphrase: "passphrase"},
			wantErr:     true,
		},
		{
			name:        "invalid passphrase",
			credentials: Credentials{Username: "user", Passphrase: "invalid"},
			wantErr:     true,
		},
		{
			name:        "empty username",
			credentials: Credentials{Username: "", Passphrase: "passphrase"},
			wantErr:     true,
		},
		{
			name:        "empty passphrase",
			credentials: Credentials{Username: "user", Passphrase: ""},
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			user, err := authenticator.Authenticate(tc.credentials)
			if err != nil != tc.wantErr {
				t.Fatalf("Authenticate() error = %v, wantErr %v", err, tc.wantErr)
			}

			if tc.wantErr {
				return
			}

			if user.Username != tc.credentials.Username {
				t.Fatalf("Authenticate() username = %v, want %v", user.Username, tc.credentials.Username)
			}
		})
	}
}

func newAuthenticator(t *testing.T) (authenticator *Authenticator, masterKey string) {
	t.Helper()

	dataDir := t.TempDir()

	logger := observer.New()

	authenticator, err := NewAuthenticator(dataDir, logger)
	if err != nil {
		t.Fatal(err)
	}

	log := logger.First()
	lastIndex := strings.LastIndex(log, " ")
	masterKey = strings.Trim(log[lastIndex+1:], `"`)

	return authenticator, masterKey
}
