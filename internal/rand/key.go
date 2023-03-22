package rand

import "crypto/rand"

const chars = `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!@#$%^&*()-_+={}[]|<>,.?;:~`

// Key generates a random key of the given length.
func Key(length int) ([]byte, error) {
	key := make([]byte, length)

	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	for i, b := range key {
		key[i] = chars[b%byte(len(chars))]
	}

	return key, nil
}
