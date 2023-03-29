package server

import (
	"strings"
	"testing"

	"github.com/KirillMironov/beaver/internal/aes"
	"github.com/KirillMironov/beaver/internal/log/observer"
	"github.com/KirillMironov/beaver/internal/rand"
)

const (
	username   = "user"
	passphrase = "passphrase"

	fileName    = "test.txt"
	file2Name   = "test2.txt"
	fileContent = "test content"
)

func TestStorage_UploadDownload(t *testing.T) {
	t.Parallel()

	authenticator := newAuthenticatorMock(t, t.TempDir())

	storage := NewStorage(authenticator, observer.New())

	if err := storage.Upload(username, passphrase, fileName, strings.NewReader(fileContent)); err != nil {
		t.Fatal(err)
	}

	if err := storage.Upload(username, passphrase, fileName, strings.NewReader(fileContent)); err == nil {
		t.Fatalf("got nil, want error on file already exists")
	}

	dst := &strings.Builder{}

	if err := storage.Download(username, passphrase, fileName, dst); err != nil {
		t.Fatal(err)
	}

	if got, want := dst.String(), fileContent; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestStorage_List(t *testing.T) {
	t.Parallel()

	authenticator := newAuthenticatorMock(t, t.TempDir())

	storage := NewStorage(authenticator, observer.New())

	filenames, err := storage.List(username, passphrase)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(filenames), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	for _, v := range []string{fileName, file2Name} {
		if err = storage.Upload(username, passphrase, v, strings.NewReader(fileContent)); err != nil {
			t.Fatal(err)
		}
	}

	filenames, err = storage.List(username, passphrase)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(filenames), 2; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	if got, want := filenames[0], fileName; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	if got, want := filenames[1], file2Name; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

type authenticatorMock struct {
	dataDir string
	key     []byte
}

func newAuthenticatorMock(t *testing.T, dataDir string) authenticator {
	key, err := rand.Key(aes.KeyLength)
	if err != nil {
		t.Fatal(err)
	}

	return authenticatorMock{
		dataDir: dataDir,
		key:     key,
	}
}

func (am authenticatorMock) Authenticate(_, _ string) (User, error) {
	return User{
		Username: "mock",
		DataDir:  am.dataDir,
		key:      am.key,
	}, nil
}
