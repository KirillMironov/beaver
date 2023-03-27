package storage

import (
	"io"
	"strings"
	"testing"

	"github.com/KirillMironov/beaver/internal/aes"
	"github.com/KirillMironov/beaver/internal/log/observer"
	"github.com/KirillMironov/beaver/internal/rand"
	"github.com/KirillMironov/beaver/internal/server/auth"
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

	var (
		authenticator = newAuthenticatorMock(t, t.TempDir())
		file          = newFileMock(fileName, fileContent)
	)

	storage := NewStorage(authenticator, observer.New())

	if err := storage.Upload(username, passphrase, file); err != nil {
		t.Fatal(err)
	}

	if err := storage.Upload(username, passphrase, file); err == nil {
		t.Fatalf("got nil, want error on file already exists")
	}

	dst := &strings.Builder{}

	if err := storage.Download(username, passphrase, file.Name(), dst); err != nil {
		t.Fatal(err)
	}

	if got, want := dst.String(), fileContent; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestStorage_List(t *testing.T) {
	t.Parallel()

	var (
		authenticator = newAuthenticatorMock(t, t.TempDir())
		file          = newFileMock(fileName, fileContent)
		file2         = newFileMock(file2Name, fileContent)
	)

	storage := NewStorage(authenticator, observer.New())

	filenames, err := storage.List(username, passphrase)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(filenames), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	for _, v := range []File{file, file2} {
		if err = storage.Upload(username, passphrase, v); err != nil {
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

	if got, want := filenames[0], file.Name(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}

	if got, want := filenames[1], file2.Name(); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

type fileMock struct {
	name string
	src  io.Reader
}

func newFileMock(name, content string) File {
	return fileMock{
		name: name,
		src:  strings.NewReader(content),
	}
}

func (fm fileMock) Read(p []byte) (n int, err error) {
	return fm.src.Read(p)
}

func (fm fileMock) Name() string {
	return fm.name
}

type authenticatorMock struct {
	dataDir string
	key     []byte
}

func newAuthenticatorMock(t *testing.T, dataDir string) Authenticator {
	key, err := rand.Key(aes.KeyLength)
	if err != nil {
		t.Fatal(err)
	}

	return authenticatorMock{
		dataDir: dataDir,
		key:     key,
	}
}

func (am authenticatorMock) Authenticate(_, _ string) (auth.User, error) {
	return auth.User{
		Username: "mock",
		DataDir:  am.dataDir,
		Key:      am.key,
	}, nil
}
