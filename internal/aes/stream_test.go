package aes

import (
	"strings"
	"testing"
)

func TestEncrypterDecrypter(t *testing.T) {
	t.Parallel()

	const (
		message = "message to encrypt"
		key     = "super secret key"
	)

	var (
		ciphertext = &strings.Builder{}
		plaintext  = &strings.Builder{}
	)

	encrypter := NewEncrypter(strings.NewReader(message), ciphertext)
	if err := encrypter.Encrypt([]byte(key)); err != nil {
		t.Fatal(err)
	}

	decrypter := NewDecrypter(strings.NewReader(ciphertext.String()), plaintext)
	if err := decrypter.Decrypt([]byte(key)); err != nil {
		t.Fatal(err)
	}

	if got := plaintext.String(); got != message {
		t.Fatalf("got %q, want %q", got, message)
	}
}
