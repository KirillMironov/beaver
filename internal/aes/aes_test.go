package aes

import (
	"bytes"
	"testing"
)

func TestEncrypt(t *testing.T) {
	t.Parallel()

	plaintext := []byte("hello world")
	key := []byte("super secret key")

	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatal(err)
	}

	if len(ciphertext) == 0 {
		t.Fatal("expected ciphertext to be non-empty")
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatalf("expected ciphertext to be different from plaintext, got %q", ciphertext)
	}

	ciphertext2, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(ciphertext, ciphertext2) {
		t.Fatalf("expected ciphertext to be different from ciphertext2, got %q, %q", ciphertext, ciphertext2)
	}
}

func TestDecrypt(t *testing.T) {
	t.Parallel()

	plaintext := []byte("hello world")
	key := []byte("super secret key")

	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatal(err)
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatalf("expected %q, got %q", plaintext, decrypted)
	}
}
