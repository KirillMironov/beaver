package jwt

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	testSecret  = "secret"
	testPayload = "payload"
)

func TestManager(t *testing.T) {
	t.Parallel()

	manager := NewManager[string](testSecret, time.Minute)

	token, err := manager.GenerateToken(testPayload)
	if err != nil {
		t.Fatal(err)
	}

	if token == "" {
		t.Fatal("got empty string, want token")
	}

	payload, err := manager.ValidateToken(token)
	if err != nil {
		t.Fatal(err)
	}

	if payload != testPayload {
		t.Fatalf("got %q, want %q", payload, testPayload)
	}
}

func TestManager_ValidateToken_Expired(t *testing.T) {
	t.Parallel()

	manager := NewManager[string](testSecret, time.Nanosecond)

	token, err := manager.GenerateToken(testPayload)
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Nanosecond)

	payload, err := manager.ValidateToken(token)
	if want := jwt.ErrTokenExpired; !errors.Is(err, want) {
		t.Fatalf("got %v, want %v", err, want)
	}

	if payload != "" {
		t.Fatalf("got %q, want empty string", payload)
	}
}
