package auth

import (
	"testing"

	"github.com/KirillMironov/beaver/pkg/log"
)

func TestService_Create(t *testing.T) {
	logger, err := log.New()
	if err != nil {
		t.Fatal(err)
	}

	service, err := NewService(t.TempDir(), logger)
	if err != nil {
		t.Fatal(err)
	}

	user, err := service.Create("passphrase", "ad")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(user)
}
