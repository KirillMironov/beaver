package server

import (
	"strings"
	"testing"
)

const (
	fileName    = "test.txt"
	file2Name   = "test2.txt"
	fileContent = "test content"
)

func TestStorage_UploadDownload(t *testing.T) {
	t.Parallel()

	storage := NewStorage()

	user := User{
		Username: "user",
		DataDir:  t.TempDir(),
		key:      deriveKey("key", "salt"),
	}

	if err := storage.Upload(user, fileName, strings.NewReader(fileContent)); err != nil {
		t.Fatal(err)
	}

	if err := storage.Upload(user, fileName, strings.NewReader(fileContent)); err == nil {
		t.Fatalf("got nil, want error on file already exists")
	}

	dst := &strings.Builder{}

	if err := storage.Download(user, fileName, dst); err != nil {
		t.Fatal(err)
	}

	if got, want := dst.String(), fileContent; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestStorage_List(t *testing.T) {
	t.Parallel()

	storage := NewStorage()

	user := User{
		Username: "user",
		DataDir:  t.TempDir(),
		key:      deriveKey("key", "salt"),
	}

	filenames, err := storage.List(user)
	if err != nil {
		t.Fatal(err)
	}

	if got, want := len(filenames), 0; got != want {
		t.Fatalf("got %d, want %d", got, want)
	}

	for _, v := range []string{fileName, file2Name} {
		if err = storage.Upload(user, v, strings.NewReader(fileContent)); err != nil {
			t.Fatal(err)
		}
	}

	filenames, err = storage.List(user)
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
