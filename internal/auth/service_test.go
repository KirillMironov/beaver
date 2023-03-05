package auth

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/KirillMironov/beaver/pkg/mock"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	dataDir := t.TempDir()

	_, err := NewService(dataDir, mock.Logger{})
	if err != nil {
		t.Fatal(err)
	}

	stat, err := os.Stat(filepath.Join(dataDir, beaverFilename))
	if err != nil {
		t.Fatal(err)
	}

	if stat.Size() == 0 {
		t.Fatal("expected non-empty file")
	}
}
