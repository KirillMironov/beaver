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

func TestGenerateKey(t *testing.T) {
	t.Parallel()

	const length = 32

	key, err := GenerateKey(length)
	if err != nil {
		t.Fatal(err)
	}

	if len(key) != length {
		t.Fatalf("expected key length to be %d, got %d", length, len(key))
	}

	key2, err := GenerateKey(length)
	if err != nil {
		t.Fatal(err)
	}

	if len(key2) != length {
		t.Fatalf("expected key2 length to be %d, got %d", length, len(key2))
	}

	if bytes.Equal(key, key2) {
		t.Fatalf("expected key to be different from key2, got %q, %q", key, key2)
	}
}
