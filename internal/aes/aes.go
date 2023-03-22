package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"

	"github.com/KirillMironov/beaver/internal/rand"
)

// KeyLength is the recommended length of the AES key.
const KeyLength = 32

// Encrypt encrypts a byte slice using the given key
// and returns the encrypted byte slice in base64-encoded form,
// including the initialization vector (IV) used during encryption.
// The IV is a unique value used in the encryption process to prevent replay attacks.
func Encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	iv, err := rand.Key(gcm.NonceSize())
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, iv, plaintext, nil)

	encoded := make([]byte, base64.StdEncoding.EncodedLen(len(iv)+len(ciphertext)))

	iv = append(iv, ciphertext...)

	base64.StdEncoding.Encode(encoded, iv)

	return encoded, nil
}

// Decrypt decrypts a base64-encoded byte slice using the given key
// and initialization vector (IV) contained in the byte slice.
func Decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(ciphertext)))

	n, err := base64.StdEncoding.Decode(decoded, ciphertext)
	if err != nil {
		return nil, err
	}

	decoded = decoded[:n]

	if len(decoded) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	iv, ciphertext := decoded[:gcm.NonceSize()], decoded[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
